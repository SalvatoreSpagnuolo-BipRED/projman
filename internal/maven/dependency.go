package maven

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// PomProject rappresenta le informazioni base di un progetto Maven
type PomProject struct {
	Name         string   // Nome del progetto (directory)
	Path         string   // Percorso completo del progetto
	GroupId      string   // GroupId del progetto
	ArtifactId   string   // ArtifactId del progetto
	Dependencies []string // Lista di chiavi "groupId:artifactId" delle dipendenze
}

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

// ParsePomFile legge e analizza un file pom.xml, restituendo le informazioni del progetto
func ParsePomFile(projectName, rootPath string) (*PomProject, error) {
	pomPath := filepath.Join(rootPath, projectName, "pom.xml")

	data, err := os.ReadFile(pomPath)
	if err != nil {
		return nil, fmt.Errorf("impossibile leggere %s: %w", pomPath, err)
	}

	var p pom
	if err := xml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("errore parsing XML di %s: %w", pomPath, err)
	}

	// Usa il groupId del parent se non specificato
	groupId := p.GroupId
	if groupId == "" && p.Parent != nil {
		groupId = p.Parent.GroupId
	}

	pomProject := &PomProject{
		Name:         projectName,
		Path:         filepath.Join(rootPath, projectName),
		GroupId:      groupId,
		ArtifactId:   p.ArtifactId,
		Dependencies: make([]string, 0),
	}

	// Estrai le dipendenze
	if p.Dependencies != nil {
		for _, dep := range p.Dependencies.Dependency {
			key := fmt.Sprintf("%s:%s", dep.GroupId, dep.ArtifactId)
			pomProject.Dependencies = append(pomProject.Dependencies, key)
		}
	}

	// Se il progetto ha moduli, analizza anche quelli per trovare dipendenze
	if p.Modules != nil && len(p.Modules.Module) > 0 {
		for _, module := range p.Modules.Module {
			moduleDeps, err := extractModuleDependencies(filepath.Join(pomProject.Path, module, "pom.xml"))
			if err == nil {
				pomProject.Dependencies = append(pomProject.Dependencies, moduleDeps...)
			}
		}
	}

	return pomProject, nil
}

// extractModuleDependencies estrae le dipendenze da un modulo figlio
func extractModuleDependencies(modulePomPath string) ([]string, error) {
	data, err := os.ReadFile(modulePomPath)
	if err != nil {
		return nil, err
	}

	var p pom
	if err := xml.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	deps := make([]string, 0)
	if p.Dependencies != nil {
		for _, dep := range p.Dependencies.Dependency {
			key := fmt.Sprintf("%s:%s", dep.GroupId, dep.ArtifactId)
			deps = append(deps, key)
		}
	}

	return deps, nil
}

// BuildDependencyGraph costruisce il grafo delle dipendenze tra i progetti selezionati
func BuildDependencyGraph(projectNames []string, rootPath string) (map[string][]string, error) {
	// Mappa: nome progetto -> PomProject
	projects := make(map[string]*PomProject)

	// Mappa: "groupId:artifactId" -> nome progetto
	artifactToProject := make(map[string]string)

	// Prima passata: parse di tutti i pom.xml
	for _, name := range projectNames {
		pomProj, err := ParsePomFile(name, rootPath)
		if err != nil {
			return nil, fmt.Errorf("errore analisi progetto %s: %w", name, err)
		}
		projects[name] = pomProj

		// Registra il progetto nella mappa degli artifact
		key := fmt.Sprintf("%s:%s", pomProj.GroupId, pomProj.ArtifactId)
		artifactToProject[key] = name
	}

	// Seconda passata: costruisci il grafo delle dipendenze
	// La mappa risultante indica: progetto -> lista di progetti da cui dipende
	graph := make(map[string][]string)

	for name, pomProj := range projects {
		graph[name] = make([]string, 0)

		for _, depKey := range pomProj.Dependencies {
			// Controlla se la dipendenza è uno dei progetti selezionati
			if depProjectName, exists := artifactToProject[depKey]; exists {
				graph[name] = append(graph[name], depProjectName)
			}
		}
	}

	return graph, nil
}

// TopologicalSort ordina i progetti in base alle loro dipendenze
// Restituisce l'ordine di esecuzione corretto o un errore se ci sono cicli
func TopologicalSort(graph map[string][]string) ([]string, error) {
	// Calcola l'in-degree (numero di progetti che dipendono da questo) per ogni nodo
	inDegree := make(map[string]int)

	// Inizializza tutti i nodi con in-degree 0
	for node := range graph {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
	}

	// Calcola l'in-degree: per ogni progetto, incrementa l'in-degree delle sue dipendenze
	// Se A dipende da B, allora dobbiamo installare B prima di A
	// Quindi quando processiamo A (che ha B nelle sue dipendenze), incrementiamo in-degree di A
	for node, deps := range graph {
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
	result := make([]string, 0)

	// Algoritmo di Kahn per il topological sort
	for len(queue) > 0 {
		// Prendi il primo nodo dalla coda
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Per ogni progetto che dipende dal nodo corrente, decrementa il suo in-degree
		for node, deps := range graph {
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
	if len(result) != len(graph) {
		return nil, fmt.Errorf("dipendenze circolari rilevate tra i progetti")
	}

	return result, nil
}
