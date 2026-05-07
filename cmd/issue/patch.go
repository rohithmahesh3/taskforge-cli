package issue

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	patchOp     string
	patchField  string
	patchValue  string
	patchDiff   string
	patchAfter  string
	patchBefore string
	patchFile   string
	patchActual string
)

func init() {
	patchCmd := &cobra.Command{
		Use:   "patch <issue-id>",
		Short: "Apply RFC 6902-style JSON Patch to an issue",
		Long: `Apply JSON Patch operations to an issue's fields.

Operations:
  - replace: Replace a field value entirely
  - diff:    Apply a unified diff to a text field
  - insert:  Insert text after an anchor point in a field
  - delete:  Delete text before an anchor point in a field

Examples:
  taskforge issue patch TF-1 --op replace --field description --value "new text"
  taskforge issue patch TF-1 --op diff --field description --diff "@@ -1,3 +1,4 @@..."
  taskforge issue patch TF-1 --op insert --field description --after "anchor" --value "text"
  taskforge issue patch TF-1 --op delete --field description --before "anchor"
  taskforge issue patch TF-1 --file patches.json`,
		Args: cobra.ExactArgs(1),
		RunE: runPatch,
	}

	patchCmd.Flags().StringVar(&patchOp, "op", "", "Patch operation: replace, diff, insert, delete")
	patchCmd.Flags().StringVar(&patchField, "field", "", "Field to patch (e.g., description)")
	patchCmd.Flags().StringVar(&patchValue, "value", "", "Value for replace/insert operations")
	patchCmd.Flags().StringVar(&patchDiff, "diff", "", "Unified diff string for diff operation")
	patchCmd.Flags().StringVar(&patchAfter, "after", "", "Anchor text after which to insert")
	patchCmd.Flags().StringVar(&patchBefore, "before", "", "Anchor text before which to delete")
	patchCmd.Flags().StringVar(&patchActual, "actual", "", "Current value of the field (safety precondition for replace/diff operations)")
	patchCmd.Flags().StringVar(&patchFile, "file", "", "Read patches from a JSON file")

	IssueCmd.AddCommand(patchCmd)
}

func runPatch(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	issueRef := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	issueID, err := resolveIssueID(client, projectID, issueRef)
	if err != nil {
		return err
	}

	var patches []taskforge.PatchOp

	if patchFile != "" {
		patches, err = loadPatchesFromFile(patchFile)
		if err != nil {
			return fmt.Errorf("failed to load patches from file: %w", err)
		}
	} else {
		patches, err = buildPatchFromFlags()
		if err != nil {
			return err
		}
	}

	if len(patches) == 0 {
		return fmt.Errorf("no patch operations specified. Use --op/--field flags or --file to provide patches")
	}

	issue, err := client.PatchIssue(projectID, issueID, patches)
	if err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	output.Success(fmt.Sprintf("Patched issue %d", issue.SequenceID))
	return nil
}

func buildPatchFromFlags() ([]taskforge.PatchOp, error) {
	if patchOp == "" {
		return nil, fmt.Errorf("--op flag is required when not using --file")
	}
	if patchField == "" {
		return nil, fmt.Errorf("--field flag is required when not using --file")
	}

	p := taskforge.PatchOp{
		Op:     patchOp,
		Field:  patchField,
		Value:  patchValue,
		Diff:   patchDiff,
		After:  patchAfter,
		Before: patchBefore,
		Old:    patchActual,
	}

	// Validate the operation has required fields
	switch patchOp {
	case "replace":
		if patchValue == "" {
			return nil, fmt.Errorf("--value is required for replace operation")
		}
	case "diff":
		if patchDiff == "" {
			return nil, fmt.Errorf("--diff is required for diff operation")
		}
	case "insert":
		if patchValue == "" {
			return nil, fmt.Errorf("--value is required for insert operation")
		}
		if patchAfter == "" {
			return nil, fmt.Errorf("--after is required for insert operation")
		}
	case "delete":
		if patchBefore == "" {
			return nil, fmt.Errorf("--before is required for delete operation")
		}
	default:
		return nil, fmt.Errorf("unsupported operation %q: must be one of replace, diff, insert, delete", patchOp)
	}

	return []taskforge.PatchOp{p}, nil
}

func loadPatchesFromFile(path string) ([]taskforge.PatchOp, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Try parsing as PatchRequest ({"patches": [...]})
	var req taskforge.PatchRequest
	if err := json.Unmarshal(data, &req); err == nil && len(req.Patches) > 0 {
		return req.Patches, nil
	}

	// Try parsing as a plain array of PatchOp
	var patches []taskforge.PatchOp
	if err := json.Unmarshal(data, &patches); err != nil {
		return nil, fmt.Errorf("file must contain a JSON array of patch operations or a JSON object with a 'patches' field")
	}

	return patches, nil
}