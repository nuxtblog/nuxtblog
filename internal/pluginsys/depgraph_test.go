package pluginsys

import (
	"strings"
	"testing"

	"github.com/nuxtblog/nuxtblog/sdk"
)

func TestTopologicalSort_Basic(t *testing.T) {
	g := newDepGraph()
	// C has no deps, B depends on C, A depends on B
	g.Add("C", nil)
	g.Add("B", []sdk.Dependency{{ID: "C"}})
	g.Add("A", []sdk.Dependency{{ID: "B"}})

	versions := func(id string) string {
		return "1.0.0"
	}

	order, err := g.TopologicalSort([]string{"A", "B", "C"}, versions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// C must come before B, B must come before A
	idx := map[string]int{}
	for i, id := range order {
		idx[id] = i
	}
	if idx["C"] >= idx["B"] {
		t.Errorf("C should come before B, got order: %v", order)
	}
	if idx["B"] >= idx["A"] {
		t.Errorf("B should come before A, got order: %v", order)
	}
}

func TestTopologicalSort_CircularDependency(t *testing.T) {
	g := newDepGraph()
	g.Add("A", []sdk.Dependency{{ID: "B"}})
	g.Add("B", []sdk.Dependency{{ID: "A"}})

	_, err := g.TopologicalSort([]string{"A", "B"}, func(string) string { return "1.0.0" })
	if err == nil {
		t.Fatal("expected circular dependency error")
	}
	if !strings.Contains(err.Error(), "circular dependency") {
		t.Errorf("expected 'circular dependency' in error, got: %v", err)
	}
}

func TestTopologicalSort_MissingRequired(t *testing.T) {
	g := newDepGraph()
	g.Add("A", []sdk.Dependency{{ID: "B"}})

	_, err := g.TopologicalSort([]string{"A"}, func(string) string { return "" })
	if err == nil {
		t.Fatal("expected missing dependency error")
	}
	if !strings.Contains(err.Error(), "missing plugin B") {
		t.Errorf("expected 'missing plugin B' in error, got: %v", err)
	}
}

func TestTopologicalSort_OptionalMissing(t *testing.T) {
	g := newDepGraph()
	g.Add("A", []sdk.Dependency{{ID: "B", Optional: true}})

	order, err := g.TopologicalSort([]string{"A"}, func(string) string { return "" })
	if err != nil {
		t.Fatalf("optional missing dep should not error: %v", err)
	}
	if len(order) != 1 || order[0] != "A" {
		t.Errorf("expected [A], got %v", order)
	}
}

func TestTopologicalSort_VersionConstraint(t *testing.T) {
	g := newDepGraph()
	g.Add("B", nil)
	g.Add("A", []sdk.Dependency{{ID: "B", Version: ">=2.0.0"}})

	_, err := g.TopologicalSort([]string{"A", "B"}, func(id string) string {
		if id == "B" {
			return "1.5.0"
		}
		return "1.0.0"
	})
	if err == nil {
		t.Fatal("expected version constraint error")
	}
	if !strings.Contains(err.Error(), "requires") {
		t.Errorf("expected version error, got: %v", err)
	}
}

func TestCascadeUnloadOrder(t *testing.T) {
	g := newDepGraph()
	g.Add("C", nil)
	g.Add("B", []sdk.Dependency{{ID: "C"}})
	g.Add("A", []sdk.Dependency{{ID: "B"}})

	// Unloading C should cascade: A, B, C
	order := g.GetCascadeUnloadOrder("C")

	if len(order) != 3 {
		t.Fatalf("expected 3 items, got %v", order)
	}
	// Last item must be C (the root being unloaded)
	if order[len(order)-1] != "C" {
		t.Errorf("last item should be C, got %v", order)
	}
	// A and B should come before C
	found := map[string]bool{}
	for _, id := range order {
		found[id] = true
	}
	if !found["A"] || !found["B"] {
		t.Errorf("expected A and B in cascade, got %v", order)
	}
}

func TestGetDependents(t *testing.T) {
	g := newDepGraph()
	g.Add("C", nil)
	g.Add("B", []sdk.Dependency{{ID: "C"}})
	g.Add("A", []sdk.Dependency{{ID: "C"}})

	deps := g.GetDependents("C")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependents, got %v", deps)
	}
}

func TestRemove(t *testing.T) {
	g := newDepGraph()
	g.Add("B", nil)
	g.Add("A", []sdk.Dependency{{ID: "B"}})
	g.Remove("A")

	deps := g.GetDependents("B")
	if len(deps) != 0 {
		t.Errorf("expected 0 dependents after remove, got %v", deps)
	}
}

func TestMatchSemverConstraint(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		want       bool
	}{
		{"1.0.0", ">=1.0.0", true},
		{"1.0.0", ">=1.0.1", false},
		{"2.0.0", ">=1.0.0", true},
		{"1.0.0", ">1.0.0", false},
		{"1.0.1", ">1.0.0", true},
		{"1.0.0", "=1.0.0", true},
		{"1.0.1", "=1.0.0", false},
		{"0.9.0", "<1.0.0", true},
		{"1.0.0", "<1.0.0", false},
		{"1.0.0", "<=1.0.0", true},
		{"1.0.1", "<=1.0.0", false},
		{"1.0.0", "1.0.0", true},  // default >= operator
		{"0.9.0", "1.0.0", false}, // default >= operator
	}
	for _, tt := range tests {
		got := matchSemverConstraint(tt.version, tt.constraint)
		if got != tt.want {
			t.Errorf("matchSemverConstraint(%q, %q) = %v, want %v", tt.version, tt.constraint, got, tt.want)
		}
	}
}
