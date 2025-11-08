# Projman

Projman è un piccolo tool da linea di comando per gestire in batch più progetti Maven/Git contemporaneamente.
È pensato per scenari in cui lavori con numerosi repository locali e vuoi eseguire operazioni ripetitive (pull, build) su tutti i progetti selezionati.

## Caratteristiche

- Scansione automatica di progetti Maven (ricerca di `pom.xml`)
- Interfaccia interattiva per selezionare i progetti da gestire
- Comandi batch per Git (es. `git pull`)
- Comandi batch per Maven (es. `mvn install`), con flag per abilitare/disabilitare i test
- Configurazione persistente (lista dei progetti selezionati + root)
- Output migliorato con formattazione (usa `pterm` per colori e tabelle)

## Comandi principali

- `projman init <directory>`

  - Scansiona `<directory>` per trovare progetti Maven e permette di selezionarli
  - Salva la configurazione (root + progetti selezionati)

- `projman git update`

  - Esegue `git pull` su tutti i progetti selezionati

- `projman mvn install [--tests|-t]`

  - Esegue `mvn install` su tutti i progetti selezionati
  - Per default i test sono disabilitati; usa `--tests` o `-t` per abilitarli

- `projman help` (o `projman help <comando>`)
  - Mostra la guida formattata (il comando `help` ora usa `pterm` per una presentazione migliore)

## Esempi

```powershell
# Inizializza la configurazione con la cartella dei progetti
projman init C:\\Users\\me\\projects

# Aggiorna (git pull) tutti i progetti selezionati
projman git update

# Esegui mvn install su tutti i progetti (skip tests)
projman mvn install

# Esegui mvn install e lancia i test
projman mvn install --tests

# Visualizza la guida formattata
projman help
```

## File di configurazione

La configurazione è salvata in formato JSON nella directory di configurazione dell'utente:

- Windows: `%APPDATA%/projman/projman_config.json`
- Linux/macOS: `~/.config/projman/projman_config.json`

Esempio di contenuto:

```json
{
  "root_of_projects": "C:\\Users\\\\me\\\\projects",
  "selected_projects": ["project-a", "project-b"]
}
```

## Requisiti

- Git (nel PATH)
- Maven (nel PATH) per i comandi `mvn`
- Go (per compilare/installare il tool)

## Contribuire

Pull request e segnalazioni di issue sono benvenute. Apri una issue per discutere cambi significativi prima di inviare una PR.

## Licenza

Copyright © 2025 Salvatore Spagnuolo - BIP RED
