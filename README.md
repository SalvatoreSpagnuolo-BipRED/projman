# Projman

[![Go Version](https://img.shields.io/badge/Go-1.25.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Projman è un tool da linea di comando per gestire in batch più progetti Maven/Git contemporaneamente.
È pensato per scenari in cui lavori con numerosi repository locali e vuoi eseguire operazioni ripetitive (pull, build) su tutti i progetti selezionati.

## ✨ Caratteristiche

- 🔍 Scansione automatica di progetti Maven (ricerca di `pom.xml`)
- 🎯 Interfaccia interattiva per selezionare i progetti da gestire
- 🔄 Comandi batch per Git con gestione intelligente dei branch
  - Stash automatico delle modifiche non committate
  - Supporto per branch develop, deploy e feature
  - Merge automatico da develop
- 🏗️ Comandi batch per Maven con controllo dei test
- 💾 Configurazione persistente (lista dei progetti selezionati + root)
- 🎨 Output migliorato con formattazione colorata (usa `pterm`)

## 📋 Comandi principali

### `projman init <directory>`

Scansiona `<directory>` per trovare progetti Maven (cerca file `pom.xml`) e permette di selezionarli interattivamente.
Salva la configurazione (root + progetti selezionati) per usi futuri.

**Esempio:**
```bash
projman init ~/progetti
```

### `projman git update`

Esegue operazioni Git su tutti i progetti selezionati:
- **Stash automatico** delle modifiche non committate
- **Cambio branch** opzionale a `develop` (con selezione interattiva)
- **Pull/Merge** in base al tipo di branch:
  - Branch `develop`: esegue `git pull origin develop`
  - Branch `deploy/*`: esegue `git pull origin <branch-corrente>`
  - Altri branch: esegue `git fetch origin develop` e `git merge origin/develop`
- **Ripristino stash** automatico dopo l'aggiornamento

### `projman mvn install [--tests|-t]`

Esegue `mvn install` su tutti i progetti selezionati.
- Per default i test sono **disabilitati** (usa `-DskipTests=true`)
- Usa il flag `--tests` o `-t` per **abilitare** l'esecuzione dei test

**Esempi:**
```bash
# Install senza test
projman mvn install

# Install con test abilitati
projman mvn install --tests
```

### `projman help [comando]`

Mostra la guida formattata con una presentazione migliorata tramite `pterm`.
Puoi specificare un comando specifico per vedere informazioni dettagliate.

## 🚀 Esempi d'uso

```bash
# 1. Inizializza la configurazione con la cartella dei progetti
projman init ~/progetti

# 2. Aggiorna (git pull/merge) tutti i progetti selezionati
projman git update

# 3. Esegui mvn install su tutti i progetti (skip tests)
projman mvn install

# 4. Esegui mvn install e lancia i test
projman mvn install --tests
# oppure
projman mvn install -t

# 5. Visualizza la guida formattata
projman help
```

### Workflow tipico

1. **Inizializzazione**: `projman init ~/workspace`
2. **Aggiornamento quotidiano**: `projman git update`
3. **Build dei progetti**: `projman mvn install`

## ⚙️ File di configurazione

La configurazione è salvata in formato JSON nella directory di configurazione dell'utente:

| Sistema | Percorso |
|---------|----------|
| **Windows** | `%APPDATA%/projman/projman_config.json` |
| **Linux/macOS** | `~/.config/projman/projman_config.json` |

### Esempio di configurazione

```json
{
  "root_of_projects": "/Users/username/progetti",
  "selected_projects": [
    "project-a",
    "project-b",
    "project-c"
  ]
}
```

## 📦 Requisiti

- **Git** (nel PATH)
- **Maven** (nel PATH) per i comandi `mvn`
- **Go 1.25+** (per compilare/installare il tool)

## 🔧 Installazione

### Da sorgente

```bash
# Clona il repository
git clone https://github.com/SalvatoreSpagnuolo-BipRED/projman.git
cd projman

# Compila e installa
go build -o projman
go install

# Oppure usa direttamente
go run main.go help
```

### Binary precompilato

Scarica il binary per il tuo sistema dalla sezione [Releases](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases) e aggiungilo al PATH.

## 🤝 Contribuire

Pull request e segnalazioni di issue sono benvenute! 

1. Fai un fork del progetto
2. Crea un branch per la tua feature (`git checkout -b feature/AmazingFeature`)
3. Committa le tue modifiche (`git commit -m 'Add some AmazingFeature'`)
4. Pusha il branch (`git push origin feature/AmazingFeature`)
5. Apri una Pull Request

Per modifiche significative, apri prima una issue per discutere cosa vorresti cambiare.

## 📄 Licenza

Questo progetto è distribuito sotto licenza MIT. Vedi il file [LICENSE](LICENSE) per maggiori dettagli.

**Copyright © 2025 Salvatore Spagnuolo**

## 👤 Autore

**Salvatore Spagnuolo**
- Organizzazione: BIP RED

## 🔗 Link utili

- [Documentazione Cobra](https://github.com/spf13/cobra) - Framework CLI utilizzato
- [Documentazione Pterm](https://github.com/pterm/pterm) - Libreria per output colorato
- [AGENTS.md](AGENTS.md) - Guida per AI agents che lavorano su questo progetto
