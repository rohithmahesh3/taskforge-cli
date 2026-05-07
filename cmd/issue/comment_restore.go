package issue

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	commentRestoreVersionNum int
)

func newCommentRestoreCmd() *cobra.Command {
	commentRestoreCmd := &cobra.Command{
		Use:   "restore <issue-id> <comment-id>",
		Short: "Restore a deleted comment",
		Long: `Restore a previously deleted comment from its version history.

If --version-num is not specified, the latest version will be restored.

Examples:
  taskforge issue comment restore TF-1 comment-uuid
  taskforge issue comment restore TF-1 comment-uuid --version-num 2`,
		Args: cobra.ExactArgs(2),
		RunE: runCommentRestore,
	}

	commentRestoreCmd.Flags().IntVar(&commentRestoreVersionNum, "version-num", 0, "Version number to restore (default: latest)")

	return commentRestoreCmd
}

func runCommentRestore(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]
	commentID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	// Confirm restoration
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to restore comment %s?", commentID),
		Default: false,
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		output.Info("Restore cancelled")
		return nil
	}

	comment, err := client.RestoreComment(projectID, issueID, commentID, commentRestoreVersionNum)
	if err != nil {
		return fmt.Errorf("failed to restore comment: %w", err)
	}

	if commentRestoreVersionNum > 0 {
		output.Success(fmt.Sprintf("Restored comment %s to version %d", comment.ID, commentRestoreVersionNum))
	} else {
		output.Success(fmt.Sprintf("Restored comment %s", comment.ID))
	}

	return nil
}