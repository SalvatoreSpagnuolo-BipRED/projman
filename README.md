# Projman

[![Go Version](https://img.shields.io/badge/Go-1.25.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.1.0-brightgreen.svg)](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases)

**Projman** è un tool CLI scritto in Go per gestire in batch multipli progetti Maven/Git.
È pensato per scenari in cui lavori con numerosi repository locali e vuoi eseguire operazioni ripetitive (pull, build, test) su tutti i progetti selezionati con un solo comando.

## ✨ Caratteristiche
- 🔍 **Scansione automatica** di progetti Maven (ricerca di `pom.xml`)
- 🎯 **Interfaccia interattiva** per selezionare i progetti da gestire
- 🔄 **Comandi batch per Git** con gestione intelligente dei branch:
- 🔄 Comandi batch per Git con gestione intelligente dei branch
  - Supporto per branch `develop`, `deploy/*` e feature branch
  - Merge automatico da develop sui feature branch
  - Ripristino automatico dello stash
- 🏗️ **Comandi batch per Maven** con controllo dei test
- 💾 **Configurazione persistente** (JSON con lista progetti e root path)
- 🎨 **Output formattato** con colori e tabelle interattive (powered by `pterm`)
- 📊 **Report dettagliati** con statistiche di successo/fallimento
- 🚀 **Architettura modulare** seguendo le best practice Go

## 🏗️ Struttura del Progetto

```
projman/
├── main.go                    # Entry point dell'applicazione
├── cmd/                       # Comandi CLI (Cobra)
│   ├── root.go               # Comando root
│   ├── init.go               # Inizializzazione configurazione
│   ├── help.go               # Help personalizzato
│   ├── git/                  # Comandi Git
│   │   ├── git.go           # Parent command
│   │   └── update.go        # Git update con gestione branch
│   └── mvn/                  # Comandi Maven
│       ├── mvn.go           # Parent command
│       └── install.go       # Maven install
└── internal/                  # Package interni
    ├── config/               # Gestione configurazione
    ├── project/              # Discovery e filtering progetti
    ├── exec/                 # Wrapper per esecuzione comandi
    └── ui/                   # Componenti UI interattivi
```

> **Nota**: Questa struttura segue le [best practice Go](https://github.com/golang-standards/project-layout) con package `internal` e organizzazione modulare dei comandi.
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
### 📋 Workflow tipico

```bash
# 1️⃣ Inizializza la configurazione (una sola volta)
projman init ~/workspace

# 2️⃣ Aggiorna tutti i progetti (quotidiano)
projman git update

# 3️⃣ Build di tutti i progetti
projman mvn install

# 4️⃣ Build con test (opzionale)
projman mvn install --tests
```

### 🎬 Demo

![Projman Demo](https://via.placeholder.com/800x400?text=Projman+CLI+Demo)

*Screenshot dell'interfaccia interattiva con selezione progetti e output colorato*
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
### Opzione 1: Binary precompilato (Consigliato)

Scarica il binary per il tuo sistema dalla sezione [Releases](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases):

```bash
# Linux/macOS
curl -L https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases/latest/download/projman -o projman
chmod +x projman
sudo mv projman /usr/local/bin/

# Windows (PowerShell)
# Scarica da GitHub Releases e aggiungi al PATH
```

### Opzione 2: Go install

Se hai Go installato:

```bash
## ❓ FAQ

### Dove vengono salvati i progetti selezionati?
La configurazione è salvata in `~/.config/projman/projman_config.json` (Linux/macOS) o `%APPDATA%/projman/projman_config.json` (Windows).

### Posso usare Projman con altri build system?
Attualmente Projman supporta solo Maven, ma è progettato per essere estensibile. Contributi per supportare Gradle, npm, ecc. sono benvenuti!

### Come funziona la gestione dei branch?
Projman usa una strategia intelligente:
- **Branch `develop`**: esegue `git pull origin develop`
- **Branch `deploy/*`**: esegue `git pull origin <branch-corrente>`
## 🗺️ Roadmap

### v1.1.0 (Corrente) ✅
- [x] Refactoring completo del codice
- [x] Struttura modulare con subpackage
- [x] Documentazione completa (AGENTS.md, README.md)
- [x] Gestione errori migliorata
- [x] Output formattato con statistiche

### v1.2.0 (Prossima)
- [ ] Comando `git status` per vedere lo stato di tutti i progetti
- [ ] Comando `git commit` batch
- [ ] Supporto per file `.projmanignore`
- [ ] Logging delle operazioni in file
- [ ] Configurazione per-progetto

### v2.0.0 (Futuro)
- [ ] Supporto Gradle
- [ ] Supporto npm/yarn
- [ ] Esecuzione parallela delle operazioni
- [ ] UI web per gestione progetti
- [ ] Plugin system

- **Altri branch**: esegue `git fetch origin develop` + `git merge origin/develop`

Questo progetto è distribuito sotto **licenza MIT**. Vedi il file [LICENSE](LICENSE) per maggiori dettagli.

```
MIT License
Projman esegue automaticamente `git stash` prima di operazioni Git e ripristina le modifiche al termine. I tuoi file locali non committati sono al sicuro!
Copyright © 2025 Salvatore Spagnuolo

Permission is hereby granted, free of charge, to any person obtaining a copy...
```
### Cosa succede se un comando fallisce su un progetto?
Projman continua con i progetti successivi e mostra un report finale con successi/fallimenti.

go install github.com/SalvatoreSpagnuolo-BipRED/projman@latest
- 🏢 Organizzazione: BIP RED
- 📧 Email: [salvatore.spagnuolo@bipred.com](mailto:salvatore.spagnuolo@bipred.com)
- 🐙 GitHub: [@SalvatoreSpagnuolo-BipRED](https://github.com/SalvatoreSpagnuolo-BipRED)

## 🙏 Riconoscimenti

Questo progetto utilizza le seguenti librerie open source:
- [Cobra](https://github.com/spf13/cobra) - Framework CLI potente e flessibile
- [Pterm](https://github.com/pterm/pterm) - Libreria per output terminal colorato e interattivo

## 🔗 Link Utili
Le contribuzioni sono benvenute! Ecco come puoi aiutare:
- 📖 [Documentazione Completa](https://github.com/SalvatoreSpagnuolo-BipRED/projman/wiki)
- 🤖 [AGENTS.md](AGENTS.md) - Guida per AI agents
- 📝 [REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md) - Storia del refactoring
- 🐛 [Issues](https://github.com/SalvatoreSpagnuolo-BipRED/projman/issues) - Segnala bug o richiedi feature
- 📦 [Releases](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases) - Download binary precompilati
- 💬 [Discussions](https://github.com/SalvatoreSpagnuolo-BipRED/projman/discussions) - Discussioni e Q&A
### 🐛 Segnala Bug
## ⭐ Supporta il Progetto

Se trovi Projman utile, considera di:
- ⭐ Mettere una stella su GitHub
- 🐛 Segnalare bug o suggerire miglioramenti
- 🔧 Contribuire con codice o documentazione
- 📢 Condividere il progetto con colleghi

---

**Made with ❤️ by Salvatore Spagnuolo**

*Ultimo aggiornamento: Novembre 2025 - Versione 1.1.0*
- Cosa è successo invece
- Output di `projman --version`

### ✨ Proponi Feature
Apri una [issue](https://github.com/SalvatoreSpagnuolo-BipRED/projman/issues) con:
- Descrizione della funzionalità
- Casi d'uso
- Esempio di utilizzo proposto

### 🔧 Contribuisci al Codice
### Opzione 3: Da sorgente
1. **Fork** il progetto
2. **Crea branch** per la tua feature:
   ```bash
   git checkout -b feature/AmazingFeature
   ```
3. **Sviluppa** seguendo le convenzioni in [AGENTS.md](AGENTS.md)
4. **Testa** le modifiche:
   ```bash
   go build -o projman
   go fmt ./...
   go vet ./...
   ./projman help
   ```
5. **Commit** con messaggi descrittivi:
   ```bash
   git commit -m 'feat: Add amazing feature'
   ```
6. **Push** al tuo fork:
   ```bash
   git push origin feature/AmazingFeature
   ```
7. **Apri Pull Request** verso `develop`
```
### 📚 Convenzioni
- Leggi [AGENTS.md](AGENTS.md) per best practice
- Segui lo stile Go standard (`gofmt`)
- Aggiungi commenti esportati per funzioni pubbliche
- Mantieni funzioni piccole (<50 righe)
- Scrivi messaggi di commit chiari

Per modifiche significative, apri prima una issue per discuterne.
### Workflow tipico

# Installa globalmente (opzionale)
sudo mv projman /usr/local/bin/
1. **Inizializzazione**: `projman init ~/workspace`
# Oppure esegui direttamente
./projman help

## ⚙️ File di configurazione
### Verifica installazione
La configurazione è salvata in formato JSON nella directory di configurazione dell'utente:
```bash
projman --version
projman help
```
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
