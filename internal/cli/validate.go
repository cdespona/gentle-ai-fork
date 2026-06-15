package cli

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/catalog"
	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

type InstallInput struct {
	Selection model.Selection
	DryRun    bool
}

func NormalizeInstallFlags(flags InstallFlags, detection system.DetectionResult) (InstallInput, error) {
	selection := model.Selection{}

	agents := defaultAgentsFromDetection(detection)
	if len(flags.Agents) > 0 {
		agents = asAgentIDs(flags.Agents)
	}
	selection.Agents = unique(agents)

	persona, err := normalizePersona(flags.Persona)
	if err != nil {
		return InstallInput{}, err
	}
	selection.Persona = persona

	preset, err := normalizePreset(flags.Preset)
	if err != nil {
		return InstallInput{}, err
	}
	selection.Preset = preset

	memoryBackend, err := normalizeMemoryBackend(flags.MemoryBackend)
	if err != nil {
		return InstallInput{}, err
	}

	components, err := normalizeComponents(flags.Components, selection.Preset)
	if err != nil {
		return InstallInput{}, err
	}
	if len(flags.Components) == 0 && strings.TrimSpace(flags.Preset) == "" && isPiOnlyAgents(selection.Agents) {
		components = piOnlyComponents()
	}
	if hasComponent(components, model.ComponentOpenCodeLeanWorkflow) {
		if strings.TrimSpace(flags.MemoryBackend) != "" && memoryBackend != model.MemoryBackendMarkdown {
			return InstallInput{}, fmt.Errorf("component %q requires --memory-backend markdown", model.ComponentOpenCodeLeanWorkflow)
		}
		memoryBackend = model.MemoryBackendMarkdown
	}
	if hasComponent(components, model.ComponentOpenCodeLayeredTDD) {
		if strings.TrimSpace(flags.MemoryBackend) != "" && memoryBackend != model.MemoryBackendMarkdown {
			return InstallInput{}, fmt.Errorf("component %q requires --memory-backend markdown", model.ComponentOpenCodeLayeredTDD)
		}
		memoryBackend = model.MemoryBackendMarkdown
	}
	if hasComponent(components, model.ComponentConductorLayeredTDD) {
		switch {
		case strings.TrimSpace(flags.MemoryBackend) == "":
			if onlyWorkspaceLocalComponents(components) {
				memoryBackend = model.MemoryBackendNone
			}
		case memoryBackend == model.MemoryBackendEngram:
			return InstallInput{}, fmt.Errorf("component %q supports only --memory-backend markdown or none", model.ComponentConductorLayeredTDD)
		}
	}
	if len(flags.Agents) == 0 && onlyWorkspaceLocalComponents(components) {
		selection.Agents = nil
	}

	selection.MemoryBackend = memoryBackend
	memoryVault, memoryNamespace, memoryProject, err := normalizeMarkdownMemoryConfig(flags, memoryBackend)
	if err != nil {
		return InstallInput{}, err
	}
	selection.MemoryVault = memoryVault
	selection.MemoryNamespace = memoryNamespace
	selection.MemoryProject = memoryProject

	selection.Components = componentsForMemoryBackend(components, memoryBackend)

	skills, err := normalizeSkills(flags.Skills)
	if err != nil {
		return InstallInput{}, err
	}
	selection.Skills = skills

	projectSkills, err := normalizeSkills(flags.ProjectSkills)
	if err != nil {
		return InstallInput{}, err
	}
	selection.ProjectSkills = projectSkills

	sddMode, err := normalizeSDDMode(flags.SDDMode)
	if err != nil {
		return InstallInput{}, err
	}
	selection.SDDMode = sddMode

	return InstallInput{Selection: selection, DryRun: flags.DryRun}, nil
}

func normalizePersona(value string) (model.PersonaID, error) {
	if strings.TrimSpace(value) == "" {
		return model.PersonaGentleman, nil
	}

	switch model.PersonaID(value) {
	case model.PersonaGentleman, model.PersonaNeutral, model.PersonaCustom:
		return model.PersonaID(value), nil
	default:
		return "", fmt.Errorf("unsupported persona %q", value)
	}
}

func normalizePreset(value string) (model.PresetID, error) {
	if strings.TrimSpace(value) == "" {
		return model.PresetFullGentleman, nil
	}

	switch model.PresetID(value) {
	case model.PresetFullGentleman, model.PresetEcosystemOnly, model.PresetMinimal, model.PresetCustom:
		return model.PresetID(value), nil
	default:
		return "", fmt.Errorf("unsupported preset %q", value)
	}
}

