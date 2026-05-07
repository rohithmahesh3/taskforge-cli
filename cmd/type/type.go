package issuetype

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	typeName        string
	typeDescription string
	typeDeleteYes   bool
)

var TypeCmd = &cobra.Command{
	Use:     "type",
	Aliases: []string{"issue-type"},
	Short:   "Manage issue types",
	Long:    `List, create, and manage custom issue types for your project.`,
}

var typeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issue types",
	RunE:    runTypeList,
}

var typeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue type",
	RunE:  runTypeCreate,
}

var typeDeleteCmd = &cobra.Command{
	Use:     "delete <type-id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an issue type",
	Args:    cobra.ExactArgs(1),
	RunE:    runTypeDelete,
}

func init() {
	typeCreateCmd.Flags().StringVarP(&typeName, "name", "n", "", "Issue type name")
	typeCreateCmd.Flags().StringVarP(&typeDescription, "description", "d", "", "Issue type description")
	typeDeleteCmd.Flags().BoolVarP(&typeDeleteYes, "yes", "y", false, "Skip confirmation")

	_ = typeCreateCmd.MarkFlagRequired("name")

	TypeCmd.AddCommand(typeListCmd)
	TypeCmd.AddCommand(typeCreateCmd)
	TypeCmd.AddCommand(typeDeleteCmd)
}

func runTypeList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	types, err := client.ListIssueTypes(projectID)
	if err != nil {
		return err
	}

	if len(types) == 0 {
		output.Info("No issue types found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type typeOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Description string `table:"DESCRIPTION" json:"description"`
		IsDefault   string `table:"DEFAULT" json:"is_default"`
	}

	var outputs []typeOutput
	for _, t := range types {
		isDefault := ""
		if t.IsDefault {
			isDefault = "✓"
		}
		outputs = append(outputs, typeOutput{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			IsDefault:   isDefault,
		})
	}

	return formatter.Print(outputs)
}

func runTypeCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	req := taskforge.CreateIssueTypeRequest{
		Name:        typeName,
		Description: typeDescription,
		IsActive:    true,
	}

	issueType, err := client.CreateIssueType(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created issue type '%s' (%s)", issueType.Name, issueType.ID))
	return nil
}

func runTypeDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	typeID := args[0]

	// Confirm deletion
	// Confirm deletion
	if !typeDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsIssueTypeEnabled {
		return fmt.Errorf("issue types are disabled for project %s", projectID)
	}

	if err := client.DeleteIssueType(projectID, typeID); err != nil {
		return err
	}

	output.Success("Issue type deleted")
	return nil
}
