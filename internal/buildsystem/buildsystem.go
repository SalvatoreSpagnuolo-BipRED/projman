// Package buildsystem fornisce interfacce comuni per diversi sistemi di build
package buildsystem

import "github.com/SalvatoreSpagnuolo-BipRED/projman/internal/graph"

// Project rappresenta un progetto generico indipendente dal build system
type Project struct {
	Name         string   // Nome del progetto (directory)
	Path         string   // Percorso completo del progetto
	Identifier   string   // Identificatore univoco del progetto (es. groupId:artifactId per Maven)
	Dependencies []string // Lista di identificatori delle dipendenze
}

// Parser è l'interfaccia per parser di build system specifici
type Parser interface {
	// ParseProject analizza un progetto e restituisce le sue informazioni
	ParseProject(projectName, rootPath string) (*Project, error)

	// GetIdentifier restituisce l'identificatore univoco del progetto
	GetIdentifier(project *Project) string

	// SupportsModules indica se il build system supporta moduli/subproject
	SupportsModules() bool
}

// ArtifactRegistry mantiene la mappatura tra identificatori di artifact e nomi di progetti
type ArtifactRegistry struct {
	artifactToProject map[string]string
}

// NewArtifactRegistry crea un nuovo registro di artifact
func NewArtifactRegistry() *ArtifactRegistry {
	return &ArtifactRegistry{
		artifactToProject: make(map[string]string),
	}
}

// Register registra un artifact associandolo a un progetto
func (r *ArtifactRegistry) Register(artifactId, projectName string) {
	r.artifactToProject[artifactId] = projectName
}

// Lookup cerca il progetto associato a un artifact
func (r *ArtifactRegistry) Lookup(artifactId string) (string, bool) {
	project, exists := r.artifactToProject[artifactId]
	return project, exists
}

// GraphBuilder costruisce grafi di dipendenze da progetti
type GraphBuilder struct {
	parser   Parser
	registry *ArtifactRegistry
}

// NewGraphBuilder crea un nuovo builder di grafi con un parser specifico
func NewGraphBuilder(parser Parser) *GraphBuilder {
	return &GraphBuilder{
		parser:   parser,
		registry: NewArtifactRegistry(),
	}
}

// BuildGraph costruisce il grafo delle dipendenze per i progetti specificati
func (b *GraphBuilder) BuildGraph(projectNames []string, rootPath string) (graph.DependencyGraph, error) {
	projects := make(map[string]*Project)

	// Prima passata: parse di tutti i progetti e registrazione degli artifact
	for _, name := range projectNames {
		project, err := b.parser.ParseProject(name, rootPath)
		if err != nil {
			return nil, err
		}
		projects[name] = project

		// Registra il progetto principale
		identifier := b.parser.GetIdentifier(project)
		b.registry.Register(identifier, name)
	}

	// Seconda passata: costruisci il grafo delle dipendenze
	dependencyGraph := graph.NewDependencyGraph()

	for name, project := range projects {
		dependencies := make([]string, 0)

		for _, depId := range project.Dependencies {
			// Controlla se la dipendenza è uno dei progetti selezionati
			if depProjectName, exists := b.registry.Lookup(depId); exists {
				// Evita self-dependency
				if depProjectName != name {
					dependencies = append(dependencies, depProjectName)
				}
			}
		}

		dependencyGraph.AddNode(name, dependencies)
	}

	return dependencyGraph, nil
}

// GetRegistry restituisce il registro degli artifact
func (b *GraphBuilder) GetRegistry() *ArtifactRegistry {
	return b.registry
}