func normalizeComponents(values []string, preset model.PresetID) ([]model.ComponentID, error) {
	if len(values) == 0 {
		return componentsForPreset(preset), nil
	}

	allowed := map[model.ComponentID]struct{}{}
	for _, component := range catalog.MVPComponents() {
		allowed[component.ID] = struct{}{}
	}

	components := []model.ComponentID{}
	for _, raw := range values {
		component := model.ComponentID(raw)
		if _, ok := allowed[component]; !ok {
			return nil, fmt.Errorf("unsupported component %q", raw)
		}
		components = append(components, component)
	}

	return unique(components), nil
}

func onlyWorkspaceLocalComponents(components []model.ComponentID) bool {
	if len(components) == 0 {
		return false
	}
	for _, component := range components {
		if component != model.ComponentConductorLayeredTDD {
			return false
		}
	}
	return true
}

func normalizeSkills(values []string) ([]model.SkillID, error) {
	if len(values) == 0 {
		return nil, nil
	}

	allowed := map[model.SkillID]struct{}{}
	for _, skill := range catalog.MVPSkills() {
		allowed[skill.ID] = struct{}{}
	}

	skills := []model.SkillID{}
	for _, raw := range values {
		skill := model.SkillID(raw)
		if _, ok := allowed[skill]; !ok {
			return nil, fmt.Errorf("unsupported skill %q", raw)
		}
		skills = append(skills, skill)
	}

	return unique(skills), nil
}

func normalizeSDDMode(value string) (model.SDDModeID, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}

	switch model.SDDModeID(value) {
	case model.SDDModeSingle, model.SDDModeMulti:
		return model.SDDModeID(value), nil
	default:
		return "", fmt.Errorf("unsupported sdd-mode %q (valid: single, multi)", value)
	}
}

func normalizeMemoryBackend(value string) (model.MemoryBackendID, error) {
	if strings.TrimSpace(value) == "" {
		return model.MemoryBackendEngram, nil
	}

	switch model.MemoryBackendID(value) {
	case model.MemoryBackendEngram, model.MemoryBackendMarkdown, model.MemoryBackendNone:
		return model.MemoryBackendID(value), nil
	default:
		return "", fmt.Errorf("unsupported memory-backend %q (valid: engram, markdown, none)", value)
	}
}

func normalizeMarkdownMemoryConfig(flags InstallFlags, backend model.MemoryBackendID) (vault, namespace, project string, err error) {
	vault = strings.TrimSpace(flags.MemoryVault)
	namespace = strings.TrimSpace(flags.MemoryNamespace)
	project = strings.TrimSpace(flags.MemoryProject)

	if backend != model.MemoryBackendMarkdown {
		return vault, namespace, project, nil
	}

	if vault == "" {
		vault = "/Users/cdespona/Documents/thoughts/central/central"
	}
	if !filepath.IsAbs(vault) {
		return "", "", "", fmt.Errorf("memory-vault must be an absolute path")
	}
	if namespace == "" {
		namespace = "machine/agent-memory"
	}
	if err := validateMemoryNamespace(namespace); err != nil {
		return "", "", "", err
	}
	if project == "" {
		project = "gentle-ai-fork"
	}
	if err := validateMemoryProject(project); err != nil {
		return "", "", "", err
	}

	return filepath.Clean(vault), filepath.ToSlash(filepath.Clean(namespace)), project, nil
}

func validateMemoryNamespace(namespace string) error {
	if filepath.IsAbs(namespace) {
		return fmt.Errorf("memory-namespace must be relative to memory-vault")
	}
	clean := filepath.ToSlash(filepath.Clean(namespace))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../") {
		return fmt.Errorf("memory-namespace %q escapes memory-vault", namespace)
	}
	if clean != "machine/agent-memory" && !strings.HasPrefix(clean, "machine/agent-memory/") {
		return fmt.Errorf("memory-namespace must stay under machine/agent-memory")
	}
	return nil
}

var memoryProjectPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

func validateMemoryProject(project string) error {
	if !memoryProjectPattern.MatchString(project) {
		return fmt.Errorf("memory-project %q must be a lowercase slug using letters, numbers, dots, underscores, or hyphens", project)
	}
	return nil
}

