package page

import (
	"fmt"

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
