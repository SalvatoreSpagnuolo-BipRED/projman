package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
)

// Costanti per la gestione della configurazione
const (
	// ConfigDirName è il nome della directory che contiene la configurazione
	ConfigDirName = "projman"
	// ConfigFileName è il nome del file di configurazione JSON
	ConfigFileName = "projman_config.json"
	// ConfigDirPermissions sono i permessi per la directory di configurazione
	ConfigDirPermissions = 0755
	// ConfigFilePermissions sono i permessi per il file di configurazione
	ConfigFilePermissions = 0644
)

// Config rappresenta la struttura della configurazione di projman
type Config struct {
	RootOfProjects   string   `json:"root_of_projects"`  // Percorso root contenente tutti i progetti
	SelectedProjects []string `json:"selected_projects"` // Lista dei progetti selezionati dall'utente
}

// SaveSettings salva la configurazione nel file JSON di sistema
func SaveSettings(c Config) error {
	// Ottiene la directory di configurazione dell'utente
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		pterm.Error.Println("Impossibile determinare la directory di configurazione:", err)
		return fmt.Errorf("impossibile ottenere la directory di configurazione: %w", err)
	}

	// Crea la directory di configurazione se non esiste
	configDirPath := filepath.Join(userConfigDir, ConfigDirName)
	if err := ensureConfigDirExists(configDirPath); err != nil {
		return err
	}

	// Serializza la configurazione in JSON con indentazione
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		pterm.Error.Println("Errore durante la serializzazione della configurazione:", err)
		return fmt.Errorf("impossibile serializzare la configurazione: %w", err)
	}

	// Salva il file di configurazione
	savingPath := filepath.Join(configDirPath, ConfigFileName)
	if err := os.WriteFile(savingPath, data, ConfigFilePermissions); err != nil {
		pterm.Error.Println("Errore durante il salvataggio della configurazione:", err)
		return fmt.Errorf("impossibile salvare il file di configurazione: %w", err)
	}

	pterm.Success.Println("Configurazione salvata in", savingPath)
	return nil
}

// LoadSettings carica la configurazione dal file JSON di sistema
func LoadSettings() (Config, error) {
	var cfg Config

	// Ottiene la directory di configurazione dell'utente
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return cfg, fmt.Errorf("impossibile ottenere la directory di configurazione: %w", err)
	}

	// Costruisce il percorso completo del file di configurazione
	loadPath := filepath.Join(userConfigDir, ConfigDirName, ConfigFileName)

	// Legge il file di configurazione
	data, err := os.ReadFile(loadPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, fmt.Errorf("file di configurazione non trovato: esegui 'projman init <directory>' per inizializzare")
		}
		return cfg, fmt.Errorf("impossibile leggere il file di configurazione: %w", err)
	}

	// Deserializza il JSON nella struttura Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("file di configurazione non valido: %w", err)
	}

	return cfg, nil
}

// CheckAndGetDirectory verifica che il percorso fornito sia una directory valida
func CheckAndGetDirectory(directory string) (string, error) {
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("la directory '%s' non esiste", directory)
		}
		return "", fmt.Errorf("impossibile accedere alla directory '%s': %w", directory, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("'%s' non è una directory", directory)
	}

	// Restituisce il percorso assoluto
	absPath, err := filepath.Abs(directory)
	if err != nil {
		return directory, nil // Fallback al percorso originale
	}

	return absPath, nil
}

// ensureConfigDirExists crea la directory di configurazione se non esiste
func ensureConfigDirExists(configDirPath string) error {
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDirPath, ConfigDirPermissions); err != nil {
			pterm.Error.Println("Impossibile creare la directory di configurazione:", configDirPath)
			return fmt.Errorf("impossibile creare la directory di configurazione: %w", err)
		}
		pterm.Info.Println("Directory di configurazione creata:", configDirPath)
	}
	return nil
}
