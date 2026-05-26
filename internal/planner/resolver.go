package planner

import (
	"fmt"

	"github.com/gentleman-programming/gentle-ai/internal/catalog"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

type dependencyResolver struct {
	graph Graph
}

func NewResolver(graph Graph) Resolver {
	return dependencyResolver{graph: graph}
}

func (r dependencyResolver) Resolve(selection model.Selection) (ResolvedPlan, error) {
	resolved := ResolvedPlan{}

	components := append([]model.ComponentID(nil), selection.Components...)
	if memoryDependency := memoryDependencyForSDD(selection); memoryDependency != "" {
		components = append(components, memoryDependency)
	}
	components = uniqueComponents(components)

	selectedSet := make(map[model.ComponentID]struct{}, len(selection.Components))
	dependencies := map[model.ComponentID][]model.ComponentID{}
	for _, selected := range selection.Components {
		selectedSet[selected] = struct{}{}
	}
	for _, selected := range components {
		if !r.graph.Has(selected) {
			return ResolvedPlan{}, fmt.Errorf("unknown component %q", selected)
		}

		if err := r.expandDependencies(selected, dependencies); err != nil {
			return ResolvedPlan{}, err
		}
	}

	orderedComponents, err := TopologicalSort(dependencies)
	if err != nil {
		return ResolvedPlan{}, err
	}

	// Apply soft ordering constraints: when BOTH components in a pair are
	// present, ensure the first appears before the second. This does NOT
	// add missing components — it only reorders what is already selected.
	orderedComponents = applySoftOrdering(orderedComponents, SoftOrderingConstraints())

	for _, component := range orderedComponents {
		if _, selected := selectedSet[component]; !selected {
			resolved.AddedDependencies = append(resolved.AddedDependencies, component)
		}
	}

	resolved.OrderedComponents = orderedComponents

	for _, agent := range selection.Agents {
		if catalog.IsSupportedAgent(agent) {
			resolved.Agents = append(resolved.Agents, agent)
			continue
		}

		resolved.UnsupportedAgents = append(resolved.UnsupportedAgents, agent)
	}

	return resolved, nil
}

func memoryDependencyForSDD(selection model.Selection) model.ComponentID {
	if !hasComponent(selection.Components, model.ComponentSDD) && !hasComponent(selection.Components, model.ComponentSkills) {
		return ""
	}
	if hasComponent(selection.Components, model.ComponentEngram) || hasComponent(selection.Components, model.ComponentMarkdownMemory) {
		return ""
	}

	switch selection.MemoryBackend {
	case model.MemoryBackendMarkdown:
		return model.ComponentMarkdownMemory
	case model.MemoryBackendNone:
		return ""
	default:
		return model.ComponentEngram
	}
}

func hasComponent(components []model.ComponentID, target model.ComponentID) bool {
	for _, component := range components {
		if component == target {
			return true
		}
	}
	return false
}

func uniqueComponents(components []model.ComponentID) []model.ComponentID {
	seen := make(map[model.ComponentID]struct{}, len(components))
	result := make([]model.ComponentID, 0, len(components))
	for _, component := range components {
		if _, ok := seen[component]; ok {
			continue
		}
		seen[component] = struct{}{}
		result = append(result, component)
	}
	return result
}

func (r dependencyResolver) expandDependencies(component model.ComponentID, dependencies map[model.ComponentID][]model.ComponentID) error {
	if _, visited := dependencies[component]; visited {
		return nil
	}

	deps := r.graph.DependenciesOf(component)
	dependencies[component] = deps
	for _, dep := range deps {
		if !r.graph.Has(dep) {
			return fmt.Errorf("component %q depends on unknown dependency %q", component, dep)
		}

		if err := r.expandDependencies(dep, dependencies); err != nil {
			return err
		}
	}

	return nil
}
