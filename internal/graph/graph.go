// Package graph fornisce strutture e algoritmi per la gestione di grafi di dipendenze
package graph

import "fmt"

// DependencyGraph rappresenta un grafo orientato di dipendenze tra progetti
// La chiave è il nome del progetto, il valore è la lista di progetti da cui dipende
type DependencyGraph map[string][]string

// NewDependencyGraph crea un nuovo grafo di dipendenze vuoto
func NewDependencyGraph() DependencyGraph {
	return make(DependencyGraph)
}

// AddNode aggiunge un nodo al grafo con una lista di dipendenze
func (g DependencyGraph) AddNode(name string, dependencies []string) {
	g[name] = dependencies
}

// AddDependency aggiunge una singola dipendenza a un nodo esistente
func (g DependencyGraph) AddDependency(name string, dependency string) {
	if _, exists := g[name]; !exists {
		g[name] = make([]string, 0)
	}
	g[name] = append(g[name], dependency)
}

// TopologicalSort ordina i nodi del grafo in base alle loro dipendenze
// Restituisce l'ordine di esecuzione corretto o un errore se ci sono cicli
// Usa l'algoritmo di Kahn per il topological sort
func (g DependencyGraph) TopologicalSort() ([]string, error) {
	// Calcola l'in-degree (numero di dipendenze) per ogni nodo
	inDegree := make(map[string]int)

	// Inizializza tutti i nodi con in-degree 0
	for node := range g {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
	}

	// Calcola l'in-degree: per ogni progetto, incrementa l'in-degree in base al numero di dipendenze
	for node, deps := range g {
		inDegree[node] += len(deps)
	}

	// Coda dei nodi senza dipendenze (in-degree = 0)
	queue := make([]string, 0)
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// Risultato dell'ordinamento topologico
	result := make([]string, 0, len(g))

	// Algoritmo di Kahn
	for len(queue) > 0 {
		// Prendi il primo nodo dalla coda
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Per ogni progetto che dipende dal nodo corrente, decrementa il suo in-degree
		for node, deps := range g {
			for _, dep := range deps {
				if dep == current {
					inDegree[node]--
					if inDegree[node] == 0 {
						queue = append(queue, node)
					}
				}
			}
		}
	}

	// Se non tutti i nodi sono stati processati, c'è un ciclo
	if len(result) != len(g) {
		return nil, fmt.Errorf("dipendenze circolari rilevate nel grafo")
	}

	return result, nil
}

// HasCycles verifica se il grafo contiene cicli
func (g DependencyGraph) HasCycles() bool {
	_, err := g.TopologicalSort()
	return err != nil
}

// GetDependencies restituisce la lista di dipendenze per un nodo
func (g DependencyGraph) GetDependencies(name string) []string {
	if deps, exists := g[name]; exists {
		return deps
	}
	return []string{}
}

// Nodes restituisce la lista di tutti i nodi nel grafo
func (g DependencyGraph) Nodes() []string {
	nodes := make([]string, 0, len(g))
	for node := range g {
		nodes = append(nodes, node)
	}
	return nodes
}
