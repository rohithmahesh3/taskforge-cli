package issue

import (
	"fmt"
	"strconv"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var versionField string

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Manage issue version history",
		Long:  `View and inspect version history for issues.`,
	}

	versionListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List version history for an issue",
		Args:    cobra.ExactArgs(1),
		RunE:    runVersionList,
	}

	versionGetCmd := &cobra.Command{
		Use:   "get <issue-id> <version-num>",
		Short: "Get a specific version of an issue field",
		Args:  cobra.ExactArgs(2),
		RunE:  runVersionGet,
	}
	versionGetCmd.Flags().StringVarP(&versionField, "field", "f", "", "Field to retrieve version for (required, e.g. name, description)")
	versionGetCmd.MarkFlagRequired("field")

	versionDiffCmd := &cobra.Command{
		Use:   "diff <issue-id> <version-num>",
		Short: "View the diff for a specific version of an issue field",
		Args:  cobra.ExactArgs(2),
		RunE:  runVersionDiff,
	}
	versionDiffCmd.Flags().StringVarP(&versionField, "field", "f", "", "Field to view diff for (required, e.g. name, description)")
	versionDiffCmd.MarkFlagRequired("field")

	versionCmd.AddCommand(versionListCmd)
	versionCmd.AddCommand(versionGetCmd)
	versionCmd.AddCommand(versionDiffCmd)

	IssueCmd.AddCommand(versionCmd)
}

func runVersionList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	versions, err := client.ListIssueVersions(projectID, issueID)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		output.Info("No version history found for this issue")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type versionOutput struct {
		Version   int    `json:"version"`
		Field     string `json:"field"`
		CreatedAt string `json:"created_at"`
		ActorID   string `json:"actor_id"`
	}

	var outputs []versionOutput
	for _, v := range versions {
		outputs = append(outputs, versionOutput{
			Version:   v.VersionNum,
			Field:     v.Field,
			CreatedAt: v.CreatedAt,
			ActorID:   v.ActorID,
		})
	}

	return formatter.Print(outputs)
}

func runVersionGet(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]
	versionNum, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[1])
	}

	field := versionField
	if field == "" {
		return fmt.Errorf("--field is required (e.g. name, description)")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	version, err := client.GetIssueVersion(projectID, issueID, field, versionNum)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(version)
}

func runVersionDiff(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]
	versionNum, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[1])
	}

	field := versionField
	if field == "" {
		return fmt.Errorf("--field is required (e.g. name, description)")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	diff, err := client.GetIssueVersionDiff(projectID, issueID, field, versionNum)
	if err != nil {
		return err
	}

	fmt.Println(diff)
	return nil
}