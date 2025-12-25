# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**typtui** is a TUI (Terminal User Interface) application for editing Garmin TYP files on Linux, built with Go and Bubbletea. It provides a native, user-friendly tool for customizing Garmin map styles without requiring Windows GUI tools or manual text file editing.

**Target Platform:** Linux (optimized for Kitty terminal, but should work in other modern terminals)

## Development Commands

### Initial Setup
```bash
# Initialize Go module (if not already done)
go mod init github.com/yourusername/typtui
go mod tidy

# Install dependencies
go get github.com/charmbracelet/bubbletea@v0.25.0
go get github.com/charmbracelet/bubbles@v0.18.0
go get github.com/charmbracelet/lipgloss@v0.9.1
go get github.com/adrg/xdg@v0.4.0
go get gopkg.in/yaml.v3@v3.0.1
```

### Build & Run
```bash
# Build the application
make build  # or: go build -o typtui ./cmd/typtui

# Run directly
go run ./cmd/typtui [file.typ]

# Install locally
go install ./cmd/typtui
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/parser/...

# Run specific test
go test -run TestParsePointType ./internal/parser/
```

### Code Quality
```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run
```

## Architecture

### Core Components

**Parser (`internal/parser/`)**: Parses TYP text format into Go structures. The TYP format has sections like `[_point]`, `[_line]`, `[_polygon]` with various properties. Key challenges include multi-line XPM data, multiple string encodings, and optional fields.

**TUI (`internal/tui/`)**: Bubbletea-based interface with multiple modes:
- **List Mode**: Browse all type definitions (points, lines, polygons) with tab switching
- **Edit Mode**: Form-based editing of type properties, colors, labels
- **Help Mode**: Keyboard shortcuts and usage help

**Compiler (`internal/compiler/`)**: Wraps mkgmap external tool to compile TYP text files into binary `.typ` files for Garmin devices.

### Data Flow

1. User selects TYP file → Parser reads and parses into `TYPFile` struct
2. TUI displays type definitions in List Mode
3. User edits properties → Modified in-memory structures
4. User saves → Serialize back to text format → Optionally compile with mkgmap

### Key Data Structures

All defined in `internal/parser/types.go`:

- **TYPFile**: Root structure containing Header, Points, Lines, Polygons, DrawOrder, Icons
- **PointType**: POI definitions with labels, colors, icons (XPM format)
- **LineType**: Road/trail definitions with width, border, colors, line style
- **PolygonType**: Area definitions with patterns, colors
- **XPMIcon**: Icon data in XPM format (multi-line bitmap with color palette)

### Bubbletea Model Structure

The main TUI model (`internal/tui/model.go`) contains:
- **Core data**: Current `TYPFile`, modification state
- **UI state**: Current mode, active tab, selected index
- **Components**: List view, editor forms, color picker, status bar
- **Navigation**: Tab switching between Points/Lines/Polygons

Update logic delegates to mode-specific handlers. View rendering is split by sections (header, content, footer).

## TYP File Format

TYP files are text-based with sections like:

```
[_id]
CodePage=1252
FID=1234
[end]

[_point]
Type=0x2f06
String1=0x04,Bank
IconXpm="16 16 2 1"
"! c #778899"
"  c none"
"!!!!!!!!!!"
[end]
```

**Parsing Challenges:**
- Comments start with `;` or `#`
- XPM data spans multiple lines within quotes
- Language codes in strings (0x04 = English, etc.)
- Many optional fields
- Various quote styles and escaping

## Configuration (XDG Compliance)

Follows XDG Base Directory Specification:

- **Config**: `$XDG_CONFIG_HOME/typtui/config.yaml` (fallback: `~/.config/typtui/config.yaml`)
- **Data**: `$XDG_DATA_HOME/typtui/recent.json` (fallback: `~/.local/share/typtui/recent.json`)
- **Cache**: `$XDG_CACHE_HOME/typtui/` (fallback: `~/.cache/typtui/`)

Config file includes mkgmap path, editor preferences, color theme, recent files.

## External Dependencies

**mkgmap**: Required for TYP compilation. Should be auto-detected in PATH or configured via config file. Target latest stable version (r4922+).

**Go packages**: See `go.mod` - primarily Charm libraries (bubbletea, bubbles, lipgloss), plus xdg and yaml support.

## Terminal Optimization

**Primary target**: Kitty terminal with true color (24-bit), GPU acceleration, Unicode support, and future image protocol support for icon previews.

**Fallback strategy**: Detect terminal capabilities on startup (`TERM` env var, `COLORTERM` for true color). Gracefully degrade to 256 colors and ASCII characters if needed.

## Testing Strategy

**Unit tests**: Parser for each section type, color utilities, validation functions, type database lookups. Use table-driven tests.

**Integration tests**: Load real TYP files (OpenHiking, OpenMTBMap), test parse → modify → save → parse round-trips, verify mkgmap compilation.

**Test data**: Located in `testdata/` directory with sample files (minimal, complex, real-world, and invalid cases).

## Implementation Phases

1. **Foundation**: Parser (read-only) + basic TUI shell
2. **List Mode**: Type browser with tab switching and details view
3. **Edit Mode**: Form-based editing of properties and colors
4. **Save & Compile**: Write TYP format + mkgmap integration
5. **Advanced Features**: XPM support, type management, Garmin type database
6. **Polish**: UX improvements, documentation, testing

Current status: Early development phase. See `typ-editor-implementation-plan.md` for detailed phase breakdown.

## Error Handling

Use explicit error handling with context:

```go
if err != nil {
    return fmt.Errorf("failed to parse point type at line %d: %w", lineNum, err)
}
```

Define custom error types (e.g., `ParseError`) with line/column information for user-friendly error messages.

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Keep functions small and focused
- Write clear comments for TYP format parsing logic
- Use meaningful variable names
