# TUI TYP Editor - Comprehensive Implementation Plan

## Project Overview

**Project Name:** `typtui` (TYP TUI Editor)

**Description:** A terminal user interface (TUI) application for editing Garmin TYP files on Linux, built with Go and Bubbletea.

**Initial Focus:**
- Linux only (other platforms later)
- Latest mkgmap version
- Optimized for Kitty terminal (but should work in others)

**Goal:** Provide Linux users with a native, user-friendly tool for customizing Garmin map styles without requiring Windows GUI tools or manual text file editing.

**Target Users:** 
- Linux-based Garmin map creators
- OpenStreetMap contributors
- Hiking/outdoor enthusiasts customizing maps
- Users of mkgmap who want easier TYP file management

---

## Technical Stack

- **Language:** Go 1.21+
- **TUI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **UI Components:** [Bubbles](https://github.com/charmbracelet/bubbles) (list, textinput, viewport, etc.)
- **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **TYP Compilation:** mkgmap (external tool, called via exec)
- **File Parsing:** Custom parser (Go standard library + regex)

---

## Project Structure

```
typtui/
├── cmd/
│   └── typtui/
│       └── main.go                 # Entry point
├── internal/
│   ├── parser/
│   │   ├── parser.go              # TYP text format parser
│   │   ├── parser_test.go
│   │   └── types.go               # Type definitions for TYP structures
│   ├── compiler/
│   │   ├── mkgmap.go              # mkgmap wrapper
│   │   └── validator.go           # Validation logic
│   ├── tui/
│   │   ├── model.go               # Main bubbletea model
│   │   ├── update.go              # Update logic
│   │   ├── view.go                # View rendering
│   │   ├── modes/
│   │   │   ├── list.go            # List mode
│   │   │   ├── edit.go            # Edit mode
│   │   │   └── help.go            # Help mode
│   │   └── components/
│   │       ├── colorpicker.go     # Color picker component
│   │       ├── typelist.go        # Type list component
│   │       └── form.go            # Form component
│   └── utils/
│       ├── colors.go              # Color conversion utilities
│       └── garmin.go              # Garmin type code database
├── testdata/
│   ├── sample.typ                 # Sample TYP files for testing
│   └── expected/                  # Expected parsing results
├── docs/
│   ├── USAGE.md                   # User guide
│   ├── TYP_FORMAT.md             # TYP format documentation
│   └── DEVELOPMENT.md             # Development guide
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── Makefile
```

---

## Data Structures

### Core Types

```go
// TYPFile represents the entire TYP file structure
type TYPFile struct {
    Header      Header
    Points      []PointType
    Lines       []LineType
    Polygons    []PolygonType
    DrawOrder   DrawOrder
    Icons       []Icon
    FilePath    string
    Modified    bool
}

// Header contains TYP file metadata
type Header struct {
    CodePage    int
    FID         int
    ProductCode int
    MapID       int
}

// PointType represents a POI definition
type PointType struct {
    Type        string                 // e.g., "0x2f06"
    SubType     string                 // Optional subtype
    Labels      map[string]string      // Language code -> label
    Icon        *XPMIcon
    DayColors   []Color
    NightColors []Color
    FontStyle   string
}

// LineType represents a line definition (roads, trails, etc.)
type LineType struct {
    Type           string
    Labels         map[string]string
    LineWidth      int
    BorderWidth    int
    DayColor       Color
    NightColor     Color
    DayBorderColor Color
    NightBorderColor Color
    UseOrientation bool
    LineStyle      string // "solid", "dashed", etc.
}

// PolygonType represents an area definition
type PolygonType struct {
    Type           string
    Labels         map[string]string
    Pattern        *XPMPattern
    DayColor       Color
    NightColor     Color
    FontStyle      string
    ExtendedLabels bool
}

// Color represents a color in hex format
type Color struct {
    Hex   string // "#RRGGBB"
    Day   bool   // true if day color, false if night
    Name  string // Optional color name
}

// XPMIcon represents icon data in XPM format
type XPMIcon struct {
    Width   int
    Height  int
    Colors  int
    Data    []string
    Palette map[string]Color
}

// DrawOrder specifies rendering order
type DrawOrder struct {
    Points   []string
    Lines    []string
    Polygons []string
}
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1-2)

#### 1.1 Project Setup
- [ ] Initialize Go module
- [ ] Set up project structure
- [ ] Add dependencies (bubbletea, bubbles, lipgloss)
- [ ] Create Makefile with build/test/run targets
- [ ] Set up basic CI (GitHub Actions)
- [ ] Create README with project description

#### 1.2 TYP Parser - Read-Only
- [ ] Implement lexer for TYP text format
- [ ] Parse header section
- [ ] Parse polygon definitions
- [ ] Parse line definitions  
- [ ] Parse point definitions
- [ ] Parse XPM data (basic structure)
- [ ] Handle comments and whitespace
- [ ] Unit tests for parser (use testdata/)
- [ ] Error handling with line numbers

**Deliverable:** Can load and parse existing TYP files into Go structures

**Test Files:**
```
testdata/
├── minimal.typ          # Bare minimum valid TYP
├── openhiking.typ       # Real-world OpenHiking TYP
├── complex.typ          # Complex with all features
└── invalid/
    ├── syntax_error.typ
    └── missing_end.typ
```

#### 1.3 Basic TUI Shell
- [ ] Create main bubbletea model
- [ ] Implement basic navigation (quit, help)
- [ ] File selection/loading screen
- [ ] Error display
- [ ] Status bar with file info
- [ ] Basic keyboard shortcuts

**Deliverable:** Can launch TUI, select file, display parsed data

---

### Phase 2: List Mode (Week 3)

#### 2.1 Type Browser
- [ ] Display list of all point types
- [ ] Display list of all line types
- [ ] Display list of all polygon types
- [ ] Tab switching between type categories
- [ ] Search/filter functionality
- [ ] Sort options (by type code, by name)
- [ ] Selection highlighting
- [ ] Pagination for long lists

#### 2.2 Type Details View
- [ ] Show selected type details in split pane
- [ ] Display color swatches (terminal approximation)
- [ ] Show labels for all languages
- [ ] Display basic properties
- [ ] Show type code with Garmin description
- [ ] Indicate if icon/pattern is defined

**Deliverable:** Can browse and view all types in a TYP file

**UI Mockup:**
```
┌─────────────────────────────────────────────────────────────┐
│ typtui - openhiking.typ                        [Modified: *] │
├─────────────────────────────────────────────────────────────┤
│ Points │ Lines │ Polygons                                    │
├─────────────┬───────────────────────────────────────────────┤
│             │ Type Details                                  │
│ 0x0100  ▸   │ ─────────────                                 │
│ 0x0200      │ Type: 0x2f06 (Bank POI)                       │
│ 0x2f06  ■   │ Label (EN): "Bank"                            │
│ 0x2f07      │ Label (HU): "Bank"                            │
│ 0x3000      │ Day Color:   ████ #778899                     │
│ 0x3001      │ Night Color: ████ #334455                     │
│             │ Icon: 16x16 (defined)                         │
│             │ Font: SmallFont                               │
│             │                                               │
│  (12/156)   │ [e] Edit  [d] Delete  [c] Clone               │
└─────────────┴───────────────────────────────────────────────┘
 [tab] Switch  [/] Search  [e] Edit  [q] Quit  [?] Help
```

---

### Phase 3: Edit Mode (Week 4-5)

#### 3.1 Type Editor - Basic Properties
- [ ] Edit mode activation (from list)
- [ ] Form-based editing interface
- [ ] Edit type code (with validation)
- [ ] Edit labels (add/remove languages)
- [ ] Edit font style (dropdown)
- [ ] Save changes to memory
- [ ] Cancel changes
- [ ] Dirty flag tracking

#### 3.2 Color Editor
- [ ] Color picker component
- [ ] Hex color input with validation
- [ ] Terminal color preview (closest match)
- [ ] Day/night color toggling
- [ ] Predefined color palette
- [ ] Color naming/descriptions
- [ ] Copy color from another type

**Color Picker Component:**
```
┌─────────────────────────────┐
│ Select Day Color            │
├─────────────────────────────┤
│ Hex: #778899 ████           │
│                             │
│ Common Colors:              │
│ [White  ] [Black  ] [Red   ]│
│ [Green  ] [Blue   ] [Yellow]│
│ [Brown  ] [Gray   ] [Orange]│
│                             │
│ Or enter hex code:          │
│ > #__________              │
│                             │
│ [OK] [Cancel]               │
└─────────────────────────────┘
```

#### 3.3 Line-Specific Editors
- [ ] Line width editor (numeric input)
- [ ] Border width editor
- [ ] Line style selector (solid, dashed, etc.)
- [ ] Orientation toggle
- [ ] Preview of line style (ASCII)

#### 3.4 Validation
- [ ] Type code range validation (Garmin specs)
- [ ] Hex color format validation
- [ ] Required field checking
- [ ] Duplicate type detection
- [ ] Real-time validation feedback

**Deliverable:** Can edit basic type properties and colors

---

### Phase 4: Save & Compile (Week 6)

#### 4.1 TYP Writer
- [ ] Serialize TYPFile back to text format
- [ ] Preserve comments where possible
- [ ] Format output consistently
- [ ] Handle all type categories
- [ ] Write XPM data correctly
- [ ] Generate proper section headers
- [ ] Save to new file
- [ ] Backup original file

#### 4.2 mkgmap Integration
- [ ] Detect mkgmap installation
- [ ] Build mkgmap command line
- [ ] Execute compilation
- [ ] Capture stdout/stderr
- [ ] Parse compilation errors
- [ ] Display results in TUI
- [ ] Handle compilation failures
- [ ] Success/failure notifications

**mkgmap Wrapper:**
```go
type MkgmapCompiler struct {
    BinaryPath string
    WorkDir    string
}

func (m *MkgmapCompiler) Compile(typFile string) (*CompileResult, error) {
    // Run: mkgmap --family-id=1234 typfile.txt
    // Parse output
    // Return result with errors/warnings
}

type CompileResult struct {
    Success  bool
    Warnings []string
    Errors   []string
    Output   string
}
```

#### 4.3 Save Workflow
- [ ] Unsaved changes warning
- [ ] Save dialog (confirm file path)
- [ ] Auto-save option
- [ ] Save as... functionality
- [ ] Compile after save option
- [ ] Show compilation results

**Deliverable:** Can save changes and compile TYP files with mkgmap

---

### Phase 5: Advanced Features (Week 7-8)

#### 5.1 XPM Support (Basic)
- [ ] Display XPM icon info (dimensions, colors)
- [ ] ASCII preview of icons (simplified)
- [ ] Import XPM from file
- [ ] Export XPM to file
- [ ] Basic validation
- [ ] Future: inline editor (stretch goal)

#### 5.2 Type Management
- [ ] Add new type (from template)
- [ ] Clone existing type
- [ ] Delete type (with confirmation)
- [ ] Reorder types
- [ ] Import types from another TYP
- [ ] Export selected types

#### 5.3 Garmin Type Database
- [ ] Built-in type code reference
- [ ] Descriptions for standard types
- [ ] Quick lookup (e.g., "What is 0x2f06?")
- [ ] Type category browsing
- [ ] Add custom type codes

**Type Database Format:**
```go
var GarminTypes = map[string]TypeInfo{
    "0x2f06": {
        Category:    "Point",
        Name:        "Bank",
        Description: "Banking and ATM facilities",
        Common:      true,
    },
    // ... hundreds more
}
```

#### 5.4 Configuration
- [ ] User config file (XDG_CONFIG_HOME/typtui/config.yaml)
- [ ] Fallback to ~/.config/typtui/config.yaml if XDG_CONFIG_HOME not set
- [ ] mkgmap path configuration
- [ ] Color preferences
- [ ] Default language for labels
- [ ] Editor preferences
- [ ] Recent files list (XDG_DATA_HOME/typtui/recent.json or ~/.local/share/typtui/recent.json)

**Deliverable:** Full-featured editor with advanced capabilities

---

### Phase 6: Polish & Documentation (Week 9)

#### 6.1 UX Improvements
- [ ] Keyboard shortcut customization
- [ ] Mouse support (optional)
- [ ] Undo/redo for edits
- [ ] Multiple file tabs
- [ ] Diff view (before/after)
- [ ] Export to HTML/Markdown report

#### 6.2 Error Handling
- [ ] Graceful handling of malformed files
- [ ] Clear error messages
- [ ] Recovery suggestions
- [ ] Debug mode with verbose logging
- [ ] Error reporting template

#### 6.3 Documentation
- [ ] Complete README.md
- [ ] USAGE.md with examples
- [ ] TYP_FORMAT.md (format specification)
- [ ] DEVELOPMENT.md (for contributors)
- [ ] Man page (optional)
- [ ] Tutorial: "Creating Your First Custom TYP"
- [ ] Video demo (screencast)

#### 6.4 Testing
- [ ] Unit tests (>80% coverage)
- [ ] Integration tests
- [ ] Test with real-world TYP files
- [ ] Performance testing (large files)
- [ ] Memory leak testing
- [ ] Platform testing:
  - [ ] Primary: Linux (Ubuntu, Arch, Fedora) with Kitty terminal
  - [ ] Secondary: Linux with gnome-terminal, alacritty, konsole
  - [ ] Future: macOS, other platforms

**Deliverable:** Production-ready v1.0 release

---

## Key Features by Priority

### Must Have (MVP)
- [x] Load and parse TYP files
- [x] Browse all type definitions
- [x] Edit basic properties (colors, labels)
- [x] Save changes back to TYP format
- [x] Call mkgmap for compilation
- [x] Keyboard navigation
- [x] Validation

### Should Have (v1.0)
- [ ] Color picker with preview
- [ ] Type management (add/delete/clone)
- [ ] Search and filter
- [ ] Unsaved changes warning
- [ ] Configuration file
- [ ] Garmin type reference
- [ ] Good documentation

### Nice to Have (v1.1+)
- [ ] XPM icon editor (inline)
- [ ] Visual diff tool
- [ ] Import/export between TYP files
- [ ] Theme customization
- [ ] Multi-file editing
- [ ] Undo/redo
- [ ] Mouse support
- [ ] Export to different formats

### Future Ideas
- [ ] Web preview of map styles
- [ ] Integration with mkgmap style files
- [ ] Batch processing
- [ ] TYP file templates library
- [ ] Cloud sync for TYP files
- [ ] Plugin system
- [ ] GUI version (Fyne/Qt)

---

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` and `golangci-lint`
- Write clear comments for complex logic
- Keep functions small and focused
- Use meaningful variable names

### Error Handling
```go
// Prefer explicit error handling
if err != nil {
    return fmt.Errorf("failed to parse point type at line %d: %w", lineNum, err)
}

// Use custom error types for specific cases
type ParseError struct {
    Line    int
    Column  int
    Message string
}
```

### Testing Strategy
```go
// Table-driven tests for parser
func TestParsePointType(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *PointType
        wantErr bool
    }{
        {"valid point", "[_point]\nType=0x2f06\n[end]", &PointType{...}, false},
        {"missing type", "[_point]\n[end]", nil, true},
        // ...
    }
    // ...
}
```

### Performance Considerations
- Use buffered I/O for file reading
- Parse lazily where possible
- Cache terminal size calculations
- Optimize rendering (only redraw changed parts)
- Profile with `pprof` if performance issues arise

---

## Kitty Terminal Optimizations

### Why Kitty?
- True color support (24-bit)
- GPU acceleration for smooth rendering
- Unicode support (including box drawing characters)
- Image protocol support (future: inline icon previews)
- Modern, actively maintained

### Kitty-Specific Features to Leverage

**Phase 1 (MVP):**
- Use full RGB color space for accurate color representation
- Leverage Unicode box drawing for better UI borders
- Utilize true color in status bar and highlights

**Future Enhancements:**
- Kitty graphics protocol for displaying actual icon previews
- Inline image display of XPM icons
- Better color picker with RGB gradients

### Fallback for Other Terminals
- Detect terminal capabilities on startup
- Graceful degradation to 256 colors
- Use ASCII fallbacks for fancy Unicode
- Skip image features if not supported

```go
package terminal

import "os"

type Capabilities struct {
    TrueColor     bool
    Images        bool
    Unicode       bool
    Name          string
}

func Detect() Capabilities {
    caps := Capabilities{
        Name: os.Getenv("TERM"),
    }
    
    // Check for Kitty
    if os.Getenv("TERM") == "xterm-kitty" {
        caps.TrueColor = true
        caps.Images = true
        caps.Unicode = true
        return caps
    }
    
    // Check COLORTERM for true color support
    if os.Getenv("COLORTERM") == "truecolor" || 
       os.Getenv("COLORTERM") == "24bit" {
        caps.TrueColor = true
    }
    
    caps.Unicode = true // Most modern terminals
    return caps
}
```

---

## TUI Architecture

### Bubbletea Model Structure

```go
type Model struct {
    // Core data
    typFile     *parser.TYPFile
    modified    bool
    
    // UI state
    mode        Mode  // list, edit, help, etc.
    width       int
    height      int
    
    // Components
    list        list.Model
    editor      *EditorModel
    colorPicker *ColorPickerModel
    statusBar   *StatusBarModel
    
    // Navigation
    activeTab   Tab   // points, lines, polygons
    selectedIdx int
    
    // Messages
    error       string
    info        string
}

type Mode int
const (
    ModeList Mode = iota
    ModeEdit
    ModeHelp
    ModeConfirm
)

type Tab int
const (
    TabPoints Tab = iota
    TabLines
    TabPolygons
)
```

### Update Logic Flow

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            if m.modified {
                return m, tea.Cmd(/* show unsaved warning */)
            }
            return m, tea.Quit
        case "tab":
            m.activeTab = (m.activeTab + 1) % 3
            return m, nil
        case "e":
            if m.mode == ModeList {
                m.mode = ModeEdit
                m.editor = NewEditorForSelectedType()
                return m, nil
            }
        // ... more cases
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }
    
    // Delegate to current mode
    return m.updateForCurrentMode(msg)
}
```

### View Rendering Strategy

```go
func (m Model) View() string {
    if m.width == 0 {
        return "Loading..."
    }
    
    var sections []string
    
    // Header
    sections = append(sections, m.renderHeader())
    
    // Main content (depends on mode)
    switch m.mode {
    case ModeList:
        sections = append(sections, m.renderListMode())
    case ModeEdit:
        sections = append(sections, m.renderEditMode())
    case ModeHelp:
        sections = append(sections, m.renderHelp())
    }
    
    // Footer
    sections = append(sections, m.renderFooter())
    
    return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
```

---

## TYP File Format Reference

### Example TYP File Structure

```
[_id]
CodePage=1252
FID=1234
ProductCode=1
[end]

[_drawOrder]
Type=0x01,1
Type=0x02,2
[end]

[_point]
Type=0x2f06
String1=0x04,Bank
String2=0x01,Banque
IconXpm="16 16 2 1"
"! c #778899"
"  c none"
"!!!!!!!!!!!!!!!!"
"!              !"
; ... more rows
[end]

[_line]
Type=0x01
String=0x04,Highway
LineWidth=5
BorderWidth=1
Xpm="0 0 2 0"
"! c #FF0000"
"  c #000000"
[end]

[_polygon]
Type=0x13
String=0x04,Park
Xpm="32 32 2 1"
"a c #90EE90"
"b c none"
; ... pattern
[end]
```

### Parsing Challenges
1. **Comments**: Lines starting with `;` or `#`
2. **Multi-line strings**: XPM data spans multiple lines
3. **String formats**: Various quote styles
4. **Optional fields**: Many fields are optional
5. **Language codes**: String1=0x04 (English), etc.

---

## Configuration & XDG Compliance

### XDG Base Directory Specification

The application follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) for storing configuration and data files.

