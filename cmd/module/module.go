package module

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	moduleName        string
	moduleDescription string
	moduleStatus      string
	showArchived      bool
	moduleDeleteYes   bool
)

var ModuleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mod"},
	Short:   "Manage modules",
	Long:    `List, create, edit, and manage TaskForge modules.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List modules",
	Long:    `List all modules in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View module details",
	Long:  `Display detailed information about a specific module.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new module",
	Long: `Create a new module in the current project.

Examples:
  taskforge module create --name "Authentication" --description "User auth features"
  taskforge module create -n "API Integration" -s "in-progress"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a module",
	Long:  `Edit an existing module.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a module",
	Long:    `Delete a module from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive a module",
	Long:  `Archive a module to hide it from active modules.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runArchive,
}

var issuesCmd = &cobra.Command{
	Use:     "issues <id>",
	Aliases: []string{"work-items"},
	Short:   "List module issues",
	Long:    `List all work items in a specific module.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runIssues,
}

var addIssuesCmd = &cobra.Command{
	Use:   "add-issues <module-id> <issue-ids...\u003e",
	Short: "Add issues to module",
	Long:  `Add work items to a module.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAddIssues,
}

var removeIssueCmd = &cobra.Command{
	Use:   "remove-issue <module-id> <issue-id>",
	Short: "Remove issue from module",
	Long:  `Remove a work item from a module.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runRemoveIssue,
}

func init() {
	ModuleCmd.AddCommand(listCmd)
	ModuleCmd.AddCommand(viewCmd)
	ModuleCmd.AddCommand(createCmd)
	ModuleCmd.AddCommand(editCmd)
	ModuleCmd.AddCommand(deleteCmd)
	ModuleCmd.AddCommand(archiveCmd)
	ModuleCmd.AddCommand(issuesCmd)
	ModuleCmd.AddCommand(addIssuesCmd)
	ModuleCmd.AddCommand(removeIssueCmd)

	// List flags
	listCmd.Flags().BoolVarP(&showArchived, "archived", "a", false, "Show archived modules")

	// Create flags
	createCmd.Flags().StringVarP(&moduleName, "name", "n", "", "Module name")
	createCmd.Flags().StringVarP(&moduleDescription, "description", "d", "", "Module description")
	createCmd.Flags().StringVarP(&moduleStatus, "status", "s", "backlog", "Module status (backlog, planned, in-progress, paused, completed, cancelled)")

	// Edit flags
	editCmd.Flags().StringVarP(&moduleName, "name", "n", "", "New module name")
	editCmd.Flags().StringVarP(&moduleDescription, "description", "d", "", "New module description")
	editCmd.Flags().StringVarP(&moduleStatus, "status", "s", "", "New module status")

	deleteCmd.Flags().BoolVarP(&moduleDeleteYes, "yes", "y", false, "Skip confirmation")

	_ = createCmd.MarkFlagRequired("name")
}

func runList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	modules, err := client.ListModules(projectID, showArchived)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		output.Info("No modules found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type moduleOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Status      string `table:"STATUS" json:"status,omitempty"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []moduleOutput
	for _, m := range modules {
		outputs = append(outputs, moduleOutput{
			ID:          m.ID,
			Name:        m.Name,
			Status:      m.Status,
			Description: m.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	module, err := client.GetModule(projectID, moduleID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(module)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := taskforge.CreateModuleRequest{
		Name:        moduleName,
		Description: moduleDescription,
		Status:      moduleStatus,
	}

	module, err := client.CreateModule(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created module '%s' (%s)", module.Name, module.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	req := taskforge.UpdateModuleRequest{}

	// Interactive mode if no flags provided
	if moduleName == "" && moduleDescription == "" && moduleStatus == "" {
		return fmt.Errorf("no edit flags provided. Available: --name, --description, --status")
	}

	// Use provided flags
	if moduleName != "" {
		req.Name = moduleName
	}
	if moduleDescription != "" {
		req.Description = moduleDescription
	}
	if moduleStatus != "" {
		req.Status = moduleStatus
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	updatedModule, err := client.UpdateModule(projectID, moduleID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated module '%s'", updatedModule.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	// Confirm deletion
	if !moduleDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeleteModule(projectID, moduleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted module %s", moduleID))
	return nil
}

func runArchive(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.ArchiveModule(projectID, moduleID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Archived module %s", moduleID))
	return nil
}

func runIssues(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issues, err := client.ListModuleIssues(projectID, moduleID)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found in this module")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type issueOutput struct {
		ID       string                `table:"ID" json:"id"`
		Sequence int                   `table:"#" json:"sequence_id"`
		Title    string                `table:"TITLE" json:"title"`
		State    taskforge.StateOutput `table:"STATE" json:"state"`
		Priority string                `table:"PRIORITY" json:"priority"`
	}

	var outputs []issueOutput
	for _, issue := range issues {
		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State:    taskforge.StateOutputFromIssue(issue),
			Priority: issue.Priority,
		})
	}

	return formatter.Print(outputs)
}

func runAddIssues(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]
	issueIDs := args[1:]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.AddIssuesToModule(projectID, moduleID, issueIDs); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Added %d issue(s) to module %s", len(issueIDs), moduleID))
	return nil
}

func runRemoveIssue(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	moduleID := args[0]
	issueID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.RemoveIssueFromModule(projectID, moduleID, issueID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Removed issue %s from module %s", issueID, moduleID))
	return nil
}
