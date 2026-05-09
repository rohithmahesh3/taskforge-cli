package page

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rohithmahesh3/taskforge-cli/internal/api"
	"github.com/rohithmahesh3/taskforge-cli/internal/config"
	"github.com/rohithmahesh3/taskforge-cli/internal/output"
	"github.com/rohithmahesh3/taskforge-cli/pkg/taskforge"
	"github.com/spf13/cobra"
)

var (
	pageTitle       string
	pageContent     string
	pageSlug        string
	pageCategoryID  string
	pageSortOrder   int
	pagePublished   bool
	pageSearchQuery string
	pageDeleteYes   bool

	pageField       string
	pageVersionNum  int
	pageRestoreYes  bool
	pagePatchOp     string
	pagePatchField  string
	pagePatchValue  string
	pagePatchDiff   string
	pagePatchAfter  string
	pagePatchBefore string
	pagePatchFile   string
	pagePatchActual string

	catName        string
	catDescription string
	catSortOrder   int
	catDeleteYes   bool
)

var PageCmd = &cobra.Command{
	Use:     "page",
	Aliases: []string{"pages", "doc"},
	Short:   "Manage project pages",
	Long:    `List, create, edit, and manage documentation pages in your project.`,
}

// ── Page subcommands ──

var pageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pages",
	Long: `List all pages in the current project.

Examples:
  taskforge page list
  taskforge page list --search "getting started"
  taskforge page list --category <category-id>`,
	RunE: runPageList,
}

var pageViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View page details",
	Long:  `Display detailed information about a specific page.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runPageView,
}

var pageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new page",
	Long: `Create a new page in the current project.

Examples:
  taskforge page create --title "Getting Started" --slug "getting-started" --category <category-id>
  taskforge page create -t "API Docs" -s "api-docs" -c <category-id> --content "..."`,
	RunE: runPageCreate,
}

var pageEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a page",
	Long:  `Edit an existing page's metadata or content.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runPageEdit,
}

var pageDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a page",
	Long:    `Soft-delete a page from the project.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runPageDelete,
}

// ── Page Category subcommands ──

var pageVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Page version operations",
	Long:  `View version history, specific versions, and diffs for page fields.`,
}

var pageVersionListCmd = &cobra.Command{
	Use:     "list <page-id>",
	Aliases: []string{"ls"},
	Short:   "List page version history",
	Long: `List version history for a page.

Examples:
  taskforge page version list <page-id>
  taskforge page version list <page-id> --field content`,
	Args: cobra.ExactArgs(1),
	RunE: runPageVersionList,
}

var pageVersionGetCmd = &cobra.Command{
	Use:   "get <page-id> <version-num>",
	Short: "Get a specific version of a page field",
	Long: `View a specific version of a page field.

Examples:
  taskforge page version get <page-id> <version-num> --field content`,
	Args: cobra.ExactArgs(2),
	RunE: runPageVersionGet,
}

var pageVersionDiffCmd = &cobra.Command{
	Use:   "diff <page-id> <version-num>",
	Short: "Get diff for a specific page version",
	Long: `View the diff for a specific version of a page field.

Examples:
  taskforge page version diff <page-id> <version-num> --field content`,
	Args: cobra.ExactArgs(2),
	RunE: runPageVersionDiff,
}

var pageRestoreCmd = &cobra.Command{
	Use:   "restore <page-id>",
	Short: "Restore a page to a specific version",
	Long: `Restore a page to a specific version. Version number is required.

Examples:
  taskforge page restore <page-id> --version-num 2
  taskforge page restore <page-id> --version-num 1 --field content`,
	Args: cobra.ExactArgs(1),
	RunE: runPageRestore,
}

var pagePatchCmd = &cobra.Command{
	Use:   "patch <page-id>",
	Short: "Apply text patch operations to a page",
	Long: `Apply JSON Patch operations to a page's fields.