**Configuration Files:**
```
$XDG_CONFIG_HOME/typtui/config.yaml
(fallback: ~/.config/typtui/config.yaml)
```

**Data Files (recent files, cache):**
```
$XDG_DATA_HOME/typtui/recent.json
(fallback: ~/.local/share/typtui/recent.json)
```

**Cache (parsed TYP data, etc.):**
```
$XDG_CACHE_HOME/typtui/
(fallback: ~/.cache/typtui/)
```

### Configuration File Example

```yaml
# ~/.config/typtui/config.yaml
mkgmap:
  path: "/usr/local/bin/mkgmap"
  default_args:
    - "--family-id=1234"

editor:
  default_language: "en"
  auto_save: true
  auto_compile: false
  
colors:
  theme: "default"  # or "dark", "light"
  
recent_files:
  max_count: 10

terminal:
  mouse_enabled: false
```

### Implementation

```go
package config

import (
    "github.com/adrg/xdg"
    "gopkg.in/yaml.v3"
    "path/filepath"
)

type Config struct {
    Mkgmap struct {
        Path        string   `yaml:"path"`
        DefaultArgs []string `yaml:"default_args"`
    } `yaml:"mkgmap"`
    Editor struct {
        DefaultLanguage string `yaml:"default_language"`
        AutoSave        bool   `yaml:"auto_save"`
        AutoCompile     bool   `yaml:"auto_compile"`
    } `yaml:"editor"`
    // ... more fields
}

func Load() (*Config, error) {
    // Try XDG_CONFIG_HOME first
    configPath, err := xdg.ConfigFile("typtui/config.yaml")
    if err != nil {
        return DefaultConfig(), nil
    }
    
    // Load and parse config
    // ...
}

func (c *Config) Save() error {
    configPath, err := xdg.ConfigFile("typtui/config.yaml")
    if err != nil {
        return err
    }
    // Save config
    // ...
}
```

