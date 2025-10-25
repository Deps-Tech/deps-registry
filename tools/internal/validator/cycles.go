package validator

import (
	"fmt"

	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
)

type CycleError struct {
	Cycle []string
}

func (e *CycleError) Error() string {
	return fmt.Sprintf("circular dependency detected: %v", e.Cycle)
}

type GraphNode struct {
	ID           string
	Dependencies []string
}

func DetectCycles(manifests map[string]*manifest.Manifest) []*CycleError {
	graph := buildGraph(manifests)
	cycles := []*CycleError{}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := []string{}

	for id := range graph {
		if !visited[id] {
			if cycle := dfs(id, graph, visited, recStack, path); cycle != nil {
				cycles = append(cycles, &CycleError{Cycle: cycle})
			}
		}
	}

	return cycles
}

func buildGraph(manifests map[string]*manifest.Manifest) map[string]*GraphNode {
	graph := make(map[string]*GraphNode)

	for id, m := range manifests {
		deps := []string{}
		for depID := range m.Dependencies {
			deps = append(deps, depID)
		}
		graph[id] = &GraphNode{
			ID:           id,
			Dependencies: deps,
		}
	}

	return graph
}

func dfs(nodeID string, graph map[string]*GraphNode, visited, recStack map[string]bool, path []string) []string {
	visited[nodeID] = true
	recStack[nodeID] = true
	path = append(path, nodeID)

	node := graph[nodeID]
	if node != nil {
		for _, depID := range node.Dependencies {
			if !visited[depID] {
				if cycle := dfs(depID, graph, visited, recStack, path); cycle != nil {
					return cycle
				}
			} else if recStack[depID] {
				cycleStart := -1
				for i, id := range path {
					if id == depID {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					return append(path[cycleStart:], depID)
				}
			}
		}
	}

	recStack[nodeID] = false
	return nil
}

