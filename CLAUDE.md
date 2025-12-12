# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Projman is a CLI tool written in Go for batch managing multiple Maven/Git projects. It allows users to:
- Create profiles with selected Maven projects from a root directory
- Execute batch Git operations (update with smart branch handling)
- Execute batch Maven operations (install with automatic dependency ordering)

The tool uses multi-profile configuration stored in JSON format at `~/.config/projman/projman_config.json` (Linux/macOS) or `%APPDATA%/projman/projman_config.json` (Windows).

## Development Commands

### Build and Test
```bash
make build      # Compile the binary (outputs 'projman' or 'projman.exe')
make test       # Run all tests
make lint       # Format code (go fmt) and run vet
make install    # Install to GOPATH/bin
make clean      # Remove build artifacts
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/graph
go test ./internal/maven

# Run with verbose output
go test -v ./...
```

### Release Process
```bash
# Test release locally (requires goreleaser)
make snapshot

# Trigger GitHub release (creates tag and triggers CI)
make release VERSION=v1.2.3
```

The release process is automated via GitHub Actions (`.github/workflows/release.yml`) which uses goreleaser to create multi-platform binaries.

## Architecture

### Command Structure (Cobra-based)
- `cmd/root.go` - Root command setup, registers `git` and `mvn` subcommands
- `cmd/init.go` - Create new profile: scan directory for Maven projects, interactive selection
- `cmd/list.go` - List all profiles, show current profile
- `cmd/use.go` - Switch active profile
- `cmd/delete.go` - Delete a profile
- `cmd/git/update.go` - Git batch operations with smart branch handling
- `cmd/mvn/install.go` - Maven install with dependency-based ordering

### Core Internal Packages

#### `internal/config`
Manages multi-profile configuration with JSON persistence:
- `Config` - Single profile (root path + selected projects)
- `ProfileConfig` - All profiles + current profile name
- Profile operations: Load/Save/List/Delete/Switch

#### `internal/project`
Maven project discovery:
- Scans directory for `pom.xml` files (one level deep)
- Returns list of discovered Maven projects
- Filtering and name extraction utilities

#### `internal/buildsystem`
Abstraction layer for build systems (designed for extensibility):
- `Parser` interface - Parse project metadata and dependencies
- `Project` - Generic project representation with identifier and dependencies
- `ArtifactRegistry` - Maps artifact IDs (e.g., `groupId:artifactId`) to project names
- `GraphBuilder` - Builds dependency graph from projects using a Parser

#### `internal/maven`
Maven-specific implementation:
- `Parser` - Implements `buildsystem.Parser` for Maven projects
- Parses `pom.xml` files to extract groupId, artifactId, dependencies, and modules
- Handles nested sub-modules recursively
- `RegisterSubModules()` - Registers all sub-modules in ArtifactRegistry to resolve transitive dependencies
- `BuildDependencyGraph()` - High-level function to build dependency graph for selected projects

#### `internal/graph`
Dependency graph data structure and algorithms:
- `DependencyGraph` - Map of project name to list of dependencies
- `TopologicalSort()` - Uses Kahn's algorithm to order projects for build
- Detects circular dependencies

#### `internal/exec`
Command execution utilities:
- `Run()` - Execute command and stream output
- `RunWithOutput()` - Execute and capture output as string
- `WaitForUserInput()` - Interactive prompt to continue/abort after errors

#### `internal/ui`
Interactive UI components using pterm:
- `MultiSelectTable` - Interactive table with multi-selection support
- Used for branch switching selection in git update

### Key Workflows

#### Git Update Workflow (cmd/git/update.go)
1. Load config and select projects interactively
2. Gather branch information for all projects
3. Show table of projects not on `develop`, allow user to select which to switch
4. For each project:
   - Stash uncommitted changes (if any tracked files modified)
   - Switch to develop if selected
   - Pull/merge based on branch type:
     - `develop`: `git pull origin develop`
     - `deploy/*`: `git pull origin <current-branch>`
     - Other branches: `git fetch origin develop && git merge origin/develop`
   - Pop stash if it was created

#### Maven Install Workflow (cmd/mvn/install.go)
1. Load config and select projects interactively
2. Parse all `pom.xml` files to extract dependencies
3. Register main artifacts and sub-modules in ArtifactRegistry
4. Build dependency graph between selected projects
5. Topological sort to determine correct build order
6. Execute `mvn clean install` (with `-DskipTests=true` by default, or `--tests`/`-t` flag to run tests)
7. If error occurs, prompt user to continue or abort

### Dependency Resolution Strategy

The tool builds a dependency graph based on Maven `pom.xml` files:
1. Parse each selected project's `pom.xml` for dependencies and modules
2. Recursively parse sub-modules to find all dependencies
3. Register all artifacts (main project + sub-modules) in `ArtifactRegistry`
4. Build graph where edges represent "depends on" relationships between selected projects
5. Topological sort ensures projects are built before their dependents

This approach correctly handles:
- Direct dependencies between selected projects
- Transitive dependencies through sub-modules
- Multi-module Maven projects
- Circular dependency detection

### Version Management

The version is set via ldflags during build:
```go
// cmd/root.go
var Version = "dev"  // Overridden at build time
```

Build command uses: `-ldflags="-s -w -X github.com/SalvatoreSpagnuolo-BipRED/projman/cmd.Version=$(VERSION)"`

## Testing

Tests are located in the same package as the code they test:
- `internal/graph/graph_test.go` - Graph operations and topological sort
- `internal/maven/dependency_test.go` - Dependency parsing

When adding new features, add tests in `*_test.go` files in the same directory.

## Code Conventions

- All user-facing messages are in Italian (this is intentional)
- Use `pterm` for all console output (colors, spinners, tables)
- Commands follow pattern: load config → interactive selection → save selection → execute operation
- Error handling: display error with `pterm.Error`, prompt user to continue with `exec.WaitForUserInput()`
- Project discovery is shallow (one level deep from root directory)
