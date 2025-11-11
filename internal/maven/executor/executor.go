package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// MavenExecutor gestisce l'esecuzione di comandi Maven con logging semplice
type MavenExecutor struct {
	projectName string
	args        []string
}

// NewMavenExecutor crea un nuovo executor Maven
func NewMavenExecutor(projectName string, args []string) *MavenExecutor {
	return &MavenExecutor{
		projectName: projectName,
		args:        args,
	}
}

// Run esegue il comando Maven con uno spinner e logging delle fasi
func (me *MavenExecutor) Run() error {
	pterm.Info.Printf("$ mvn %s\n", strings.Join(me.args, " "))

	// Avvia lo spinner
	spinner, _ := pterm.DefaultSpinner.Start("Esecuzione Maven in corso...")

	// Prepara il comando
	cmd := exec.Command("mvn", me.args...)
	cmd.Env = os.Environ()

	// Setup pipe per catturare output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		spinner.Fail("Errore setup comando")
		return fmt.Errorf("errore stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		spinner.Fail("Errore setup comando")
		return fmt.Errorf("errore stderr pipe: %w", err)
	}

	// Avvia il comando
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		spinner.Fail("Errore avvio Maven")
		return fmt.Errorf("errore avvio comando: %w", err)
	}

	// Canale per comunicare tra goroutine
	done := make(chan bool, 2)

	// Leggi stdout in background
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			me.processLine(scanner.Text(), spinner)
		}
		done <- true
	}()

	// Leggi stderr in background
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			me.processLine(scanner.Text(), spinner)
		}
		done <- true
	}()

	// Attendi fine lettura
	<-done
	<-done

	// Attendi fine comando
	cmdErr := cmd.Wait()

	// Calcola tempo trascorso
	elapsed := time.Since(startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	// Mostra risultato
	if cmdErr != nil {
		spinner.Fail(fmt.Sprintf("Build fallita dopo %dm %ds", minutes, seconds))
		return fmt.Errorf("comando maven fallito: %w", cmdErr)
	}

	spinner.Success(fmt.Sprintf("Build completata in %dm %ds", minutes, seconds))
	return nil
}

// processLine processa una riga di output Maven e logga le fasi importanti
func (me *MavenExecutor) processLine(line string, spinner *pterm.SpinnerPrinter) {
	line = strings.TrimSpace(line)

	if line == "" {
		return
	}

	// Rileva le fasi Maven (plugin execution)
	if strings.Contains(line, "---") && strings.Contains(line, "maven-") {
		phase := me.detectPhase(line)
		if phase != "" {
			// Ferma temporaneamente lo spinner per stampare la fase
			_ = spinner.Stop()
			pterm.Info.Printfln("▶ Fase: %s", phase)
			// Riavvia lo spinner
			newSpinner, _ := pterm.DefaultSpinner.Start("Esecuzione Maven in corso...")
			*spinner = *newSpinner
		}
		return
	}

	// Rileva classe di test in esecuzione
	if strings.HasPrefix(line, "[INFO] Running ") {
		testClass := strings.TrimPrefix(line, "[INFO] Running ")
		_ = spinner.Stop()
		pterm.Info.Printfln("  🧪 Test: %s", testClass)
		newSpinner, _ := pterm.DefaultSpinner.Start("Esecuzione Maven in corso...")
		*spinner = *newSpinner
		return
	}

	// Rileva risultati test
	if strings.HasPrefix(line, "[INFO] Tests run:") {
		_ = spinner.Stop()
		me.logTestResults(line)
		newSpinner, _ := pterm.DefaultSpinner.Start("Esecuzione Maven in corso...")
		*spinner = *newSpinner
		return
	}

	// Rileva build success/failure
	if strings.Contains(line, "BUILD SUCCESS") || strings.Contains(line, "BUILD FAILURE") {
		return // Verrà gestito dal cmdErr
	}

	// Logga errori Maven
	if strings.HasPrefix(line, "[ERROR]") && !strings.Contains(line, "ERROR") {
		_ = spinner.Stop()
		pterm.Error.Println(strings.TrimPrefix(line, "[ERROR] "))
		newSpinner, _ := pterm.DefaultSpinner.Start("Esecuzione Maven in corso...")
		*spinner = *newSpinner
	}
}

// detectPhase rileva la fase Maven dalla riga di plugin execution
func (me *MavenExecutor) detectPhase(line string) string {
	if strings.Contains(line, "maven-clean-plugin") {
		return "CLEAN"
	}
	if strings.Contains(line, "maven-compiler-plugin") {
		if strings.Contains(line, "compile") {
			return "COMPILE"
		}
		if strings.Contains(line, "testCompile") {
			return "TEST-COMPILE"
		}
	}
	if strings.Contains(line, "maven-surefire-plugin") || strings.Contains(line, "maven-failsafe-plugin") {
		return "TEST"
	}
	if strings.Contains(line, "maven-jar-plugin") || strings.Contains(line, "maven-war-plugin") {
		return "PACKAGE"
	}
	if strings.Contains(line, "maven-install-plugin") {
		return "INSTALL"
	}
	return ""
}

// logTestResults logga i risultati dei test in modo semplice
func (me *MavenExecutor) logTestResults(line string) {
	// Esempio: [INFO] Tests run: 5, Failures: 0, Errors: 0, Skipped: 0
	summary := strings.TrimPrefix(line, "[INFO] ")

	var run, failures, errors, skipped int
	_, _ = fmt.Sscanf(summary, "Tests run: %d, Failures: %d, Errors: %d, Skipped: %d",
		&run, &failures, &errors, &skipped)

	passed := run - failures - errors - skipped

	if failures > 0 || errors > 0 {
		pterm.Error.Printf("    ✗ %d passed, %d failed, %d errors\n", passed, failures, errors)
	} else if skipped > 0 {
		pterm.Success.Printf("    ✓ %d passed, %d skipped\n", passed, skipped)
	} else if run > 0 {
		pterm.Success.Printf("    ✓ %d test(s) passed\n", passed)
	}
}