---

## External Dependencies

### Required Tools
- **mkgmap**: For TYP compilation
  - Target: Latest stable version (r4922+ as of 2024)
  - Auto-detect in PATH
  - Configurable path via XDG config
  - Version compatibility check on startup
  
### Go Dependencies
```go
require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/adrg/xdg v0.4.0  // For XDG Base Directory support
    gopkg.in/yaml.v3 v3.0.1  // For config files
)
```

---

## Testing Plan

### Unit Tests
- Parser for each section type
- Color conversion utilities
- Validation functions
- Garmin type database lookups

### Integration Tests
- Load real TYP files (OpenHiking, etc.)
- Parse -> Modify -> Save -> Parse again
- mkgmap compilation of generated files
- Round-trip testing (input == output)

### Test Data Sources
- OpenHiking TYP files
- OpenMTBMap TYP files
- mkgmap example TYP files
- Custom test cases for edge cases

### Manual Testing Checklist
- [ ] Load various TYP files without crashes
- [ ] Edit and save changes correctly
- [ ] Compilation with mkgmap succeeds
- [ ] Terminal resize handling
- [ ] Keyboard navigation feels natural
- [ ] Error messages are helpful
- [ ] Performance with large files (1000+ types)

---

## Distribution & Packaging

### Release Artifacts (Phase 1 - Linux Only)
- Source code (GitHub releases)
- Pre-built binaries (Linux x86_64, ARM64)
- AUR package (Arch Linux)
- Debian/Ubuntu package (.deb)
- Future: RPM, Flatpak, other platforms

