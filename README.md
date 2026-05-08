# TaskForge CLI

A powerful command-line interface for [TaskForge](https://github.com/rohithmahesh3/taskforge) - the open-source project management tool.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- 🔐 **Secure Authentication** - API key stored in OS keyring
- 📊 **Consistent YAML Output** - Structured YAML output for scripting and readability
- 🎯 **Interactive Mode** - Prompts for missing required fields
- 🔍 **Work Item Filtering** - Filter work items by state and assignee
- ⚡ **Fast & Lightweight** - Single binary, no dependencies
- 🎨 **Beautiful Output** - Colored and formatted tables
- 🔧 **Shell Completion** - Bash, Zsh, Fish, and PowerShell support
- ⏱️ **Time Tracking** - Log and manage time spent on issues
- 💬 **Comments & Activity** - Full comment thread and activity history support
- 📎 **Attachments** - Upload and manage file attachments
- 🔗 **Issue Links** - Add external links to issues
- 🤖 **AI Context Generation** - Generate CLI reference for AI agents

## Installation

### From Source

```bash
go install github.com/rohithmahesh3/taskforge@latest
```

### Download Binary

Download the latest release from the [releases page](https://github.com/rohithmahesh3/taskforge/releases).

### Homebrew (macOS/Linux)

```bash
brew tap rohithmahesh3/taskforge
brew install taskforge
```

## Quick Start

### 1. Authenticate

```bash
taskforge auth login
```

You'll need an API key from TaskForge:
1. Go to your TaskForge workspace
2. Navigate to Profile Settings → Personal Access Tokens
3. Create a new token
4. Use it when prompted

### 2. Initialize Project

```bash
taskforge init
```

This creates a `.taskforge/settings.yaml` file in your current directory with workspace and project settings.

### 3. Inject Context into Agent Files (Optional)

```bash
taskforge inject
```

This injects CLI command documentation into agent files (AGENTS.md, GEMINI.md, CLAUDE.md, CURSOR.md) for AI assistants.

### 4. Check Workspace Access

```bash
taskforge workspace info
```

### 5. List Projects

```bash
taskforge project list
```

### 6. List Issues

```bash
taskforge issue list
```

## Configuration

Configuration is stored in `~/.config/taskforge/config.yaml`:

```yaml
version: "1.0"
default_workspace: my-workspace
default_project: my-project-id
api_host: https://api.taskforge.app
```

### Environment Variables

- `TASKFORGE_API_KEY` - Your TaskForge API key
- `TASKFORGE_WORKSPACE` - Default workspace slug
- `TASKFORGE_PROJECT` - Default project ID

## Usage Examples

### Authentication

```bash
# Interactive login
taskforge auth login

# Login with flags
taskforge auth login --token YOUR_API_KEY --workspace my-workspace

# Check authentication status
taskforge auth status

# Logout
taskforge auth logout
```

### Project Initialization

```bash
# Initialize project settings interactively
taskforge init

# Initialize with specific workspace and project
taskforge init --workspace my-workspace --project PROJECT_ID

# Create a new project during init
taskforge init --workspace my-workspace --create-project --project-name "My New Project"

# Update existing settings
taskforge init --upgrade
```

### Context Injection

```bash
# Inject context into all agent files
taskforge inject

# Inject into specific file
taskforge inject --file AGENTS.md

# Preview changes without modifying files
taskforge inject --dry-run

# Force update even if unchanged
taskforge inject --force
```

### Workspaces

```bash
# Show workspace details
taskforge workspace info

# Switch default workspace
taskforge workspace switch my-workspace

# List workspace members
taskforge workspace members

# Search workspace members before assigning issues
taskforge workspace members --search roh
```

### Projects

```bash
# List projects
taskforge project list

# Create a new project
taskforge project create

# View project details
taskforge project info PROJECT_ID

# Delete a project
taskforge project delete PROJECT_ID

# List project members
taskforge project members PROJECT_ID
```

### Issues (Work Items)

```bash
# List issues
taskforge issue list

# List with filters
taskforge issue list --state <state-id>
taskforge issue list --assignee <assignee-id>

# View issue details (supports sequence ID or UUID)
taskforge issue view 123
taskforge issue view uuid-here

# Create an issue
taskforge issue create --title "Bug fix" --priority high
taskforge issue create -t "Feature request" -d "Description" -p medium -a <assignee-id>

# Edit an issue
taskforge issue edit 123 --state <state-id>
taskforge issue edit 123 --priority urgent --assignee <assignee-id>
taskforge issue edit 123 --assignee <assignee-id> --label bug,urgent

# Delete an issue
taskforge issue delete 123

# Search issues across workspace
taskforge issue search "login bug"
```

### Issue Comments

```bash
# List comments on an issue
taskforge issue comment list 123

# Add a comment
taskforge issue comment add 123 --text "Fixed in PR #42"
taskforge issue comment add 123 --text "Customer feedback" --access EXTERNAL

# Delete a comment
taskforge issue comment delete 123 COMMENT_ID
```

### Issue Time Tracking

```bash
# List time logs
taskforge issue time list 123

# Log time (supports multiple formats)
taskforge issue time log 123 2h30m --description "Fixed the bug"
taskforge issue time log 123 90 --description "Code review"

# Show total time logged
taskforge issue time total 123

# Edit a time log
taskforge issue time edit 123 WORKLOG_ID --duration 3h

# Delete a time log
taskforge issue time delete 123 WORKLOG_ID
```

### Issue Links

```bash
# List links on an issue
taskforge issue link list 123

# Add a link
taskforge issue link add 123 https://github.com/repo/pull/42 --title "Related PR"

# Delete a link
taskforge issue link delete 123 LINK_ID
```

### Issue Attachments

```bash
# List attachments
taskforge issue attachment list 123

# Upload a file
taskforge issue attachment upload 123 ./screenshot.png

# Edit attachment metadata
taskforge issue attachment edit 123 ATTACHMENT_ID --name "new-name.png"

# Archive/unarchive
taskforge issue attachment edit 123 ATTACHMENT_ID --archive

# Delete an attachment
taskforge issue attachment delete 123 ATTACHMENT_ID
```

### Issue Activity History

```bash
# List activity history
taskforge issue activity list 123

# View specific activity details
taskforge issue activity view 123 ACTIVITY_ID
```

### Cycles (Sprints)

```bash
# List cycles
taskforge cycle list

# List including archived
taskforge cycle list --archived

# View cycle details
taskforge cycle view CYCLE_ID

# Create a cycle
taskforge cycle create --name "Sprint 1" --start-date 2024-01-01 --end-date 2024-01-14

# Edit a cycle
taskforge cycle edit CYCLE_ID --name "Sprint 1 (Revised)"

# Delete a cycle
taskforge cycle delete CYCLE_ID

# Archive
taskforge cycle archive CYCLE_ID

# List issues in a cycle
taskforge cycle issues CYCLE_ID

# Add issues to a cycle
taskforge cycle add-issues CYCLE_ID ISSUE_ID_1 ISSUE_ID_2

# Remove an issue from a cycle
taskforge cycle remove-issue CYCLE_ID ISSUE_ID
```

### Modules

```bash
# List modules
taskforge module list

# List including archived
taskforge module list --archived

# View module details
taskforge module view MODULE_ID

# Create a module
taskforge module create --name "Authentication" --description "Auth features" --status <status-id>

# Edit a module
taskforge module edit MODULE_ID --status <status-id>

# Delete a module
taskforge module delete MODULE_ID

# Archive
taskforge module archive MODULE_ID

# List issues in a module
taskforge module issues MODULE_ID

# Add issues to a module
taskforge module add-issues MODULE_ID ISSUE_ID_1 ISSUE_ID_2

# Remove an issue from a module
taskforge module remove-issue MODULE_ID ISSUE_ID
```

### Pages (Documentation)

```bash
# View page details
taskforge page view PAGE_ID

# Create a page
taskforge page create --name "API Documentation"
taskforge page create --name "Team Guidelines" --workspace

```

### States (Workflow)

```bash
# List states
taskforge state list

# View state details
taskforge state view STATE_ID

# Create a state
taskforge state create --name "In Review" --color "#F59E0B" --group started

# Edit a state
taskforge state edit STATE_ID --name "Code Review"

# Delete a state
taskforge state delete STATE_ID
```

### Labels

```bash
# List labels
taskforge label list

# View label details
taskforge label view LABEL_ID

# Create a label
taskforge label create --name "Bug" --color "#EF4444" --description "Something is broken"

# Edit a label
taskforge label edit LABEL_ID --color "#3B82F6"

# Delete a label
taskforge label delete LABEL_ID
```

### Intake (Inbox)

```bash
# List intake issues
taskforge intake list

# View intake issue
taskforge intake view INTAKE_ID

# Create an intake issue
taskforge intake create --name "Feature Request" --priority high

# Delete an intake issue
taskforge intake delete INTAKE_ID
```

### Issue Types

```bash
# List issue types
taskforge type list

# Create an issue type
taskforge type create --name "Bug" --description "Bug reports"

# Delete an issue type
taskforge type delete TYPE_ID
```

### AI Context Generation

Generate CLI command reference for AI agents:

```bash
# Default modules (issue, module, page, state, label, intake, type)
taskforge context

# Include additional modules
taskforge context --cycle
taskforge context --project
taskforge context --workspace

# Include all modules
taskforge context --all
```

### Output

```bash
# YAML is the default output format
taskforge project list

# No colors
taskforge issue list --no-color
```

## Shell Completion

### Bash

```bash
source <(taskforge completion bash)
# Add to ~/.bashrc for persistence
```

### Zsh

```bash
source <(taskforge completion zsh)
# Add to ~/.zshrc for persistence
```

### Fish

```bash
taskforge completion fish | source
# Save for persistence:
taskforge completion fish > ~/.config/fish/completions/taskforge.fish
```

### PowerShell

```powershell
taskforge completion powershell | Out-String | Invoke-Expression
```

## API Support

This CLI uses the [TaskForge REST API](https://github.com/rohithmahesh3/taskforge) v1.

Supported features:
- ✅ Workspaces (info, switch, members)
- ✅ Projects (list, create, info, delete, members)
- ✅ Issues/Work Items (list, create, edit, delete, search)
- ✅ Issue Comments (list, add, delete)
- ✅ Issue Links (list, add, delete)
- ✅ Issue Time Tracking (list, log, edit, delete, total)
- ✅ Issue Attachments (list, upload, edit, delete)
- ✅ Issue Activity History (list, view)
- ✅ Cycles (list, create, edit, delete, archive, issues management)
- ✅ Modules (list, create, edit, delete, archive, issues management)
- ✅ Pages/Documentation (create, view)
- ✅ States/Workflow (list, create, edit, delete)
- ✅ Labels (list, create, edit, delete)
- ✅ Intake/Inbox (list, create, view, delete)
- ✅ Issue Types (list, create, delete)
- ✅ AI Context Generation

## Development

### Prerequisites

- Go 1.21 or higher
- Git
- pre-commit (optional but recommended)

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Install Locally

```bash
make install
```

### Setup Pre-commit Hooks

We use pre-commit hooks to ensure code quality. Install pre-commit and the hooks:

```bash
# Install pre-commit (if not already installed)
pip install pre-commit

# Install the git hooks
make setup-hooks
```

The pre-commit hooks will automatically run on every commit and check:
- Code formatting (`go fmt`)
- Static analysis (`go vet`)
- Linting (`golangci-lint`)
- Tests (`go test`)

You can also run all checks manually:

```bash
make check  # Runs fmt, vet, lint, and test
```

## Project Structure

```
taskforge/
├── cmd/                # Command definitions
│   ├── auth/           # Authentication commands
│   ├── config/         # Config commands
│   ├── context/        # AI context generation
│   ├── cycle/          # Cycle/sprint commands
│   ├── intake/         # Intake/inbox commands
│   ├── issue/          # Issue commands + subcommands
│   │   ├── activity.go # Activity history
│   │   ├── attachment.go # File attachments
│   │   ├── comment.go  # Comments
│   │   ├── link.go     # External links
│   │   └── time.go     # Time tracking
│   ├── label/          # Label commands
│   ├── module/         # Module commands
│   ├── page/           # Page/documentation commands
│   ├── project/        # Project commands
│   ├── state/          # State/workflow commands
│   ├── type/           # Issue type commands
│   └── workspace/      # Workspace commands
├── internal/
│   ├── api/            # API client and endpoints
│   ├── config/         # Configuration management
│   └── output/         # Output formatting (json/yaml)
├── pkg/taskforge/          # TaskForge API types and models
├── main.go
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [TaskForge](https://github.com/rohithmahesh3/taskforge) - The open-source project management tool
- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Survey](https://github.com/AlecAivazis/survey) - Interactive prompts

## Support

- 🐛 [Report bugs](../../issues)
- 💡 [Request features](../../issues)
- 💬 [Discussions](../../discussions)

---

Made with ❤️ for the TaskForge community
