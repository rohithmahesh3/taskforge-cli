package issue

import (
	"fmt"
	"strconv"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCommentVersionCmd() *cobra.Command {
	commentVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Manage comment version history",
		Long:  `View and inspect version history for issue comments.`,
	}

	commentVersionListCmd := &cobra.Command{
		Use:     "list <issue-id> <comment-id>",
		Aliases: []string{"ls"},
		Short:   "List version history for a comment",
		Args:    cobra.ExactArgs(2),
		RunE:    runCommentVersionList,
	}

	commentVersionGetCmd := &cobra.Command{
		Use:   "get <issue-id> <comment-id> <version-num>",
		Short: "Get a specific version of a comment",
		Args:  cobra.ExactArgs(3),
		RunE:  runCommentVersionGet,
	}

	commentVersionDiffCmd := &cobra.Command{
		Use:   "diff <issue-id> <comment-id> <version-num>",
		Short: "View the diff for a specific version of a comment",
		Args:  cobra.ExactArgs(3),
		RunE:  runCommentVersionDiff,
	}

	commentVersionCmd.AddCommand(commentVersionListCmd)
	commentVersionCmd.AddCommand(commentVersionGetCmd)
	commentVersionCmd.AddCommand(commentVersionDiffCmd)

	return commentVersionCmd
}

func runCommentVersionList(cmd *cobra.Command, args []string) error {
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

	versions, err := client.ListCommentVersions(projectID, issueID, commentID)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		output.Info("No version history found for this comment")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type versionOutput struct {
		Version   int    `json:"version"`
		Field     string `json:"field"`
		Created   string `json:"created_at"`
		CreatedBy string `json:"created_by"`
	}

	var outputs []versionOutput
	for _, v := range versions {
		createdBy := v.ActorID

		outputs = append(outputs, versionOutput{
			Version:   v.VersionNum,
			Field:     v.Field,
			Created:   v.CreatedAt,
			CreatedBy: createdBy,
		})
	}

	return formatter.Print(outputs)
}

func runCommentVersionGet(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]
	commentID := args[1]
	versionNum, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[2])
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	version, err := client.GetCommentVersion(projectID, issueID, commentID, versionNum)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(version)
}

func runCommentVersionDiff(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]
	commentID := args[1]
	versionNum, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[2])
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	diff, err := client.GetCommentVersionDiff(projectID, issueID, commentID, versionNum)
	if err != nil {
		return err
	}

	fmt.Println(diff)
	return nil
}