### Installation Methods
```bash
# From source
go install github.com/yourusername/typtui@latest

# From releases
curl -LO https://github.com/yourusername/typtui/releases/latest/download/typtui-linux-amd64
chmod +x typtui-linux-amd64
sudo mv typtui-linux-amd64 /usr/local/bin/typtui

# Via package manager (future)
yay -S typtui-bin  # AUR
```

---

## Success Metrics

### Technical Metrics
- Parse 99%+ of real-world TYP files
- <100ms startup time
- <50MB memory usage
- Zero crashes during testing
- 80%+ code coverage

### User Experience Metrics
- Can complete basic edit in <2 minutes
- Help system answers common questions
- Error messages guide users to solutions
- Compilation success rate >95%

---

## Risk Mitigation

### Potential Risks

1. **TYP Format Complexity**
   - Risk: Format has undocumented edge cases
   - Mitigation: Test with many real files, add verbose error logging

2. **mkgmap Compatibility**
   - Risk: mkgmap versions have different requirements
   - Mitigation: Test with multiple versions, document requirements

3. **Terminal Compatibility**
   - Risk: Doesn't work in all terminals
   - Mitigation: Initially optimize for Kitty terminal, test in other common terminals (gnome-terminal, alacritty, etc.) as secondary priority
   - Initial focus: Kitty on Linux, other terminals best-effort

