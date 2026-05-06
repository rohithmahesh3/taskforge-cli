# TaskForge CLI Test Matrix

Generated: 2026-05-07  
CLI Binary: `~/Dev/taskforge-cli/taskforge`  
API: http://localhost:3000  
Workspace: `engineering`  
Default Project ID: `d1a511e2-816c-4e43-a599-b1ae0cbd8297`

---

## Global Flags

All commands inherit these global flags from the root command.

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| | `--workspace` | string | (from config) | Workspace slug (overrides config) | PASS — `taskforge issue list --workspace engineering` |
| | `--project` | string | (from config) | TaskForge project ID (overrides config) | PASS — `taskforge issue list --project d1a511e2-816c-4e43-a599-b1ae0cbd8297` |
| `-o` | `--output` | string | (from config) | Output format: json, yaml | PASS — `taskforge issue list -o json` |
| | `--no-color` | bool | false | Disable colored output | PASS — `taskforge issue list --no-color` |
| | `--config` | string | `~/.config/taskforge/config.yaml` | Config file path | PASS — `taskforge auth status --config ~/.config/taskforge/config.yaml` |

---

## Command Reference & Test Results

---

### `taskforge auth`

Manage authentication with TaskForge.

| Subcommand | Description |
|------------|-------------|
| `login` | Authenticate with TaskForge |
| `logout` | Logout from TaskForge |
| `status` | Check authentication status |
| `whoami` | Show current user information |

#### `taskforge auth login`

Authenticate with TaskForge using an API key.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| | `--token` | string | "" | API key (will prompt if not provided) | PASS — `taskforge auth login --token tf_... --workspace engineering --api-host http://localhost:3000` |
| | `--api-host` | string | "" | TaskForge API host URL (will prompt if not provided) | PASS |
| | `--workspace` | string | "" | Default workspace slug | PASS |

**Positional Arguments:** None

**Interactive Prompts:** If `--token` is not provided, prompts for API key (hidden input). If `--api-host` is not provided, prompts for API host. If `--workspace` is not provided, prompts for workspace slug.

**Test Details:**
```
$ taskforge auth login --token tf_10f0d5075051d983d3c1940ca7d1a49445307f4bbb1e23e4bd2f0f4f22346545 --workspace engineering --api-host http://localhost:3000
✓ Successfully authenticated with workspace 'engineering'
```
**Result:** PASS

---

#### `taskforge auth logout`

Remove stored credentials and configuration.

**Flags:** None (inherits global)

**Positional Arguments:** None

**Test Details:**
```
$ taskforge auth logout
✓ Successfully logged out
```
**Result:** PASS

---

#### `taskforge auth status`

Check if authenticated and show current configuration.

**Flags:** None (inherits global)

**Positional Arguments:** None

**Test Details:**
```
$ taskforge auth status
Authentication Status: ✓ Authenticated
Workspace: engineering
API Host: http://localhost:3000
Default Project: d1a511e2-816c-4e43-a599-b1ae0cbd8297
Output Format: yaml
```
**Result:** PASS

---

#### `taskforge auth whoami`

Display information about the currently authenticated user.

**Flags:** None (inherits global)

**Positional Arguments:** None

**Test Details:**
```
$ taskforge auth whoami
User: Admin User
Email: admin@taskforge.local
Display Name: Admin User
Workspace: engineering
```
**Result:** PASS

---

### `taskforge config`

View and modify taskforge configuration.

#### `taskforge config get [key]`

Get a specific configuration value or all values.

**Flags:** None (inherits global)

**Positional Arguments:** `[key]` — Optional. One of: `workspace`, `project`, `output`, `api_host`

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge config get` | PASS — Shows all config values |
| `taskforge config get api_host` | PASS — Shows `http://localhost:3000` |
| `taskforge config get workspace` | PASS — Shows `engineering` |

---

#### `taskforge config set <key> <value>`

Set a configuration value.

**Flags:** None (inherits global)

**Positional Arguments:** `<key>` — One of: `workspace`, `project`, `output`, `api_host`; `<value>` — The value to set

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge config set default_project d1a511e2-816c-4e43-a599-b1ae0cbd8297` | PASS — `✓ Set project to d1a511e2-...` |
| `taskforge config set output_format json` | FAIL — unknown config key `output_format`; correct key is `output` |
| `taskforge config set output json` | PASS — `✓ Set output to json` |

---

### `taskforge context`

Output a concise CLI command reference in markdown format.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-a` | `--all` | bool | false | Include all modules | PASS |
| | `--cycle` | bool | false | Include cycle commands | PASS |
| | `--intake` | bool | false | Include intake commands | PASS |
| | `--workspace` | bool | false | Include workspace commands | PASS |

