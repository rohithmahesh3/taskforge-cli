package taskforge

import (
	"encoding/json"
	"time"
)

type Workspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	LogoURL     string    `json:"logo_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Project struct {
	ID                    string      `json:"id"`
	Name                  string      `json:"name"`
	Identifier            string      `json:"identifier"`
	Description           string      `json:"description,omitempty"`
	CoverImage            string      `json:"cover_image,omitempty"`
	IconProp              IconProp    `json:"icon_prop,omitempty"`
	ModuleView            bool        `json:"module_view,omitempty"`
	CycleView             bool        `json:"cycle_view,omitempty"`
	PageView              bool        `json:"page_view,omitempty"`
	IntakeView            bool        `json:"intake_view,omitempty"`
	IsTimeTrackingEnabled bool        `json:"is_time_tracking_enabled,omitempty"`
	IsIssueTypeEnabled    bool        `json:"is_issue_type_enabled,omitempty"`
	Estimate              interface{} `json:"estimate,omitempty"`
	CreatedAt             time.Time   `json:"created_at"`
	UpdatedAt             time.Time   `json:"updated_at"`
}

type IconProp struct {
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

type Issue struct {
	ID            string          `json:"id"`
	Identifier    string          `json:"identifier,omitempty"`
	SequenceID    int             `json:"sequence_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description,omitempty"`
	State         FlexibleState   `json:"state"`
	Priority      string          `json:"priority"`
	Assignees     []FlexibleUser  `json:"assignees,omitempty"`
	Labels        []FlexibleLabel `json:"labels,omitempty"`
	CycleID       string          `json:"cycle_id,omitempty"`
	ModuleID      string          `json:"module_id,omitempty"`
	Parent        string          `json:"parent,omitempty"`
	StartDate     string          `json:"start_date,omitempty"`
	TargetDate    string          `json:"target_date,omitempty"`
	EstimatePoint int             `json:"estimate_point,omitempty"`
	Type          string          `json:"type,omitempty"`
	CompletedAt   *time.Time      `json:"completed_at,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	CreatedBy     string          `json:"created_by,omitempty"`
	UpdatedBy     string          `json:"updated_by,omitempty"`
	ProjectID     string          `json:"project_id,omitempty"`
	Project       string          `json:"project,omitempty"`
	WorkspaceID   string          `json:"workspace_id,omitempty"`
	Workspace     string          `json:"workspace,omitempty"`
	IsDraft       bool            `json:"is_draft,omitempty"`
	ArchivedAt    string          `json:"archived_at,omitempty"`
}

type State struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Color         string    `json:"color"`
	WorkspaceSlug string    `json:"workspace_slug,omitempty"`
	Sequence      float64   `json:"sequence,omitempty"`
	Group         string    `json:"group,omitempty"`
	IsDefault     bool      `json:"default,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type CreateStateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color"`
	Group       string `json:"group,omitempty"`
}

type UpdateStateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Group       string `json:"group,omitempty"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Role        int    `json:"role,omitempty"`
}

type Label struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateLabelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type UpdateLabelRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

type Cycle struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitempty"`
	Status      string    `json:"status,omitempty"`
	OwnedBy     string    `json:"owned_by,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ArchivedAt  string    `json:"archived_at,omitempty"`
}

type CreateCycleRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	OwnedBy        string `json:"owned_by,omitempty"`
	ExternalSource string `json:"external_source,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
}

type UpdateCycleRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	OwnedBy        string `json:"owned_by,omitempty"`
	ExternalSource string `json:"external_source,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
}

type Module struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	TargetDate  string    `json:"target_date,omitempty"`
	Status      string    `json:"status,omitempty"`
	Lead        string    `json:"lead,omitempty"`
	Members     []string  `json:"members,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ArchivedAt  string    `json:"archived_at,omitempty"`
}

type CreateModuleRequest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	StartDate      string   `json:"start_date,omitempty"`
	TargetDate     string   `json:"target_date,omitempty"`
	Status         string   `json:"status,omitempty"`
	Lead           string   `json:"lead,omitempty"`
	Members        []string `json:"members,omitempty"`
	ExternalSource string   `json:"external_source,omitempty"`
	ExternalID     string   `json:"external_id,omitempty"`
}

type UpdateModuleRequest struct {
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	StartDate      string   `json:"start_date,omitempty"`
	TargetDate     string   `json:"target_date,omitempty"`
	Status         string   `json:"status,omitempty"`
	Lead           string   `json:"lead,omitempty"`
	Members        []string `json:"members,omitempty"`
	ExternalSource string   `json:"external_source,omitempty"`
	ExternalID     string   `json:"external_id,omitempty"`
}

type CreateIssueRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Priority      string   `json:"priority,omitempty"`
	Assignees     []string `json:"assignees,omitempty"`
	Labels        []string `json:"labels,omitempty"`
	Parent        string   `json:"parent,omitempty"`
	EstimatePoint int      `json:"estimate_point,omitempty"`
	Type          string   `json:"type,omitempty"`
	Module        string   `json:"module,omitempty"`
	StartDate     string   `json:"start_date,omitempty"`
	TargetDate    string   `json:"target_date,omitempty"`
}