4. **XPM Editing Complexity**
   - Risk: Inline XPM editor is very hard in TUI
   - Mitigation: Phase 1 just imports/exports, inline editor is future work

5. **Color Representation**
   - Risk: Terminal colors don't match hex accurately
   - Mitigation: Show both terminal approximation and hex value

---

## Community Engagement

### Target Communities
- OpenHiking forums/mailing list
- mkgmap users
- OpenStreetMap community
- /r/openstreetmap
- Linux hiking/GPS communities

### Launch Strategy
1. Soft launch: Post to mkgmap mailing list for feedback
2. Beta testing with OpenHiking contributors
3. Public release: GitHub, Reddit, OpenStreetMap forums
4. Article: "Customizing Garmin Maps on Linux"
5. Video demo on YouTube

---

## Timeline Summary

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| 1. Foundation | 2 weeks | Parser + Basic TUI |
| 2. List Mode | 1 week | Type browser |
| 3. Edit Mode | 2 weeks | Property editing |
| 4. Save & Compile | 1 week | Write + mkgmap integration |
| 5. Advanced Features | 2 weeks | Full feature set |
| 6. Polish | 1 week | v1.0 release |
| **Total** | **9 weeks** | **Production ready** |

---

## Next Steps

1. Set up project repository
2. Initialize Go module and dependencies
3. Create sample TYP files in testdata/
4. Start with parser implementation
5. Set up basic TUI scaffold
6. Iterate from there!

