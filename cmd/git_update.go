package cmd

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/executil"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/tableutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ProjectInfo struct {
	Name            string
	Path            string
	CurrentBranch   string
	IsDevelop       bool
	IsDeploy        bool
	switchToDevelop bool
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Aggiorna tutti i progetti selezionati con git pull",
	Long: `Esegue git pull su tutti i progetti precedentemente selezionati.
Questo comando mantiene aggiornati tutti i progetti scaricando le ultime modifiche
dai rispettivi repository remoti.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Carica la configurazione
		cfg, err := config.LoadSettings()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento della configurazione:", err)
			return
		}

		// Chiedi all'utente se vuole passare al ramo develop i progetti non su develop
		projectInfos := make([]ProjectInfo, len(cfg.SelectedProjects))
		for _, p := range cfg.SelectedProjects {
			path := filepath.Join(cfg.RootOfProjects, p)
			isDevelop, isDeploy, currBranch, err := branchInformation(path)
			if err != nil {
				pterm.Error.Println("Errore nel recuperare le informazioni del ramo per il progetto " + p + ": " + err.Error())
				return
			}
			projectInfos = append(projectInfos, ProjectInfo{
				Name:            p,
				Path:            path,
				IsDevelop:       isDevelop,
				IsDeploy:        isDeploy,
				CurrentBranch:   currBranch,
				switchToDevelop: !isDevelop,
			})
		}

		// Mostra la conferma interattiva
		multiSelectTable := tableutil.MultiSelectTable{
			Headers: []string{"Progetto", "Ramo Attuale"},
			Rows: func() [][]string {
				rows := make([][]string, 0, len(projectInfos))
				for _, pInfo := range projectInfos {
					rows = append(rows, []string{pInfo.Name, pInfo.CurrentBranch})
				}
				return rows
			}(),
			Message: "Seleziona i progetti da passare a 'develop':",
			DefaultIndices: func() []int {
				indices := []int{}
				for i, pInfo := range projectInfos {
					if pInfo.switchToDevelop {
						indices = append(indices, i)
					}
				}
				return indices
			}(),
		}
		selectedIndices, err := multiSelectTable.Show()
		if err != nil {
			pterm.Error.Println("Errore nella selezione interattiva:", err)
			return
		}
		// Aggiorna le informazioni sui progetti in base alla selezione dell'utente
		for i := range projectInfos {
			if slices.Contains(selectedIndices, i) {
				projectInfos[i].switchToDevelop = true
			} else {
				projectInfos[i].switchToDevelop = false
			}
		}

		// Elabora ogni progetto individualmente
		for _, p := range projectInfos {
			pterm.DefaultSection.Println("Elaborazione progetto: " + p.Name)
			hasStash := false

			pterm.Info.Println("Ramo corrente:", p.CurrentBranch)

			// 1. Verifica e stasha i cambiamenti non committati
			status, err := executil.RunWithOutput("git", "-C", p.Path, "status", "--porcelain")
			if err := errorHandling("Impossibile versificare esistenza stash", err); err != nil {
				continue
			}

			if status != "" {
				pterm.Info.Println("Rilevati cambiamenti non committati, eseguo stash...")
				if err := errorHandling("Errore durante lo stash", executil.Run("git", "-C", p.Path, "stash", "push", "-m", "projman auto-stash")); err != nil {
					continue
				}

				hasStash = true
				pterm.Success.Println("Stash eseguito con successo")

			}

			// 2. Cambio ramo se richiesto
			if p.switchToDevelop {
				pterm.Info.Println("Cambio ramo a develop...")
				if err := errorHandling("Errore durante il cambio ramo", executil.Run("git", "-C", p.Path, "checkout", "develop")); err != nil {
					continue
				}

				pterm.Success.Println("Ramo cambiato a develop")
				p.IsDevelop, p.IsDeploy, p.CurrentBranch = true, false, "develop"
			}

			// 3. Esegui git pull
			pterm.Info.Println("Eseguo git pull...")
			if p.IsDevelop {
				// git pull per develop
				if err := errorHandling("Errore durante il git pull", executil.Run("git", "-C", p.Path, "pull", "origin", "develop")); err != nil {
					continue
				}
				pterm.Success.Println("Git pull eseguito con successo sul ramo develop")
			} else if p.IsDeploy {
				// git pull per deploy
				if err := errorHandling("Errore durante il git pull", executil.Run("git", "-C", p.Path, "pull", "origin", p.CurrentBranch)); err != nil {
					continue
				}
				pterm.Success.Println("Git pull eseguito con successo sul ramo " + p.CurrentBranch)
			} else {
				// git fetch + git merge (develop) per altri rami
				if err := errorHandling("Errore durante il git fetch", executil.Run("git", "-C", p.Path, "fetch", "origin", "develop")); err != nil {
					continue
				}
				if err := errorHandling("Errore durante il git merge", executil.Run("git", "-C", p.Path, "merge", "origin/develop")); err != nil {
					continue
				}
				pterm.Success.Println("Fetch e merge di develop eseguito con successo sul ramo " + p.CurrentBranch)
			}

			// 4. Ripristina lo stash se era stato fatto
			if hasStash {
				pterm.Info.Println("Ripristino modifiche locali...")
				if err := errorHandling("Errore durante il ripristino dello stash", executil.Run("git", "-C", p.Path, "stash", "pop")); err != nil {
					continue
				}
				pterm.Success.Println("Modifiche locali ripristinate con successo")
			}

			pterm.Println() // Riga vuota per separare i progetti
		}

		pterm.Success.Println("Operazioni completate per tutti i progetti")

	},
}

func init() {
	gitCmd.AddCommand(updateCmd)
}

func branchInformation(projectPath string) (bool, bool, string, error) {
	currentBranch, err := executil.RunWithOutput("git", "-C", projectPath, "branch", "--show-current")
	if err != nil {
		return false, false, "", err
	}

	isDevelopBranch := currentBranch == "develop"
	isDeployBranch := strings.HasPrefix(currentBranch, "deploy")

	return isDevelopBranch, isDeployBranch, currentBranch, nil

}

func errorHandling(msg string, err error) error {
	if err != nil {
		pterm.Warning.Printfln("Si è verificato un errore: %s, Risolvi manualmente il problema prima di continuare.", msg)
		pterm.Error.Println(err.Error())
		pterm.DefaultInteractiveConfirm.Show("Premi Invio per continuare con il prossimo progetto")
		return err
	}
	return nil
}
