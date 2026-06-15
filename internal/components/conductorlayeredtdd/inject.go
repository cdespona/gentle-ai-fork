package conductorlayeredtdd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

const (
	assetRoot    = "conductor/layered-tdd"
	workflowDir  = "workflows/conductor"
	skillRootDir = ".github/skills"
)

var defaultSkillIDs = []model.SkillID{
	model.SkillTDD,
	model.SkillTestTypeClass,
	model.SkillCognitiveDoc,
	model.SkillWorkUnitCommits,
	model.SkillCommentWriter,
}

var memorySkillIDs = []model.SkillID{
	model.SkillMemoryRecall,
	model.SkillMemoryCapture,
	model.SkillMemoryConsolidate,
}

type InjectOptions struct {
	IncludeMemorySkills bool
}

type InjectionResult struct {
	Changed bool
	Files   []string
}

func Inject(workspaceDir string, opts InjectOptions) (InjectionResult, error) {
	if strings.TrimSpace(workspaceDir) == "" {
		return InjectionResult{}, nil
	}

	files := []string{}
	changed := false

	workflowResult, err := copyWorkflowAssets(workspaceDir)
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || workflowResult.Changed
	files = append(files, workflowResult.Files...)

	skillResult, err := copySkillsIfMissing(filepath.Join(workspaceDir, skillRootDir), skillIDs(opts))
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || skillResult.Changed
	files = append(files, skillResult.Files...)

	gitignorePath := filepath.Join(workspaceDir, ".gitignore")
	gitignoreChanged, err := ensureGitignoreEntries(gitignorePath, []string{
		"workflows/conductor/",
		".github/plans/",
		".github/skills/conductor-*/",
	})
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || gitignoreChanged
	files = append(files, gitignorePath)

	return InjectionResult{Changed: changed, Files: files}, nil
}

func Paths(workspaceDir string, opts InjectOptions) []string {
	if strings.TrimSpace(workspaceDir) == "" {
		return nil
	}

	paths := []string{
		filepath.Join(workspaceDir, workflowDir, "layered-tdd.yaml"),
		filepath.Join(workspaceDir, ".gitignore"),
	}
	for _, file := range []string{
		"final-reviewer.md",
		"implementor.md",
		"layer-mapper.md",
		"layer-reviewer.md",
		"layer-todo-generator.md",
		"memory-capturer.md",
		"requirements-griller.md",
		"slice-run-starter.md",
	} {
		paths = append(paths, filepath.Join(workspaceDir, workflowDir, "prompts", file))
	}
	for _, id := range skillIDs(opts) {
		paths = append(paths, filepath.Join(workspaceDir, skillRootDir, conductorSkillDirName(id), "SKILL.md"))
	}
	return paths
}

func skillIDs(opts InjectOptions) []model.SkillID {
	ids := append([]model.SkillID{}, defaultSkillIDs...)
	if opts.IncludeMemorySkills {
		ids = append(ids, memorySkillIDs...)
	}
	return ids
}

func copyWorkflowAssets(workspaceDir string) (InjectionResult, error) {
	files := []string{}
	changed := false

	err := fs.WalkDir(assets.FS, assetRoot, func(assetPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		content, err := assets.Read(assetPath)
		if err != nil {
			return fmt.Errorf("read Conductor layered TDD asset %q: %w", assetPath, err)
		}
		relPath, err := filepath.Rel(filepath.FromSlash(assetRoot), filepath.FromSlash(assetPath))
		if err != nil {
			return fmt.Errorf("resolve Conductor layered TDD relative path for %q: %w", assetPath, err)
		}
		path := filepath.Join(workspaceDir, workflowDir, relPath)
		result, err := writeFileIfMissing(path, []byte(content))
		if err != nil {
			return err
		}
		changed = changed || result.Changed
		files = append(files, path)
		return nil
	})
	if err != nil {
		return InjectionResult{}, fmt.Errorf("copy Conductor layered TDD workflow assets: %w", err)
	}

	return InjectionResult{Changed: changed, Files: files}, nil
}

func copySkillsIfMissing(skillDir string, skillIDs []model.SkillID) (InjectionResult, error) {
	files := []string{}
	changed := false

	for _, id := range skillIDs {
		embedDir := "skills/" + string(id)
		err := fs.WalkDir(assets.FS, embedDir, func(assetPath string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}

			content, err := assets.Read(assetPath)
			if err != nil {
				return fmt.Errorf("read project skill asset %q: %w", assetPath, err)
			}
			relPath, err := filepath.Rel(filepath.FromSlash(embedDir), filepath.FromSlash(assetPath))
			if err != nil {
				return fmt.Errorf("resolve project skill relative path for %q: %w", assetPath, err)
			}
			path := filepath.Join(skillDir, conductorSkillDirName(id), relPath)
			result, err := writeFileIfMissing(path, []byte(content))
			if err != nil {
				return err
			}
			changed = changed || result.Changed
			files = append(files, path)
			return nil
		})
		if err != nil {
			return InjectionResult{}, fmt.Errorf("copy project skill %q: %w", id, err)
		}
	}

	return InjectionResult{Changed: changed, Files: files}, nil
}

func conductorSkillDirName(id model.SkillID) string {
	return "conductor-" + string(id)
}

func writeFileIfMissing(path string, content []byte) (filemerge.WriteResult, error) {
	if _, err := os.Lstat(path); err == nil {
		return filemerge.WriteResult{}, nil
	} else if !os.IsNotExist(err) {
		return filemerge.WriteResult{}, fmt.Errorf("stat Conductor layered TDD file %q: %w", path, err)
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
	b.WriteString("# Gentle AI Conductor layered TDD workflow\n")
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