type UpdateIssueRequest struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	State         string   `json:"state,omitempty"`
	Priority      string   `json:"priority,omitempty"`
	Assignees     []string `json:"assignees_list,omitempty"`
	Labels        []string `json:"labels_list,omitempty"`
	Parent        string   `json:"parent,omitempty"`
	EstimatePoint int      `json:"estimate_point,omitempty"`
	Type          string   `json:"type,omitempty"`
	Module        string   `json:"module,omitempty"`
	StartDate     string   `json:"start_date,omitempty"`
	TargetDate    string   `json:"target_date,omitempty"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"description,omitempty"`
}

type Worklog struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	CreatedBy   string    `json:"created_by,omitempty"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	ProjectID   string    `json:"project_id,omitempty"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	LoggedBy    string    `json:"logged_by,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateWorklogRequest struct {
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

type UpdateWorklogRequest struct {
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"`
}

type WorklogTotal struct {
	TotalTime int `json:"total_time"`
}

type Attachment struct {
	ID              string                 `json:"id"`
	FileName        string                 `json:"file_name,omitempty"`
	FileType        string                 `json:"file_type,omitempty"`
	FileSize        FlexibleInt64          `json:"file_size,omitempty"`
	Asset           string                 `json:"asset,omitempty"`
	EntityType      string                 `json:"entity_type,omitempty"`
	IsDeleted       bool                   `json:"is_deleted,omitempty"`
	IsArchived      bool                   `json:"is_archived,omitempty"`
	Size            FlexibleInt64          `json:"size,omitempty"`
	IsUploaded      bool                   `json:"is_uploaded,omitempty"`
	StorageMetadata map[string]interface{} `json:"storage_metadata,omitempty"`
	CreatedBy       string                 `json:"created_by,omitempty"`
	UpdatedBy       string                 `json:"updated_by,omitempty"`
	WorkspaceID     string                 `json:"workspace_id,omitempty"`
	Workspace       string                 `json:"workspace,omitempty"`
	ProjectID       string                 `json:"project_id,omitempty"`
	Project         string                 `json:"project,omitempty"`
	IssueID         string                 `json:"issue_id,omitempty"`
	Issue           string                 `json:"issue,omitempty"`
	CreatedAt       time.Time              `json:"created_at,omitempty"`
	UpdatedAt       time.Time              `json:"updated_at,omitempty"`
}

type AttachmentAttributes struct {
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
	Type string `json:"type,omitempty"`
}

type UploadData struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields,omitempty"`
}

type UploadCredentials struct {
	ID         string            `json:"id,omitempty"`
	Upload     UploadData        `json:"upload,omitempty"`
	UploadData UploadData        `json:"upload_data,omitempty"`
	AssetID    string            `json:"asset_id,omitempty"`
	Attachment Attachment        `json:"attachment,omitempty"`
	AssetURL   string            `json:"asset_url,omitempty"`
	URL        string            `json:"url,omitempty"`
	Fields     map[string]string `json:"fields,omitempty"`
}

func (u UploadCredentials) UploadTarget() UploadData {
	if u.Upload.URL != "" {
		return u.Upload
	}

	if u.UploadData.URL != "" {
		return u.UploadData
	}

	return UploadData{
		URL:    u.URL,
		Fields: u.Fields,
	}
}

