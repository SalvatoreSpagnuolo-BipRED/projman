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
│   ├── git/                            # Sottocomandi Git
│   │   ├── git.go                      # Parent command
│   │   └── update.go                   # Git pull/merge intelligente
│   └── mvn/                            # Sottocomandi Maven
│       ├── mvn.go                      # Parent command
│       └── install.go                  # Maven install con dep sorting
└── internal/                           # Package interni
    ├── config/config.go                # Gestione configurazione JSON
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

### Configurazione
- **Path**: `~/.config/projman/projman_config.json` (Unix) | `%APPDATA%/projman/projman_config.json` (Win)
- **Formato**: `{"root_of_projects": "/path", "selected_projects": ["proj-a", "proj-b"]}`

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
- Output dettagliato con parsing fasi Maven:
  - Riconoscimento fasi (clean, compile, test, package, install)
  - Logging test con risultati per classe/metodo
  - Riepilogo finale delle fasi eseguite
- Flag `--tests/-t` (default: test disabilitati)

## 🛠️ Task Comuni

### Aggiungere Comando Git/Maven
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
3. **Test** manuale: `go build && ./projman <comando>`
4. **Format**: `go fmt ./...`
5. **Lint**: `go vet ./...`
6. **Commit** con messaggio chiaro
7. **Update docs** se necessario (README.md, AGENTS.md)

## 🐛 Troubleshooting

| Problema | Soluzione |
|----------|-----------|
| Config non trovata | Esegui `projman init <dir>` |
| Git fallisce | Verifica Git nel PATH e repo inizializzati |
| Maven fallisce | Verifica Maven nel PATH e connessione internet |
| Comando non registrato | Importa subpackage in `cmd/root.go` e verifica `init()` |

## 📚 Reference

- [Cobra Docs](https://github.com/spf13/cobra/blob/main/user_guide.md)
- [Pterm Examples](https://github.com/pterm/pterm/tree/master/_examples)
- [Effective Go](https://golang.org/doc/effective_go)

---

**v1.1.0** | Novembre 2025 | Salvatore Spagnuolo

