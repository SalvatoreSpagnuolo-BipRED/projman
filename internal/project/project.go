package project

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// MavenProjectFile è il nome del file che identifica un progetto Maven
	MavenProjectFile = "pom.xml"
)

// Project rappresenta un progetto Maven con nome e percorso
type Project struct {
	Name string // Nome del progetto (nome della directory)
	Path string // Percorso assoluto del progetto
}

// Discover scansiona la directory root e restituisce tutti i progetti Maven trovati.
// Un progetto Maven viene identificato dalla presenza di un file pom.xml nella directory.
func Discover(root string) ([]Project, error) {
	// Valida che la directory root esista
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("la directory root '%s' non esiste", root)
		}
		return nil, fmt.Errorf("impossibile accedere alla directory root '%s': %w", root, err)
	}

	// Legge le entry nella directory root
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("impossibile leggere la directory root '%s': %w", root, err)
	}

	// Scansiona ogni entry alla ricerca di progetti Maven
	projects := make([]Project, 0)
	for _, entry := range entries {
		// Considera solo le directory
		if !entry.IsDir() {
			continue
		}

		// Verifica se la directory contiene un pom.xml
		projectPath := filepath.Join(root, entry.Name())
		if isMavenProject(projectPath) {
			projects = append(projects, Project{
				Name: entry.Name(),
				Path: projectPath,
			})
		}
	}

	return projects, nil
}

// Names estrae i nomi di tutti i progetti dalla lista fornita
func Names(projs []Project) []string {
	names := make([]string, len(projs))
	for i, p := range projs {
		names[i] = p.Name
	}
	return names
}

// Filter filtra i progetti in base ai nomi selezionati.
// Se la lista selected è vuota, restituisce tutti i progetti.
func Filter(all []Project, selected []string) []Project {
	// Se non ci sono filtri, restituisce tutti i progetti
	if len(selected) == 0 {
		return all
	}

	// Crea una mappa per lookup O(1) dei nomi selezionati
	selectedMap := make(map[string]bool, len(selected))
	for _, name := range selected {
		selectedMap[name] = true
	}

	// Filtra i progetti che sono nella lista dei selezionati
	filtered := make([]Project, 0, len(selected))
	for _, proj := range all {
		if selectedMap[proj.Name] {
			filtered = append(filtered, proj)
		}
	}

	return filtered
}

// isMavenProject verifica se una directory contiene un file pom.xml
func isMavenProject(projectPath string) bool {
	pomPath := filepath.Join(projectPath, MavenProjectFile)
	_, err := os.Stat(pomPath)
	return err == nil
}
