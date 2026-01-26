# Stacks CLI - GitHub Copilot Instructions

## Project Overview

The Stacks CLI is a Go-based project scaffolding tool that generates complete application templates with infrastructure, pipelines, and supporting files. It supports multiple frameworks (dotnet, java, nx, infra, data) and cloud platforms (Azure, AWS) with sophisticated Go templating.

**Key Facts for AI Agents:**

- This is NOT a simple CLI - it's a sophisticated templating and scaffolding system
- Uses **eirctl** (not taskctl) for builds - see [eirctl.yaml](../eirctl.yaml) for pipeline definitions
- Templates are downloaded from external sources (Git/NuGet/Filesystem) and processed with Go templates
- Integration tests are the primary validation method
- Go version: **1.25.5** (see [go.mod](../go.mod))

## Key Architecture Patterns

### Command Structure

- **Entry point**: [stacks-cli.go](../stacks-cli.go) → `cmd/` package with Cobra commands
- **Core commands**: `scaffold` (main), `interactive`, `export`, `setup` (with subcommands: `update`, `list`, `latest`), `version`
- **Configuration binding**: Uses Viper with flag binding pattern in [cmd/scaffold.go](../cmd/scaffold.go):
  ```go
  viper.BindPFlag("input.project.name", scaffoldCmd.Flags().Lookup("name"))
  ```

### Package Organization

```
├── cmd/                    # Cobra CLI commands and flag definitions
│   ├── root.go            # Root command, global flags, preRun configuration
│   ├── scaffold.go        # Main scaffolding command
│   ├── interactive.go     # Interactive configuration wizard
│   ├── export.go          # Export internal files
│   ├── setup.go           # Setup commands (update, list, latest)
│   └── version.go         # Version display
├── pkg/                    # Public API packages
│   ├── config/            # Configuration structs and parsing (40+ files)
│   ├── scaffold/          # Core scaffolding engine
│   ├── downloaders/       # Git, NuGet, Filesystem downloaders
│   ├── interactive/       # Interactive prompts
│   ├── export/            # File export functionality
│   ├── setup/             # Setup operations
│   ├── filter/            # Struct filtering utilities
│   └── interfaces/        # Shared interfaces
├── internal/               # Private packages
│   ├── config/staticFiles/ # Embedded config, templates, help
│   ├── constants/         # Application constants
│   ├── models/            # Core models (App, Command, Log, etc.)
│   ├── interfaces/        # Internal interfaces
│   └── util/              # Utility functions
└── testing/integration/    # Integration test suite
```

**Configuration flows**: CLI flags → Viper → `config.Config` struct → template rendering

### Scaffolding Engine (`pkg/scaffold/`)

Core scaffolding logic in [scaffold.go](../pkg/scaffold/scaffold.go) processes projects through these phases:

1. **Validation**: Check configuration, validate frameworks, ensure required binaries exist
2. **Download**: Templates fetched via `pkg/downloaders/`:
   - `git.go` - Git repositories (GitHub, etc.)
   - `nuget.go` - NuGet packages (.NET templates)
   - `filesystem.go` - Local filesystem paths
3. **Phases**: Projects have two phases defined in `stackscli.yml`:
   - `init` - Operations run in the temp/clone directory
   - `setup` - Operations run in the working directory
4. **Operations**: Two action types:
   - `copy` - Direct file copying from template to target directory
   - `cmd` - Execute framework commands with Go template expansion

**Critical**: All user inputs and file paths go through Go template processing - ensure proper escaping!

### Configuration System

- **Hierarchical config**: CLI flags → config files → defaults
- **Template variables**: Available via `config.Replacements` struct:
  - `.Input.*` - All CLI inputs (Business, Cloud, Network, etc.)
  - `.Project.*` - Project-specific data (Name, Framework, SourceControl, etc.)
- **Framework definitions**: Defined in [internal/config/staticFiles/config.yml](../internal/config/staticFiles/config.yml)
- **Supported frameworks**: `dotnet`, `java`, `nx`, `infra`, `data`
- **Project settings**: `stackscli.yml` files define per-project operations and pipeline configs

