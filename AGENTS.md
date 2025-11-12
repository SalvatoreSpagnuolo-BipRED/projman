# AGENTS.md - Guida per AI Agents

Guida tecnica per AI agents che lavorano su **Projman**, un tool CLI Go per la gestione batch di progetti Maven/Git.

## ⚠️ REGOLE OPERATIVE

**Esegui SOLO quanto richiesto:**

1. ⛔ NO funzionalità extra non richieste
2. ⛔ NO file di report .md (REFACTORING_SUMMARY, CHANGELOG, ecc.)
3. ⛔ NO over-engineering - mantieni la semplicità
4. ⛔ NO riorganizzazione struttura senza richiesta esplicita
5. ✅ Aggiorna README e AGENTS solo per modifiche a funzionalità/struttura

## 🏗️ Architettura

```
projman/
├── main.go                             # Entry point
├── cmd/                                # Comandi CLI (Cobra)
│   ├── root.go, init.go, help.go       # Comandi base
│   ├── list.go, use.go, delete.go      # Gestione profili
│   ├── git/                            # Sottocomandi Git
│   │   ├── git.go                      # Parent command
│   │   └── update.go                   # Git pull/merge intelligente
│   └── mvn/                            # Sottocomandi Maven
│       ├── mvn.go                      # Parent command
│       └── install.go                  # Maven install con dep sorting
└── internal/                           # Package interni
    ├── config/config.go                # Gestione configurazione JSON multi-profilo
    ├── project/project.go              # Discovery progetti Maven
    ├── graph/graph.go                  # Ordinamento topologico (Kahn)
    ├── buildsystem/buildsystem.go      # Interfacce comuni
    ├── maven/                          # Parser Maven
    │   ├── parser.go                   # Parsing pom.xml
    │   └── executor/                   # Executor Maven-specific
    │       └── executor.go             # Parsing output Maven + logging
    ├── exec/exec.go                    # Esecuzione comandi generici
    └── ui/multiselect.go               # Selezione interattiva
```

