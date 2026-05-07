package state

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	stateName        string
	stateDescription string
	stateColor       string
	stateGroup       string
	stateDeleteYes   bool
)

var StateCmd = &cobra.Command{
	Use:     "state",
	Aliases: []string{"states"},
	Short:   "Manage project states",
	Long:    `List, create, edit, and manage workflow states for your project.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List states",
	Long:    `List all states in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View state details",
	Long:  `Display detailed information about a specific state.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new state",
	Long: `Create a new workflow state in the current project.

Examples:
  taskforge state create --name "In Review" --color "#F59E0B" --group started`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a state",
	Long:  `Edit an existing workflow state.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a state",
	Long:    `Delete a state from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	StateCmd.AddCommand(listCmd)
	StateCmd.AddCommand(viewCmd)
	StateCmd.AddCommand(createCmd)
	StateCmd.AddCommand(editCmd)
	StateCmd.AddCommand(deleteCmd)

	// Create flags
	createCmd.Flags().StringVarP(&stateName, "name", "n", "", "State name")
	createCmd.Flags().StringVarP(&stateDescription, "description", "d", "", "State description")
	createCmd.Flags().StringVarP(&stateColor, "color", "c", "#F59E0B", "State color (hex code, e.g., #F59E0B)")
	createCmd.Flags().StringVarP(&stateGroup, "group", "g", "backlog", "State group (backlog, unstarted, started, completed, cancelled)")

	// Edit flags
	editCmd.Flags().StringVarP(&stateName, "name", "n", "", "New state name")
	editCmd.Flags().StringVarP(&stateDescription, "description", "d", "", "New state description")
	editCmd.Flags().StringVarP(&stateColor, "color", "c", "", "New state color")
	editCmd.Flags().StringVarP(&stateGroup, "group", "g", "", "New state group")

	deleteCmd.Flags().BoolVarP(&stateDeleteYes, "yes", "y", false, "Skip confirmation")

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

	states, err := client.ListStates(projectID)
	if err != nil {
		return err
	}

	if len(states) == 0 {
		output.Info("No states found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type stateOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Group       string `table:"GROUP" json:"group"`
		Color       string `table:"COLOR" json:"color"`
		IsDefault   string `table:"DEFAULT" json:"is_default"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []stateOutput
	for _, s := range states {
		isDefault := ""
		if s.IsDefault {
			isDefault = "✓"
		}
		outputs = append(outputs, stateOutput{
			ID:          s.ID,
			Name:        s.Name,
			Group:       s.Group,
			Color:       s.Color,
			IsDefault:   isDefault,
			Description: s.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	stateID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	state, err := client.GetState(projectID, stateID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(state)
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

	req := taskforge.CreateStateRequest{
		Name:        stateName,
		Description: stateDescription,
		Color:       stateColor,
		Group:       stateGroup,
	}

	state, err := client.CreateState(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created state '%s' (%s)", state.Name, state.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	stateID := args[0]

	req := taskforge.UpdateStateRequest{}

	// Interactive mode if no flags provided
	if stateName == "" && stateDescription == "" && stateColor == "" && stateGroup == "" {
		return fmt.Errorf("no edit flags provided. Available: --name, --description, --color, --group")
	}

	// Use provided flags
	if stateName != "" {
		req.Name = stateName
	}
	if stateDescription != "" {
		req.Description = stateDescription
	}
	if stateColor != "" {
		req.Color = stateColor
	}
	if stateGroup != "" {
		req.Group = stateGroup
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	updatedState, err := client.UpdateState(projectID, stateID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated state '%s'", updatedState.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	stateID := args[0]

	// Confirm deletion
	if !stateDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeleteState(projectID, stateID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted state %s", stateID))
	return nil
}
