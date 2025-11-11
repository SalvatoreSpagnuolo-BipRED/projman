# Projman

[![Go Version](https://img.shields.io/badge/Go-1.25.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-0.8.0-brightgreen.svg)](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases)

**Projman** è un tool CLI scritto in Go per gestire in batch multipli progetti Maven/Git.

## ✨ Caratteristiche
- 🔍 Scansione automatica di progetti Maven (ricerca di `pom.xml`)
- 🎯 Interfaccia interattiva per selezionare i progetti da gestire
- 🔄 Comandi batch per Git con gestione intelligente dei branch (develop, deploy/*, feature)
- 🏗️ Comandi batch per Maven con ordinamento automatico delle dipendenze
- 💾 Configurazione persistente (JSON)
- 🎨 Output formattato con colori e tabelle interattive

## 📋 Comandi

### `projman init <directory>`
Scansiona `<directory>` per trovare progetti Maven e permette di selezionarli interattivamente.
Salva la configurazione per usi futuri.

```bash
projman init ~/progetti
```

### `projman git update`
Esegue operazioni Git su tutti i progetti selezionati:
- Stash automatico delle modifiche
- Cambio branch opzionale
- Pull/Merge in base al tipo di branch:
  - `develop`: `git pull origin develop`
  - `deploy/*`: `git pull origin <branch-corrente>`
  - Altri: `git fetch origin develop` + `git merge origin/develop`
- Ripristino stash automatico

### `projman mvn install [--tests|-t]`
Esegue `mvn install` con ordinamento automatico delle dipendenze.
Di default i test sono disabilitati. Usa `--tests` o `-t` per abilitarli.

```bash
# Install senza test
projman mvn install

# Install con test
projman mvn install --tests
```

## 📦 Requisiti
- **Git** (nel PATH)
- **Maven** (nel PATH)
- **Go 1.25+** (solo per compilare da sorgente)

## 🔧 Installazione

### Binary precompilato (Consigliato)
Scarica il binary dalla sezione [Releases](https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases):

```bash
# Linux/macOS
curl -L https://github.com/SalvatoreSpagnuolo-BipRED/projman/releases/latest/download/projman -o projman
chmod +x projman
sudo mv projman /usr/local/bin/

# Windows (PowerShell)
# Scarica da GitHub Releases e aggiungi al PATH
```

### Da sorgente
```bash
git clone https://github.com/SalvatoreSpagnuolo-BipRED/projman.git
cd projman
go build -o projman
sudo mv projman /usr/local/bin/
```

### Verifica installazione
```bash
projman --version
projman help
```

## ⚙️ Configurazione
La configurazione è salvata in:
- **Windows**: `%APPDATA%/projman/projman_config.json`
- **Linux/macOS**: `~/.config/projman/projman_config.json`

Esempio:
```json
{
  "root_of_projects": "/Users/username/progetti",
  "selected_projects": ["project-a", "project-b", "project-c"]
}
```

## 📄 Licenza
Questo progetto è distribuito sotto licenza MIT. Vedi il file [LICENSE](LICENSE) per maggiori dettagli.

**Copyright © 2025 Salvatore Spagnuolo**

## 👤 Autore
**Salvatore Spagnuolo**
- Organizzazione: BIP RED
- Email: [salvatore.spagnuolo@bipred.com](mailto:salvatore.spagnuolo@bipred.com)
- GitHub: [@SalvatoreSpagnuolo-BipRED](https://github.com/SalvatoreSpagnuolo-BipRED)

---

**Made with ❤️ by Salvatore Spagnuolo**

