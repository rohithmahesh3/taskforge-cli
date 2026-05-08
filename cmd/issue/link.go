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
	linkTitle     string
	linkDeleteYes bool
)

func init() {
	// Add link subcommand
	linkCmd := &cobra.Command{
		Use:   "link",
		Short: "Manage issue links",
		Long:  `Add, list, and manage external links attached to issues.`,
	}

	linkListCmd := &cobra.Command{
		Use:     "list <issue-id>",
		Aliases: []string{"ls"},
		Short:   "List links for an issue",
		Args:    cobra.ExactArgs(1),
		RunE:    runLinkList,
	}

	linkAddCmd := &cobra.Command{
		Use:   "add <issue-id> <url>",
		Short: "Add a link to an issue",
		Args:  cobra.ExactArgs(2),
		RunE:  runLinkAdd,
	}

	linkDeleteCmd := &cobra.Command{
		Use:     "delete <issue-id> <link-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove a link from an issue",
		Args:    cobra.ExactArgs(2),
		RunE:    runLinkDelete,
	}

	linkAddCmd.Flags().StringVarP(&linkTitle, "title", "t", "", "Link title")
	linkDeleteCmd.Flags().BoolVarP(&linkDeleteYes, "yes", "y", false, "Skip confirmation")

	linkCmd.AddCommand(linkListCmd)
	linkCmd.AddCommand(linkAddCmd)
	linkCmd.AddCommand(linkDeleteCmd)

	IssueCmd.AddCommand(linkCmd)
}

func runLinkList(cmd *cobra.Command, args []string) error {
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

	links, err := client.ListLinks(projectID, issueID)
	if err != nil {
		return err
	}

	if len(links) == 0 {
		output.Info("No links found for this issue")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type linkOutput struct {
		ID    string `table:"ID" json:"id"`
		Title string `table:"TITLE" json:"title"`
		URL   string `table:"URL" json:"url"`
	}

	var outputs []linkOutput
	for _, link := range links {
		outputs = append(outputs, linkOutput{
			ID:    link.ID,
			Title: link.Title,
			URL:   link.URL,
		})
	}

	return formatter.Print(outputs)
}

func runLinkAdd(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]

	url := args[1]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	projectID, issueID, err = resolveIssueContext(client, projectID, issueID)
	if err != nil {
		return err
	}

	req := taskforge.CreateLinkRequest{
		Title: linkTitle,
		URL:   url,
	}

	link, err := client.CreateLink(projectID, issueID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Added link '%s' to issue", link.Title))
	return nil
}

func runLinkDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueID := args[0]
	linkID := args[1]

	// Confirm deletion
	if !linkDeleteYes {
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

	if err := client.DeleteLink(projectID, issueID, linkID); err != nil {
		return err
	}

	output.Success("Link deleted")
	return nil
}
