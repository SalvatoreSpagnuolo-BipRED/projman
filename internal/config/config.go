package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
)

const ConfigDirName = "projman"
const ConfigFileName = "projman_config.json"

type Config struct {
	SelectedProjects []string `json:"selected_projects"`
}

func SaveSettings(c Config) error {
	rootPath, _ := os.UserConfigDir()

	// Creo la directory di configurazione se non esiste
	configDirPath := filepath.Join(rootPath, ConfigDirName)
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(configDirPath, 0755); err != nil {
			pterm.Error.Println("Errore durante la creazione della directory di configurazione in " + configDirPath)
			return nil
		}
	}

	// Salvo il file di configurazione
	savingPath := filepath.Join(rootPath, ConfigDirName, ConfigFileName)
	data, _ := json.MarshalIndent(c, "", "  ")
	if err := os.WriteFile(savingPath, data, 0644); err != nil {
		pterm.Error.Println("Errore durante il salvataggio della configurazione in " + savingPath)
		return nil
	}

	pterm.Success.Println("Configurazione salvata in " + savingPath)
	return nil

}

func LoadSettings() (Config, error) {
	rootPath, _ := os.UserConfigDir()
	loadPath := filepath.Join(rootPath, ConfigDirName, ConfigFileName)

	var cfg Config
	data, err := os.ReadFile(loadPath)
	if err != nil {
		return cfg, err
	}
	json.Unmarshal(data, &cfg)
	return cfg, nil
}

func CheckAndGetDirectory(directory string) (string, error) {
	info, err := os.Stat(directory)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", os.ErrInvalid
	}
	return directory, nil
}
