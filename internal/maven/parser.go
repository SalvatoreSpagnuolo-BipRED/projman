// Package maven fornisce funzionalità per l'analisi di progetti Maven
package maven

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/buildsystem"
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

// Parser implementa l'interfaccia buildsystem.Parser per progetti Maven
type Parser struct{}

// NewParser crea un nuovo parser Maven
func NewParser() *Parser {
	return &Parser{}
}

// ParseProject analizza un progetto Maven e restituisce le sue informazioni
func (p *Parser) ParseProject(projectName, rootPath string) (*buildsystem.Project, error) {
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

// GetIdentifier restituisce l'identificatore univoco del progetto (groupId:artifactId)
func (p *Parser) GetIdentifier(project *buildsystem.Project) string {
	return project.Identifier
}

// SupportsModules indica che Maven supporta i moduli
func (p *Parser) SupportsModules() bool {
	return true
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
