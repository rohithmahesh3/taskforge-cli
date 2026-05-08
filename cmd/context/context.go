package context

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	includeAll       bool
	includeCycle     bool
	includeWorkspace bool
	includeIntake    bool
)

var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Generate CLI command reference for AI agents",
	Long: `Output a concise CLI command reference in markdown format.
Use flags to include additional modules beyond the default set.

Default modules: issue, project, module, state, label, type
Optional modules: --cycle, --workspace, --intake, --all`,
	RunE: runContext,
}

func init() {
	ContextCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all modules")
	ContextCmd.Flags().BoolVar(&includeCycle, "cycle", false, "Include cycle commands")
	ContextCmd.Flags().BoolVar(&includeWorkspace, "workspace", false, "Include workspace commands")
	ContextCmd.Flags().BoolVar(&includeIntake, "intake", false, "Include intake commands")
}

func runContext(cmd *cobra.Command, args []string) error {
	output := "# TaskForge CLI Commands\n\n"

	output += getGlobalFlags()
	output += getIssueCommands()
	output += getProjectCommands()
	output += getModuleCommands()
	output += getStateCommands()
	output += getLabelCommands()
	output += getTypeCommands()

	if includeAll || includeCycle {
		output += getCycleCommands()
	}
	if includeAll || includeWorkspace {
		output += getWorkspaceCommands()
	}
	if includeAll || includeIntake {
		output += getIntakeCommands()
	}

	fmt.Print(output)
	return nil
}

func getGlobalFlags() string {
	return `## Global Flags
` + "```" + `
--workspace <slug:text>     Workspace slug (overrides config)
--project <id:text>         Project ID (overrides config)
--no-color                  Disable colored output
--config <path:text>        Config file path
` + "```" + `

`
}

func getIssueCommands() string {
	return `## Issue (aliases: i, issues, ticket)
` + "```" + `
` + getIssueQuickStartCommands() + `
taskforge issue delete <id:seq_id|uuid>
taskforge issue search <query:text>
taskforge issue patch <issue-id:seq_id|uuid> --op <replace|diff|insert|delete> --field <name|description> [--actual <text>] [--value <text>] [--after <anchor>] [--diff <unified-diff>] [--file <patches.json>]

# Issue Versions & Restore
taskforge issue version list <issue-id:seq_id|uuid>
taskforge issue version get <issue-id:seq_id|uuid> <version-num:int> --field <name|description>
taskforge issue version diff <issue-id:seq_id|uuid> <version-num:int> --field <name|description>
taskforge issue restore <issue-id:seq_id|uuid> --version-num <int>

# Issue Comments
taskforge issue comment list <issue-id:seq_id|uuid>
taskforge issue comment add <issue-id:seq_id|uuid> [--text <markdown:text>]
                        [--access <enum:INTERNAL|EXTERNAL>]
taskforge issue comment delete <issue-id:seq_id|uuid> <comment-id:uuid>
taskforge issue comment version list <issue-id:seq_id|uuid> <comment-id:uuid>
taskforge issue comment version get <issue-id:seq_id|uuid> <comment-id:uuid> <version-num:int>
taskforge issue comment version diff <issue-id:seq_id|uuid> <comment-id:uuid> <version-num:int>
taskforge issue comment restore <issue-id:seq_id|uuid> <comment-id:uuid> --version-num <int>

# Issue Links
taskforge issue link list <issue-id:seq_id|uuid>
taskforge issue link add <issue-id:seq_id|uuid> <url:text> [--title <text>]
taskforge issue link delete <issue-id:seq_id|uuid> <link-id:uuid>

# Issue Time Tracking
taskforge issue time list <issue-id:seq_id|uuid>
taskforge issue time log <issue-id:seq_id|uuid> <duration:minutes|1h30m>
                     [--description <markdown:text>]
taskforge issue time total <issue-id:seq_id|uuid>
taskforge issue time edit <issue-id:seq_id|uuid> <worklog-id:uuid>
                      [--description <markdown:text>] [--duration <minutes|1h30m>]
taskforge issue time delete <issue-id:seq_id|uuid> <worklog-id:uuid>

# Issue Attachments
taskforge issue attachment list <issue-id:seq_id|uuid>
taskforge issue attachment upload <issue-id:seq_id|uuid> <file-path:text>
taskforge issue attachment edit <issue-id:seq_id|uuid> <attachment-id:uuid>
                           [--name <text>] [--archive | --unarchive]
taskforge issue attachment delete <issue-id:seq_id|uuid> <attachment-id:uuid>

# Issue Activity
taskforge issue activity list <issue-id:seq_id|uuid>
taskforge issue activity view <issue-id:seq_id|uuid> <activity-id:uuid>
` + "```" + `

`
}