type IntakeIssue struct {
	ID          string        `json:"id"`
	Status      int           `json:"status"`
	SnoozedTill string        `json:"snoozed_till,omitempty"`
	Source      string        `json:"source,omitempty"`
	Inbox       string        `json:"inbox,omitempty"`
	Issue       FlexibleIssue `json:"issue,omitempty"`
	DuplicateTo string        `json:"duplicate_to,omitempty"`
	ProjectID   string        `json:"project_id,omitempty"`
	Project     string        `json:"project,omitempty"`
	WorkspaceID string        `json:"workspace_id,omitempty"`
	Workspace   string        `json:"workspace,omitempty"`
	CreatedBy   string        `json:"created_by,omitempty"`
	UpdatedBy   string        `json:"updated_by,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

type CreateIntakeIssueRequest struct {
	Name     string `json:"name"`
	Priority string `json:"priority,omitempty"`
	Source   string `json:"source,omitempty"`
}

type Activity struct {
	ID             string    `json:"id"`
	Verb           string    `json:"verb"`
	Field          string    `json:"field,omitempty"`
	OldValue       string    `json:"old_value,omitempty"`
	NewValue       string    `json:"new_value,omitempty"`
	Comment        string    `json:"comment,omitempty"`
	Attachments    []string  `json:"attachments,omitempty"`
	OldIdentifier  string    `json:"old_identifier,omitempty"`
	NewIdentifier  string    `json:"new_identifier,omitempty"`
	Epoch          float64   `json:"epoch,omitempty"`
	ProjectID      string    `json:"project_id,omitempty"`
	Project        string    `json:"project,omitempty"`
	WorkspaceID    string    `json:"workspace_id,omitempty"`
	Workspace      string    `json:"workspace,omitempty"`
	IssueID        string    `json:"issue_id,omitempty"`
	Issue          string    `json:"issue,omitempty"`
	IssueCommentID string    `json:"issue_comment_id,omitempty"`
	IssueComment   string    `json:"issue_comment,omitempty"`
	ActorID        string    `json:"actor_id,omitempty"`
	Actor          string    `json:"actor,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

type UpdateAttachmentRequest struct {
	Attributes       AttachmentAttributes   `json:"attributes,omitempty"`
	Asset            string                 `json:"asset,omitempty"`
	EntityType       string                 `json:"entity_type,omitempty"`
	EntityIdentifier string                 `json:"entity_identifier,omitempty"`
	IsDeleted        bool                   `json:"is_deleted,omitempty"`
	IsArchived       bool                   `json:"is_archived,omitempty"`
	ExternalID       string                 `json:"external_id,omitempty"`
	ExternalSource   string                 `json:"external_source,omitempty"`
	Size             int64                  `json:"size,omitempty"`
	IsUploaded       bool                   `json:"is_uploaded,omitempty"`
	StorageMetadata  map[string]interface{} `json:"storage_metadata,omitempty"`
}

// PatchOp represents a single TaskForge text patch operation.
type PatchOp struct {
	Op      string `json:"op"`                // replace, diff, insert, delete
	Field   string `json:"field"`             // e.g. "description"
	Old     string `json:"old,omitempty"`     // required for replace/delete precondition checks
	New     string `json:"new,omitempty"`     // required for replace
	Content string `json:"content,omitempty"` // required for insert
	Unified string `json:"unified,omitempty"` // required for diff
	After   string `json:"after,omitempty"`   // required for insert

	// Legacy input compatibility fields. CLI normalizes these before sending.
	Value  string `json:"value,omitempty"`  // legacy replace/insert payload
	Diff   string `json:"diff,omitempty"`   // legacy diff payload
	Before string `json:"before,omitempty"` // legacy delete anchor payload
}

// PatchRequest represents a patch request body containing a list of operations
type PatchRequest struct {
	Patches []PatchOp `json:"patches"`
}

// VersionHistory represents a version history entry for an issue or comment
type VersionHistory struct {
	ID         string          `json:"id"`
	VersionNum int             `json:"version"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	Field      string          `json:"field"`
	Content    string          `json:"content"`
	Diff       string          `json:"diff"`
	PatchOps   json.RawMessage `json:"patch_ops"`
	CreatedAt  string          `json:"created_at"`
	ActorID    string          `json:"actor_id"`
	ActorType  string          `json:"actor_type"`
	CreatedBy  FlexibleUser    `json:"created_by"`
}

// RestoreRequest represents a request to restore a deleted entity
type RestoreRequest struct {
	Version int `json:"version"`
}

type IssueRestoreRequest struct {
	Version int `json:"version"`
}

type CommentRestoreRequest struct {
	Version int `json:"version"`
}

// ── Pages ──

type PageCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreatePageCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

type UpdatePageCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

type Page struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content,omitempty"`
	CategoryID  string    `json:"category_id,omitempty"`
	SortOrder   float64   `json:"sort_order,omitempty"`
	Published   bool      `json:"published,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreatePageRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content,omitempty"`
	Slug       string `json:"slug"`
	CategoryID string `json:"category_id"`
	SortOrder  int    `json:"sort_order,omitempty"`
	Published  bool   `json:"published,omitempty"`
}

type UpdatePageRequest struct {
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Slug       string `json:"slug,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	SortOrder  int    `json:"sort_order,omitempty"`
	Published  *bool  `json:"published,omitempty"`
}

// ContentVersion represents a version snapshot of a field value.
type ContentVersion struct {
	Version   int        `json:"version"`
	Content   string     `json:"content,omitempty"`
	Diff      string     `json:"diff,omitempty"`
	Field     string     `json:"field,omitempty"`
	ActorID   string     `json:"actor_id,omitempty"`
	ActorType string     `json:"actor_type,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
}

// VersionDiff represents a diff between two versions.
type VersionDiff struct {
	Version         int    `json:"version"`
	Field           string `json:"field"`
	Diff            string `json:"diff"`
	PreviousContent string `json:"previous_content,omitempty"`
	CurrentContent  string `json:"current_content,omitempty"`
}

// RestorePageRequest represents a request to restore a page to a specific version.
type RestorePageRequest struct {
	Version int    `json:"version"`
	Field   string `json:"field,omitempty"`
}

// DependencyViewItem represents an issue dependency in a project
type DependencyViewItem struct {
	ID          string `json:"id"`
	FromIssueID string `json:"from_issue_id"`
	ToIssueID   string `json:"to_issue_id"`
	Project     string `json:"project"`
	Workspace   string `json:"workspace"`
	CreatedBy   string `json:"created_by,omitempty"`
	UpdatedBy   string `json:"updated_by,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	RelatedIssue struct {
		ID         string `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
	} `json:"related_issue"`
}

// GroupedDependencies represents issue dependencies grouped by direction
type GroupedDependencies struct {
	DependsOn []DependencyViewItem `json:"depends_on"`
	Blocks    []DependencyViewItem `json:"blocks"`
}
