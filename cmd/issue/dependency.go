package issue

import (
	"fmt"
	"strings"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:     "dependency",
	Aliases: []string{"deps", "dep"},
	Short:   "Manage issue dependencies",
	Long:    `List, add, or remove dependencies for an issue.`,
}

var isBlocks bool

var depsListCmd = &cobra.Command{
	Use:     "list <issue-id>",
	Aliases: []string{"ls"},
	Short:   "List dependencies for an issue",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		issueID := args[0]
		if strings.Contains(issueID, "-") {
			issue, err := client.GetIssueByIdentifier(issueID)
			if err != nil {
				return err
			}
			issueID = issue.ID
		}

		deps, err := client.ListIssueDependencies(projectID, issueID)
		if err != nil {
			return err
		}

		formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
		formatter.Print(deps)
		return nil
	},
}

var depsAddCmd = &cobra.Command{
	Use:   "add <issue-id> <target-identifier>",
	Short: "Add a dependency",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		issueID := args[0]
		if strings.Contains(issueID, "-") {
			issue, err := client.GetIssueByIdentifier(issueID)
			if err != nil {
				return err
			}
			issueID = issue.ID
		}

		targetID := args[1]
		if strings.Contains(targetID, "-") {
			target, err := client.GetIssueByIdentifier(targetID)
			if err != nil {
				return err
			}
			targetID = target.ID
		}

		if err := client.AddIssueDependency(projectID, issueID, targetID, isBlocks); err != nil {
			return err
		}

		output.Success("Successfully added dependency")
		return nil
	},
}

var depsRemoveCmd = &cobra.Command{
	Use:     "remove <issue-id> <target-identifier>",
	Aliases: []string{"rm"},
	Short:   "Remove a dependency",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := config.Cfg.DefaultProject
		if projectID == "" {
			return fmt.Errorf("no project specified")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		issueID := args[0]
		if strings.Contains(issueID, "-") {
			issue, err := client.GetIssueByIdentifier(issueID)
			if err != nil {
				return err
			}
			issueID = issue.ID
		}

		targetID := args[1]
		if strings.Contains(targetID, "-") {
			target, err := client.GetIssueByIdentifier(targetID)
			if err != nil {
				return err
			}
			targetID = target.ID
		}

		deps, err := client.ListIssueDependencies(projectID, issueID)
		if err != nil {
			return err
		}

		var depID string
		for _, dep := range deps.Blocks {
			if dep.FromIssueID == issueID && dep.ToIssueID == targetID {
				depID = dep.ID
				break
			}
			if dep.FromIssueID == targetID && dep.ToIssueID == issueID {
				depID = dep.ID
				break
			}
		}
		if depID == "" {
			for _, dep := range deps.DependsOn {
				if dep.FromIssueID == issueID && dep.ToIssueID == targetID {
					depID = dep.ID
					break
				}
				if dep.FromIssueID == targetID && dep.ToIssueID == issueID {
					depID = dep.ID
					break
				}
			}
		}

		if depID == "" {
			return fmt.Errorf("dependency link not found")
		}

		if err := client.RemoveIssueDependency(projectID, issueID, depID); err != nil {
			return err
		}

		output.Success("Successfully removed dependency")
		return nil
	},
}

func init() {
	dependencyCmd.AddCommand(depsListCmd)
	dependencyCmd.AddCommand(depsAddCmd)
	dependencyCmd.AddCommand(depsRemoveCmd)

	depsAddCmd.Flags().BoolVar(&isBlocks, "blocks", false, "Mark dependency as 'blocks' (default is 'depends_on')")
}
