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
// Usa l'algoritmo di Kahn ottimizzato con reverse graph (O(n+m) invece di O(n²))
func (g DependencyGraph) TopologicalSort() ([]string, error) {
	// Inizializza in-degree e reverse graph
	inDegree := make(map[string]int)
	reverseGraph := make(map[string][]string) // chi dipende da me

	// Inizializza tutti i nodi con in-degree 0
	for node := range g {
		inDegree[node] = 0
	}

	// Calcola in-degree e costruisci reverse graph in una sola passata
	for node, deps := range g {
		inDegree[node] += len(deps)
		for _, dep := range deps {
			reverseGraph[dep] = append(reverseGraph[dep], node)
		}
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

	// Algoritmo di Kahn ottimizzato
	for len(queue) > 0 {
		// Prendi il primo nodo dalla coda
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// O(1) lookup: trova tutti i nodi che dipendono dal nodo corrente
		for _, dependent := range reverseGraph[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
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