### Key Configuration Types (`pkg/config/`)

```go
// Main configuration struct
type Config struct {
    Commands      Commands       // Git commands for project init
    FrameworkDefs []FrameworkDef // Framework binary requirements
    Input         InputConfig    // All user inputs
    Internal      Internal       // Internal configuration
    Help          Help           // Help URLs
    Stacks        Stacks         // Component definitions
}

// Template rendering context
type Replacements struct {
    Input   InputConfig
    Project Project
}
```

## Key Patterns

### Logging

- Uses logrus for structured logging,
- Generate methods using the `App.Help` singleton,

```go
msg := App.Help.GetMessage("INT001", overrideConfig)
App.Logger.Infof(msg)
```

## Critical Development Workflows

### Building ⚠️ IMPORTANT

**DO NOT use `go build`** - this project uses eirctl with Docker containers:

```bash
# Full build pipeline (compile + test + docs)
eirctl build

# Just compile binaries (multi-platform)
eirctl compile

# Build documentation only
eirctl docs

# Run integration tests
eirctl inttest
```

**Build contexts** (defined in [build/eirctl/contexts.yaml](../build/eirctl/contexts.yaml)):

- `buildenv` - Go compilation (ensono/eir-golang container)
- `inttestenv` - Integration tests (ensono/eir-dotnet container)
- `docsenv` - Documentation build (ensono/eir-asciidoctor container)

Build scripts in `build/scripts/` are PowerShell-based, producing multi-platform binaries in `outputs/bin/`:

- `stacks-cli-linux-amd64-{version}`
- `stacks-cli-windows-amd64-{version}.exe`
- `stacks-cli-darwin-amd64-{version}`
- `stacks-cli-darwin-arm64-{version}`

### Testing Strategy

**Unit tests**: Standard Go tests (`*_test.go`)

```bash
# Via eirctl
eirctl build  # Includes unit tests

# Direct Go (for local development)
go test -v ./...
```

Key test files:

- [pkg/config/\*\_test.go](../pkg/config/) - Configuration parsing tests
- [pkg/scaffold/scaffold_test.go](../pkg/scaffold/scaffold_test.go) - Scaffolding logic tests
- [pkg/downloaders/\*\_test.go](../pkg/downloaders/) - Downloader tests

**Integration tests**: Located in [testing/integration/](../testing/integration/) with build tag `//go:build integration`

```bash
# Run via eirctl (recommended - uses proper containers)
eirctl inttest

# Direct Go (requires compiled binary)
go test -tags=integration -v ./testing/integration/...
```

Integration tests use the actual CLI binary and validate complete scaffolding workflows across frameworks.

### Documentation Workflow

- **Source format**: AsciiDoc files in [docs/](../docs/) directory
- **Configuration**: [docs/conf/docs.json](../docs/conf/docs.json)
- **Build**: `eirctl docs`
- **Live preview**:
  ```bash
  docker run --rm -it -v ${PWD}/docs:/hugo-project/content/docs -v ${PWD}:/repo -p 1313:1313 russellseymour/hugo-docker
  ```

## Framework Integration Patterns

### Supported Frameworks and Components

Defined in [internal/config/staticFiles/config.yml](../internal/config/staticFiles/config.yml):

| Framework | Options                         | Package Type |
| --------- | ------------------------------- | ------------ |
| `dotnet`  | `webapi`                        | NuGet        |
| `java`    | `webapi`, `cqrs`, `events`      | Git          |
| `nx`      | `next`, `apps`                  | Git          |
| `infra`   | `aca`, `aks`, `eks`, `template` | Git          |
| `data`    | `fabric`, `azure`               | Git          |

### Template Processing

All framework commands use Go template expansion:

```go
replacements := config.Replacements{
    Input: s.Config.Input,     // CLI inputs
    Project: *project,         // Project-specific data
}

// Template example in stackscli.yml:
// args: new stacks-webapi -n {{ .Input.Business.Company }}.{{ .Input.Business.Domain }}
```