func componentsForMemoryBackend(components []model.ComponentID, backend model.MemoryBackendID) []model.ComponentID {
	filtered := make([]model.ComponentID, 0, len(components)+1)
	hasSDD := false
	hasMarkdown := false
	hasEngram := false
	hasConductorLayeredTDD := false
	for _, component := range components {
		if component == model.ComponentSDD {
			hasSDD = true
		}
		if component == model.ComponentEngram {
			hasEngram = true
		}
		if component == model.ComponentMarkdownMemory {
			hasMarkdown = true
		}
		if component == model.ComponentConductorLayeredTDD {
			hasConductorLayeredTDD = true
		}
		if component == model.ComponentEngram || component == model.ComponentMarkdownMemory {
			continue
		}
		filtered = append(filtered, component)
	}

	switch backend {
	case model.MemoryBackendMarkdown:
		if hasSDD || hasEngram || hasMarkdown || hasConductorLayeredTDD {
			filtered = append([]model.ComponentID{model.ComponentMarkdownMemory}, filtered...)
		}
	case model.MemoryBackendNone:
		return filtered
	default:
		if hasEngram || hasSDD || hasMarkdown {
			filtered = append([]model.ComponentID{model.ComponentEngram}, filtered...)
		}
	}
	return unique(filtered)
}

func componentsForPreset(preset model.PresetID) []model.ComponentID {
	switch preset {
	case model.PresetMinimal:
		return []model.ComponentID{model.ComponentEngram}
	case model.PresetEcosystemOnly:
		return []model.ComponentID{model.ComponentEngram, model.ComponentSDD, model.ComponentSkills, model.ComponentContext7, model.ComponentGGA}
	case model.PresetCustom:
		return nil
	default:
		return []model.ComponentID{
			model.ComponentEngram,
			model.ComponentSDD,
			model.ComponentSkills,
			model.ComponentContext7,
			model.ComponentPersona,
			model.ComponentPermission,
			model.ComponentGGA,
			model.ComponentClaudeTheme,
			model.ComponentOpenCodeGentleLogo,
		}
	}
}

func defaultAgentsFromDetection(detection system.DetectionResult) []model.AgentID {
	agents := []model.AgentID{}
	for _, state := range detection.Configs {
		if !state.Exists {
			continue
		}

		switch strings.TrimSpace(state.Agent) {
		case string(model.AgentClaudeCode):
			agents = append(agents, model.AgentClaudeCode)
		case string(model.AgentOpenCode):
			agents = append(agents, model.AgentOpenCode)
		case string(model.AgentKilocode):
			agents = append(agents, model.AgentKilocode)
		case string(model.AgentGeminiCLI):
			agents = append(agents, model.AgentGeminiCLI)
		case string(model.AgentCursor):
			agents = append(agents, model.AgentCursor)
		case string(model.AgentVSCodeCopilot):
			agents = append(agents, model.AgentVSCodeCopilot)
		case string(model.AgentCodex):
			agents = append(agents, model.AgentCodex)
		case string(model.AgentAntigravity):
			agents = append(agents, model.AgentAntigravity)
		case string(model.AgentWindsurf):
			agents = append(agents, model.AgentWindsurf)
		case string(model.AgentKimi):
			agents = append(agents, model.AgentKimi)
		case string(model.AgentQwenCode):
			agents = append(agents, model.AgentQwenCode)
		case string(model.AgentKiroIDE):
			agents = append(agents, model.AgentKiroIDE)
		case string(model.AgentOpenClaw):
			agents = append(agents, model.AgentOpenClaw)
		case string(model.AgentPi):
			agents = append(agents, model.AgentPi)
		}
	}

	if len(agents) > 0 {
		return agents
	}

	catalogAgents := catalog.AllAgents()
	agents = make([]model.AgentID, 0, len(catalogAgents))
	for _, agent := range catalogAgents {
		agents = append(agents, agent.ID)
	}

	return agents
}

func asAgentIDs(values []string) []model.AgentID {
	agents := make([]model.AgentID, 0, len(values))
	for _, value := range values {
		agents = append(agents, model.AgentID(value))
	}

	return agents
}

func isPiOnlyAgents(agents []model.AgentID) bool {
	return len(agents) == 1 && agents[0] == model.AgentPi
}

func piOnlyComponents() []model.ComponentID {
	return []model.ComponentID{model.ComponentEngram}
}

func unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}, len(items))
	result := make([]T, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}
