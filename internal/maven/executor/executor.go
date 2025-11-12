package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
)

// Pattern per il matching dell'output Maven
var (
	// Pattern per rilevare l'esecuzione di un plugin Maven
	// Esempio: [INFO] --- compiler:3.14.1:compile (default-compile) @ demo-1 ---
	// Esempio: [INFO] --- spring-boot:3.5.7:repackage (repackage) @ demo-1 ---
	// Cattura: plugin, goal, module (ignoriamo la versione)
	pluginPattern = regexp.MustCompile(`\[INFO] --- ([a-z-]+):([^:]+):([a-zA-Z]+).*?@\s+(\S+)\s+---`)

	// Pattern per i test in esecuzione
	// Esempio: [INFO] Running com.example.MyTest
	testRunningPattern = regexp.MustCompile(`\[INFO] Running (.+)`)

	// Pattern per i risultati dei test
	// Esempio: [INFO] Tests run: 5, Failures: 0, Errors: 0, Skipped: 0
	testResultsPattern = regexp.MustCompile(`\[INFO] Tests run: (\d+), Failures: (\d+), Errors: (\d+), Skipped: (\d+)`)

	// Pattern per BUILD SUCCESS/FAILURE
	buildResultPattern = regexp.MustCompile(`\[INFO] BUILD (SUCCESS|FAILURE)`)
)

// MavenPhase rappresenta una fase Maven riconosciuta
type MavenPhase struct {
	Name        string // Nome leggibile della fase (es: "COMPILE")
	Plugin      string // Nome del plugin Maven
	Goal        string // Goal del plugin
	Module      string // Nome del modulo Maven
	Description string // Descrizione breve per lo spinner
}

// MavenExecutor gestisce l'esecuzione di comandi Maven
type MavenExecutor struct {
	projectName    string
	args           []string
	currentSpinner *pterm.SpinnerPrinter
	currentPhase   *MavenPhase
}

// NewMavenExecutor crea un nuovo executor Maven
func NewMavenExecutor(projectName string, args []string) *MavenExecutor {
	return &MavenExecutor{
		projectName: projectName,
		args:        args,
	}
}

// Run esegue il comando Maven mostrando le fasi con spinner
func (mavenExec *MavenExecutor) Run() error {
	// Mostra comando con Info (senza spinner che si chiude subito)
	pterm.Info.Printf("$ mvn %s\n", strings.Join(mavenExec.args, " "))

	// Prepara ed esegui il comando
	cmd := exec.Command("mvn", mavenExec.args...)
	// Disabilita buffering Maven per output in tempo reale
	env := os.Environ()
	env = append(env, "MAVEN_OPTS=-Djansi.force=true")
	cmd.Env = env

	// Unisci stdout e stderr per catturare tutto l'output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("errore stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("errore stderr pipe: %w", err)
	}

	// Avvia il comando Maven
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("errore avvio comando: %w", err)
	}

	// Info spinner iniziale
	mavenExec.currentSpinner, _ = pterm.DefaultSpinner.Start("Starting Maven build...")

	// Leggi output in goroutine
	done := make(chan bool, 2)
	go mavenExec.readOutput(stdout, done)
	go mavenExec.readOutput(stderr, done)

	// Attendi completamento lettura
	<-done
	<-done

	// Attendi fine comando
	return cmd.Wait()

}

// readOutput legge l'output da uno stream
func (mavenExec *MavenExecutor) readOutput(stream io.Reader, done chan bool) {
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		mavenExec.processOutputLine(scanner.Text())
	}
	done <- true
}

// processOutputLine processa una riga di output Maven
func (mavenExec *MavenExecutor) processOutputLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// 1. Rileva inizio nuova fase Maven (plugin execution)
	if matches := pluginPattern.FindStringSubmatch(line); matches != nil {
		// matches[1] = plugin, matches[2] = version (ignorata), matches[3] = goal, matches[4] = module
		mavenExec.handlePhaseStart(matches[1], matches[2], matches[3], matches[4])
		return
	}

	// 2. Rileva esecuzione test
	if matches := testRunningPattern.FindStringSubmatch(line); matches != nil {
		mavenExec.handleTestStart(matches[1])
		return
	}

	// 3. Rileva risultati test
	if matches := testResultsPattern.FindStringSubmatch(line); matches != nil {
		mavenExec.handleTestResults(matches[1], matches[2], matches[3], matches[4])
		return
	}

	// 4. Rileva fine build (implicitamente completa ultima fase)
	if buildResultPattern.MatchString(line) {
		if mavenExec.currentSpinner != nil {
			mavenExec.currentSpinner.Success("Build completed")
			mavenExec.currentSpinner = nil
		}
		return
	}
}

// handlePhaseStart gestisce l'inizio di una nuova fase Maven
func (mavenExec *MavenExecutor) handlePhaseStart(plugin, version, goal, module string) {
	// Mappa plugin+goal a fase leggibile
	phase := mavenExec.identifyPhase(plugin, version, goal, module)
	if phase == nil {
		return // Fase non riconosciuta, ignora
	}

	// Se stessa fase, non fare nulla
	if mavenExec.currentPhase != nil && mavenExec.currentPhase.Name == phase.Name {
		return
	}

	// Completa fase precedente con Success
	if mavenExec.currentSpinner != nil {
		mavenExec.currentSpinner.Success()
	}

	// Avvia nuovo spinner per questa fase
	mavenExec.currentPhase = phase
	mavenExec.currentSpinner, _ = pterm.DefaultSpinner.Start(phase.Description)
}