### Framework Commands

Each framework defines required binaries in the config:

```yaml
frameworks:
  - name: dotnet
    commands:
      - name: dotnet
        version:
          arguments: --version
          pattern: (?P<version>...)
      - name: git
```

### Pipeline Integration

Generates build pipeline files using embedded templates in [internal/config/staticFiles/](../internal/config/staticFiles/):

- `ado_variable_template.yml` - Azure DevOps variables

## Key Conventions

### Error Handling

- Uses logrus for structured logging with configurable levels (trace, debug, info, warn, error, fatal)
- Fatal errors exit via `App.Logger.Fatal()` or `App.Logger.Fatalf()`
- Validation errors accumulate and display together
- Help messages loaded from embedded [help.yml](../internal/config/staticFiles/help.yml)

### File Structure

- **Project settings**: `stackscli.yml` in template repositories
- **CLI config**: `.stackscli/config.yml` in working directory or home
- **Config constant**: `ConfigFileDir = ".stackscli"`, `ConfigName = "config"`
- **Temp/cache directories**: Configurable via `--tempdir` and `--cachedir`
- **Output structure**: Multi-platform binaries in `outputs/bin/`

### Version Management

- Version injection during build: `-ldflags "-X github.com/Ensono/stacks-cli/cmd.version={version}"`
- Default version `0.0.1-workstation` indicates local build (see [internal/constants/constants.go](../internal/constants/constants.go))
- GitHub integration via `GitHubRef = "ensono/stacks-cli"`

## Development Loop

1. Make code changes
2. Run `eirctl build` to compile, test, and build docs
3. For quick iteration, run `go test -v ./...` for unit tests
4. Update documentation in `docs/` if functionality changes
5. Create a conventional commit message and push changes
6. Create a pull request for review

## Security Considerations

- **Never override or bypass security controls** such as GPG signing - prompt the user to correct issues
- **Template rendering security**: All user inputs go through Go templates - ensure proper escaping and validation
- **Command execution**: Only allowed commands per framework are executed (whitelist in config)
- **External inputs**: CLI flags and config files must be sanitized and validated
- **Dependencies**: Keep up to date and monitor for vulnerabilities
- **Least privilege**: Use minimal permissions when downloading templates or executing commands
- **Logging**: Never log sensitive information such as tokens or secrets

## Agent-Specific Guidance

### When Making Changes

1. **Understand the flow**: CLI flags → Viper → Config → Scaffold → Operations
2. **Check the config package first**: Most functionality depends on [pkg/config/](../pkg/config/)
3. **Template changes**: Test with multiple frameworks - templates are shared
4. **New frameworks**: Update [internal/config/staticFiles/config.yml](../internal/config/staticFiles/config.yml)

### Common Tasks

**Adding a new CLI flag**:

1. Add flag definition in appropriate `cmd/*.go` file
2. Bind to Viper: `viper.BindPFlag("input.path.to.value", cmd.Flags().Lookup("flagname"))`
3. Add corresponding field in `pkg/config/` structs
4. Update documentation in `docs/runtime_config/`

**Adding a new framework component**:

1. Add to `stacks.components` in [config.yml](../internal/config/staticFiles/config.yml)
2. Ensure framework binary requirements are defined
3. Create template repository with `stackscli.yml`

**Adding a new downloader**:

1. Implement `interfaces.Downloader` interface
2. Add to switch in [pkg/scaffold/scaffold.go](../pkg/scaffold/scaffold.go) `processProject` method

### Files to Check Before Committing

- [ ] Unit tests pass: `go test -v ./...`
- [ ] No linting errors
- [ ] Documentation updated if user-facing changes
- [ ] Config structs have proper mapstructure tags
- [ ] New flags bound to Viper correctly

When modifying this codebase, pay special attention to the template rendering system and ensure all user inputs are properly escaped and validated before processing.
