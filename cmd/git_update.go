package cmd

import (
	"path/filepath"
	"strings"

	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/config"
	"github.com/SalvatoreSpagnuolo-BipRED/projman/internal/executil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Aggiorna tutti i progetti selezionati con git pull",
	Long: `Esegue git pull su tutti i progetti precedentemente selezionati.
Questo comando mantiene aggiornati tutti i progetti scaricando le ultime modifiche
dai rispettivi repository remoti.`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.DefaultHeader.Println("Git Update")

		// Carica la configurazione
		cfg, err := config.LoadSettings()
		if err != nil {
			pterm.Error.Println("Errore nel caricamento della configurazione:", err)
			return
		}

		// Chiedi all'utente se vuole passare al ramo develop (solo se necessario)

		// Elabora ogni progetto individualmente
		for _, p := range cfg.SelectedProjects {
			pterm.DefaultSection.Println("Elaborazione progetto: " + p)
			projectPath := filepath.Join(cfg.RootOfProjects, p)
			hasStash := false
			switchToDevelop := false

			// Controlla il ramo corrente
			isDevelopBranch, isDeployBranch, currBranch, err := branchInformation(projectPath)
			if err := errorHandling("Impossibilre recuperare info sul bransh", err); err != nil {
				continue
			}

			// Chiedi se si vuole cambiare ramo solo se non è già develop
			if !isDevelopBranch {
				pterm.Info.Println("Ramo corrente:", currBranch)
				switchToDevelop, _ = pterm.DefaultInteractiveConfirm.Show("Vuoi passare tutti i progetti al ramo 'develop' prima di eseguire il pull?")
			}

			// 1. Verifica e stasha i cambiamenti non committati
			status, err := executil.RunWithOutput("git", "-C", projectPath, "status", "--porcelain")
			if err := errorHandling("Impossibile versificare esistenza stash", err); err != nil {
				continue
			}

			if status != "" {
				pterm.Info.Println("Rilevati cambiamenti non committati, eseguo stash...")
				if err := errorHandling("Errore durante lo stash", executil.Run("git", "-C", projectPath, "stash", "push", "-m", "projman auto-stash")); err != nil {
					continue
				}

				hasStash = true
				pterm.Success.Println("Stash eseguito con successo")

			}

			// 2. Cambio ramo se richiesto
			if switchToDevelop {
				pterm.Info.Println("Cambio ramo a develop...")
				if err := errorHandling("Errore durante il cambio ramo", executil.Run("git", "-C", projectPath, "checkout", "develop")); err != nil {
					continue
				}

				pterm.Success.Println("Ramo cambiato a develop")
				isDevelopBranch, isDeployBranch, currBranch = true, false, "develop"
			}

			// 3. Esegui git pull
			pterm.Info.Println("Eseguo git pull...")
			if isDevelopBranch {
				// git pull per develop
				if err := errorHandling("Errore durante il git pull", executil.Run("git", "-C", projectPath, "pull", "origin/develop")); err != nil {
					continue
				}
				pterm.Success.Println("Git pull eseguito con successo sul ramo develop")
			} else if isDeployBranch {
				// git pull per deploy
				if err := errorHandling("Errore durante il git pull", executil.Run("git", "-C", projectPath, "pull", "origin/"+currBranch)); err != nil {
					continue
				}
				pterm.Success.Println("Git pull eseguito con successo sul ramo " + currBranch)
			} else {
				// git fetch + git merge (develop) per altri rami
				if err := errorHandling("Errore durante il git fetch", executil.Run("git", "-C", projectPath, "fetch", "origin/develop")); err != nil {
					continue
				}
				if err := errorHandling("Errore durante il git merge", executil.Run("git", "-C", projectPath, "merge", "origin/develop")); err != nil {
					continue
				}
				pterm.Success.Println("Fetch e merge di develop eseguito con successo sul ramo " + currBranch)
			}

			// 4. Ripristina lo stash se era stato fatto
			if hasStash {
				pterm.Info.Println("Ripristino modifiche locali...")
				if err := errorHandling("Errore durante il ripristino dello stash", executil.Run("git", "-C", projectPath, "stash", "pop")); err != nil {
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
