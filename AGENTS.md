# AGENTS.md - Guida per AI Agents

Questo documento fornisce informazioni per AI agents che lavorano sul progetto **Projman**.

## 📋 Panoramica del Progetto

**Projman** è un tool CLI scritto in Go per la gestione batch di progetti Maven/Git. Consente di eseguire operazioni su multipli repository contemporaneamente.

### Tecnologie Utilizzate
- **Linguaggio**: Go 1.25.4+
- **Framework CLI**: [Cobra](https://github.com/spf13/cobra)
- **UI Terminal**: [Pterm](https://github.com/pterm/pterm)

## 🏗️ Architettura del Progetto

```
projman/
├── main.go                    # Entry point dell'applicazione
│
├── cmd/                       # Comandi CLI (Cobra)
│   ├── root.go               # Comando root + esportazione RootCmd
│   ├── init.go               # Inizializzazione configurazione
│   ├── help.go               # Comando help personalizzato
│   │
│   ├── git/                  # Comandi Git (subpackage)
│   │   ├── git.go           # Comando parent per operazioni Git
│   │   └── update.go        # Aggiornamento Git con gestione branch
│   │
│   └── mvn/                  # Comandi Maven (subpackage)
│       ├── mvn.go           # Comando parent per operazioni Maven
│       └── install.go       # Maven install con ordinamento dipendenze
│
├── internal/                  # Pacchetti interni (non esportabili)
│   ├── config/               # Gestione configurazione
│   │   └── config.go         # Load/Save config JSON
│   │
│   ├── project/              # Gestione progetti
│   │   └── project.go        # Discovery e filtering progetti Maven
│   │
│   ├── maven/                # Gestione dipendenze Maven
│   │   └── dependency.go     # Parsing pom.xml e ordinamento topologico
│   │
│   ├── exec/                 # Utility per esecuzione comandi
│   │   ├── exec.go           # Wrapper per os/exec (Run, RunWithOutput, RunWithScrollableOutput)
│   │   └── scrollable/       # Componenti per output scrollabile
│   │       ├── area.go       # Area scrollabile con spinner e tempo
│   │       ├── buffer.go     # Buffer circolare thread-safe per righe
│   │       ├── executor.go   # Gestione esecuzione comandi
│   │       ├── streamer.go   # Streaming stdout/stderr in tempo reale
│   │       └── updater.go    # Aggiornamento periodico dell'area
│   │
│   └── ui/                   # Utility per UI interattiva
│       └── multiselect.go    # Selezione multipla con tabelle
│
├── go.mod                     # Dipendenze Go
├── go.sum                     # Checksum dipendenze
│
├── AGENTS.md                  # Questa guida
├── README.md                  # Documentazione utente
├── REFACTORING_SUMMARY.md     # Storia del refactoring
└── LICENSE                    # Licenza MIT
```

### 📂 Organizzazione dei Package

#### Package `cmd`
- **root.go**: esporta `RootCmd` per permettere ai subpackage di registrare comandi
- **init.go**, **help.go**: comandi diretti sotto root

#### Subpackage `cmd/git` e `cmd/mvn`
- Ogni subpackage ha il proprio `init()` che registra i comandi su `cmd.RootCmd`
- File nominati senza prefissi ridondanti (es. `update.go` invece di `git_update.go`)
- Organizzazione logica: un file per comando

#### Package `internal`
- Nomi concisi: `exec`, `ui`, `config`, `project`
- Non esportabili fuori da projman (best practice Go)
- Ogni package ha responsabilità singola

## 🔑 Concetti Chiave

### 1. Configurazione
- **Percorso**: `~/.config/projman/projman_config.json` (Linux/macOS) o `%APPDATA%/projman/projman_config.json` (Windows)
- **Struttura**:
  ```json
  {
    "root_of_projects": "/path/to/projects",
    "selected_projects": ["project-a", "project-b"]
  }
  ```

### 2. Discovery Progetti
- Scansiona la directory root
- Identifica progetti Maven tramite presenza di `pom.xml`
- Permette selezione interattiva con `pterm.DefaultInteractiveMultiselect`

### 3. Gestione Git (package `cmd/git`)
- **Stash automatico**: salva modifiche non committate prima delle operazioni
- **Branch intelligenti**: gestisce diversamente `develop`, `deploy/*` e feature branch
- **Merge da develop**: sui feature branch esegue fetch+merge invece di pull
- Modularizzato in ~18 funzioni per leggibilità

### 4. Gestione Maven (package `cmd/mvn`)
- Esegue `mvn install` su tutti i progetti selezionati
- **Ordinamento automatico**: analizza le dipendenze Maven e ordina i progetti topologicamente
- **Analisi dipendenze**: parse di `pom.xml` per identificare `groupId:artifactId` e dipendenze
- **Rilevamento sottomoduli**: considera anche le dipendenze nei sottomoduli Maven
- **Rilevamento cicli**: identifica dipendenze circolari e notifica l'errore
- **Output scrollabile**: mostra l'output Maven in un'area di 15 righe con altezza fissa
- **Indicatore progresso**: spinner animato e tempo di esecuzione incorporati nell'area (aggiornamento ogni 500ms)
- **Buffer circolare**: mantiene solo le ultime N righe visibili, scrollando automaticamente
- Flag `--tests/-t` per abilitare/disabilitare test
- Default: test disabilitati (`-DskipTests=true`)
- Report finale con statistiche successi/fallimenti

### 5. Analisi Dipendenze Maven (package `internal/maven`)
- **Parsing pom.xml**: estrae `groupId`, `artifactId` e dipendenze da file Maven
- **Grafo dipendenze**: costruisce un grafo delle dipendenze tra progetti selezionati
- **Ordinamento topologico**: usa l'algoritmo di Kahn per ordinare i progetti
- **Gestione moduli**: analizza anche i sottomoduli per dipendenze annidate
- Garantisce che le dipendenze siano installate prima dei progetti che le utilizzano

## 🛠️ Modifiche Comuni

### Aggiungere un nuovo comando Git

1. Crea file `cmd/git/<comando>.go` (senza prefisso `git_`)
2. Definisci un `cobra.Command` con package `git`
3. Implementa la logica nel campo `Run` o `RunE`
4. Aggiungi il comando come subcommand di `gitCmd`:
   ```go
   package git
   
   import (
       "github.com/spf13/cobra"
   )
   
   var nuovoCmd = &cobra.Command{
       Use:   "nuovo",
       Short: "Descrizione breve",
       Run: func(cmd *cobra.Command, args []string) {
           // Logica del comando
       },
   }
   
   func init() {
       gitCmd.AddCommand(nuovoCmd)
   }
   ```

### Aggiungere un nuovo comando Maven

Stessa procedura dei comandi Git, ma in `cmd/mvn/`:
```go
package mvn

import (
    "github.com/spf13/cobra"
)

var nuovoCmd = &cobra.Command{
    Use:   "nuovo",
    Short: "Descrizione breve",
    Run: func(cmd *cobra.Command, args []string) {
        // Logica del comando
    },
}

func init() {
    mvnCmd.AddCommand(nuovoCmd)
}
```

### Modificare il comando help

Il file `cmd/help.go` usa `pterm` per formattare l'output. Modifica:
- `pterm.DefaultHeader` per titoli
- `pterm.DefaultTable` per tabelle
- `pterm.DefaultBulletList` per liste
- `pterm.DefaultBox` per box informativi

### Aggiungere nuove funzionalità al discovery

Modifica `internal/project/project.go`:
- `Discover()`: cambia logica di rilevamento progetti
- `Filter()`: modifica criteri di filtraggio
- `Names()`: cambia estrazione nomi progetti
- `isMavenProject()`: modifica criteri di identificazione Maven

### Usare l'output scrollabile per comandi esterni

Per eseguire comandi esterni con output che scrolla in un'area fissa:

```go
import "github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"

// Output scrollabile in 15 righe (ideale per Maven, Git, ecc.)
if err := exec.RunWithScrollableOutput("mvn", []string{"clean", "install"}, 15); err != nil {
    pterm.Error.Println("Errore:", err)
}

// Output normale (stampa tutto in tempo reale)
if err := exec.Run("git", "status"); err != nil {
    pterm.Error.Println("Errore:", err)
}
```

L'area scrollabile mantiene lo spinner e altri elementi visibili, migliorando l'UX durante operazioni lunghe.

**Nota**: I componenti per l'output scrollabile sono organizzati nel subpackage `internal/exec/scrollable/` per una migliore modularità e separazione delle responsabilità.

## 📝 Convenzioni di Codice

### Organizzazione File
- **Un file per comando**: evita file monolitici
- **Nomi senza prefissi ridondanti**: `update.go` invece di `git_update.go`
- **Subpackage per raggruppamento logico**: `cmd/git/`, `cmd/mvn/`

### Stile
- Usa `gofmt` per formattare il codice
- Nomi variabili: camelCase (es. `rootPath`, `selectedProjects`)
- Nomi costanti: PascalCase (es. `ConfigDirName`)
- Commenti esportati: iniziano con il nome dell'elemento
- Funzioni piccole: max 50 righe (responsabilità singola)

### Messaggi Output
- **Info**: `pterm.Info.Println()` - informazioni generiche
- **Success**: `pterm.Success.Println()` - operazioni completate
- **Warning**: `pterm.Warning.Println()` - avvisi
- **Error**: `pterm.Error.Println()` - errori
- **Section**: `pterm.DefaultSection.Println()` - separatori sezioni
- **Header**: `pterm.DefaultHeader.Println()` - intestazioni principali

### Gestione Errori
```go
if err != nil {
    pterm.Error.Println("Descrizione errore:", err)
    return fmt.Errorf("contesto: %w", err)  // Wrapping con %w
}
```

Per errori Git durante operazioni batch, usa la funzione `handleGitError()` in `git/update.go` che permette all'utente di continuare con i progetti successivi.

## 🧪 Testing

### Comandi di Test Rapido

```bash
# Build
go build -o projman

# Formattazione
go fmt ./...

# Linting
go vet ./...

# Test comando init (usa directory temporanea)
./projman init /path/to/test/projects

# Test comando help
./projman help

# Test git update (richiede configurazione esistente)
./projman git update

# Test mvn install
./projman mvn install --tests
```

### Debug

Usa `pterm.Debug.Println()` per output di debug:
```go
pterm.Debug.Println("Valore variabile:", variabile)
```

## 🔄 Workflow di Sviluppo

1. **Branch**: crea un feature branch da `develop`
2. **Modifica**: implementa la funzionalità
3. **Test**: verifica manualmente i comandi modificati
4. **Format**: esegui `go fmt ./...`
5. **Lint**: esegui `go vet ./...`
6. **Build**: verifica che compili con `go build`
7. **Commit**: commit con messaggio descrittivo
8. **PR**: crea Pull Request verso `develop`

## 🐛 Debug Comuni

### Configurazione non trovata
- Verifica che `projman init` sia stato eseguito
- Controlla il percorso: `~/.config/projman/projman_config.json`

### Git operations falliscono
- Verifica che Git sia nel PATH
- Controlla che i progetti abbiano repository Git inizializzati
- Verifica permessi su directory progetti

### Maven install fallisce
- Verifica che Maven sia nel PATH
- Controlla che i progetti abbiano file `pom.xml` validi
- Verifica connessione internet (per download dipendenze)

### Comando non registrato
- Verifica che il subpackage sia importato in `cmd/root.go`
- Controlla che `init()` chiami correttamente `cmd.RootCmd.AddCommand()`
- Assicurati che il package name sia corretto

## 📚 Risorse Utili

- [Cobra User Guide](https://github.com/spf13/cobra/blob/main/user_guide.md)
- [Pterm Examples](https://github.com/pterm/pterm/tree/master/_examples)
- [Go Package os/exec](https://pkg.go.dev/os/exec)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## 🎯 Obiettivi Futuri

Possibili estensioni del progetto:
- [ ] Supporto per altri build system (Gradle, npm, etc.)
- [ ] Comandi Git aggiuntivi (commit, push, status)
- [ ] Logging strutturato delle operazioni
- [ ] Configurazione globale e per-progetto
- [ ] Supporto per gruppi di progetti
- [ ] Esecuzione parallela delle operazioni
- [ ] Report dettagliati post-esecuzione
- [ ] Unit tests completi

## 💡 Suggerimenti per AI Agents

Quando lavori su questo progetto:

1. **Mantieni la struttura**: rispetta l'organizzazione cmd/git, cmd/mvn
2. **Nomi file semplici**: evita prefissi ridondanti (il package dice già di cosa si tratta)
3. **Funzioni piccole**: max 50 righe, responsabilità singola
4. **Commenti esportati**: tutti i tipi e funzioni esportati devono avere commenti
5. **Gestisci errori gracefully**: permetti all'utente di continuare dopo errori
6. **Output chiaro**: usa colori e formattazione `pterm` per migliorare UX
7. **Interattività**: usa prompt interattivi quando appropriato
8. **Configurazione**: rispetta il formato JSON della configurazione
9. **Cross-platform**: testa su Windows, Linux e macOS quando possibile
10. **Documentazione**: aggiorna README.md e AGENTS.md per nuove funzionalità

## 🔍 Best Practice Applicate

### Organizzazione
- ✅ Subpackage per raggruppamento logico
- ✅ Package `internal` per codice non esportabile
- ✅ Nomi concisi e significativi

### Codice
- ✅ Funzioni piccole e modulari
- ✅ Commenti esportati completi
- ✅ Gestione errori con wrapping (`%w`)
- ✅ Costanti per magic numbers

### Documentazione
- ✅ AGENTS.md per AI agents
- ✅ README.md per utenti finali
- ✅ REFACTORING_SUMMARY.md per storia

## 📞 Contatti

Per domande o suggerimenti sul progetto, apri una issue su GitHub.

---

**Ultimo aggiornamento**: Novembre 2025
**Versione Progetto**: 1.1.0
**Maintainer**: Salvatore Spagnuolo

