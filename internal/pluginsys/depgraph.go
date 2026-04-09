package pluginsys

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// depGraph tracks plugin dependency relationships for topological ordering
// and cascade unloading.
type depGraph struct {
	mu         sync.RWMutex
	deps       map[string][]sdk.Dependency   // plugin ID -> its declared dependencies
	dependents map[string]map[string]struct{} // plugin ID -> set of IDs that depend on it
}

func newDepGraph() *depGraph {
	return &depGraph{
		deps:       make(map[string][]sdk.Dependency),
		dependents: make(map[string]map[string]struct{}),
	}
}

// Add registers a plugin and its dependencies, updating both forward and reverse indexes.
func (g *depGraph) Add(id string, deps []sdk.Dependency) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Remove old reverse edges if re-registering
	if old, ok := g.deps[id]; ok {
		for _, d := range old {
			if s, ok := g.dependents[d.ID]; ok {
				delete(s, id)
			}
		}
	}

	g.deps[id] = deps

	for _, d := range deps {
		if g.dependents[d.ID] == nil {
			g.dependents[d.ID] = make(map[string]struct{})
		}
		g.dependents[d.ID][id] = struct{}{}
	}
}

// Remove deletes a node and cleans up both forward and reverse edges.
func (g *depGraph) Remove(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Remove forward edges (this plugin's deps)
	if deps, ok := g.deps[id]; ok {
		for _, d := range deps {
			if s, ok := g.dependents[d.ID]; ok {
				delete(s, id)
				if len(s) == 0 {
					delete(g.dependents, d.ID)
				}
			}
		}
		delete(g.deps, id)
	}

	// Remove reverse edges (plugins that depend on this one)
	delete(g.dependents, id)
}

// GetDependents returns the IDs of plugins that directly depend on id.
func (g *depGraph) GetDependents(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	s := g.dependents[id]
	result := make([]string, 0, len(s))
	for pid := range s {
		result = append(result, pid)
	}
	return result
}

// GetCascadeUnloadOrder returns all plugins that must be unloaded when id is unloaded,
// ordered from deepest dependent to id itself (i.e. unload order).
func (g *depGraph) GetCascadeUnloadOrder(id string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// BFS collecting all transitive dependents
	visited := map[string]bool{id: true}
	queue := []string{id}
	var order []string

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)

		for dep := range g.dependents[cur] {
			if !visited[dep] {
				visited[dep] = true
				queue = append(queue, dep)
			}
		}
	}

	// Reverse: deepest dependents first, id last
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order
}

// TopologicalSort returns a load order for the given plugin IDs using Kahn's algorithm.
// Non-optional dependencies missing from ids cause an error.
// Version constraints are checked via getVersion.
// Circular dependencies are detected and reported.
func (g *depGraph) TopologicalSort(ids []string, getVersion func(string) string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}

	// Build in-degree map (only for edges within idSet)
	inDegree := make(map[string]int, len(ids))
	for _, id := range ids {
		inDegree[id] = 0
	}

	for _, id := range ids {
		for _, dep := range g.deps[id] {
			if !idSet[dep.ID] {
				if !dep.Optional {
					return nil, fmt.Errorf("plugin %s requires missing plugin %s", id, dep.ID)
				}
				continue
			}
			// Check version constraint
			if dep.Version != "" {
				ver := getVersion(dep.ID)
				if ver != "" && !matchSemverConstraint(ver, dep.Version) {
					return nil, fmt.Errorf("plugin %s requires %s %s, but found %s", id, dep.ID, dep.Version, ver)
				}
			}
			inDegree[id]++
		}
	}

	// Kahn's algorithm
	var queue []string
	for _, id := range ids {
		if inDegree[id] == 0 {
			queue = append(queue, id)
		}
	}

	var sorted []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		sorted = append(sorted, node)

		// For each plugin that depends on node, decrement in-degree
		for dependent := range g.dependents[node] {
			if !idSet[dependent] {
				continue
			}
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	if len(sorted) < len(ids) {
		// Circular dependency detected — find the cycle via DFS
		cycle := g.findCycle(ids, idSet)
		if cycle != "" {
			return nil, fmt.Errorf("circular dependency: %s", cycle)
		}
		return nil, fmt.Errorf("circular dependency detected among plugins")
	}

	return sorted, nil
}

// findCycle uses DFS to find and format a cycle path among the given IDs.
func (g *depGraph) findCycle(ids []string, idSet map[string]bool) string {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int, len(ids))
	parent := make(map[string]string, len(ids))

	var dfs func(string) string
	dfs = func(u string) string {
		color[u] = gray
		for _, dep := range g.deps[u] {
			if !idSet[dep.ID] {
				continue
			}
			if color[dep.ID] == gray {
				// Found cycle: trace back from u to dep.ID
				path := []string{dep.ID, u}
				cur := u
				for cur != dep.ID {
					cur = parent[cur]
					path = append(path, cur)
				}
				// Reverse to get forward direction
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return strings.Join(path, " → ")
			}
			if color[dep.ID] == white {
				parent[dep.ID] = u
				if result := dfs(dep.ID); result != "" {
					return result
				}
			}
		}
		color[u] = black
		return ""
	}

	for _, id := range ids {
		if color[id] == white {
			if result := dfs(id); result != "" {
				return result
			}
		}
	}
	return ""
}

// matchSemverConstraint checks if version satisfies a semver constraint string.
// Supported prefixes: >=, >, =, <, <= (default is >= if no prefix).
func matchSemverConstraint(version, constraint string) bool {
	constraint = strings.TrimSpace(constraint)
	if constraint == "" {
		return true
	}

	var op string
	var target string

	for _, prefix := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(constraint, prefix) {
			op = prefix
			target = strings.TrimSpace(constraint[len(prefix):])
			break
		}
	}
	if op == "" {
		op = ">="
		target = constraint
	}

	cmp := compareSemver(version, target)
	switch op {
	case ">=":
		return cmp >= 0
	case ">":
		return cmp > 0
	case "=":
		return cmp == 0
	case "<=":
		return cmp <= 0
	case "<":
		return cmp < 0
	}
	return false
}

// compareSemver returns -1, 0, or 1 comparing version a to b.
func compareSemver(a, b string) int {
	parse := func(s string) (major, minor, patch int, ok bool) {
		_, err := fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
		return major, minor, patch, err == nil
	}

	aMaj, aMin, aPat, ok1 := parse(a)
	bMaj, bMin, bPat, ok2 := parse(b)
	if !ok1 || !ok2 {
		return 0
	}

	if aMaj != bMaj {
		if aMaj > bMaj {
			return 1
		}
		return -1
	}
	if aMin != bMin {
		if aMin > bMin {
			return 1
		}
		return -1
	}
	if aPat != bPat {
		if aPat > bPat {
			return 1
		}
		return -1
	}
	return 0
}
