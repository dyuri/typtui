# typtui - TUI TYP Editor

A terminal user interface (TUI) application for editing Garmin TYP files on Linux, built with Go and Bubbletea.

## Overview

**typtui** provides Linux users with a native, user-friendly tool for customizing Garmin map styles without requiring Windows GUI tools or manual text file editing. It's designed for:

- Linux-based Garmin map creators
- OpenStreetMap contributors
- Hiking/outdoor enthusiasts customizing maps
- Users of mkgmap who want easier TYP file management

## Features (Phase 1 - MVP)

✅ **Current Features:**
- Parse and load TYP text files
- Browse point, line, and polygon type definitions
- Tab-based navigation between type categories
- Keyboard-driven interface
- Detailed type information display

🚧 **In Development:**
- Edit type properties (colors, labels, dimensions)
- Save changes back to TYP format
- mkgmap compilation integration
- Color picker with terminal preview
- Type management (add/delete/clone)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/dyuri/typtui.git
cd typtui

# Build the application
make build

# Run the application
./bin/typtui path/to/file.typ
```

### Using Go Install

```bash
go install github.com/dyuri/typtui/cmd/typtui@latest
```

## Usage

```bash
# Open a TYP file
typtui mymap.typ

# The application will launch in your terminal
```

### Keyboard Shortcuts

- **Tab** - Switch between Points/Lines/Polygons tabs
- **↑/k** - Move selection up
- **↓/j** - Move selection down
- **?** - Toggle help screen
- **q** or **Ctrl+C** - Quit

## Requirements

- Go 1.21 or later (for building from source)
- A modern terminal (optimized for Kitty, but works in others)
- mkgmap (for TYP compilation - future feature)

## Development

### Building

```bash
make build        # Build the application
make run          # Run without building
make run-file FILE=path/to/file.typ  # Run with specific file
```

### Testing

```bash
make test         # Run all tests
make test-parser  # Run parser tests only
make test-coverage # Generate coverage report
```

### Code Quality

```bash
make fmt          # Format code
make lint         # Run linter
make tidy         # Tidy dependencies
```

## Project Structure

```
typtui/
├── cmd/typtui/           # Main entry point
├── internal/
│   ├── parser/           # TYP file parser
│   ├── tui/              # Bubbletea TUI components
│   ├── compiler/         # mkgmap wrapper (future)
│   └── utils/            # Utilities (future)
├── testdata/             # Test TYP files
├── docs/                 # Documentation
└── Makefile              # Build tasks
```

## TYP File Format

TYP files are text-based configuration files used by mkgmap to create custom map styles for Garmin devices. They define:

- **Points**: POI (Point of Interest) icons and labels
- **Lines**: Road and trail styles (width, color, borders)
- **Polygons**: Area fills and patterns (forests, parks, water)

Example:

```
[_id]
CodePage=1252
FID=1234
[end]

[_point]
Type=0x2f06
String=0x04,Bank
DayXpm="16 16 2 1"
"! c #778899"
"  c none"
"!!!!!!!!!!"
[end]
```

## Contributing

Contributions are welcome! This project is in early development (Phase 1 of the implementation plan).

### Current Focus

We're currently working on:
1. ✅ Parser implementation (read-only)
2. ✅ Basic TUI shell
3. 🚧 Edit mode for type properties
4. 🚧 Save functionality
5. 🚧 mkgmap integration

See `typ-editor-implementation-plan.md` for the full roadmap.

## License

MIT License - See LICENSE file for details

## Resources

- [mkgmap TYP Documentation](https://www.mkgmap.org.uk/doc/typ-compiler)
- [OSM Garmin Type Codes](https://wiki.openstreetmap.org/wiki/OSM_Map_On_Garmin/POI_Types)
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework

## Acknowledgments

- Built with [Charm's](https://charm.sh/) excellent TUI libraries (Bubbletea, Bubbles, Lipgloss)
- Designed for the OpenStreetMap and Garmin mapping community
- Optimized for Kitty terminal with fallback support for other terminals
