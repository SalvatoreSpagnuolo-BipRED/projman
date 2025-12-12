// Package buildsystem fornisce interfacce comuni per diversi sistemi di build
package buildsystem

// Project rappresenta un progetto generico indipendente dal build system
type Project struct {
	Name         string   // Nome del progetto (directory)
	Path         string   // Percorso completo del progetto
	Identifier   string   // Identificatore univoco del progetto (es. groupId:artifactId per Maven)
	Dependencies []string // Lista di identificatori delle dipendenze
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