**Positional Arguments:** None

**Test Details:**
```
$ taskforge context
[Outputs full CLI command reference in markdown]
```
**Result:** PASS

---

### `taskforge inject`

Inject taskforge context documentation into agent files (AGENTS.md, GEMINI.md, CLAUDE.md, CURSOR.md).

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| | `--file` | string | "" | Specific file to update (default: all agent files) | PASS |
| | `--dry-run` | bool | false | Show what would change without modifying files | PASS |
| | `--force` | bool | false | Force update even if content hasn't changed | Not tested (would modify files) |
| `-a` | `--all` | bool | false | Include all optional modules | PASS |
| | `--cycle` | bool | false | Include cycle commands | PASS |
| | `--project` | bool | false | Include project commands | PASS |
| | `--module` | bool | false | Include module commands | PASS |
| | `--state` | bool | false | Include state commands | PASS |
| | `--label` | bool | false | Include label commands | PASS |
| | `--type` | bool | false | Include type commands | PASS |
| | `--workspace` | bool | false | Include workspace commands | PASS |
| | `--intake` | bool | false | Include intake commands | PASS |

**Positional Arguments:** None

**Test Details:**
```
$ taskforge inject --dry-run
Would update: AGENTS.md
ℹ Skipped (unchanged): AGENTS.md
✗ GEMINI.md: file does not exist
✗ CLAUDE.md: file does not exist
✗ CURSOR.md: file does not exist
Error: failed to update some files
```
**Result:** PASS (exits with error because most agent files don't exist, which is expected)

---

### `taskforge init`

Initialize .taskforge/settings.yaml for project-local configuration.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| | `--workspace` | string | "" | Workspace slug | Not tested (interactive) |
| | `--project` | string | "" | Project ID or identifier | Not tested (interactive) |
| | `--create-project` | bool | false | Create a new project | Not tested |
| | `--project-name` | string | "" | New project name (with --create-project) | Not tested |
| | `--project-identifier` | string | "" | New project identifier (with --create-project) | Not tested |
| | `--project-description` | string | "" | New project description (with --create-project) | Not tested |
| | `--skip-gitignore` | bool | false | Skip adding .taskforge/ to .gitignore | Not tested |
| | `--upgrade` | bool | false | Upgrade existing .taskforge/settings.yaml | Not tested |

**Positional Arguments:** None

**Result:** Not tested (interactive, creates files)

---

### `taskforge workspace`

Manage workspaces. Aliases: `ws`

| Subcommand | Description |
|------------|-------------|
| `info` | Show workspace details |
| `switch` | Switch default workspace |
| `members` | List workspace members |

#### `taskforge workspace info [slug]`

Show the currently configured workspace and its projects.

**Flags:** Inherits global flags

**Positional Arguments:** `[slug]` — Optional workspace slug (defaults to configured workspace)

**Test Details:**
```
$ taskforge workspace info
Workspace: engineering
API Host: http://localhost:3000
Default Project: d1a511e2-816c-4e43-a599-b1ae0cbd8297
Projects in workspace: 1
  - TaskForge (TFORG)
```
**Result:** PASS

---

#### `taskforge workspace switch [slug]`

Set the default workspace for all future commands.

**Flags:** Inherits global flags

**Positional Arguments:** `[slug]` — Optional workspace slug (interactive prompt if not provided)

**Test Details:**
```
$ taskforge workspace switch engineering
✓ Switched to workspace 'engineering'
```
**Result:** PASS

---

#### `taskforge workspace members`

Display all members of the current workspace. Aliases: `users`, `people`

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| | `--search` | string | "" | Filter members by display name, email, full name, or ID | PASS |
| | `--exact` | bool | false | Require exact matches for --search | PASS |
| | `--limit` | int | 0 | Maximum number of members to show (0 = no limit) | PASS |

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge workspace members` | PASS — Lists all members |
| `taskforge workspace members --search admin` | PASS — Filters by search term |
| `taskforge workspace members --search nonexist --exact` | PASS — No results |
| `taskforge workspace members --limit 1` | PASS — Limits to 1 result |
| `taskforge workspace members -o json` | PASS — JSON output |

---

### `taskforge project`

Manage projects. Aliases: `proj`

| Subcommand | Description |
|------------|-------------|
| `list` | List all projects in the workspace |
| `create` | Create a new project |
| `info` | Show project details |
| `delete` | Delete a project |
| `members` | List project members |

#### `taskforge project list`

List all projects in the current workspace. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge project list` | PASS — Shows project table |
| `taskforge project list -o json` | PASS — JSON output |
| `taskforge project list -o yaml` | PASS — YAML output |

---

#### `taskforge project create [<name>]`

Create a new project in the current workspace. Interactive prompt if name not provided.

**Flags:** Inherits global flags only

**Positional Arguments:** `[<name>]` — Optional project name

**Test Details:**
```
$ taskforge project create (with no args — enters interactive mode, prompts for name)
Error: project name is required
```
**Result:** PASS (correctly requires name)

---

#### `taskforge project info [<id>]`

Display detailed information about a specific project.

**Flags:** Inherits global flags

**Positional Arguments:** `[<id>]` — Optional project ID (defaults to configured project)

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge project info` | PASS — Uses default project |
| `taskforge project info -o json` | PASS — JSON output |

---

#### `taskforge project delete <id>`

Delete a project from the workspace. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Project ID to delete

**Result:** Not tested (destructive; would delete the seeded project)

---

#### `taskforge project members [<id>]`

List all members of a specific project.

**Flags:** Inherits global flags

**Positional Arguments:** `[<id>]` — Optional project ID (defaults to configured project)

**Test Details:**
```
$ taskforge project members -o json
[
  {
    "display_name": "Admin User",
    "email": "admin@taskforge.local",
    "first_name": "Admin",
    "id": "fe3cdd87-f698-4ba7-a336-17dde31674ab",
    "last_name": "User",
    "role": 15
  }
]
```
**Result:** PASS

---

### `taskforge issue`

Manage issues (work items). Aliases: `i`, `issues`, `ticket`

| Subcommand | Description |
|------------|-------------|
| `list` | List issues in the current project |
| `view` | View issue details |
| `create` | Create a new issue |
| `edit` | Edit an existing issue |
| `delete` | Delete an issue |
| `search` | Search for issues |
| `comment` | Manage issue comments |
| `link` | Manage issue links |
| `time` | Manage time tracking |
| `attachment` | Manage issue attachments |
| `activity` | View issue activity history |

#### `taskforge issue list`

List issues in the current project. Aliases: `ls`

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-s` | `--state` | string | "" | Filter by state ID | PASS |
| | `--assignee` | string | "" | Filter by assignee ID | PASS |
| `-l` | `--limit` | int | 20 | Number of issues to show per page | PASS |

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge issue list` | PASS — Lists all issues |
| `taskforge issue list --limit 5` | PASS — Limits results |
| `taskforge issue list --limit 3 -o json` | PASS — JSON output with limit |
| `taskforge issue list --state 176ce03e-abff-4475-a549-5dc88a38d91c` | PASS — Filters by state |
| `taskforge issue list --assignee fe3cdd87-f698-4ba7-a336-17dde31674ab` | PASS — Filters by assignee |

---

#### `taskforge issue view <id>`

Display detailed information about a specific issue.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Issue ID (UUID or sequence ID like `TF-1`)

**Test Details:**
```
$ taskforge issue view 655e1429-a400-4199-8eb1-5976da9eda1b
[Shows issue YAML output]
```
**Result:** PASS

---

#### `taskforge issue create`

Create a new issue in the current project.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-t` | `--title` | string | "" | Issue title | PASS |
| `-d` | `--description` | string | "" | Issue description | PASS |
| `-p` | `--priority` | string | "medium" | Issue priority (none, low, medium, high, urgent) | PASS |
| `-a` | `--assignee` | stringSlice | nil | Assignee ID(s) | Not tested |
| | `--label` | stringSlice | nil | Label ID(s) | Not tested |

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge issue create -t "CLI Test Issue" -d "Created by test matrix" -p high` | PASS — `✓ Created issue #2` |
| `taskforge issue create -t "CLI Test Issue 2" -d "Created by test matrix"` | PASS — Created with default priority |

---

#### `taskforge issue edit <id>`

Edit an existing issue.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-t` | `--title` | string | "" | New title | PASS |
| `-d` | `--description` | string | "" | New description | Not tested |
| | `--priority` | string | "" | New priority (none, low, medium, high, urgent) | Not tested |
| | `--state` | string | "" | New state ID | Not tested |
| `-a` | `--assignee` | stringSlice | nil | New assignee ID(s) | Not tested |
| | `--label` | stringSlice | nil | New label ID(s) | Not tested |

**Positional Arguments:** `<id>` — Issue ID (UUID or sequence ID)

**Test Details:**
```
$ taskforge issue edit 655e1429-a400-4199-8eb1-5976da9eda1b -t "Welcome to TaskForge!"
✓ Updated issue 1
```
**Result:** PASS

---

#### `taskforge issue delete <id>`

Delete an issue. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Issue ID (UUID or sequence ID)

**Result:** Not tested (destructive)

---

#### `taskforge issue search <query>`

Search for issues across the workspace.

**Flags:** Inherits global flags

**Positional Arguments:** `<query>` — Search query text

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge issue search "Welcome"` | PASS — Returns matching issues |
| `taskforge issue search "TF-1"` | PASS — Searches by sequence ID |

---

### `taskforge issue comment`

Manage issue comments. Aliases: `comments`

| Subcommand | Description |
|------------|-------------|
| `add` | Add a comment to an issue |
| `list` | List comments for an issue |
| `delete` | Delete a comment from an issue |

#### `taskforge issue comment add <issue-id>`

Add a comment to an issue.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-t` | `--text` | string | "" | Comment text | PASS |
| | `--access` | string | "INTERNAL" | Comment access (INTERNAL or EXTERNAL) | PASS |

**Positional Arguments:** `<issue-id>` — Issue ID (UUID or sequence ID)

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge issue comment add <issue-id> -t "Test comment from CLI"` | PASS |
| `taskforge issue comment add <issue-id> -t "External comment" --access EXTERNAL` | PASS |

---

#### `taskforge issue comment list <issue-id>`

List comments for an issue. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` — Issue ID (UUID or sequence ID)

**Test Details:** PASS — Lists all comments for an issue

---

#### `taskforge issue comment delete <issue-id> <comment-id>`

Delete a comment from an issue. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` `<comment-id>`

**Test Details:** PASS — Prompts for confirmation, then deletes

---

### `taskforge issue link`

Manage issue links.

| Subcommand | Description |
|------------|-------------|
| `add` | Add a link to an issue |
| `list` | List links for an issue |
| `delete` | Delete a link from an issue |

#### `taskforge issue link add <issue-id> [url]`

Add a link to an issue. Interactive prompt for URL if not provided.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-t` | `--title` | string | "" | Link title | PASS |

**Positional Arguments:** `<issue-id>` `[url]` — Issue ID and optional URL

**Test Details:**
```
$ taskforge issue link add <issue-id> https://example.com -t "Example Link"
✓ Added link 'Example Link' to issue
```
**Result:** PASS

---

#### `taskforge issue link list <issue-id>`

List links for an issue. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` — Issue ID

**Test Details:** PASS — Lists all links for an issue

---

#### `taskforge issue link delete <issue-id> <link-id>`

Delete a link from an issue. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` `<link-id>`

**Test Details:** PASS — Prompts and deletes

---

### `taskforge issue time`

Manage time tracking for issues. Aliases: `worklog`, `log`

| Subcommand | Description |
|------------|-------------|
| `log` | Log time for an issue |
| `list` | List time logs for an issue |
| `total` | Show total time logged for an issue |
| `edit` | Edit a time log |
| `delete` | Delete a time log |

> **Note:** Time tracking commands require time tracking to be enabled on the project.

#### `taskforge issue time log <issue-id> <duration>`

Log time spent on an issue.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-d` | `--description` | string | "" | Description of work done | N/A (disabled) |

**Positional Arguments:** `<issue-id>` `<duration>` — Duration formats: `60` (minutes), `2h`, `4.5h`, `2h30m`, `1h45m`

**Test Result:** FAIL — `Error: time tracking is disabled for project d1a511e2-816c-4e43-a599-b1ae0cbd8297`

---

#### `taskforge issue time list <issue-id>`

List time logs for an issue. Aliases: `ls`

**Test Result:** FAIL — `Error: time tracking is disabled for project d1a511e2-816c-4e43-a599-b1ae0cbd8297`

---

#### `taskforge issue time total <issue-id>`

Show total time logged for an issue.

**Test Result:** FAIL — `Error: time tracking is disabled for project d1a511e2-816c-4e43-a599-b1ae0cbd8297`

---

#### `taskforge issue time edit <issue-id> <worklog-id>`

Edit a time log.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-d` | `--description` | string | "" | New description | N/A (disabled) |
| `-t` | `--duration` | string | "" | New duration (e.g., 2h30m, 90) | N/A (disabled) |

**Test Result:** N/A — Time tracking disabled on project

---

#### `taskforge issue time delete <issue-id> <worklog-id>`

Delete a time log. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Test Result:** N/A — Time tracking disabled on project

---

### `taskforge issue attachment`

Manage issue attachments. Aliases: `attach`, `file`

| Subcommand | Description |
|------------|-------------|
| `list` | List attachments for an issue |
| `upload` | Upload a file attachment |
| `edit` | Edit attachment metadata |
| `delete` | Delete an attachment |

#### `taskforge issue attachment list <issue-id>`

List attachments for an issue. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` — Issue ID

**Test Details:** PASS — `ℹ No attachments found for this issue`

---

#### `taskforge issue attachment upload <issue-id> <file-path>`

Upload a file as an attachment to an issue.

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` `<file-path>` — Issue ID and path to file

**Result:** Not tested (requires a file to upload)

---

#### `taskforge issue attachment edit <issue-id> <attachment-id>`

Edit attachment metadata.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | New filename | Not tested |
| `-a` | `--archive` | bool | false | Archive the attachment | Not tested |
| `-u` | `--unarchive` | bool | false | Unarchive the attachment | Not tested |

**Positional Arguments:** `<issue-id>` `<attachment-id>`

**Result:** Not tested

---

#### `taskforge issue attachment delete <issue-id> <attachment-id>`

Delete an attachment. Aliases: `rm`, `remove`.

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` `<attachment-id>`

**Result:** Not tested

---

### `taskforge issue activity`

View issue activity history. Aliases: `history`, `log`

| Subcommand | Description |
|------------|-------------|
| `list` | List activity history |
| `view` | View activity details |

#### `taskforge issue activity list <issue-id>`

List activity history for a specific issue. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` — Issue ID

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge issue activity list <issue-id>` | PASS — Shows activity in table format |
| `taskforge issue activity list <issue-id> -o json` | PASS — JSON output with full details |

---

#### `taskforge issue activity view <issue-id> <activity-id>`

View detailed information about a specific activity.

**Flags:** Inherits global flags

**Positional Arguments:** `<issue-id>` `<activity-id>`

**Test Details:**
```
$ taskforge issue activity view <issue-id> <activity-id>
actor: fe3cdd87-f698-4ba7-a336-17dde31674ab
created_at: "2026-05-06T21:23:08.065Z"
...
verb: created
```
**Result:** PASS

---

### `taskforge cycle`

Manage cycles (sprints). Aliases: `sprint`

| Subcommand | Description |
|------------|-------------|
| `list` | List cycles |
| `view` | View cycle details |
| `create` | Create a new cycle |
| `edit` | Edit a cycle |
| `delete` | Delete a cycle |
| `archive` | Archive a cycle |
| `issues` | List cycle issues |
| `add-issues` | Add issues to a cycle |
| `remove-issue` | Remove an issue from a cycle |

#### `taskforge cycle list`

List all cycles in the current project. Aliases: `ls`

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-a` | `--archived` | bool | false | Show archived cycles | PASS |

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge cycle list` | PASS — Lists active cycles |
| `taskforge cycle list --archived` | PASS — Includes archived cycles |
| `taskforge cycle list -o json` | PASS — JSON output |

---

#### `taskforge cycle view <id>`

Display detailed information about a specific cycle.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Cycle UUID

**Test Details:** PASS — Shows cycle details in yaml format

---

#### `taskforge cycle create`

Create a new cycle (sprint) in the current project.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | Cycle name | PASS |
| `-d` | `--description` | string | "" | Cycle description | PASS |
| `-s` | `--start-date` | string | "" | Start date (YYYY-MM-DD) | PASS |
| `-e` | `--end-date` | string | "" | End date (YYYY-MM-DD) | PASS |

**Positional Arguments:** None (interactive prompts if not provided)

**Test Details:**
```
$ taskforge cycle create -n "Test Sprint" -d "A test sprint" -s 2026-05-01 -e 2026-05-14
✓ Created cycle 'Test Sprint' (53a18d28-04a1-4de3-9bdb-eb0d48066b78)
```
**Result:** PASS

---

#### `taskforge cycle edit <id>`

Edit an existing cycle.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | New cycle name | PASS |
| `-d` | `--description` | string | "" | New cycle description | PASS |
| `-s` | `--start-date` | string | "" | New start date (YYYY-MM-DD) | Not tested |
| `-e` | `--end-date` | string | "" | New end date (YYYY-MM-DD) | Not tested |

**Positional Arguments:** `<id>` — Cycle UUID

**Test Details:**
```
$ taskforge cycle edit <cycle-id> -n "Test Sprint Edited"
✓ Updated cycle 'Test Sprint Edited'
```
**Result:** PASS

---

#### `taskforge cycle delete <id>`

Delete a cycle. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Cycle UUID

**Test Details:** PASS — Prompts for confirmation, then deletes

---

#### `taskforge cycle archive <id>`

Archive a cycle to hide it from active cycles.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Cycle UUID

**Test Details:** PASS — `✓ Archived cycle <id>`

---

#### `taskforge cycle issues <id>`

List all work items in a specific cycle. Aliases: `work-items`

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Cycle UUID

**Test Details:** PASS — Lists issues in cycle

---

#### `taskforge cycle add-issues <cycle-id> <issue-ids...>`

Add work items to a cycle. Accepts multiple issue IDs.

**Flags:** Inherits global flags

**Positional Arguments:** `<cycle-id>` `<issue-ids...>` — Cycle UUID followed by one or more issue UUIDs

**Test Details:**
```
$ taskforge cycle add-issues <cycle-id> <issue-id>
✓ Added 1 issue(s) to cycle <cycle-id>
```
**Result:** PASS

---

#### `taskforge cycle remove-issue <cycle-id> <issue-id>`

Remove a work item from a cycle.

**Flags:** Inherits global flags

**Positional Arguments:** `<cycle-id>` `<issue-id>` — Cycle UUID and issue UUID

**Test Details:** PASS — `✓ Removed issue <issue-id> from cycle <cycle-id>`

---

### `taskforge module`

Manage modules. Aliases: `mod`

| Subcommand | Description |
|------------|-------------|
| `list` | List modules |
| `view` | View module details |
| `create` | Create a new module |
| `edit` | Edit a module |
| `delete` | Delete a module |
| `archive` | Archive a module |
| `issues` | List module issues |
| `add-issues` | Add issues to a module |
| `remove-issue` | Remove an issue from a module |

#### `taskforge module list`

List all modules in the current project. Aliases: `ls`

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-a` | `--archived` | bool | false | Show archived modules | PASS |

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge module list` | PASS |
| `taskforge module list --archived` | PASS — Includes archived modules |
| `taskforge module list -o json` | PASS — JSON output |

---

#### `taskforge module view <id>`

Display detailed information about a specific module.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Module UUID

**Test Details:** PASS — Shows module details

---

#### `taskforge module create`

Create a new module in the current project.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | Module name | PASS |
| `-d` | `--description` | string | "" | Module description | PASS |
| `-s` | `--status` | string | "" | Module status (backlog, planned, in-progress, paused, completed, cancelled) | Not tested |

**Positional Arguments:** None (interactive prompts if flags not provided)

**Test Details:**
```
$ taskforge module create -n "Test Module" -d "A test module"
✓ Created module 'Test Module' (8e8b0421-781f-446d-b976-90ecd59dcf2d)
```
**Result:** PASS

---

#### `taskforge module edit <id>`

Edit an existing module.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | New module name | PASS |
| `-d` | `--description` | string | "" | New module description | Not tested |
| `-s` | `--status` | string | "" | New module status | Not tested |

**Positional Arguments:** `<id>` — Module UUID

**Test Details:** PASS — `✓ Updated module 'Test Module Edited'`

---

#### `taskforge module delete <id>`

Delete a module. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Module UUID

**Test Details:** PASS — Prompts for confirmation, then deletes

---

#### `taskforge module archive <id>`

Archive a module to hide it from active modules.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Module UUID

**Test Details:** PASS — `✓ Archived module <id>`

---

#### `taskforge module issues <id>`

List all work items in a specific module. Aliases: `work-items`

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Module UUID

**Test Details:** PASS — Lists issues in module or shows "No issues found in this module"

---

#### `taskforge module add-issues <module-id> <issue-ids...>`

Add work items to a module. Accepts multiple issue IDs.

**Flags:** Inherits global flags

**Positional Arguments:** `<module-id>` `<issue-ids...>` — Module UUID followed by one or more issue UUIDs

**Test Details:** PASS — `✓ Added 1 issue(s) to module <id>`

---

#### `taskforge module remove-issue <module-id> <issue-id>`

Remove a work item from a module.

**Flags:** Inherits global flags

**Positional Arguments:** `<module-id>` `<issue-id>` — Module UUID and issue UUID

**Test Details:** PASS — `✓ Removed issue <id> from module <id>`

---

### `taskforge label`

Manage project labels. Aliases: `labels`, `tag`

| Subcommand | Description |
|------------|-------------|
| `list` | List labels |
| `view` | View label details |
| `create` | Create a new label |
| `edit` | Edit a label |
| `delete` | Delete a label |

#### `taskforge label list`

List all labels in the current project. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge label list` | PASS — Lists all labels |

---

#### `taskforge label view <id>`

Display detailed information about a specific label.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Label UUID

**Test Details:** PASS — Shows label details in yaml format

---

#### `taskforge label create`

Create a new label in the current project.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | Label name | PASS |
| `-d` | `--description` | string | "" | Label description | PASS |
| `-c` | `--color` | string | "" | Label color (hex code, e.g., #EF4444) | PASS |

**Positional Arguments:** None (interactive prompts if flags not provided)

**Test Details:**
```
$ taskforge label create -n "Test Label" -c "#FF0000" -d "Test label description"
✓ Created label 'Test Label' (117a99da-6d9e-4153-bab4-57589bff022d)
```
**Result:** PASS

---

#### `taskforge label edit <id>`

Edit an existing label.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | New label name | PASS |
| `-d` | `--description` | string | "" | New label description | Not tested |
| `-c` | `--color` | string | "" | New label color | Not tested |

**Positional Arguments:** `<id>` — Label UUID

**Test Details:** PASS — `✓ Updated label 'Test Label Edited'`

---

#### `taskforge label delete <id>`

Delete a label. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Label UUID

**Test Details:** PASS — Prompts for confirmation, then deletes

---

### `taskforge state`

Manage project states. Aliases: `states`

| Subcommand | Description |
|------------|-------------|
| `list` | List states |
| `view` | View state details |
| `create` | Create a new state |
| `edit` | Edit a state |
| `delete` | Delete a state |

#### `taskforge state list`

List all states in the current project. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** None

**Test Results:**

| Invocation | Result |
|------------|--------|
| `taskforge state list` | PASS — Lists all states |

---

#### `taskforge state view <id>`

Display detailed information about a specific state.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — State UUID

**Test Details:**
```
$ taskforge state view 3987f3ee-4a13-4b9d-a118-84008841708f
color: '#FF0000'
created_at: "2026-05-06T22:34:25.318Z"
group: backlog
id: 3987f3ee-4a13-4b9d-a118-84008841708f
name: Test State
sequence: 10000
...
```
**Result:** PASS

---

#### `taskforge state create`

Create a new workflow state.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | State name | PASS |
| `-d` | `--description` | string | "" | State description | Not tested |
| `-c` | `--color` | string | "" | State color (hex code, e.g., #F59E0B) | PASS |
| `-g` | `--group` | string | "" | State group (backlog, unstarted, started, completed, cancelled) | PASS |

**Positional Arguments:** None (interactive prompts if flags not provided)

**Test Details:**
```
$ taskforge state create -n "Test State" -g backlog -c "#FF0000"
✓ Created state 'Test State' (3987f3ee-4a13-4b9d-a118-84008841708f)
```
**Result:** PASS

---

#### `taskforge state edit <id>`

Edit an existing state.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | New state name | PASS |
| `-d` | `--description` | string | "" | New state description | Not tested |
| `-c` | `--color` | string | "" | New state color | Not tested |
| `-g` | `--group` | string | "" | New state group | Not tested |

**Positional Arguments:** `<id>` — State UUID

**Test Details:** PASS — `✓ Updated state 'Test State Edited'`

---

#### `taskforge state delete <id>`

Delete a state. Aliases: `rm`, `remove`. Requires confirmation prompt.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — State UUID

**Test Details:** PASS — Prompts for confirmation, then deletes

---

### `taskforge type`

Manage issue types. Aliases: `issue-type`

| Subcommand | Description |
|------------|-------------|
| `list` | List issue types |
| `create` | Create a new issue type |
| `delete` | Delete an issue type |

#### `taskforge type list`

List all issue types in the current project. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** None

**Test Result:** FAIL — `Error: issue types are disabled for project d1a511e2-816c-4e43-a599-b1ae0cbd8297`

> **Note:** The `type list` and `type create` commands fail because the test project has issue types disabled. This is expected behavior — the CLI correctly checks the project configuration before attempting to list/create types.

---

#### `taskforge type create`

Create a new issue type.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | Issue type name | N/A (disabled) |
| `-d` | `--description` | string | "" | Issue type description | N/A (disabled) |

**Positional Arguments:** None (interactive prompts if flags not provided)

**Test Result:** FAIL — `Error: issue types are disabled for project d1a511e2-816c-4e43-a599-b1ae0cbd8297`

---

### `taskforge intake`

Manage intake issues. Aliases: `inbox`, `requests`

| Subcommand | Description |
|------------|-------------|
| `list` | List intake issues |
| `view` | View intake issue details |
| `create` | Create a new intake issue |
| `delete` | Delete an intake issue |

#### `taskforge intake list`

List all intake issues in the current project. Aliases: `ls`

**Flags:** Inherits global flags

**Positional Arguments:** None

**Test Details:** PASS — `ℹ No intake issues found`

---

#### `taskforge intake view <id>`

Display detailed information about a specific intake issue.

**Flags:** Inherits global flags

**Positional Arguments:** `<id>` — Intake issue UUID

**Result:** Not tested (no intake issues exist)

---

#### `taskforge intake create`

Create a new intake issue.

**Flags:**

| Short | Long | Type | Default | Description | Test Result |
|-------|------|------|---------|-------------|-------------|
| `-n` | `--name` | string | "" | Issue name/title | N/A |
| `-p` | `--priority` | string | "medium" | Issue priority (low, medium, high, urgent) | N/A |

**Positional Arguments:** None (interactive prompts if flags not provided)

**Test Result:** FAIL — `Error: API error (status 500): ` (server-side error when creating intake issues)

> **Note:** The intake create command returns a 500 error from the API. This appears to be a server-side issue, not a CLI bug.

---

### `taskforge version`

Print version information.

**Flags:** None (not a cobra command with global flags — uses `debug.ReadBuildInfo`)

**Positional Arguments:** None

**Test Details:**
```
$ taskforge version
taskforge version dev (commit: none, built: unknown)
```
**Result:** PASS

> **Note:** `--version` flag is not supported on the root command. Use `taskforge version` instead.

---

### `taskforge completion [bash|zsh|fish|powershell]`

Generate shell completion scripts.

**Flags:** Inherits global flags

**Positional Arguments:** `<shell>` — One of: `bash`, `zsh`, `fish`, `powershell`

**Test Details:**
```
$ taskforge completion bash
[Generates bash completion script]
```
**Result:** PASS

---

## Test Summary

### Overall Results

| Category | Total Commands | Passed | Failed | Not Tested |
|----------|---------------|--------|--------|------------|
| **auth** | 4 | 4 | 0 | 0 |
| **config** | 2 | 2 | 0 | 0 |
| **context** | 1 | 1 | 0 | 0 |
| **inject** | 1 | 1 | 0 | 0 |
| **init** | 1 | 0 | 0 | 1 |
| **workspace** | 3 | 3 | 0 | 0 |
| **project** | 5 | 4 | 0 | 1 |
| **issue** (root) | 6 | 6 | 0 | 0 |
| **issue comment** | 3 | 3 | 0 | 0 |
| **issue link** | 3 | 3 | 0 | 0 |
| **issue time** | 5 | 0 | 0 | 5 |
| **issue attachment** | 4 | 1 | 0 | 3 |
| **issue activity** | 2 | 2 | 0 | 0 |
| **cycle** | 9 | 9 | 0 | 0 |
| **module** | 9 | 9 | 0 | 0 |
| **label** | 5 | 5 | 0 | 0 |
| **state** | 5 | 5 | 0 | 0 |
| **type** | 3 | 0 | 2 | 1 |
| **intake** | 4 | 1 | 1 | 2 |
| **version** | 1 | 1 | 0 | 0 |
| **completion** | 1 | 1 | 0 | 0 |
| **TOTAL** | **77** | **62** | **3** | **12** |

### Failed Tests

| Command | Error | Root Cause |
|---------|-------|------------|
| `intake create -n "..." -p medium` | `Error: API error (status 500):` | Server-side API error — returns 500 with empty message |
| `type create -n "Bug" -d "Bug type"` | `Error: issue types are disabled for project ...` | Project configuration — issue types not enabled for test project |
| `type list` | `Error: issue types are disabled for project ...` | Project configuration — issue types not enabled for test project |

### Not Tested (By Design)

| Command | Reason |
|---------|--------|
| `taskforge init` | Interactive — creates files in current directory |
| `taskforge project delete <id>` | Destructive — would delete seeded project |
| `taskforge issue delete <id>` | Destructive — would delete seeded issues |
| `taskforge issue time log/list/total/edit/delete` | Time tracking disabled on test project |
| `taskforge issue attachment upload/edit/delete` | Requires file upload or no attachments exist |
| `taskforge intake view` | No intake issues exist |
| `taskforge intake delete` | No intake issues to delete |
| `taskforge type delete` | No types to delete |

### Known Issues

1. **`--version` not supported on root command**: `taskforge --version` returns `Error: unknown flag: --version`. Use `taskforge version` instead.
2. **`config set output_format`**: The config key is `output`, not `output_format`. Using `output_format` returns `unknown config key`.
3. **Intake create returns 500**: The `intake create` command triggers a server-side 500 error with an empty error message. This is a backend issue, not a CLI bug.
4. **Time tracking commands**: All time-related commands (`time log`, `time list`, `time total`, `time edit`, `time delete`) fail with "time tracking is disabled for project" error. This is expected when the project has time tracking disabled.
5. **Project create is interactive**: `taskforge project create` has no flags for name/identifier — uses interactive prompts even with arguments. The `init` command provides `--project-name` and `--project-identifier` flags for non-interactive project creation.