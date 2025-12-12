// Package maven fornisce funzionalità per l'analisi di progetti Maven
package maven

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/buildsystem"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/graph"
)

// pom rappresenta la struttura XML di un file pom.xml (semplificata)
type pom struct {
	XMLName      xml.Name    `xml:"project"`
	GroupId      string      `xml:"groupId"`
	ArtifactId   string      `xml:"artifactId"`
	Parent       *pomParent  `xml:"parent"`
	Dependencies *pomDeps    `xml:"dependencies"`
	Modules      *pomModules `xml:"modules"`
}

type pomParent struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
}

type pomDeps struct {
	Dependency []pomDependency `xml:"dependency"`
}

type pomDependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
}

type pomModules struct {
	Module []string `xml:"module"`
}

// ParseProject analizza un progetto Maven e restituisce le sue informazioni
func ParseProject(projectName, rootPath string) (*buildsystem.Project, error) {
	pomPath := filepath.Join(rootPath, projectName, "pom.xml")

	pomData, err := parsePomFile(pomPath)
	if err != nil {
		return nil, fmt.Errorf("errore analisi progetto %s: %w", projectName, err)
	}

	project := &buildsystem.Project{
		Name:         projectName,
		Path:         filepath.Join(rootPath, projectName),
		Identifier:   makeIdentifier(pomData.GroupId, pomData.ArtifactId),
		Dependencies: make([]string, 0),
	}

	// Estrai le dipendenze dirette
	project.Dependencies = append(project.Dependencies, extractDependencies(pomData)...)

	// Se il progetto ha moduli, analizza anche quelli per trovare dipendenze
	if pomData.Modules != nil && len(pomData.Modules.Module) > 0 {
		for _, module := range pomData.Modules.Module {
			modulePomPath := filepath.Join(project.Path, module, "pom.xml")
			moduleDeps, err := extractModuleDependenciesRecursive(modulePomPath)
			if err == nil {
				project.Dependencies = append(project.Dependencies, moduleDeps...)
			}
		}
	}

	return project, nil
}

// parsePomFile legge e analizza un file pom.xml
func parsePomFile(pomPath string) (*pom, error) {
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return nil, fmt.Errorf("impossibile leggere %s: %w", pomPath, err)
	}

	var p pom
	if err := xml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("errore parsing XML di %s: %w", pomPath, err)
	}

	// Usa il groupId del parent se non specificato
	if p.GroupId == "" && p.Parent != nil {
		p.GroupId = p.Parent.GroupId
	}

	return &p, nil
}

// extractDependencies estrae le dipendenze da un pom
func extractDependencies(p *pom) []string {
	deps := make([]string, 0)
	if p.Dependencies != nil {
		for _, dep := range p.Dependencies.Dependency {
			deps = append(deps, makeIdentifier(dep.GroupId, dep.ArtifactId))
		}
	}
	return deps
}

// extractModuleDependenciesRecursive estrae le dipendenze da un modulo figlio ricorsivamente
func extractModuleDependenciesRecursive(modulePomPath string) ([]string, error) {
	pomData, err := parsePomFile(modulePomPath)
	if err != nil {
		return nil, err
	}

	deps := extractDependencies(pomData)

	// Gestione ricorsiva dei moduli annidati
	if pomData.Modules != nil && len(pomData.Modules.Module) > 0 {
		for _, module := range pomData.Modules.Module {
			nestedPomPath := filepath.Join(filepath.Dir(modulePomPath), module, "pom.xml")
			nestedDeps, err := extractModuleDependenciesRecursive(nestedPomPath)
			if err == nil {
				deps = append(deps, nestedDeps...)
			}
		}
	}

	return deps, nil
}

// makeIdentifier crea l'identificatore univoco Maven (groupId:artifactId)
func makeIdentifier(groupId, artifactId string) string {
	return fmt.Sprintf("%s:%s", groupId, artifactId)
}

// RegisterSubModules registra ricorsivamente tutti i sub-modules di un progetto nel registry
// Questo permette di riconoscere dipendenze indirette tra progetti attraverso i loro sub-modules
func RegisterSubModules(projectPath string, parentProjectName string, registry *buildsystem.ArtifactRegistry) error {
	pomPath := filepath.Join(projectPath, "pom.xml")

	pomData, err := parsePomFile(pomPath)
	if err != nil {
		return err
	}

	// Se questo pom ha moduli, registrali ricorsivamente
	if pomData.Modules != nil && len(pomData.Modules.Module) > 0 {
		for _, module := range pomData.Modules.Module {
			modulePath := filepath.Join(projectPath, module)
			modulePomPath := filepath.Join(modulePath, "pom.xml")

			// Leggi il pom del sub-module
			modulePomData, err := parsePomFile(modulePomPath)
			if err != nil {
				continue // Ignora moduli non leggibili
			}

			// Registra il sub-module mappandolo al progetto parent
			if modulePomData.GroupId != "" && modulePomData.ArtifactId != "" {
				moduleId := makeIdentifier(modulePomData.GroupId, modulePomData.ArtifactId)
				registry.Register(moduleId, parentProjectName)
			}

			// Registra ricorsivamente eventuali sub-modules annidati
			_ = RegisterSubModules(modulePath, parentProjectName, registry)
		}
	}

	return nil
}

// BuildDependencyGraph costruisce il grafo delle dipendenze tra i progetti Maven selezionati
// Questa funzione sostituisce l'implementazione in compat.go eliminando il doppio parsing
func BuildDependencyGraph(projectNames []string, rootPath string) (map[string][]string, error) {
	registry := buildsystem.NewArtifactRegistry()
	projects := make(map[string]*buildsystem.Project)

	// Parse tutti i progetti in una sola passata
	for _, name := range projectNames {
		project, err := ParseProject(name, rootPath)
		if err != nil {
			return nil, err
		}
		projects[name] = project

		// Registra il progetto principale
		registry.Register(project.Identifier, name)

		// Registra i sub-modules per gestire dipendenze indirette
		if err := RegisterSubModules(project.Path, name, registry); err != nil {
			// Ignora errori di registrazione sub-modules (già gestito con Warning in compat.go)
			// Questo è accettabile perché alcuni sub-modules potrebbero non essere accessibili
		}
	}

	// Costruisci il grafo delle dipendenze tra i progetti selezionati
	dependencyGraph := make(map[string][]string)
	for name, project := range projects {
		dependencies := make([]string, 0)

		for _, depId := range project.Dependencies {
			// Controlla se la dipendenza è uno dei progetti selezionati
			if depProjectName, exists := registry.Lookup(depId); exists {
				// Evita self-dependency
				if depProjectName != name {
					dependencies = append(dependencies, depProjectName)
				}
			}
		}

		dependencyGraph[name] = dependencies
	}

	return dependencyGraph, nil
}

// TopologicalSort ordina i progetti in base alle loro dipendenze
// Restituisce l'ordine di esecuzione corretto o un errore se ci sono cicli
func TopologicalSort(graphMap map[string][]string) ([]string, error) {
	dependencyGraph := graph.DependencyGraph(graphMap)
	return dependencyGraph.TopologicalSort()
}
