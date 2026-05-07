package config

import (
	"fmt"

	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify taskforge configuration.`,
}

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Long:  `Get a specific configuration value or all values.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runGet,
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long:  `Set a configuration value.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runSet,
}

func init() {
	ConfigCmd.AddCommand(getCmd)
	ConfigCmd.AddCommand(setCmd)
}

func normalizeConfigKey(key string) (string, error) {
	switch key {
	case "workspace", "default_workspace":
		return "workspace", nil
	case "project", "default_project":
		return "project", nil
	case "output", "output_format", "format":
		return "output", nil
	case "api_host", "api_host_url", "host":
		return "api_host", nil
	default:
		return "", fmt.Errorf("unknown config key: %q (valid keys: workspace, project, output, api_host)", key)
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Show all config
		fmt.Printf("default_workspace: %s\n", config.Cfg.DefaultWorkspace)
		fmt.Printf("default_project: %s\n", config.Cfg.DefaultProject)
		fmt.Printf("output_format: %s\n", config.Cfg.OutputFormat)
		fmt.Printf("api_host: %s\n", config.Cfg.APIHost)
		return nil
	}

	normalized, err := normalizeConfigKey(args[0])
	if err != nil {
		return err
	}

	switch normalized {
	case "workspace":
		fmt.Println(config.Cfg.DefaultWorkspace)
	case "project":
		fmt.Println(config.Cfg.DefaultProject)
	case "output":
		fmt.Println(config.Cfg.OutputFormat)
	case "api_host":
		fmt.Println(config.Cfg.APIHost)
	}

	return nil
}

func runSet(cmd *cobra.Command, args []string) error {
	normalized, err := normalizeConfigKey(args[0])
	if err != nil {
		return err
	}

	value := args[1]

	switch normalized {
	case "workspace":
		config.Cfg.DefaultWorkspace = value
	case "project":
		config.Cfg.DefaultProject = value
	case "output":
		if err := output.ValidateFormat(value); err != nil {
			return err
		}
		value = output.NormalizeFormat(value)
		config.Cfg.OutputFormat = value
	case "api_host":
		config.Cfg.APIHost = value
	}

	if err := config.SaveConfig(); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Set %s to %s", args[0], value))
	return nil
}
