package git

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/exec"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/ui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// ProjectInfo contiene le informazioni di stato di un progetto Git
type ProjectInfo struct {
	Name            string // Nome del progetto
	Path            string // Percorso assoluto del progetto
	CurrentBranch   string // Nome del branch corrente
	IsDevelop       bool   // true se il branch corrente è 'develop'
	IsDeploy        bool   // true se il branch corrente inizia con 'deploy/'
	switchToDevelop bool   // true se l'utente vuole passare a 'develop'
}

// updateCmd rappresenta il comando per aggiornare i progetti con git pull/merge
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Aggiorna tutti i progetti selezionati con git pull",
	Long: `Esegue git pull su tutti i progetti precedentemente selezionati.
Questo comando mantiene aggiornati tutti i progetti scaricando le ultime modifiche
dai rispettivi repository remoti.

Il comportamento varia in base al branch corrente:
  - Branch 'develop': esegue git pull da origin/develop
  - Branch 'deploy/*': esegue git pull dal branch corrente
  - Altri branch: esegue git fetch + git merge di origin/develop

Prima delle operazioni, eventuali modifiche non committate vengono salvate in stash
e automaticamente ripristinate al termine.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Carica e valida la configurazione
		cfg, err := config.LoadAndValidateConfig()
		if err != nil {
			return
		}

		// Permette all'utente di selezionare i progetti da aggiornare
		selectedProjects, err := config.SelectProjectsToUpdate(cfg)
		if err != nil {
			return
		}

		// Salva la selezione aggiornata
		cfg.SelectedProjects = selectedProjects
		if err := config.SaveSettings(*cfg); err != nil {
			pterm.Error.Println("Errore nel salvataggio della configurazione:", err)
			return
		}

		// Raccoglie informazioni sui branch di tutti i progetti
		projectInfos, err := gatherProjectsInfo(cfg)
		if err != nil {
			return
		}

		// Gestisce la selezione interattiva dei progetti da passare a develop
		if err := handleBranchSwitching(projectInfos); err != nil {
			return
		}

		// Processa ogni progetto individualmente
		processProjects(projectInfos)

		pterm.Success.Println("Operazioni completate per tutti i progetti")
	},
}

// gatherProjectsInfo raccoglie informazioni sui branch di tutti i progetti selezionati
func gatherProjectsInfo(cfg *config.Config) ([]ProjectInfo, error) {
	projectInfos := make([]ProjectInfo, 0, len(cfg.SelectedProjects))

	for _, projectName := range cfg.SelectedProjects {
		path := filepath.Join(cfg.RootOfProjects, projectName)
		isDevelop, isDeploy, currBranch, err := branchInformation(path)
		if err != nil {
			pterm.Error.Printf("Errore nel recuperare le informazioni del branch per '%s': %v\n", projectName, err)
			return nil, err
		}

		projectInfos = append(projectInfos, ProjectInfo{
			Name:            projectName,
			Path:            path,
			IsDevelop:       isDevelop,
			IsDeploy:        isDeploy,
			CurrentBranch:   currBranch,
			switchToDevelop: !isDevelop, // Default: suggerisci di passare a develop se non ci sei già
		})
	}

	return projectInfos, nil
}

// handleBranchSwitching gestisce la selezione interattiva dei progetti da passare a develop
func handleBranchSwitching(projectInfos []ProjectInfo) error {
	// Filtra i progetti che non sono su develop
	nonDevelopProjects := filterNonDevelopProjects(projectInfos)

	if len(nonDevelopProjects) == 0 {
		pterm.Info.Println("Tutti i progetti sono già sul branch 'develop'")
		return nil
	}

	// Mostra i progetti non su develop e permette la selezione
	pterm.Info.Println("Alcuni progetti non sono attualmente sul branch 'develop'")

	selectedIndices, err := showBranchSwitchingTable(nonDevelopProjects)
	if err != nil {
		pterm.Error.Println("Errore nella selezione interattiva:", err)
		return err
	}

	// Aggiorna le informazioni in base alla selezione dell'utente
	updateBranchSwitchingChoices(projectInfos, nonDevelopProjects, selectedIndices)

	return nil
}

// filterNonDevelopProjects filtra i progetti che non sono sul branch develop
func filterNonDevelopProjects(projectInfos []ProjectInfo) []ProjectInfo {
	nonDevelopProjects := make([]ProjectInfo, 0)
	for _, pInfo := range projectInfos {
		if !pInfo.IsDevelop {
			nonDevelopProjects = append(nonDevelopProjects, pInfo)
		}
	}
	return nonDevelopProjects
}

// showBranchSwitchingTable mostra una tabella interattiva per selezionare i progetti da passare a develop
func showBranchSwitchingTable(nonDevelopProjects []ProjectInfo) ([]int, error) {
	// Prepara i dati per la tabella
	rows := make([][]string, 0, len(nonDevelopProjects))
	defaultIndices := make([]int, 0)

	for i, pInfo := range nonDevelopProjects {
		rows = append(rows, []string{pInfo.Name, pInfo.CurrentBranch})
		if pInfo.switchToDevelop {
			defaultIndices = append(defaultIndices, i)
		}
	}

	// Mostra la tabella multiselect
	multiSelectTable := ui.MultiSelectTable{
		Headers:        []string{"Progetto", "Branch Attuale"},
		Rows:           rows,
		Message:        "Seleziona i progetti da passare a 'develop':",
		DefaultIndices: defaultIndices,
	}

	return multiSelectTable.Show()
}

// updateBranchSwitchingChoices aggiorna le scelte di cambio branch in base alla selezione utente
func updateBranchSwitchingChoices(projectInfos []ProjectInfo, nonDevelopProjects []ProjectInfo, selectedIndices []int) {
	// Crea una mappa per lookup rapido
	projectIndexMap := make(map[string]int)
	for i := range projectInfos {
		if !projectInfos[i].IsDevelop {
			for idx, p := range nonDevelopProjects {
				if p.Name == projectInfos[i].Name {
					projectIndexMap[p.Name] = idx
					break
				}
			}
		}
	}

	// Aggiorna i flag switchToDevelop
	for i := range projectInfos {
		if !projectInfos[i].IsDevelop {
			idx, exists := projectIndexMap[projectInfos[i].Name]
			if exists {
				projectInfos[i].switchToDevelop = slices.Contains(selectedIndices, idx)
			}
		}
	}
}

// processProjects elabora tutti i progetti eseguendo le operazioni git
func processProjects(projectInfos []ProjectInfo) {
	for _, pInfo := range projectInfos {
		pterm.DefaultSection.Printf("Elaborazione progetto: %s", pInfo.Name)

		if err := processProject(&pInfo); err != nil {
			// L'errore è già stato gestito, continua con il prossimo progetto
			pterm.Println() // Riga vuota per separare i progetti
			continue
		}

		pterm.Println() // Riga vuota per separare i progetti
	}
}

// processProject elabora un singolo progetto eseguendo tutte le operazioni git necessarie
func processProject(pInfo *ProjectInfo) error {
	pterm.Info.Println("Branch corrente:", pInfo.CurrentBranch)

	// 1. Gestisci lo stash delle modifiche non committate
	hasStash, err := stashUncommittedChanges(pInfo.Path)
	if err != nil {
		return err
	}

	// 2. Cambia branch se richiesto
	if pInfo.switchToDevelop {
		if err := switchToDevelopBranch(pInfo); err != nil {
			return err
		}
	}

	// 3. Esegui git pull/merge in base al tipo di branch
	if err := updateRepository(pInfo); err != nil {
		return err
	}

	// 4. Ripristina lo stash se era stato fatto
	if hasStash {
		if err := popStash(pInfo.Path); err != nil {
			return err
		}
	}

	return nil
}

// stashUncommittedChanges salva in stash eventuali modifiche non committate
func stashUncommittedChanges(projectPath string) (bool, error) {
	status, err := exec.RunWithOutput("git", "-C", projectPath, "status", "--porcelain")
	if err := handleGitError("Impossibile verificare lo status del repository", err); err != nil {
		return false, err
	}

	if status == "" {
		return false, nil // Nessuna modifica da salvare
	}

	pterm.Info.Println("Rilevati cambiamenti non committati, eseguo stash...")
	if err := handleGitError("Errore durante lo stash",
		exec.Run("git", "-C", projectPath, "stash", "push", "-m", "projman auto-stash")); err != nil {
		return false, err
	}

	pterm.Success.Println("Stash eseguito con successo")
	return true, nil
}

// switchToDevelopBranch cambia il branch corrente a develop
func switchToDevelopBranch(pInfo *ProjectInfo) error {
	pterm.Info.Println("Cambio branch a 'develop'...")

	if err := handleGitError("Errore durante il cambio branch",
		exec.Run("git", "-C", pInfo.Path, "checkout", "develop")); err != nil {
		return err
	}

	pterm.Success.Println("Branch cambiato a 'develop'")

	// Aggiorna le informazioni del progetto
	pInfo.IsDevelop = true
	pInfo.IsDeploy = false
	pInfo.CurrentBranch = "develop"

	return nil
}

// updateRepository esegue git pull o git fetch+merge in base al tipo di branch
func updateRepository(pInfo *ProjectInfo) error {
	pterm.Info.Println("Aggiornamento del repository...")

	switch {
	case pInfo.IsDevelop:
		return updateDevelopBranch(pInfo)
	case pInfo.IsDeploy:
		return updateDeployBranch(pInfo)
	default:
		return updateFeatureBranch(pInfo)
	}
}

// updateDevelopBranch esegue git pull per il branch develop
func updateDevelopBranch(pInfo *ProjectInfo) error {
	if err := handleGitError("Errore durante il git pull",
		exec.Run("git", "-C", pInfo.Path, "pull", "origin", "develop")); err != nil {
		return err
	}
	pterm.Success.Println("Git pull eseguito con successo sul branch 'develop'")
	return nil
}

// updateDeployBranch esegue git pull per un branch deploy/*
func updateDeployBranch(pInfo *ProjectInfo) error {
	if err := handleGitError("Errore durante il git pull",
		exec.Run("git", "-C", pInfo.Path, "pull", "origin", pInfo.CurrentBranch)); err != nil {
		return err
	}
	pterm.Success.Printf("Git pull eseguito con successo sul branch '%s'\n", pInfo.CurrentBranch)
	return nil
}

// updateFeatureBranch esegue git fetch + git merge di develop per feature branch
func updateFeatureBranch(pInfo *ProjectInfo) error {
	// Fetch delle modifiche da develop
	if err := handleGitError("Errore durante il git fetch",
		exec.Run("git", "-C", pInfo.Path, "fetch", "origin", "develop")); err != nil {
		return err
	}

	// Merge di develop nel branch corrente
	if err := handleGitError("Errore durante il git merge",
		exec.Run("git", "-C", pInfo.Path, "merge", "origin/develop")); err != nil {
		return err
	}

	pterm.Success.Printf("Fetch e merge di 'develop' eseguiti con successo sul branch '%s'\n", pInfo.CurrentBranch)
	return nil
}

// popStash ripristina le modifiche salvate in stash
func popStash(projectPath string) error {
	pterm.Info.Println("Ripristino modifiche locali...")

	if err := handleGitError("Errore durante il ripristino dello stash",
		exec.Run("git", "-C", projectPath, "stash", "pop")); err != nil {
		return err
	}

	pterm.Success.Println("Modifiche locali ripristinate con successo")
	return nil
}

// branchInformation recupera le informazioni sul branch corrente di un progetto
func branchInformation(projectPath string) (isDevelop, isDeploy bool, currentBranch string, err error) {
	currentBranch, err = exec.RunWithOutput("git", "-C", projectPath, "branch", "--show-current")
	if err != nil {
		return false, false, "", fmt.Errorf("impossibile recuperare il branch corrente: %w", err)
	}

	isDevelop = currentBranch == "develop"
	isDeploy = strings.HasPrefix(currentBranch, "deploy/")

	return isDevelop, isDeploy, currentBranch, nil
}

// handleGitError gestisce gli errori Git permettendo all'utente di continuare
func handleGitError(msg string, err error) error {
	if err != nil {
		pterm.Warning.Printfln("Si è verificato un errore: %s", msg)
		pterm.Error.Println("  ", err.Error())
		pterm.Info.Println("Risolvi manualmente il problema prima di continuare")
		_, _ = pterm.DefaultInteractiveConfirm.Show("Premi Invio per continuare con il prossimo progetto")
		return err
	}
	return nil
}

func init() {
	gitCmd.AddCommand(updateCmd)
}
