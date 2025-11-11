package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/project"
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

// ProfileConfig rappresenta la struttura che contiene tutti i profili e il profilo corrente
type ProfileConfig struct {
	CurrentProfile string            `json:"current_profile"` // Nome del profilo attualmente attivo
	Profiles       map[string]Config `json:"profiles"`        // Mappa nome_profilo -> Config
}

// SaveSettings salva la configurazione nel file JSON di sistema
func SaveSettings(c Config) error {
	return SaveProfile("default", c)
}

// SaveProfile salva un profilo specifico nella configurazione
func SaveProfile(profileName string, c Config) error {
	// Carica il ProfileConfig esistente o creane uno nuovo
	profileCfg, err := loadProfileConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Inizializza la mappa se non esiste
	if profileCfg.Profiles == nil {
		profileCfg.Profiles = make(map[string]Config)
	}

	// Salva il profilo
	profileCfg.Profiles[profileName] = c

	// Se non c'è un profilo corrente, imposta questo come corrente
	if profileCfg.CurrentProfile == "" {
		profileCfg.CurrentProfile = profileName
	}

	return saveProfileConfig(profileCfg)
}

// LoadSettings carica la configurazione dal file JSON di sistema
func LoadSettings() (Config, error) {
	profileCfg, err := loadProfileConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("file di configurazione non trovato: esegui 'projman init <nome-profilo> <directory>' per inizializzare")
		}
		return Config{}, fmt.Errorf("impossibile leggere il file di configurazione: %w", err)
	}

	if profileCfg.CurrentProfile == "" {
		return Config{}, fmt.Errorf("nessun profilo attivo: esegui 'projman profile use <nome-profilo>' per selezionare un profilo")
	}

	cfg, exists := profileCfg.Profiles[profileCfg.CurrentProfile]
	if !exists {
		return Config{}, fmt.Errorf("profilo corrente '%s' non trovato", profileCfg.CurrentProfile)
	}

	return cfg, nil
}

// GetCurrentProfile restituisce il nome del profilo corrente
func GetCurrentProfile() (string, error) {
	profileCfg, err := loadProfileConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file di configurazione non trovato")
		}
		return "", err
	}
	return profileCfg.CurrentProfile, nil
}

// SetCurrentProfile imposta il profilo corrente
func SetCurrentProfile(profileName string) error {
	profileCfg, err := loadProfileConfig()
	if err != nil {
		return err
	}

	if _, exists := profileCfg.Profiles[profileName]; !exists {
		return fmt.Errorf("profilo '%s' non trovato", profileName)
	}

	profileCfg.CurrentProfile = profileName
	return saveProfileConfig(profileCfg)
}

// ListProfiles restituisce la lista di tutti i profili disponibili
func ListProfiles() ([]string, string, error) {
	profileCfg, err := loadProfileConfig()
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, "", nil
		}
		return nil, "", err
	}

	profiles := make([]string, 0, len(profileCfg.Profiles))
	for name := range profileCfg.Profiles {
		profiles = append(profiles, name)
	}

	return profiles, profileCfg.CurrentProfile, nil
}

// DeleteProfile elimina un profilo dalla configurazione
func DeleteProfile(profileName string) error {
	profileCfg, err := loadProfileConfig()
	if err != nil {
		return err
	}

	if _, exists := profileCfg.Profiles[profileName]; !exists {
		return fmt.Errorf("profilo '%s' non trovato", profileName)
	}

	delete(profileCfg.Profiles, profileName)

	// Se era il profilo corrente, resetta il current profile
	if profileCfg.CurrentProfile == profileName {
		profileCfg.CurrentProfile = ""
		// Se ci sono altri profili, imposta il primo disponibile
		for name := range profileCfg.Profiles {
			profileCfg.CurrentProfile = name
			break
		}
	}

	return saveProfileConfig(profileCfg)
}

// loadProfileConfig carica il ProfileConfig dal file di sistema
func loadProfileConfig() (ProfileConfig, error) {
	var profileCfg ProfileConfig

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return profileCfg, fmt.Errorf("impossibile ottenere la directory di configurazione: %w", err)
	}

	loadPath := filepath.Join(userConfigDir, ConfigDirName, ConfigFileName)

	data, err := os.ReadFile(loadPath)
	if err != nil {
		return profileCfg, err
	}

	if err := json.Unmarshal(data, &profileCfg); err != nil {
		return profileCfg, fmt.Errorf("file di configurazione non valido: %w", err)
	}

	return profileCfg, nil
}

// saveProfileConfig salva il ProfileConfig nel file di sistema
func saveProfileConfig(profileCfg ProfileConfig) error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		pterm.Error.Println("Impossibile determinare la directory di configurazione:", err)
		return fmt.Errorf("impossibile ottenere la directory di configurazione: %w", err)
	}

	configDirPath := filepath.Join(userConfigDir, ConfigDirName)
	if err := ensureConfigDirExists(configDirPath); err != nil {
		return err
	}

	data, err := json.MarshalIndent(profileCfg, "", "  ")
	if err != nil {
		pterm.Error.Println("Errore durante la serializzazione della configurazione:", err)
		return fmt.Errorf("impossibile serializzare la configurazione: %w", err)
	}

	savingPath := filepath.Join(configDirPath, ConfigFileName)
	if err := os.WriteFile(savingPath, data, ConfigFilePermissions); err != nil {
		pterm.Error.Println("Errore durante il salvataggio della configurazione:", err)
		return fmt.Errorf("impossibile salvare il file di configurazione: %w", err)
	}

	return nil
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

// LoadAndValidateConfig carica e valida la configurazione
func LoadAndValidateConfig() (*Config, error) {
	cfg, err := LoadSettings()
	if err != nil {
		pterm.Error.Println("Errore nel caricamento della configurazione:", err)
		pterm.Info.Println("Esegui prima 'projman init <directory>' per configurare i progetti")
		return nil, err
	}
	return &cfg, nil
}

// SelectProjectsToUpdate mostra un prompt per selezionare i progetti da aggiornare
func SelectProjectsToUpdate(cfg *Config) ([]string, error) {
	projUri, err := project.Discover(cfg.RootOfProjects)
	if err != nil {
		pterm.Error.Println("Errore nella scansione dei progetti:", err)
		return nil, err
	}

	names := project.Names(projUri)
	selectedNames, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(names).
		WithDefaultOptions(cfg.SelectedProjects).
		Show("Seleziona i progetti da aggiornare:")
	if err != nil {
		pterm.Error.Println("Errore nella selezione dei progetti:", err)
		return nil, err
	}

	return selectedNames, nil
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
