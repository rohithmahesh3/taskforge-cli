package issue

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

func ensureTimeTrackingEnabled(client *api.Client, projectID string) error {
	project, err := client.GetProject(projectID)
	if err != nil {
		return err
	}
	if !project.IsTimeTrackingEnabled {
		return fmt.Errorf("time tracking is disabled for project %q (%s). Enable it in the project settings to use time tracking commands", project.Name, projectID)
	}
	return nil
}

var (
	worklogDescription string
	worklogDuration    string
	worklogDeleteYes   bool
)

func init() {
	// Add time subcommand
	timeCmd := &cobra.Command{
		Use:     "time",
		Aliases: []string{"worklog", "log"},
		Short:   "Manage time tracking for issues",
		Long:    `Log, view, and manage time spent on work items.`,
	}

	timeListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List time logs for an issue",
		Args:    cobra.ExactArgs(1),
		RunE:    runTimeList,
	}

	timeLogCmd := &cobra.Command{
		Use:   "log <issue-id> <duration>",
		Short: "Log time for an issue",
		Long: `Log time spent on an issue.

Duration formats:
  - Minutes: 60, 120
  - Hours: 2h, 4.5h
  - Hours and minutes: 2h30m, 1h45m

Examples:
  taskforge issue time log ISS-123 2h30m
  taskforge issue time log ISS-123 90 -d "Fixed the bug"`,
		Args: cobra.RangeArgs(2, 2),
		RunE: runTimeLog,
	}

	timeTotalCmd := &cobra.Command{
		Use:   "total <issue-id>",
		Short: "Show total time logged for an issue",
		Args:  cobra.ExactArgs(1),
		RunE:  runTimeTotal,
	}

	timeEditCmd := &cobra.Command{
		Use:   "edit <issue-id> <worklog-id>",
		Short: "Edit a time log",
		Args:  cobra.ExactArgs(2),
		RunE:  runTimeEdit,
	}

	timeDeleteCmd := &cobra.Command{
		Use:     "delete <issue-id> <worklog-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a time log",
		Args:    cobra.ExactArgs(2),
		RunE:    runTimeDelete,
	}

	timeLogCmd.Flags().StringVarP(&worklogDescription, "description", "d", "", "Description of work done")
	timeEditCmd.Flags().StringVarP(&worklogDescription, "description", "d", "", "New description")
	timeEditCmd.Flags().StringVarP(&worklogDuration, "duration", "t", "", "New duration (e.g., 2h30m, 90)")
	timeDeleteCmd.Flags().BoolVarP(&worklogDeleteYes, "yes", "y", false, "Skip confirmation")

	timeCmd.AddCommand(timeListCmd)
	timeCmd.AddCommand(timeLogCmd)
	timeCmd.AddCommand(timeTotalCmd)
	timeCmd.AddCommand(timeEditCmd)
	timeCmd.AddCommand(timeDeleteCmd)

	IssueCmd.AddCommand(timeCmd)
}

func runTimeList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}
	if err := ensureTimeTrackingEnabled(client, projectID); err != nil {
		return err
	}

	worklogs, err := client.ListWorklogs(projectID, issueID)
	if err != nil {
		return err
	}

	if len(worklogs) == 0 {
		output.Info("No time logs found for this issue")
		return nil
	}

	// Get total time
	totalMinutes, err := client.GetTotalWorklogTime(projectID, issueID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type worklogOutput struct {
		ID          string `table:"ID" json:"id"`
		Duration    string `table:"DURATION" json:"duration_formatted"`
		Description string `table:"DESCRIPTION" json:"description"`
		LoggedAt    string `table:"LOGGED AT" json:"created_at_formatted"`
	}

	var outputs []worklogOutput
	for _, w := range worklogs {
		outputs = append(outputs, worklogOutput{
			ID:          w.ID,
			Duration:    formatDuration(w.Duration),
			Description: truncateString(w.Description, 40),
			LoggedAt:    w.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	if err := formatter.Print(outputs); err != nil {
		return err
	}

	output.Info(fmt.Sprintf("\nTotal time: %s", formatDuration(totalMinutes)))
	return nil
}

func runTimeLog(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	durationStr := args[1]

	// Parse duration
	duration, err := parseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	// Prompt for description if flag not provided
	if worklogDescription == "" {
		worklogDescription = fmt.Sprintf("Worked for %s", formatDuration(duration))
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}
	if err := ensureTimeTrackingEnabled(client, projectID); err != nil {
		return err
	}

	req := taskforge.CreateWorklogRequest{
		Description: worklogDescription,
		Duration:    duration,
	}

	worklog, err := client.CreateWorklog(projectID, issueID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Logged %s for issue %s", formatDuration(worklog.Duration), issueID))
	return nil
}

func runTimeTotal(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}
	if err := ensureTimeTrackingEnabled(client, projectID); err != nil {
		return err
	}

	totalMinutes, err := client.GetTotalWorklogTime(projectID, issueID)
	if err != nil {
		return err
	}

	fmt.Printf("Total time logged for issue %s: %s\n", issueID, formatDuration(totalMinutes))
	return nil
}

func runTimeEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	worklogID := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}
	if err := ensureTimeTrackingEnabled(client, projectID); err != nil {
		return err
	}

	req := taskforge.UpdateWorklogRequest{}

	// Use provided flags
	if worklogDescription == "" && worklogDuration == "" {
		return fmt.Errorf("no edit flags provided. Available: --description, --duration")
	}
	if worklogDescription != "" {
		req.Description = worklogDescription
	}
	if worklogDuration != "" {
		duration, err := parseDuration(worklogDuration)
		if err != nil {
			return fmt.Errorf("invalid duration format: %w", err)
		}
		req.Duration = duration
	}

	updatedWorklog, err := client.UpdateWorklog(projectID, issueID, worklogID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated time log to %s", formatDuration(updatedWorklog.Duration)))
	return nil
}

func runTimeDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	worklogID := args[1]

	// Confirm deletion
	if !worklogDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}
	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}
	if err := ensureTimeTrackingEnabled(client, projectID); err != nil {
		return err
	}

	if err := client.DeleteWorklog(projectID, issueID, worklogID); err != nil {
		return err
	}

	output.Success("Time log deleted")
	return nil
}

// parseDuration parses duration strings like "2h30m", "90", "1.5h"
func parseDuration(s string) (int, error) {
	// Try parsing as plain minutes
	if minutes, err := strconv.Atoi(s); err == nil {
		return minutes, nil
	}

	// Try parsing as duration string (e.g., "2h30m", "1h", "30m")
	duration, err := time.ParseDuration(s)
	if err == nil {
		return int(duration.Minutes()), nil
	}

	// Try parsing "XhYm" format manually
	re := regexp.MustCompile(`^(?:(\d+(?:\.\d+)?)h)?(?:(\d+)m?)?$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s (use formats like: 90, 2h, 1h30m, 2.5h)", s)
	}

	var totalMinutes int

	if matches[1] != "" {
		hours, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}
		totalMinutes += int(hours * 60)
	}

	if matches[2] != "" {
		minutes, err := strconv.Atoi(matches[2])
		if err != nil {
			return 0, err
		}
		totalMinutes += minutes
	}

	if totalMinutes == 0 {
		return 0, fmt.Errorf("duration cannot be zero")
	}

	return totalMinutes, nil
}

// formatDuration formats minutes into human-readable string
func formatDuration(minutes int) string {
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}

// truncateString truncates a string to max length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
