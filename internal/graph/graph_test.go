package graph

import (
	"testing"
)

func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name          string
		graph         DependencyGraph
		expectedOrder []string
		expectError   bool
	}{
		{
			name: "Grafo semplice senza dipendenze",
			graph: DependencyGraph{
				"projectA": {},
				"projectB": {},
				"projectC": {},
			},
			expectError: false,
		},
		{
			name: "Grafo con dipendenze lineari",
			graph: DependencyGraph{
				"projectA": {},
				"projectB": {"projectA"},
				"projectC": {"projectB"},
			},
			expectedOrder: []string{"projectA", "projectB", "projectC"},
			expectError:   false,
		},
		{
			name: "Grafo con dipendenze multiple",
			graph: DependencyGraph{
				"projectA": {},
				"projectB": {"projectA"},
				"projectC": {"projectA"},
				"projectD": {"projectB", "projectC"},
			},
			expectError: false,
		},
		{
			name: "Grafo con dipendenza circolare",
			graph: DependencyGraph{
				"projectA": {"projectB"},
				"projectB": {"projectC"},
				"projectC": {"projectA"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.graph.TopologicalSort()

			if tt.expectError {
				if err == nil {
					t.Error("Ci si aspettava un errore ma non è stato restituito")
				}
				return
			}

			if err != nil {
				t.Errorf("Errore non previsto: %v", err)
				return
			}

			if len(result) != len(tt.graph) {
				t.Errorf("Lunghezza risultato errata: atteso %d, ottenuto %d", len(tt.graph), len(result))
			}

			// Se specificato un ordine, verifica che sia corretto
			if len(tt.expectedOrder) > 0 {
				for i, expected := range tt.expectedOrder {
					if i >= len(result) || result[i] != expected {
						t.Errorf("Ordine errato alla posizione %d: atteso %s, ottenuto %s", i, expected, result[i])
					}
				}
			}

			// Verifica che le dipendenze siano rispettate
			positions := make(map[string]int)
			for i, name := range result {
				positions[name] = i
			}

			for node, deps := range tt.graph {
				nodePos := positions[node]
				for _, dep := range deps {
					depPos := positions[dep]
					if depPos >= nodePos {
						t.Errorf("Dipendenza non rispettata: %s dipende da %s ma appare prima nell'ordine", node, dep)
					}
				}
			}
		})
	}
}

func TestDependencyGraphOperations(t *testing.T) {
	t.Run("NewDependencyGraph", func(t *testing.T) {
		g := NewDependencyGraph()
		if g == nil {
			t.Error("NewDependencyGraph ha restituito nil")
		}
		if len(g) != 0 {
			t.Error("Il grafo dovrebbe essere vuoto")
		}
	})

	t.Run("AddNode", func(t *testing.T) {
		g := NewDependencyGraph()
		g.AddNode("projectA", []string{"projectB"})

		if len(g) != 1 {
			t.Error("Il grafo dovrebbe contenere 1 nodo")
		}

		deps := g.GetDependencies("projectA")
		if len(deps) != 1 || deps[0] != "projectB" {
			t.Error("Le dipendenze non sono state aggiunte correttamente")
		}
	})

	t.Run("AddDependency", func(t *testing.T) {
		g := NewDependencyGraph()
		g.AddDependency("projectA", "projectB")
		g.AddDependency("projectA", "projectC")

		deps := g.GetDependencies("projectA")
		if len(deps) != 2 {
			t.Errorf("Numero di dipendenze errato: atteso 2, ottenuto %d", len(deps))
		}
	})

	t.Run("HasCycles", func(t *testing.T) {
		// Grafo senza cicli
		g1 := DependencyGraph{
			"A": {},
			"B": {"A"},
		}
		if g1.HasCycles() {
			t.Error("HasCycles dovrebbe restituire false per grafi senza cicli")
		}

		// Grafo con cicli
		g2 := DependencyGraph{
			"A": {"B"},
			"B": {"A"},
		}
		if !g2.HasCycles() {
			t.Error("HasCycles dovrebbe restituire true per grafi con cicli")
		}
	})

	t.Run("Nodes", func(t *testing.T) {
		g := DependencyGraph{
			"A": {},
			"B": {},
			"C": {},
		}

		nodes := g.Nodes()
		if len(nodes) != 3 {
			t.Errorf("Numero di nodi errato: atteso 3, ottenuto %d", len(nodes))
		}
	})
}
