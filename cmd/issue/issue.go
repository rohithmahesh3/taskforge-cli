package issue

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	stateFilter    string
	assigneeFilter string
	perPage        int

	issueTitle       string
	issueDescription string
	issuePriority    string
	issueState       string
	issueAssignees   []string
	issueLabels      []string
	deleteYes        bool
)

var IssueCmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i", "issues", "ticket"},
	Short:   "Manage issues (work items)",
	Long:    `List, create, edit, and manage TaskForge issues/work items.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List issues",
	Long: `List issues in the current project.

Examples:
  taskforge issue list
  taskforge issue list --state <state-id>
  taskforge issue list --assignee <assignee-id>`,
	RunE: runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View issue details",
	Long:  `Display detailed information about a specific issue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long: `Create a new issue in the current project.

Examples:
  taskforge issue create --title "Bug fix" --priority high
  taskforge issue create -t "Feature request" -d "Description here"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an issue",
	Long:  `Edit an existing issue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete an issue",
	Long:    `Delete an issue from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for issues",
	Long:  `Search for issues across the workspace.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	IssueCmd.AddCommand(listCmd)
	IssueCmd.AddCommand(viewCmd)
	IssueCmd.AddCommand(createCmd)
	IssueCmd.AddCommand(editCmd)
	IssueCmd.AddCommand(deleteCmd)
	IssueCmd.AddCommand(searchCmd)
	IssueCmd.AddCommand(dependencyCmd)

	// List flags
	listCmd.Flags().StringVarP(&stateFilter, "state", "s", "", "Filter by state ID")
	listCmd.Flags().StringVar(&assigneeFilter, "assignee", "", "Filter by assignee ID")
	listCmd.Flags().IntVarP(&perPage, "limit", "l", 20, "Number of issues to show per page")

	// Create flags
	createCmd.Flags().StringVarP(&issueTitle, "title", "t", "", "Issue title")
	createCmd.Flags().StringVarP(&issueDescription, "description", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&issuePriority, "priority", "p", "medium", "Issue priority (none, low, medium, high, urgent)")
	createCmd.Flags().StringSliceVarP(&issueAssignees, "assignee", "a", nil, "Assignee ID(s)")
	createCmd.Flags().StringSliceVar(&issueLabels, "label", nil, "Label ID(s)")

	// Edit flags
	editCmd.Flags().StringVarP(&issueTitle, "title", "t", "", "New title")
	editCmd.Flags().StringVarP(&issueDescription, "description", "d", "", "New description")
	editCmd.Flags().StringVar(&issuePriority, "priority", "", "New priority (none, low, medium, high, urgent)")
	editCmd.Flags().StringVar(&issueState, "state", "", "New state ID")
	editCmd.Flags().StringSliceVarP(&issueAssignees, "assignee", "a", nil, "New assignee ID(s)")
	editCmd.Flags().StringSliceVar(&issueLabels, "label", nil, "New label ID(s)")

	// Delete flags
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip interactive confirmation")
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

	opts := api.IssueListOptions{
		State:    stateFilter,
		Assignee: assigneeFilter,
		Limit:    perPage,
	}

	issues, _, err := client.ListIssues(projectID, opts)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type issueOutput struct {
		ID       string                `table:"ID" json:"id"`
		Sequence int                   `table:"#" json:"sequence_id"`
		Title    string                `table:"TITLE" json:"title"`
		State    taskforge.StateOutput `table:"STATE" json:"state"`
		Priority string                `table:"PRIORITY" json:"priority"`
		Assignee string                `table:"ASSIGNEE" json:"assignee"`
	}

	var outputs []issueOutput
	for _, issue := range issues {
		assignee := "-"
		if len(issue.Assignees) > 0 {
			u := issue.Assignees[0]
			if u.DisplayName != "" {
				assignee = "@" + u.DisplayName
			} else if u.Email != "" {
				assignee = u.Email
			}
		}

		outputs = append(outputs, issueOutput{
			ID:       issue.ID,
			Sequence: issue.SequenceID,
			Title:    issue.Name,
			State:    taskforge.StateOutputFromIssue(issue),
			Priority: issue.Priority,
			Assignee: assignee,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issue, err := resolveIssue(client, projectID, issueID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(issue)
}

func runCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	if issueTitle == "" {
		return fmt.Errorf("issue title is required (use --title / -t flag)")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Use default assignee from config if no assignees specified via flags
	assignees := issueAssignees
	if len(assignees) == 0 && config.Cfg.DefaultAssignee != "" {
		assignees = []string{config.Cfg.DefaultAssignee}
	}

	req := taskforge.CreateIssueRequest{
		Name:        issueTitle,
		Description: issueDescription,
		Priority:    issuePriority,
		Assignees:   assignees,
		Labels:      issueLabels,
	}

	issue, err := client.CreateIssue(projectID, req)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(issue)
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	// Get current issue
	issue, err := resolveIssue(client, projectID, issueRef)
	if err != nil {
		return err
	}
	if issue.ProjectID != "" {
		projectID = issue.ProjectID
	}

	req := taskforge.UpdateIssueRequest{}

	hasFlags := issueTitle != "" || issueDescription != "" || issuePriority != "" || issueState != "" || len(issueAssignees) > 0 || len(issueLabels) > 0

	if !hasFlags {
		return fmt.Errorf("no edit flags provided. Available: --title, --description, --priority, --state, --assignee, --label")
	}

	// Use provided flags
	if issueTitle != "" {
		req.Name = issueTitle
	}
	if issueDescription != "" {
		req.Description = issueDescription
	}
	if issuePriority != "" {
		req.Priority = issuePriority
	}
	if issueState != "" {
		req.State = issueState
	}
	if len(issueAssignees) > 0 {
		req.Assignees = issueAssignees
	}
	if len(issueLabels) > 0 {
		req.Labels = issueLabels
	}

	updatedIssue, err := client.UpdateIssue(projectID, issue.ID, req)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(updatedIssue)
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	projectID, issueID, err := resolveIssueContext(client, projectID, issueRef)
	if err != nil {
		return err
	}

	if !deleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}

	if err := client.DeleteIssue(projectID, issueID); err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(map[string]string{"status": "deleted", "issue_ref": issueRef, "issue_id": issueID, "project_id": projectID})
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issues, err := client.SearchIssues(query)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		output.Info("No issues found")
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