func getIssueQuickStartCommands() string {
	return `taskforge issue list [--state <id:uuid>] [--assignee <id:uuid>]
                 [--limit <count:int>]

taskforge issue view <id:seq_id|uuid>

taskforge issue create [--title <text>] [--description <markdown:text>]
                   [--priority <enum:none|low|medium|high|urgent>]
                   [--assignee <id:uuid>...] [--label <id:uuid>...]

taskforge issue edit <id:seq_id|uuid> [--title <text>] [--description <markdown:text>]
                 [--priority <enum:none|low|medium|high|urgent>]
                 [--state <id:uuid>]
                 [--assignee <id:uuid>...] [--label <id:uuid>...]
`
}

func GetIssueQuickStartCommands() string {
	return getIssueQuickStartCommands()
}

func getModuleCommands() string {
	return `## Module (aliases: mod)
` + "```" + `
taskforge module list [--archived]
taskforge module view <id:uuid>
taskforge module create [--name <text>] [--description <markdown:text>]
                    [--status <enum:backlog|planned|in-progress|paused|completed|cancelled>]
taskforge module edit <id:uuid> [--name <text>] [--description <markdown:text>] [--status <enum:...>]
taskforge module delete <id:uuid>
taskforge module archive <id:uuid>
taskforge module issues <id:uuid>
taskforge module add-issues <module-id:uuid> <issue-ids:uuid...>
taskforge module remove-issue <module-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getStateCommands() string {
	return `## State (aliases: states)
` + "```" + `
taskforge state list
taskforge state view <id:uuid>
taskforge state create [--name <text>] [--description <markdown:text>]
                   [--color <hex:#RRGGBB>]
                   [--group <enum:backlog|unstarted|started|completed|cancelled>]
taskforge state edit <id:uuid> [--name <text>] [--description <markdown:text>]
                 [--color <hex>] [--group <enum:...>]
taskforge state delete <id:uuid>
` + "```" + `

`
}

func getLabelCommands() string {
	return `## Label (aliases: labels, tag)
` + "```" + `
taskforge label list
taskforge label view <id:uuid>
taskforge label create [--name <text>] [--description <markdown:text>] [--color <hex:#RRGGBB>]
taskforge label edit <id:uuid> [--name <text>] [--description <markdown:text>] [--color <hex>]
taskforge label delete <id:uuid>
` + "```" + `

`
}

func getIntakeCommands() string {
	return `## Intake (aliases: inbox, requests)
` + "```" + `
taskforge intake list
taskforge intake view <id:uuid>
taskforge intake create [--name <text>] [--priority <enum:low|medium|high|urgent>]
taskforge intake delete <id:uuid>
` + "```" + `

`
}

func getTypeCommands() string {
	return `## Type (aliases: issue-type)
` + "```" + `
taskforge type list
taskforge type create [--name <text>] [--description <markdown:text>]
taskforge type delete <id:uuid>
` + "```" + `

`
}

func getCycleCommands() string {
	return `## Cycle (aliases: sprint)
` + "```" + `
taskforge cycle list [--archived]
taskforge cycle view <id:uuid>
taskforge cycle create [--name <text>] [--description <markdown:text>]
                   [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
taskforge cycle edit <id:uuid> [--name <text>] [--description <markdown:text>]
                 [--start-date <YYYY-MM-DD>] [--end-date <YYYY-MM-DD>]
taskforge cycle delete <id:uuid>
taskforge cycle archive <id:uuid>
taskforge cycle issues <id:uuid>
taskforge cycle add-issues <cycle-id:uuid> <issue-ids:uuid...>
taskforge cycle remove-issue <cycle-id:uuid> <issue-id:uuid>
` + "```" + `

`
}

func getProjectCommands() string {
	return `## Project (aliases: proj)
` + "```" + `
taskforge project list
taskforge project create [<name:text>]
taskforge project info [<id:uuid>]
taskforge project delete <id:uuid>
taskforge project members [<id:uuid>]
` + "```" + `

`
}

func getWorkspaceCommands() string {
	return `## Workspace (aliases: ws)
` + "```" + `
taskforge workspace info [<slug:text>]
taskforge workspace switch [<slug:text>]
taskforge workspace members [--search <text>] [--exact] [--limit <count:int>]
` + "```" + `

`
}

// Exported wrapper functions for use by inject command

func GetGlobalFlags() string       { return getGlobalFlags() }
func GetIssueCommands() string     { return getIssueCommands() }
func GetProjectCommands() string   { return getProjectCommands() }
func GetModuleCommands() string    { return getModuleCommands() }
func GetStateCommands() string     { return getStateCommands() }
func GetLabelCommands() string     { return getLabelCommands() }
func GetTypeCommands() string      { return getTypeCommands() }
func GetCycleCommands() string     { return getCycleCommands() }
func GetWorkspaceCommands() string { return getWorkspaceCommands() }
func GetIntakeCommands() string    { return getIntakeCommands() }