---

## Resources

### Documentation
- [mkgmap TYP File Documentation](https://www.mkgmap.org.uk/doc/typ-compiler)
- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss)
- [OSM Garmin Type Codes](https://wiki.openstreetmap.org/wiki/OSM_Map_On_Garmin/POI_Types)

### Inspiration
- [lazygit](https://github.com/jesseduffield/lazygit) - Great TUI for Git
- [k9s](https://github.com/derailed/k9s) - Kubernetes TUI
- [glow](https://github.com/charmbracelet/glow) - Markdown TUI

### Sample TYP Files
- OpenHiking: https://openhiking.eu/en/downloads/garmin-maps
- mkgmap examples: In mkgmap distribution
- OpenMTBMap: https://openmtbmap.org/

---

## License

Recommend: **MIT License** (permissive, widely adopted, OSM-compatible)

---

## Contact & Contribution

- GitHub Issues for bug reports
- Discussions for feature requests
- Pull requests welcome
- Code of Conduct: Contributor Covenant

---

**Ready to start coding!** 🚀

Begin with Phase 1, focus on getting the parser working with real TYP files, then build the TUI layer on top. The foundation is the most important part - once you can reliably parse and write TYP files, the rest is "just" UI work.

Good luck with the implementation! This could genuinely be a valuable tool for the OpenStreetMap and Garmin mapping community.
