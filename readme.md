## Development Setup

### Prerequisites

- Go 1.18 or higher
- Python 3.6 or higher
- Git

### Setting Up Your Development Environment

1. Clone the repository:
   git clone https://github.com/your-org/your-repo.git
   cd your-repo
   Copy
2. Run the setup script:
   ./setup.sh
   Copy
   This script will:

- Install pre-commit framework
- Install required Go tools (goimports, golangci-lint)
- Configure pre-commit hooks

### Workflow

- The pre-commit hooks will automatically run on each commit
- They will format your code according to Go standards
- Failed checks will prevent commits until fixed
- Our CI pipeline provides a backup check for any formatting issues

If you need to run formatting manually:

- `go fmt ./...` (basic formatting)
- `goimports -w .` (formatting + import organization)
