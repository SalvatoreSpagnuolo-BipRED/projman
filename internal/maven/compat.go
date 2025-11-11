// Package maven fornisce funzionalità per l'analisi di progetti Maven
// Questo file fornisce funzioni di compatibilità con l'API precedente
package maven

import (
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/buildsystem"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/graph"
	"github.com/pterm/pterm"
)

// BuildDependencyGraph costruisce il grafo delle dipendenze tra i progetti Maven selezionati
// Questa funzione mantiene la compatibilità con l'API precedente
func BuildDependencyGraph(projectNames []string, rootPath string) (map[string][]string, error) {
	parser := NewParser()
	builder := buildsystem.NewGraphBuilder(parser)

	// Costruisci il grafo base
	dependencyGraph, err := builder.BuildGraph(projectNames, rootPath)
	if err != nil {
		return nil, err
	}

	// Registra i sub-modules per gestire dipendenze indirette
	registry := builder.GetRegistry()
	for _, name := range projectNames {
		project, err := parser.ParseProject(name, rootPath)
		if err != nil {
			continue
		}

		if err := RegisterSubModules(project.Path, name, registry); err != nil {
			pterm.Warning.Println(err)
		}
	}

	// Ricostruisci il grafo considerando i sub-modules
	dependencyGraph, err = builder.BuildGraph(projectNames, rootPath)
	if err != nil {
		return nil, err
	}

	// Converti il grafo nel formato mappa per compatibilità
	return map[string][]string(dependencyGraph), nil
}

// TopologicalSort ordina i progetti in base alle loro dipendenze
// Restituisce l'ordine di esecuzione corretto o un errore se ci sono cicli
// Questa funzione mantiene la compatibilità con l'API precedente
func TopologicalSort(graphMap map[string][]string) ([]string, error) {
	dependencyGraph := graph.DependencyGraph(graphMap)
	return dependencyGraph.TopologicalSort()
}
