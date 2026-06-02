package layeredtdd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

const (
	assetRoot    = "opencode/layered-tdd"
	promptsDir   = assetRoot + "/prompts"
	templatesDir = assetRoot + "/templates"
)

var promptFiles = []string{
	"layered-tdd-orchestrator.md",
	"layered-requirements-griller.md",
	"layer-mapper.md",
	"layer-todo-generator.md",
	"layered-implementor.md",
	"layered-final-reviewer.md",
}

var templateFiles = []string{
	"00-requirements.md",
	"01-layer-map.md",
	"slice-selection.md",
	"layer-todo.md",
	"99-final-review.md",
}

type InjectionResult struct {
	Changed bool
	Files   []string
}

type InjectOptions struct {
	OpenCodeModelAssignments map[string]model.ModelAssignment
}

func Inject(homeDir, workspaceDir string, adapter agents.Adapter, options ...InjectOptions) (InjectionResult, error) {
	if adapter.Agent() != model.AgentOpenCode {
		return InjectionResult{}, nil
	}
	opts := InjectOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	files := []string{}
	changed := false

	configRoot := filepath.Join(homeDir, ".config", "opencode", "layered-tdd")
	promptRoot := filepath.Join(configRoot, "prompts")
	for _, file := range promptFiles {
		content, err := assets.Read(promptsDir + "/" + file)
		if err != nil {
			return InjectionResult{}, fmt.Errorf("read layered TDD prompt %q: %w", file, err)
		}
		path := filepath.Join(promptRoot, file)
		result, err := filemerge.WriteFileAtomic(path, []byte(content), 0o644)
		if err != nil {
			return InjectionResult{}, fmt.Errorf("write layered TDD prompt %q: %w", path, err)
		}
		changed = changed || result.Changed
		files = append(files, path)
	}

	readme, err := assets.Read(assetRoot + "/README.md")
	if err != nil {
		return InjectionResult{}, fmt.Errorf("read layered TDD README: %w", err)
	}
	readmePath := filepath.Join(configRoot, "README.md")
	readmeResult, err := filemerge.WriteFileAtomic(readmePath, []byte(readme), 0o644)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("write layered TDD README: %w", err)
	}
	changed = changed || readmeResult.Changed
	files = append(files, readmePath)

	if strings.TrimSpace(workspaceDir) != "" {
		templateRoot := filepath.Join(workspaceDir, ".github", "layered-tdd-templates")
		for _, file := range templateFiles {
			content, err := assets.Read(templatesDir + "/" + file)
			if err != nil {
				return InjectionResult{}, fmt.Errorf("read layered TDD template %q: %w", file, err)
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
			".github/layered-tdd-templates/",
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
	rootModelID, err := readOpenCodeRootModel(settingsPath)
	if err != nil {
		return InjectionResult{}, err
	}
	existingAgentKeys, err := readExistingAgentKeys(settingsPath)
	if err != nil {
		return InjectionResult{}, err
	}
	if len(opts.OpenCodeModelAssignments) > 0 || rootModelID != "" {
		overlay, err = injectModelAssignments(overlay, opts.OpenCodeModelAssignments, rootModelID, existingAgentKeys)
		if err != nil {
			return InjectionResult{}, err
		}
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
	configRoot := filepath.Join(homeDir, ".config", "opencode", "layered-tdd")
	promptRoot := filepath.Join(configRoot, "prompts")
	paths := []string{filepath.Join(configRoot, "README.md"), adapter.SettingsPath(homeDir)}
	for _, file := range promptFiles {
		paths = append(paths, filepath.Join(promptRoot, file))
	}
	if strings.TrimSpace(workspaceDir) != "" {
		templateRoot := filepath.Join(workspaceDir, ".github", "layered-tdd-templates")
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
		"layered-tdd-orchestrator": map[string]any{
			"mode":        "primary",
			"description": "Layered TDD workflow orchestrator with requirements, slice, layer, red-test, and review gates",
			"prompt":      promptRef("layered-tdd-orchestrator"),
			"permission": map[string]any{
				"task": map[string]any{
					"__replace__": map[string]any{
						"*":                            "deny",
						"layered-requirements-griller": "allow",
						"layer-mapper":                 "allow",
						"layer-todo-generator":         "allow",
						"layered-implementor":          "allow",
						"layered-final-reviewer":       "allow",
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
		{"layered-requirements-griller", "Create or revise requirements and slice-selection artifacts"},
		{"layer-mapper", "Map confirmed slice requirements into human-approved implementation layers"},
		{"layer-todo-generator", "Create skeleton and detailed layer todos with top-level test proposals"},
		{"layered-implementor", "Implement approved layer todos with TDD inside the selected layer boundary"},
		{"layered-final-reviewer", "Review each layer and final slice completion, including waivers and memory candidates"},
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

func injectModelAssignments(overlayBytes []byte, assignments map[string]model.ModelAssignment, rootModelID string, existingAgentKeys map[string]bool) ([]byte, error) {
	var overlay map[string]any
	if err := json.Unmarshal(overlayBytes, &overlay); err != nil {
		return nil, fmt.Errorf("unmarshal layered TDD overlay for model injection: %w", err)
	}

	agentsRaw, ok := overlay["agent"]
	if !ok {
		return overlayBytes, nil
	}
	agents, ok := agentsRaw.(map[string]any)
	if !ok {
		return overlayBytes, nil
	}

	for name, agentDef := range agents {
		agentMap, ok := agentDef.(map[string]any)
		if !ok {
			continue
		}

		assignment, hasExplicitAssignment := assignments[name]
		switch {
		case hasExplicitAssignment && assignment.ProviderID != "" && assignment.ModelID != "":
			agentMap["model"] = assignment.FullID()
			if assignment.Effort != "" {
				agentMap["variant"] = assignment.Effort
			} else {
				agentMap["variant"] = ""
			}
		case existingAgentKeys[name]:
			continue
		case rootModelID != "":
			agentMap["model"] = rootModelID
			agentMap["variant"] = ""
		}
	}

	result, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal layered TDD overlay after model injection: %w", err)
	}
	return append(result, '\n'), nil
}

func readOpenCodeRootModel(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read OpenCode root model from %q: %w", path, err)
	}

	root := map[string]any{}
	if err := json.Unmarshal(data, &root); err != nil {
		return "", nil
	}

	modelID, _ := root["model"].(string)
	return modelID, nil
}

func readExistingAgentKeys(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]bool{}, nil
		}
		return nil, fmt.Errorf("read existing OpenCode agent keys from %q: %w", path, err)
	}

	root := map[string]any{}
	if err := json.Unmarshal(data, &root); err != nil {
		return map[string]bool{}, nil
	}
	agentsRaw, ok := root["agent"]
	if !ok {
		return map[string]bool{}, nil
	}
	agents, ok := agentsRaw.(map[string]any)
	if !ok {
		return map[string]bool{}, nil
	}

	result := make(map[string]bool, len(agents))
	for name := range agents {
		result[name] = true
	}
	return result, nil
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
		return filemerge.WriteResult{}, fmt.Errorf("stat layered TDD template %q: %w", path, err)
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
	b.WriteString("# Gentle AI layered TDD OpenCode workflow\n")
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
