package leanworkflow

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

const (
	assetRoot    = "opencode/lean-workflow"
	promptsDir   = assetRoot + "/prompts"
	templatesDir = assetRoot + "/templates"
)

var promptFiles = []string{
	"lean-workflow-orchestrator.md",
	"requirements-griller.md",
	"planner.md",
	"test-guidance-planner.md",
	"todo-generator.md",
	"implementor.md",
}

var templateFiles = []string{
	"01-requirements.md",
	"02-plan.md",
	"03-test-guidance.md",
	"04-todo.md",
	"05-checkpoint.md",
}

type InjectionResult struct {
	Changed bool
	Files   []string
}

func Inject(homeDir, workspaceDir string, adapter agents.Adapter) (InjectionResult, error) {
	if adapter.Agent() != model.AgentOpenCode {
		return InjectionResult{}, nil
	}

	files := []string{}
	changed := false

	configRoot := filepath.Join(homeDir, ".config", "opencode", "lean-workflow")
	promptRoot := filepath.Join(configRoot, "prompts")
	for _, file := range promptFiles {
		content, err := assets.Read(promptsDir + "/" + file)
		if err != nil {
			return InjectionResult{}, fmt.Errorf("read lean workflow prompt %q: %w", file, err)
		}
		path := filepath.Join(promptRoot, file)
		result, err := filemerge.WriteFileAtomic(path, []byte(content), 0o644)
		if err != nil {
			return InjectionResult{}, fmt.Errorf("write lean workflow prompt %q: %w", path, err)
		}
		changed = changed || result.Changed
		files = append(files, path)
	}

	readme, err := assets.Read(assetRoot + "/README.md")
	if err != nil {
		return InjectionResult{}, fmt.Errorf("read lean workflow README: %w", err)
	}
	readmePath := filepath.Join(configRoot, "README.md")
	readmeResult, err := filemerge.WriteFileAtomic(readmePath, []byte(readme), 0o644)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("write lean workflow README: %w", err)
	}
	changed = changed || readmeResult.Changed
	files = append(files, readmePath)

	if strings.TrimSpace(workspaceDir) != "" {
		templateRoot := filepath.Join(workspaceDir, ".github", "lean-workflow-templates")
		for _, file := range templateFiles {
			content, err := assets.Read(templatesDir + "/" + file)
			if err != nil {
				return InjectionResult{}, fmt.Errorf("read lean workflow template %q: %w", file, err)
			}
			path := filepath.Join(templateRoot, file)
			result, err := writeFileIfMissing(path, []byte(content))
			if err != nil {
				return InjectionResult{}, err
			}
			changed = changed || result.Changed
			files = append(files, path)
		}

		gitignorePath := filepath.Join(workspaceDir, ".gitignore")
		gitignoreChanged, err := ensureGitignoreEntries(gitignorePath, []string{
			".github/plans/",
			".github/lean-workflow-templates/",
		})
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || gitignoreChanged
		files = append(files, gitignorePath)
	}

	settingsPath := adapter.SettingsPath(homeDir)
	overlay, err := renderOpenCodeOverlay(promptRoot)
	if err != nil {
		return InjectionResult{}, err
	}
	settingsResult, err := mergeJSONFile(settingsPath, overlay)
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || settingsResult.Changed
	files = append(files, settingsPath)

	return InjectionResult{Changed: changed, Files: files}, nil
}

func Paths(homeDir, workspaceDir string, adapter agents.Adapter) []string {
	if adapter.Agent() != model.AgentOpenCode {
		return nil
	}
	configRoot := filepath.Join(homeDir, ".config", "opencode", "lean-workflow")
	promptRoot := filepath.Join(configRoot, "prompts")
	paths := []string{filepath.Join(configRoot, "README.md"), adapter.SettingsPath(homeDir)}
	for _, file := range promptFiles {
		paths = append(paths, filepath.Join(promptRoot, file))
	}
	if strings.TrimSpace(workspaceDir) != "" {
		templateRoot := filepath.Join(workspaceDir, ".github", "lean-workflow-templates")
		for _, file := range templateFiles {
			paths = append(paths, filepath.Join(templateRoot, file))
		}
		paths = append(paths, filepath.Join(workspaceDir, ".gitignore"))
	}
	return paths
}

