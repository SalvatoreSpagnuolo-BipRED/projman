package maven

import (
	"testing"
)

// TestTopologicalSort verifica l'ordinamento topologico di un grafo di dipendenze
func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name          string
		graph         map[string][]string
		expectedOrder []string
		expectError   bool
	}{
		{
			name: "Grafo semplice senza dipendenze",
			graph: map[string][]string{
				"projectA": {},
				"projectB": {},
				"projectC": {},
			},
			expectedOrder: []string{"projectA", "projectB", "projectC"},
			expectError:   false,
		},
		{
			name: "Grafo con dipendenze lineari",
			graph: map[string][]string{
				"projectA": {},
				"projectB": {"projectA"},
				"projectC": {"projectB"},
			},
			expectedOrder: []string{"projectA", "projectB", "projectC"},
			expectError:   false,
		},
		{
			name: "Grafo con dipendenze multiple",
			graph: map[string][]string{
				"projectA": {},
				"projectB": {},
				"projectC": {"projectA", "projectB"},
				"projectD": {"projectC"},
			},
			expectedOrder: nil, // L'ordine esatto dipende dall'implementazione, ma A e B devono venire prima di C
			expectError:   false,
		},
		{
			name: "Grafo con dipendenza circolare",
			graph: map[string][]string{
				"projectA": {"projectB"},
				"projectB": {"projectA"},
			},
			expectedOrder: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TopologicalSort(tt.graph)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.graph) {
				t.Errorf("Expected %d projects, got %d", len(tt.graph), len(result))
				return
			}

			// Verifica che le dipendenze siano rispettate
			position := make(map[string]int)
			for i, proj := range result {
				position[proj] = i
			}

			for proj, deps := range tt.graph {
				projPos := position[proj]
				for _, dep := range deps {
					depPos := position[dep]
					if depPos >= projPos {
						t.Errorf("Dependency %s should come before %s, but positions are %d and %d", dep, proj, depPos, projPos)
					}
				}
			}
		})
	}
}
