# AGENTS.md - Guida per AI Agents

Questo documento fornisce informazioni per AI agents che lavorano sul progetto **Projman**.

## 📋 Panoramica del Progetto

**Projman** è un tool CLI scritto in Go per la gestione batch di progetti Maven/Git. Consente di eseguire operazioni su multipli repository contemporaneamente.

### Tecnologie Utilizzate
- **Linguaggio**: Go 1.25.4
- **Framework CLI**: [Cobra](https://github.com/spf13/cobra)
- **UI Terminal**: [Pterm](https://github.com/pterm/pterm)

## 🏗️ Architettura del Progetto

```
projman/
├── main.go                    # Entry point
├── cmd/                       # Comandi CLI (Cobra)
│   ├── root.go               # Comando root
│   ├── help.go               # Comando help custom
│   ├── init.go               # Inizializzazione configurazione
│   ├── git.go                # Comando parent per operazioni Git
│   ├── git_update.go         # Aggiornamento Git con gestione branch
│   ├── mvn.go                # Comando parent per operazioni Maven
│   └── mvn_install.go        # Maven install con controllo test
├── internal/                  # Pacchetti interni
│   ├── config/               # Gestione configurazione
│   │   └── config.go         # Load/Save config JSON
│   ├── executil/             # Utility per esecuzione comandi
│   │   └── exec.go           # Wrapper per exec.Command
│   ├── project/              # Gestione progetti
│   │   └── project.go        # Discovery e filtering progetti
│   └── tableutil/            # Utility per tabelle interattive
│       └── multiselect.go    # Selezione multipla con tabelle
└── go.mod                     # Dipendenze Go

```

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

### 3. Gestione Git
- **Stash automatico**: salva modifiche non committate prima delle operazioni
- **Branch intelligenti**: gestisce diversamente `develop`, `deploy/*` e feature branch
- **Merge da develop**: sui feature branch esegue fetch+merge invece di pull

### 4. Gestione Maven
- Esegue `mvn install` su tutti i progetti selezionati
- Flag `--tests/-t` per abilitare/disabilitare test
- Default: test disabilitati (`-DskipTests=true`)

## 🛠️ Modifiche Comuni

### Aggiungere un nuovo comando Git

1. Crea file `cmd/git_<comando>.go`
2. Definisci un `cobra.Command`
3. Implementa la logica nel campo `Run` o `RunE`
4. Aggiungi il comando come subcommand di `gitCmd`:
   ```go
   func init() {
       gitCmd.AddCommand(nuovoCmd)
   }
   ```

### Aggiungere un nuovo comando Maven

Stessa procedura dei comandi Git, ma usa `mvnCmd` come parent:
```go
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

## 📝 Convenzioni di Codice

### Stile
- Usa `gofmt` per formattare il codice
- Nomi variabili: camelCase (es. `rootPath`, `selectedProjects`)
- Nomi costanti: PascalCase (es. `ConfigDirName`)
- Errori: gestiti con `pterm.Error.Println()` per output user-friendly

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
    return err
}
```

Per errori Git durante operazioni batch, usa la funzione `errorHandling()` in `git_update.go` che permette all'utente di continuare con i progetti successivi.

## 🧪 Testing

### Comandi di Test Rapido

```bash
# Build
go build -o projman

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
5. **Build**: verifica che compili con `go build`
6. **Commit**: commit con messaggio descrittivo
7. **PR**: crea Pull Request verso `develop`

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

## 📚 Risorse Utili

- [Cobra User Guide](https://github.com/spf13/cobra/blob/main/user_guide.md)
- [Pterm Examples](https://github.com/pterm/pterm/tree/master/_examples)
- [Go Package os/exec](https://pkg.go.dev/os/exec)

## 🎯 Obiettivi Futuri

Possibili estensioni del progetto:
- [ ] Supporto per altri build system (Gradle, npm, etc.)
- [ ] Comandi Git aggiuntivi (commit, push, status)
- [ ] Logging delle operazioni
- [ ] Configurazione globale e per-progetto
- [ ] Supporto per gruppi di progetti
- [ ] Esecuzione parallela delle operazioni
- [ ] Report dettagliati post-esecuzione

## 💡 Suggerimenti per AI Agents

Quando lavori su questo progetto:

1. **Mantieni lo stile esistente**: usa `pterm` per output, `cobra` per CLI
2. **Gestisci errori gracefully**: permetti all'utente di continuare dopo errori
3. **Output chiaro**: usa colori e formattazione per migliorare UX
4. **Interattività**: usa prompt interattivi quando appropriato
5. **Configurazione**: rispetta il formato JSON della configurazione
6. **Cross-platform**: testa su Windows, Linux e macOS quando possibile
7. **Documentazione**: aggiorna README.md per nuove funzionalità

## 📞 Contatti

Per domande o suggerimenti sul progetto, apri una issue su GitHub.

---

**Ultimo aggiornamento**: Novembre 2025
**Versione Progetto**: 1.0.0
**Maintainer**: Salvatore Spagnuolo

