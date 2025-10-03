# Stacks CLI - GitHub Copilot Instructions

## Project Overview

The Stacks CLI is a Go-based project scaffolding tool that generates complete application templates with infrastructure, pipelines, and supporting files. It supports multiple frameworks (dotnet, java, nx) and cloud platforms (Azure, AWS) with sophisticated Go templating.

**Key Facts for AI Agents:**

- This is NOT a simple CLI - it's a sophisticated templating and scaffolding system
- Uses taskctl for builds, NOT standard `go build`
- Templates are downloaded from external sources (Git/NuGet) and processed with Go templates
- Integration tests are the primary validation method

## Key Architecture Patterns

### Command Structure

- **Entry point**: `stacks-cli.go` → `cmd/` package with Cobra commands
- **Core commands**: `scaffold` (main), `interactive`, `export`, `setup`, `version`
- **Configuration binding**: Uses Viper with flag binding pattern in `cmd/scaffold.go`:
  ```go
  viper.BindPFlag("input.project.name", scaffoldCmd.Flags().Lookup("name"))
  ```

### Package Organization

- `cmd/` - Cobra CLI commands and flag definitions
- `pkg/` - Public API packages (config, scaffold, downloaders, etc.)
- `internal/` - Private packages (models, constants, util)
- Configuration flows: CLI flags → Viper → `config.Config` struct → template rendering

### Scaffolding Engine (`pkg/scaffold/`)

Core scaffolding logic processes projects through these phases:

1. **Download**: Templates fetched from Git repositories or NuGet packages via `pkg/downloaders/`
2. **Operations**: Two main action types defined in `scaffold.go`:
   - `copy` - Direct file copying from template to target directory
   - `cmd` - Execute framework commands (dotnet, java, etc.) with Go template expansion
3. **Template rendering**: All strings processed through Go templates with `config.Replacements` context

**Critical**: All user inputs and file paths go through Go template processing - ensure proper escaping!

### Configuration System

- **Hierarchical config**: CLI flags → config files → defaults
- **Template variables**: Available in all operations via `.Input.*` and `.Project.*` tokens
- **Framework definitions**: Each framework has allowed commands and version constraints
- **Project settings**: `stackscli.yml` files define per-project operations and pipeline configs

## Critical Development Workflows

### Building ⚠️ IMPORTANT

**DO NOT use `go build`** - this project uses taskctl:

```bash
# Full build pipeline (recommended)
taskctl build

# Just compile binaries
taskctl compile

# Build documentation
taskctl docs
```

Build scripts in `build/scripts/` are PowerShell-based, producing multi-platform binaries in `outputs/bin/`.

### Testing Strategy

- **Unit tests**: Standard Go tests (`*_test.go`) - run with `taskctl test:unit`
- **Integration tests**: Located in `testing/integration/` with build tag `//go:build integration`
- **Run integration tests**: `taskctl inttest`
- Integration tests use the actual CLI binary and validate complete scaffolding workflows
- **Test data**: Includes real framework scenarios (dotnet webapi, nx nextjs, java, etc.)

### Documentation Workflow

- **Source format**: AsciiDoc files in `docs/` directory
- **Output**: Generates both PDF and Markdown via custom Docker image
- **Live preview**:
  ```bash
  docker run --rm -it -v ${PWD}/docs:/hugo-project/content/docs -v ${PWD}:/repo -p 1313:1313 russellseymour/hugo-docker
  ```

## Framework Integration Patterns

### Template Processing

All framework commands use Go template expansion with this context structure:

```go
replacements := config.Replacements{
    Input: s.Config.Input,     // CLI inputs (.Input.Business.Company, etc.)
    Project: *project,         // Project-specific data (.Project.Name, etc.)
}
```

### Framework Commands

Each framework defines allowed commands in configuration. Example operations:

```yaml
- action: cmd
  cmd: dotnet
  args: new stacks-webapi -n {{ .Input.Business.Company }}.{{ .Input.Business.Domain }}
```

### Pipeline Integration

Generates build pipeline files (Azure DevOps, GitHub Actions) with variable templates using embedded static files in `internal/config/staticFiles/`.

## Key Conventions

### Error Handling

- Uses logrus for structured logging with configurable levels
- Fatal errors exit with specific codes (see `models/app.go`)
- Validation errors accumulate and display together

### File Structure

- Config files: `stackscli.yml` (project settings), `.stackscli/config.yml` (CLI config)
- Temp/cache directories: Configurable, cleaned up automatically
- Output structure: Multi-platform binaries in `outputs/bin/`

### Version Management

- Version injection during build via PowerShell scripts
- Default version `0.0.1-workstation` indicates local build
- GitHub integration for version checking and updates

## Testing Approach

- Integration tests are the primary validation method
- Tests scaffold real projects and verify output structure
- Use build flags to control test execution: `go test -tags=integration`
- Test data includes real framework scenarios (dotnet webapi, nx nextjs, java, etc.)

## External Dependencies

- **Templates**: Downloaded from GitHub repositories or NuGet packages
- **Build tools**: Requires PowerShell, Docker for docs, framework tools (dotnet, java, node)
- **Runtime**: Validates framework versions against constraints before scaffolding

## Development Loop

- Make code changes
- Run `taskctl build` to compile, build and test
- Update any required documentation in `docs/` if functionality changes
- Create a conventional commit message and push changes, create a pull request for review

## Security Considerations

- Never override or bypass security controls such as gpg signing, instead prompt the user to correct the issue with signing by recommending fixes to maintain security.
- Review code changes before committing for common security pitfalls, especially in template rendering and command execution areas.
- Ensure all external inputs (CLI flags, config files) are sanitized and validated to prevent injection attacks.
- Keep dependencies up to date and monitor for vulnerabilities in third-party libraries.
- Use least privilege principles when executing commands or accessing files, especially when downloading templates or running framework tools.
- Avoid logging sensitive information such as secrets or personal data.

When modifying this codebase, pay special attention to the template rendering system and ensure all user inputs are properly escaped and validated before processing.