// handleTestStart aggiorna lo spinner con la classe di test in esecuzione
func (mavenExec *MavenExecutor) handleTestStart(testClass string) {
	if mavenExec.currentSpinner == nil {
		return
	}

	// Estrai nome breve della classe
	shortName := testClass
	if idx := strings.LastIndex(testClass, "."); idx != -1 {
		shortName = testClass[idx+1:]
	}

	mavenExec.currentSpinner.UpdateText(fmt.Sprintf("%s - %s", mavenExec.currentPhase.Description, shortName))
}

// handleTestResults aggiorna lo spinner con i risultati dei test
func (mavenExec *MavenExecutor) handleTestResults(runStr, failStr, errStr, skipStr string) {
	if mavenExec.currentSpinner == nil {
		return
	}

	var run, failures, errors, skipped int
	_, _ = fmt.Sscanf(runStr, "%d", &run)
	_, _ = fmt.Sscanf(failStr, "%d", &failures)
	_, _ = fmt.Sscanf(errStr, "%d", &errors)
	_, _ = fmt.Sscanf(skipStr, "%d", &skipped)

	passed := run - failures - errors - skipped

	var message string
	if failures > 0 || errors > 0 {
		message = fmt.Sprintf("%s - %d passed, %d failed", mavenExec.currentPhase.Description, passed, failures+errors)
	} else {
		message = fmt.Sprintf("%s - %d passed", mavenExec.currentPhase.Description, passed)
	}

	mavenExec.currentSpinner.UpdateText(message)
}

// identifyPhase identifica la fase Maven dal plugin e goal
func (mavenExec *MavenExecutor) identifyPhase(plugin, _, goal, module string) *MavenPhase {
	// Mappa (plugin, goal) -> fase
	// Supporta sia nomi brevi (compiler) che completi (maven-compiler-plugin)

	switch {
	// CLEAN
	case (plugin == "maven-clean-plugin" || plugin == "clean") && goal == "clean":
		return &MavenPhase{
			Name:        "CLEAN",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("CLEAN @ %s", module),
		}

	// RESOURCES (copia risorse)
	case (plugin == "maven-resources-plugin" || plugin == "resources") && goal == "resources":
		return &MavenPhase{
			Name:        "RESOURCES",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("RESOURCES @ %s", module),
		}

	// COMPILE (codice main)
	case (plugin == "maven-compiler-plugin" || plugin == "compiler") && goal == "compile":
		return &MavenPhase{
			Name:        "COMPILE",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("COMPILE @ %s", module),
		}

	// TEST-RESOURCES
	case (plugin == "maven-resources-plugin" || plugin == "resources") && goal == "testResources":
		return &MavenPhase{
			Name:        "TEST-RESOURCES",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("TEST-RESOURCES @ %s", module),
		}

	// TEST-COMPILE (codice test)
	case (plugin == "maven-compiler-plugin" || plugin == "compiler") && goal == "testCompile":
		return &MavenPhase{
			Name:        "TEST-COMPILE",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("TEST-COMPILE @ %s", module),
		}

	// TEST (esecuzione test)
	case (plugin == "maven-surefire-plugin" || plugin == "surefire") && goal == "test":
		return &MavenPhase{
			Name:        "TEST",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("TEST @ %s", module),
		}

	case (plugin == "maven-failsafe-plugin" || plugin == "failsafe") && goal == "integration-test":
		return &MavenPhase{
			Name:        "INTEGRATION-TEST",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("INTEGRATION-TEST @ %s", module),
		}

	// PACKAGE
	case (plugin == "maven-jar-plugin" || plugin == "jar") && goal == "jar":
		return &MavenPhase{
			Name:        "PACKAGE",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("PACKAGE @ %s", module),
		}

	case (plugin == "maven-war-plugin" || plugin == "war") && goal == "war":
		return &MavenPhase{
			Name:        "PACKAGE",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("PACKAGE @ %s", module),
		}

	// SPRING BOOT REPACKAGE
	case (plugin == "spring-boot-maven-plugin" || plugin == "spring-boot") && goal == "repackage":
		return &MavenPhase{
			Name:        "REPACKAGE",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("REPACKAGE @ %s", module),
		}

	// INSTALL
	case (plugin == "maven-install-plugin" || plugin == "install") && goal == "install":
		return &MavenPhase{
			Name:        "INSTALL",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("INSTALL @ %s", module),
		}

	// DEPLOY
	case (plugin == "maven-deploy-plugin" || plugin == "deploy") && goal == "deploy":
		return &MavenPhase{
			Name:        "DEPLOY",
			Plugin:      plugin,
			Goal:        goal,
			Module:      module,
			Description: fmt.Sprintf("DEPLOY @ %s", module),
		}

	default:
		// Fase non riconosciuta
		return nil
	}
}