**Tecnologie**: Go 1.25.4+ | [Cobra](https://github.com/spf13/cobra) | [Pterm](https://github.com/pterm/pterm)

## 🔑 Funzionalità Chiave

### Configurazione Multi-Profilo

- **Path**: `~/.config/projman/projman_config.json` (Unix) | `%APPDATA%/projman/projman_config.json` (Win)
- **Formato**:

```json
{
  "current_profile": "default",
  "profiles": {
    "default": {
      "root_of_projects": "/path",
      "selected_projects": ["proj-a", "proj-b"]
    }
  }
}
```

- **Comandi**:
  - `list`: elenca profili disponibili
  - `use <nome>`: imposta profilo corrente
  - `delete <nome>`: elimina profilo
  - `init <nome-profilo> <path>`: crea/aggiorna profilo

### Discovery Progetti

- Scansiona directory per trovare `pom.xml`
- Selezione interattiva con `pterm.DefaultInteractiveMultiselect`

### Git Update (`cmd/git/update.go`)

- **Stash automatico** → operazioni → ripristino stash
- **Branch intelligenti**:
  - `develop`: `git pull origin develop`
  - `deploy/*`: `git pull origin <branch>`
  - Altri: `git fetch + merge origin/develop`

### Maven Install (`cmd/mvn/install.go`)

- Analisi dipendenze tra progetti (parsing `pom.xml`)
- **Ordinamento topologico** con algoritmo di Kahn (`internal/graph`)
- Rilevamento cicli e sottomoduli
- **Executor Maven dedicato** (`internal/maven/executor/`):
  - Parsing preciso con regex dell'output Maven
  - Riconoscimento fasi (clean, compile, test-compile, test, package, install, deploy)
  - Spinner per ogni fase che viene completato quando la fase finisce
  - Tracking test in esecuzione con risultati (passed/failed/errors)
  - Supporto WAR, JAR, integration tests
- Flag `--tests/-t` (default: test disabilitati)

## 🛠️ Task Comuni

### Aggiungere Comando di Primo Livello

```go
// cmd/nuovo.go
package cmd

import "github.com/spf13/cobra"

var nuovoCmd = &cobra.Command{
    Use:   "nuovo",
    Short: "Descrizione",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementazione
    },
}

func init() {
    RootCmd.AddCommand(nuovoCmd)
}
```

### Aggiungere Sottocomando Git/Maven

```go
// cmd/git/nuovo.go o cmd/mvn/nuovo.go
package git  // o mvn

import "github.com/spf13/cobra"

var nuovoCmd = &cobra.Command{
    Use:   "nuovo",
    Short: "Descrizione",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementazione
    },
}

func init() {
    gitCmd.AddCommand(nuovoCmd)  // o mvnCmd
}
```

### Eseguire Comandi

**Comandi generici:**

```go
import "github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"

// Con spinner e timer
exec.RunWithSpinner("comando", []string{"arg1", "arg2"}, 0)

// Output normale in tempo reale
exec.Run("git", "status")

// Catturare output come stringa
output, err := exec.RunWithOutput("git", "rev-parse", "HEAD")
```

**Comandi Maven con parsing output:**

```go
import "github.com/SalvatoreSpagnuolo-BipRED/projman/internal/maven/executor"

// Esegue Maven con logging dettagliato delle fasi e test
mavenExec := executor.NewMavenExecutor("project-name", []string{"-f", "pom.xml", "install"})
err := mavenExec.Run()
```

## 📝 Convenzioni

### Codice

- **Formatting**: `gofmt`
- **Naming**: camelCase (var), PascalCase (const/export)
- **Commenti**: obbligatori per export, iniziano con nome elemento
- **Funzioni**: max 50 righe, singola responsabilità
- **Errori**: wrapping con `fmt.Errorf("contesto: %w", err)`

### Output Pterm

```go
pterm.Info.Println()      // Informazioni
pterm.Success.Println()   // Operazione riuscita
pterm.Warning.Println()   // Avvisi
pterm.Error.Println()     // Errori
pterm.DefaultHeader.Println()  // Intestazioni
```

## 🔄 Workflow

1. **Branch** feature da `develop`
2. **Modifica** codice
3. **Test** manuale: `make build && ./projman <comando>`
4. **Format**: `make lint` (esegue `go fmt` e `go vet`)
5. **Commit** con messaggio chiaro (usa Conventional Commits: feat:, fix:, docs:, etc.)
6. **Update docs** se necessario (README.md, AGENTS.md)
7. **Release**: `make release VERSION=v1.2.3` (crea tag e triggera GitHub Actions)

### Processo di Release

Projman usa [GoReleaser](https://goreleaser.com/) per gestire le release:

1. **Locale**: `make release VERSION=v1.2.3`

   - Crea e pusha tag Git
   - Triggera GitHub Actions

2. **GitHub Actions** (automatico):

   - Aggiorna versione in `cmd/root.go` e README badge
   - Compila con GoReleaser (Windows, macOS Intel/ARM, Linux)
   - Genera changelog automatico dai commit
   - Crea GitHub Release con archivi
   - Copia binari in `builds/` directory
   - Commit delle modifiche su `main`

3. **Test locale** (senza pubblicare):
   ```bash
   make snapshot  # Crea build in dist/ senza pubblicare
   ```

### Comandi Make Utili

```bash
make build      # Compila projman
make test       # Esegue test
make lint       # Formatta e analizza codice
make install    # Installa in GOPATH/bin
make clean      # Rimuove artifacts
make snapshot   # Test release locale (richiede goreleaser)
make ci         # Esegue tutti i check CI
```

## 🐛 Troubleshooting

| Problema               | Soluzione                                               |
| ---------------------- | ------------------------------------------------------- |
| Config non trovata     | Esegui `projman init <nome-profilo> <dir>`              |
| Nessun profilo attivo  | Esegui `projman use <nome-profilo>`                     |
| Git fallisce           | Verifica Git nel PATH e repo inizializzati              |
| Maven fallisce         | Verifica Maven nel PATH e connessione internet          |
| Comando non registrato | Importa subpackage in `cmd/root.go` e verifica `init()` |

## 📚 Reference

- [Cobra Docs](https://github.com/spf13/cobra/blob/main/user_guide.md)
- [Pterm Examples](https://github.com/pterm/pterm/tree/master/_examples)
- [GoReleaser Docs](https://goreleaser.com/intro/)
- [Effective Go](https://golang.org/doc/effective_go)

---

**v1.2.0** | Novembre 2025 | Salvatore Spagnuolo
