package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/rohithmahesh3/taskforge-cli/cmd/auth"
	"github.com/rohithmahesh3/taskforge-cli/cmd/config"
	"github.com/rohithmahesh3/taskforge-cli/cmd/context"
	"github.com/rohithmahesh3/taskforge-cli/cmd/cycle"
	"github.com/rohithmahesh3/taskforge-cli/cmd/inject"
	"github.com/rohithmahesh3/taskforge-cli/cmd/intake"
	"github.com/rohithmahesh3/taskforge-cli/cmd/issue"
	"github.com/rohithmahesh3/taskforge-cli/cmd/label"
	"github.com/rohithmahesh3/taskforge-cli/cmd/module"
	"github.com/rohithmahesh3/taskforge-cli/cmd/project"
	"github.com/rohithmahesh3/taskforge-cli/cmd/state"
	issuetype "github.com/rohithmahesh3/taskforge-cli/cmd/type"
	"github.com/rohithmahesh3/taskforge-cli/cmd/workspace"
	cfg "github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	workspaceSlug string
	projectID     string
	outputFmt     string
	noColor       bool
	configFile    string
)

var rootCmd = &cobra.Command{
	Use:   "taskforge",
	Short: "A CLI tool for managing TaskForge projects",
	Long: `taskforge is a command-line interface for TaskForge project management.

It allows you to manage workspaces, projects, issues, cycles, and modules
from the comfort of your terminal.

Get started:
  taskforge auth login                    # Authenticate with TaskForge
  taskforge workspace info                # Show configured workspace access
  taskforge workspace members --search alice # Find assignable workspace users by name/email
  taskforge project list                  # List projects in current workspace
  taskforge issue list                    # List issues in current project`,
	Version: version,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config initialization for certain commands
		if cmd.Name() == "login" || cmd.Name() == "version" || cmd.Name() == "completion" || cmd.Name() == "init" {
			return nil
		}

		if configFile != "" {
			cfg.SetConfigFile(configFile)
		}

		allowInvalidOutputConfig := shouldAllowInvalidOutputConfig(cmd, args)
		initConfig := cfg.InitConfig
		if allowInvalidOutputConfig {
			initConfig = cfg.InitConfigAllowInvalidOutput
		}

		if err := initConfig(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// Override config with flags
		if workspaceSlug != "" {
			cfg.Cfg.DefaultWorkspace = workspaceSlug
		}
		if projectID != "" {
			cfg.Cfg.DefaultProject = projectID
		}
		if outputFmt != "" {
			cfg.Cfg.OutputFormat = outputFmt
		}
		if !allowInvalidOutputConfig {
			if err := output.ValidateFormat(cfg.Cfg.OutputFormat); err != nil {
				return err
			}
			cfg.Cfg.OutputFormat = output.NormalizeFormat(cfg.Cfg.OutputFormat)
		}

		return nil
	},
}

func shouldAllowInvalidOutputConfig(cmd *cobra.Command, args []string) bool {
	return cmd.CommandPath() == "taskforge config set" && len(args) >= 1 && args[0] == "output"
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&workspaceSlug, "workspace", "", "TaskForge workspace slug (overrides config)")
	rootCmd.PersistentFlags().StringVar(&projectID, "project", "", "TaskForge project ID (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "", "Output format: json, yaml (overrides config)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path (default: ~/.config/taskforge/config.yaml)")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(workspace.WorkspaceCmd)
	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(issue.IssueCmd)
	rootCmd.AddCommand(state.StateCmd)
	rootCmd.AddCommand(label.LabelCmd)
	rootCmd.AddCommand(cycle.CycleCmd)
	rootCmd.AddCommand(module.ModuleCmd)
	rootCmd.AddCommand(intake.IntakeCmd)
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(context.ContextCmd)
	rootCmd.AddCommand(issuetype.TypeCmd)
	rootCmd.AddCommand(inject.InjectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to get version from Go module info (works with go install)
		// Version is set when installed via: go install github.com/rohithmahesh3/taskforge-cli@v1.0.1
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			fmt.Printf("taskforge version %s\n", info.Main.Version)
			return
		}

		// Fall back to ldflags for Makefile builds or local development
		fmt.Printf("taskforge version %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(taskforge completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ taskforge completion bash > /etc/bash_completion.d/taskforge
  # macOS:
  $ taskforge completion bash > $(brew --prefix)/etc/bash_completion.d/taskforge

Zsh:
  $ source <(taskforge completion zsh)
  # To load completions for each session, execute once:
  $ taskforge completion zsh > "${fpath[1]}/_taskforge"

Fish:
  $ source <(taskforge completion fish)
  # To load completions for each session, execute once:
  $ taskforge completion fish > ~/.config/fish/completions/taskforge.fish

PowerShell:
  PS> taskforge completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> taskforge completion powershell > taskforge.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}
