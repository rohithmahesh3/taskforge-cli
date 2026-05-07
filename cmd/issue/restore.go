package issue

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	restoreVersionNum int
	restoreYes        bool
)

func init() {
	restoreCmd := &cobra.Command{
		Use:   "restore <issue-id>",
		Short: "Restore a deleted issue",
		Long: `Restore a previously deleted issue from its version history.

If --version-num is not specified, the latest version will be restored.

Examples:
  taskforge issue restore TF-1
  taskforge issue restore TF-1 --version-num 3`,
		Args: cobra.ExactArgs(1),
		RunE: runRestore,
	}

	restoreCmd.Flags().IntVar(&restoreVersionNum, "version-num", 0, "Version number to restore (default: latest)")
	restoreCmd.Flags().BoolVarP(&restoreYes, "yes", "y", false, "Skip interactive confirmation")

	IssueCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	confirm := restoreYes
	if !confirm {
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return fmt.Errorf("interactive confirmation required in non-tty mode; rerun with --yes")
		}
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to restore issue %s?", issueRef),
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}
		if !confirm {
			output.Info("Restore cancelled")
			return nil
		}
	}

	issue, err := client.RestoreIssue(projectID, issueRef, restoreVersionNum)
	if err != nil {
		return fmt.Errorf("failed to restore issue: %w", err)
	}

	if restoreVersionNum > 0 {
		output.Success(fmt.Sprintf("Restored issue %d to version %d", issue.SequenceID, restoreVersionNum))
	} else {
		output.Success(fmt.Sprintf("Restored issue %d", issue.SequenceID))
	}

	return nil
}
