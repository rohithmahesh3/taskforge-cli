package label

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	labelName        string
	labelDescription string
	labelColor       string
	labelDeleteYes   bool
)

var LabelCmd = &cobra.Command{
	Use:     "label",
	Aliases: []string{"labels", "tag"},
	Short:   "Manage project labels",
	Long:    `List, create, edit, and manage labels for categorizing work items in your project.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List labels",
	Long:    `List all labels in the current project.`,
	RunE:    runList,
}

var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View label details",
	Long:  `Display detailed information about a specific label.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Long: `Create a new label in the current project.

Examples:
  taskforge label create --name "Bug" --color "#EF4444"
  taskforge label create -n "Feature" -c "#3B82F6" -d "New features"`,
	RunE: runCreate,
}

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a label",
	Long:  `Edit an existing label.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a label",
	Long:    `Delete a label from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

func init() {
	LabelCmd.AddCommand(listCmd)
	LabelCmd.AddCommand(viewCmd)
	LabelCmd.AddCommand(createCmd)
	LabelCmd.AddCommand(editCmd)
	LabelCmd.AddCommand(deleteCmd)

	// Create flags
	createCmd.Flags().StringVarP(&labelName, "name", "n", "", "Label name")
	createCmd.Flags().StringVarP(&labelDescription, "description", "d", "", "Label description")
	createCmd.Flags().StringVarP(&labelColor, "color", "c", "#EF4444", "Label color (hex code, e.g., #EF4444)")

	// Edit flags
	editCmd.Flags().StringVarP(&labelName, "name", "n", "", "New label name")
	editCmd.Flags().StringVarP(&labelDescription, "description", "d", "", "New label description")
	editCmd.Flags().StringVarP(&labelColor, "color", "c", "", "New label color")

	deleteCmd.Flags().BoolVarP(&labelDeleteYes, "yes", "y", false, "Skip confirmation")

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

	labels, err := client.ListLabels(projectID)
	if err != nil {
		return err
	}

	if len(labels) == 0 {
		output.Info("No labels found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type labelOutput struct {
		ID          string `table:"ID" json:"id"`
		Name        string `table:"NAME" json:"name"`
		Color       string `table:"COLOR" json:"color"`
		Description string `table:"DESCRIPTION" json:"description,omitempty"`
	}

	var outputs []labelOutput
	for _, l := range labels {
		outputs = append(outputs, labelOutput{
			ID:          l.ID,
			Name:        l.Name,
			Color:       l.Color,
			Description: l.Description,
		})
	}

	return formatter.Print(outputs)
}

func runView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	label, err := client.GetLabel(projectID, labelID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(label)
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

	req := taskforge.CreateLabelRequest{
		Name:        labelName,
		Description: labelDescription,
		Color:       labelColor,
	}

	label, err := client.CreateLabel(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created label '%s' (%s)", label.Name, label.ID))
	return nil
}

func runEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	req := taskforge.UpdateLabelRequest{}

	// Interactive mode if no flags provided
	if labelName == "" && labelDescription == "" && labelColor == "" {
		return fmt.Errorf("no edit flags provided. Available: --name, --description, --color")
	}

	// Use provided flags
	if labelName != "" {
		req.Name = labelName
	}
	if labelDescription != "" {
		req.Description = labelDescription
	}
	if labelColor != "" {
		req.Color = labelColor
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	updatedLabel, err := client.UpdateLabel(projectID, labelID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated label '%s'", updatedLabel.Name))
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	labelID := args[0]

	// Confirm deletion
	if !labelDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}
	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeleteLabel(projectID, labelID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted label %s", labelID))
	return nil
}
