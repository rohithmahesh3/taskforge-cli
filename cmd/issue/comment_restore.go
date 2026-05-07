package issue

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	commentRestoreVersionNum int
	commentRestoreYes        bool
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
	commentRestoreCmd.Flags().BoolVarP(&commentRestoreYes, "yes", "y", false, "Skip interactive confirmation")

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

	issueID := issueRef
	if !looksLikeUUID(issueRef) {
		var resolveErr error
		projectID, issueID, resolveErr = resolveIssueContext(client, projectID, issueRef)
		if resolveErr != nil {
			return resolveErr
		}
	}

	confirm := commentRestoreYes
	if !confirm {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm restoration")
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