func renderOpenCodeOverlay(promptRoot string) ([]byte, error) {
	promptRef := func(name string) string {
		return "{file:" + filepath.ToSlash(filepath.Join(promptRoot, name+".md")) + "}"
	}
	agent := map[string]any{
		"lean-workflow-orchestrator": map[string]any{
			"mode":        "primary",
			"description": "Lean workflow pilot orchestrator with requirements, plan, acceptance-test, todo, and implementation gates",
			"prompt":      promptRef("lean-workflow-orchestrator"),
			"permission": map[string]any{
				"task": map[string]any{
					"__replace__": map[string]any{
						"*":                     "deny",
						"requirements-griller":  "allow",
						"planner":               "allow",
						"test-guidance-planner": "allow",
						"todo-generator":        "allow",
						"implementor":           "allow",
					},
				},
			},
			"tools": map[string]any{
				"read": true, "write": true, "edit": true, "bash": true,
				"delegate": true, "delegation_read": true, "delegation_list": true,
			},
		},
	}
	for _, subagent := range []struct {
		Name        string
		Description string
	}{
		{"requirements-griller", "Write and revise bounded requirements-grill artifacts"},
		{"planner", "Create confirmed-requirements implementation plans with test inventories"},
		{"test-guidance-planner", "Write optional acceptance-test guidance without writing tests"},
		{"todo-generator", "Create compact implementation todos from confirmed inputs"},
		{"implementor", "Implement confirmed lean workflow todos and verification only"},
	} {
		agent[subagent.Name] = map[string]any{
			"mode":        "subagent",
			"hidden":      true,
			"description": subagent.Description,
			"prompt":      promptRef(subagent.Name),
			"tools": map[string]any{
				"read": true, "write": true, "edit": true, "bash": true,
			},
		}
	}

	root := map[string]any{
		"agent": agent,
		"permission": map[string]any{
			"bash": map[string]any{"*": "ask"},
		},
	}
	return json.MarshalIndent(root, "", "  ")
}

func mergeJSONFile(path string, overlay []byte) (filemerge.WriteResult, error) {
	baseJSON, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return filemerge.WriteResult{}, fmt.Errorf("read OpenCode settings %q: %w", path, err)
		}
		baseJSON = nil
	}
	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return filemerge.WriteResult{}, fmt.Errorf("merge OpenCode settings %q: %w", path, err)
	}
	return filemerge.WriteFileAtomic(path, merged, 0o644)
}

func writeFileIfMissing(path string, content []byte) (filemerge.WriteResult, error) {
	if _, err := os.Lstat(path); err == nil {
		return filemerge.WriteResult{}, nil
	} else if !os.IsNotExist(err) {
		return filemerge.WriteResult{}, fmt.Errorf("stat lean workflow template %q: %w", path, err)
	}
	return filemerge.WriteFileAtomic(path, content, 0o644)
}

func ensureGitignoreEntries(path string, entries []string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read gitignore %q: %w", path, err)
	}

	content := string(data)
	lines := map[string]struct{}{}
	for _, line := range strings.Split(content, "\n") {
		lines[strings.TrimSpace(line)] = struct{}{}
	}

	missing := []string{}
	for _, entry := range entries {
		if _, ok := lines[entry]; !ok {
			missing = append(missing, entry)
		}
	}
	if len(missing) == 0 {
		return false, nil
	}

	var b strings.Builder
	b.WriteString(content)
	if content != "" && !strings.HasSuffix(content, "\n") {
		b.WriteByte('\n')
	}
	if content != "" {
		b.WriteByte('\n')
	}
	b.WriteString("# Gentle AI lean OpenCode workflow pilot\n")
	for _, entry := range missing {
		b.WriteString(entry)
		b.WriteByte('\n')
	}

	result, err := filemerge.WriteFileAtomic(path, []byte(b.String()), 0o644)
	if err != nil {
		return false, err
	}
	return result.Changed, nil
}

func AssetFiles() ([]string, error) {
	files := []string{}
	for _, dir := range []string{promptsDir, templatesDir} {
		err := fs.WalkDir(assets.FS, dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