Operations:
  - replace: Replace exact old text with new text
  - diff:    Apply a unified diff to a text field
  - insert:  Insert text after an anchor point in a field
  - delete:  Delete exact old text from a field

Examples:
  taskforge page patch <page-id> --op replace --field content --actual "old text" --value "new text"
  taskforge page patch <page-id> --op diff --field content --diff "@@ -1,3 +1,4 @@..."
  taskforge page patch <page-id> --op insert --field content --after "anchor" --value "text"
  taskforge page patch <page-id> --op delete --field content --actual "text to delete"
  taskforge page patch <page-id> --file patches.json`,
	Args: cobra.ExactArgs(1),
	RunE: runPagePatch,
}

// ── Page Category subcommands ──

var pageCategoryCmd = &cobra.Command{
	Use:     "category",
	Aliases: []string{"cat", "categories"},
	Short:   "Manage page categories",
	Long:    `List, create, edit, and manage page categories.`,
}

var catListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List page categories",
	Long:    `List all page categories in the current project.`,
	RunE:    runCategoryList,
}

var catCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a page category",
	Long: `Create a new page category.

Examples:
  taskforge page category create --name "Guides"
  taskforge page category create -n "API Reference" -d "API documentation"`,
	RunE: runCategoryCreate,
}

var catEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a page category",
	Long:  `Edit an existing page category.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCategoryEdit,
}

var catDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a page category",
	Long:    `Delete a page category.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runCategoryDelete,
}

func init() {
	// Page subcommands
	PageCmd.AddCommand(pageListCmd)
	PageCmd.AddCommand(pageViewCmd)
	PageCmd.AddCommand(pageCreateCmd)
	PageCmd.AddCommand(pageEditCmd)
	PageCmd.AddCommand(pageDeleteCmd)
	PageCmd.AddCommand(pageCategoryCmd)
	PageCmd.AddCommand(pageVersionCmd)
	PageCmd.AddCommand(pageRestoreCmd)
	PageCmd.AddCommand(pagePatchCmd)

	// List flags
	pageListCmd.Flags().StringVarP(&pageSearchQuery, "search", "s", "", "Search pages by keyword")
	pageListCmd.Flags().StringVarP(&pageCategoryID, "category", "c", "", "Filter by category ID")

	// Create flags
	pageCreateCmd.Flags().StringVarP(&pageTitle, "title", "t", "", "Page title")
	pageCreateCmd.Flags().StringVarP(&pageContent, "content", "b", "", "Page content (body)")
	pageCreateCmd.Flags().StringVarP(&pageSlug, "slug", "s", "", "URL slug")
	pageCreateCmd.Flags().StringVarP(&pageCategoryID, "category", "c", "", "Category ID")
	pageCreateCmd.Flags().IntVarP(&pageSortOrder, "sort-order", "o", 0, "Sort order")
	pageCreateCmd.Flags().BoolVar(&pagePublished, "published", false, "Publish immediately")
	_ = pageCreateCmd.MarkFlagRequired("title")
	_ = pageCreateCmd.MarkFlagRequired("slug")
	_ = pageCreateCmd.MarkFlagRequired("category")

	// Edit flags
	pageEditCmd.Flags().StringVarP(&pageTitle, "title", "t", "", "New title")
	pageEditCmd.Flags().StringVarP(&pageContent, "content", "b", "", "New content")
	pageEditCmd.Flags().StringVarP(&pageSlug, "slug", "s", "", "New URL slug")
	pageEditCmd.Flags().StringVarP(&pageCategoryID, "category", "c", "", "New category ID")
	pageEditCmd.Flags().IntVarP(&pageSortOrder, "sort-order", "o", 0, "New sort order")
	pageEditCmd.Flags().BoolVar(&pagePublished, "published", false, "Publish status")
	pageEditCmd.Flags().BoolVar(&pagePublished, "unpublished", false, "Unpublish")
	pageEditCmd.MarkFlagsMutuallyExclusive("published", "unpublished")

	// Delete flags
	pageDeleteCmd.Flags().BoolVarP(&pageDeleteYes, "yes", "y", false, "Skip confirmation")

	// --- Page Category subcommand group ---
	pageCategoryCmd.AddCommand(catListCmd)
	pageCategoryCmd.AddCommand(catCreateCmd)
	pageCategoryCmd.AddCommand(catEditCmd)
	pageCategoryCmd.AddCommand(catDeleteCmd)

	// Category create flags
	catCreateCmd.Flags().StringVarP(&catName, "name", "n", "", "Category name")
	catCreateCmd.Flags().StringVarP(&catDescription, "description", "d", "", "Category description")
	catCreateCmd.Flags().IntVarP(&catSortOrder, "sort-order", "o", 0, "Sort order")
	_ = catCreateCmd.MarkFlagRequired("name")

	// Category edit flags
	catEditCmd.Flags().StringVarP(&catName, "name", "n", "", "New name")
	catEditCmd.Flags().StringVarP(&catDescription, "description", "d", "", "New description")
	catEditCmd.Flags().IntVarP(&catSortOrder, "sort-order", "o", 0, "New sort order")

	// Category delete flags
	catDeleteCmd.Flags().BoolVarP(&catDeleteYes, "yes", "y", false, "Skip confirmation")

	// --- Version subcommand group ---
	pageVersionCmd.AddCommand(pageVersionListCmd)
	pageVersionCmd.AddCommand(pageVersionGetCmd)
	pageVersionCmd.AddCommand(pageVersionDiffCmd)

	// Version list flags
	pageVersionListCmd.Flags().StringVarP(&pageField, "field", "f", "", "Field to get versions for (content or title)")

	// Version get flags
	pageVersionGetCmd.Flags().StringVarP(&pageField, "field", "f", "", "Field to get version for (required)")
	_ = pageVersionGetCmd.MarkFlagRequired("field")

	// Version diff flags
	pageVersionDiffCmd.Flags().StringVarP(&pageField, "field", "f", "", "Field to get diff for (required)")
	_ = pageVersionDiffCmd.MarkFlagRequired("field")

	// Restore flags
	pageRestoreCmd.Flags().IntVar(&pageVersionNum, "version-num", 0, "Version number to restore to (required)")
	pageRestoreCmd.Flags().StringVarP(&pageField, "field", "f", "", "Specific field to restore (default: all versionable fields)")
	pageRestoreCmd.Flags().BoolVarP(&pageRestoreYes, "yes", "y", false, "Skip confirmation")
	_ = pageRestoreCmd.MarkFlagRequired("version-num")

	// Patch flags
	pagePatchCmd.Flags().StringVar(&pagePatchOp, "op", "", "Patch operation: replace, diff, insert, delete")
	pagePatchCmd.Flags().StringVar(&pagePatchField, "field", "", "Field to patch (content or title)")
	pagePatchCmd.Flags().StringVar(&pagePatchValue, "value", "", "Value for replace/insert operations")
	pagePatchCmd.Flags().StringVar(&pagePatchDiff, "diff", "", "Unified diff string for diff operation")
	pagePatchCmd.Flags().StringVar(&pagePatchAfter, "after", "", "Anchor text after which to insert")
	pagePatchCmd.Flags().StringVar(&pagePatchBefore, "before", "", "Legacy delete anchor input (auto-mapped to old text when possible)")
	pagePatchCmd.Flags().StringVar(&pagePatchActual, "actual", "", "Current/expected text precondition (required for replace/delete)")
	pagePatchCmd.Flags().StringVar(&pagePatchFile, "file", "", "Read patches from a JSON file")
}

// ── Page run functions ──

func runPageList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	pages, err := client.ListPages(projectID, pageSearchQuery, pageCategoryID)
	if err != nil {
		return err
	}

	if len(pages) == 0 {
		output.Info("No pages found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type pageOutput struct {
		ID         string `table:"ID" json:"id"`
		Title      string `table:"TITLE" json:"title"`
		Slug       string `table:"SLUG" json:"slug"`
		CategoryID string `table:"CATEGORY" json:"category_id,omitempty"`
		Content    string `table:"CONTENT" json:"content,omitempty"`
		Published  bool   `table:"PUBLISHED" json:"published"`
	}

	var outputs []pageOutput
	for _, p := range pages {
		contentPreview := p.Content
		if len(contentPreview) > 80 {
			contentPreview = contentPreview[:80] + "..."
		}
		outputs = append(outputs, pageOutput{
			ID:         p.ID,
			Title:      p.Title,
			Slug:       p.Slug,
			CategoryID: p.CategoryID,
			Content:    contentPreview,
			Published:  p.Published,
		})
	}

	return formatter.Print(outputs)
}

func runPageView(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	page, err := client.GetPage(projectID, pageID)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(page)
}

func runPageCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := taskforge.CreatePageRequest{
		Title:      pageTitle,
		Content:    pageContent,
		Slug:       pageSlug,
		CategoryID: pageCategoryID,
		SortOrder:  pageSortOrder,
		Published:  pagePublished,
	}

	page, err := client.CreatePage(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created page '%s' (%s)", page.Title, page.ID))
	return nil
}

func runPageEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]

	// Check that at least one flag was provided
	hasFlags := cmd.Flags().Changed("title") ||
		cmd.Flags().Changed("content") ||
		cmd.Flags().Changed("slug") ||
		cmd.Flags().Changed("category") ||
		cmd.Flags().Changed("sort-order") ||
		cmd.Flags().Changed("published") ||
		cmd.Flags().Changed("unpublished")

	if !hasFlags {
		return fmt.Errorf("no edit flags provided. Available: --title, --content, --slug, --category, --sort-order, --published/--unpublished")
	}

	req := taskforge.UpdatePageRequest{}

	if cmd.Flags().Changed("title") {
		req.Title = pageTitle
	}
	if cmd.Flags().Changed("content") {
		req.Content = pageContent
	}
	if cmd.Flags().Changed("slug") {
		req.Slug = pageSlug
	}
	if cmd.Flags().Changed("category") {
		req.CategoryID = pageCategoryID
	}
	if cmd.Flags().Changed("sort-order") {
		req.SortOrder = pageSortOrder
	}
	if cmd.Flags().Changed("published") {
		pub := true
		req.Published = &pub
	}
	if cmd.Flags().Changed("unpublished") {
		pub := false
		req.Published = &pub
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	updatedPage, err := client.UpdatePage(projectID, pageID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated page '%s'", updatedPage.Title))
	return nil
}

func runPageDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]

	if !pageDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeletePage(projectID, pageID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted page %s", pageID))
	return nil
}

// ── Category run functions ──

func runCategoryList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	categories, err := client.ListPageCategories(projectID)
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		output.Info("No page categories found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type catOutput struct {
		ID          string  `table:"ID" json:"id"`
		Name        string  `table:"NAME" json:"name"`
		Description string  `table:"DESCRIPTION" json:"description,omitempty"`
		SortOrder   float64 `table:"SORT" json:"sort_order,omitempty"`
	}

	var outputs []catOutput
	for _, c := range categories {
		outputs = append(outputs, catOutput{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			SortOrder:   c.SortOrder,
		})
	}

	return formatter.Print(outputs)
}

func runCategoryCreate(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified. Use --project flag or set default project")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := taskforge.CreatePageCategoryRequest{
		Name:        catName,
		Description: catDescription,
		SortOrder:   catSortOrder,
	}

	category, err := client.CreatePageCategory(projectID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Created page category '%s' (%s)", category.Name, category.ID))
	return nil
}

func runCategoryEdit(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	categoryID := args[0]

	hasFlags := cmd.Flags().Changed("name") || cmd.Flags().Changed("description") || cmd.Flags().Changed("sort-order")
	if !hasFlags {
		return fmt.Errorf("no edit flags provided. Available: --name, --description, --sort-order")
	}

	req := taskforge.UpdatePageCategoryRequest{}

	if cmd.Flags().Changed("name") {
		req.Name = catName
	}
	if cmd.Flags().Changed("description") {
		req.Description = catDescription
	}
	if cmd.Flags().Changed("sort-order") {
		req.SortOrder = catSortOrder
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	updated, err := client.UpdatePageCategory(projectID, categoryID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Updated page category '%s'", updated.Name))
	return nil
}

func runCategoryDelete(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	categoryID := args[0]

	if !catDeleteYes {
		return fmt.Errorf("confirmation required; use --yes / -y flag to confirm deletion")
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	if err := client.DeletePageCategory(projectID, categoryID); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Deleted page category %s", categoryID))
	return nil
}

// ── Page Version run functions ──

func runPageVersionList(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]
	if err := validatePageField(pageField, true); err != nil {
		return err
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	versions, err := client.GetPageVersions(projectID, pageID, pageField)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		output.Info("No versions found")
		return nil
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)

	type versionOutput struct {
		Version   int    `table:"VER" json:"version"`
		Field     string `table:"FIELD" json:"field"`
		Content   string `table:"CONTENT" json:"content,omitempty"`
		ActorType string `table:"ACTOR" json:"actor_type,omitempty"`
	}

	var outputs []versionOutput
	for _, v := range versions {
		content := v.Content
		if len(content) > 80 {
			content = content[:80] + "..."
		}
		outputs = append(outputs, versionOutput{
			Version:   v.Version,
			Field:     v.Field,
			Content:   content,
			ActorType: v.ActorType,
		})
	}

	return formatter.Print(outputs)
}

func runPageVersionGet(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]
	versionNum := args[1]
	if err := validatePageField(pageField, false); err != nil {
		return err
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	version, err := client.GetPageVersion(projectID, pageID, versionNum, pageField)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(version)
}

func runPageVersionDiff(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]
	versionNum := args[1]
	if err := validatePageField(pageField, false); err != nil {
		return err
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	diff, err := client.GetPageVersionDiff(projectID, pageID, versionNum, pageField)
	if err != nil {
		return err
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(diff)
}

func runPageRestore(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]

	if !cmd.Flags().Changed("version-num") {
		return fmt.Errorf("--version-num is required")
	}
	if err := validatePageField(pageField, true); err != nil {
		return err
	}
	if pageVersionNum < 1 {
		return fmt.Errorf("must be a positive integer (got %d)", pageVersionNum)
	}

	if !pageRestoreYes {
		output.Info(fmt.Sprintf("Will restore page %s to version %d. Use --yes to confirm.", pageID, pageVersionNum))
		return nil
	}

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	req := taskforge.RestorePageRequest{
		Version: pageVersionNum,
	}
	if pageField != "" {
		req.Field = pageField
	}

	page, err := client.RestorePage(projectID, pageID, req)
	if err != nil {
		return err
	}

	output.Success(fmt.Sprintf("Restored page '%s' to version %d", page.Title, pageVersionNum))
	return nil
}

func runPagePatch(cmd *cobra.Command, args []string) error {
	projectID := config.Cfg.DefaultProject
	if projectID == "" {
		return fmt.Errorf("no project specified")
	}

	pageID := args[0]

	client, err := api.NewClient()
	if err != nil {
		return err
	}

	var patches []taskforge.PatchOp
	if pagePatchFile != "" {
		patches, err = loadPagePatchesFromFile(pagePatchFile)
		if err != nil {
			return fmt.Errorf("failed to load patches from file: %w", err)
		}
	} else {
		patches, err = buildPagePatchFromFlags()
		if err != nil {
			return err
		}
	}

	if len(patches) == 0 {
		return fmt.Errorf("no patch operations specified. Use --op/--field flags or --file to provide patches")
	}

	page, err := client.PatchPage(projectID, pageID, patches)
	if err != nil {
		var apiErr *api.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 409 && len(apiErr.Conflicts) > 0 {
			return fmt.Errorf("failed to apply patch: %s (conflicts: %s)", apiErr.Message, string(apiErr.Conflicts))
		}
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	formatter := output.NewFormatter(config.Cfg.OutputFormat, false)
	return formatter.Print(page)
}

func buildPagePatchFromFlags() ([]taskforge.PatchOp, error) {
	if pagePatchOp == "" {
		return nil, fmt.Errorf("--op flag is required when not using --file")
	}
	if pagePatchField == "" {
		return nil, fmt.Errorf("--field flag is required when not using --file")
	}
	if err := validatePageField(pagePatchField, false); err != nil {
		return nil, err
	}

	p := taskforge.PatchOp{
		Op:      pagePatchOp,
		Field:   pagePatchField,
		After:   pagePatchAfter,
		Old:     pagePatchActual,
		New:     pagePatchValue,
		Content: pagePatchValue,
		Unified: pagePatchDiff,
		Value:   pagePatchValue,
		Diff:    pagePatchDiff,
		Before:  pagePatchBefore,
	}

	switch pagePatchOp {
	case "replace":
		if pagePatchValue == "" {
			return nil, fmt.Errorf("--value is required for replace operation")
		}
		if pagePatchActual == "" {
			return nil, fmt.Errorf("--actual is required for replace operation")
		}
	case "diff":
		if pagePatchDiff == "" {
			return nil, fmt.Errorf("--diff is required for diff operation")
		}
	case "insert":
		if pagePatchValue == "" {
			return nil, fmt.Errorf("--value is required for insert operation")
		}
		if pagePatchAfter == "" {
			return nil, fmt.Errorf("--after is required for insert operation")
		}
	case "delete":
		if pagePatchActual == "" && pagePatchBefore == "" {
			return nil, fmt.Errorf("--actual is required for delete operation")
		}
		if p.Old == "" {
			p.Old = pagePatchBefore
		}
	default:
		return nil, fmt.Errorf("unsupported operation %q: must be one of replace, diff, insert, delete", pagePatchOp)
	}

	return normalizePagePatches([]taskforge.PatchOp{p}), nil
}

func loadPagePatchesFromFile(path string) ([]taskforge.PatchOp, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var req taskforge.PatchRequest
	if err := json.Unmarshal(data, &req); err == nil && len(req.Patches) > 0 {
		return req.Patches, nil
	}

	var patches []taskforge.PatchOp
	if err := json.Unmarshal(data, &patches); err != nil {
		return nil, fmt.Errorf("file must contain a JSON array of patch operations or a JSON object with a 'patches' field")
	}

	return normalizePagePatches(patches), nil
}

func normalizePagePatches(patches []taskforge.PatchOp) []taskforge.PatchOp {
	out := make([]taskforge.PatchOp, 0, len(patches))
	for _, p := range patches {
		switch p.Op {
		case "replace":
			if p.New == "" && p.Value != "" {
				p.New = p.Value
			}
		case "insert":
			if p.Content == "" && p.Value != "" {
				p.Content = p.Value
			}
		case "delete":
			if p.Old == "" && p.Before != "" {
				p.Old = p.Before
			}
		case "diff":
			if p.Unified == "" && p.Diff != "" {
				p.Unified = p.Diff
			}
		}
		out = append(out, p)
	}
	return out
}

func validatePageField(field string, allowEmpty bool) error {
	normalized := strings.TrimSpace(field)
	if normalized == "" {
		if allowEmpty {
			return nil
		}
		return fmt.Errorf("--field is required")
	}
	if normalized != "content" && normalized != "title" {
		return fmt.Errorf("invalid --field %q: must be one of content, title", field)
	}
	return nil
}
