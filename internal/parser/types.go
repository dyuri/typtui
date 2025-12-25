package parser

// TYPFile represents the entire TYP file structure
type TYPFile struct {
	Header   Header
	Points   []PointType
	Lines    []LineType
	Polygons []PolygonType
	DrawOrder DrawOrder
	FilePath string
	Modified bool
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
	Type        string            // e.g., "0x2f06"
	SubType     string            // Optional subtype
	Labels      map[string]string // Language code -> label
	DayXpm      *XPMIcon
	NightXpm    *XPMIcon
	DayColors   []Color
	NightColors []Color
	FontStyle   string
}

// LineType represents a line definition (roads, trails, etc.)
type LineType struct {
	Type              string
	Labels            map[string]string
	LineWidth         int
	BorderWidth       int
	DayXpm            *XPMIcon
	NightXpm          *XPMIcon
	UseOrientation    bool
	LineStyle         string // "solid", "dashed", etc.
}

// PolygonType represents an area definition
type PolygonType struct {
	Type           string
	Labels         map[string]string
	DayXpm         *XPMIcon
	NightXpm       *XPMIcon
	FontStyle      string
	ExtendedLabels bool
}

// Color represents a color in hex format
type Color struct {
	Hex  string // "#RRGGBB"
	Day  bool   // true if day color, false if night
	Name string // Optional color name
}

// XPMIcon represents icon/pattern data in XPM format
type XPMIcon struct {
	Width   int
	Height  int
	Colors  int
	CharsPerPixel int
	Data    []string
	Palette map[string]Color
}

// DrawOrder specifies rendering order
type DrawOrder struct {
	Points   []string
	Lines    []string
	Polygons []string
}

// ParseError represents a parsing error with location information
type ParseError struct {
	Line    int
	Column  int
	Message string
	File    string
}

func (e *ParseError) Error() string {
	if e.File != "" {
		return e.File + ":" + string(rune(e.Line)) + ":" + string(rune(e.Column)) + ": " + e.Message
	}
	return "line " + string(rune(e.Line)) + ", col " + string(rune(e.Column)) + ": " + e.Message
}
