package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parser handles TYP file parsing
type Parser struct {
	scanner  *bufio.Scanner
	lineNum  int
	filePath string
}

// NewParser creates a new parser for the given file
func NewParser(filePath string) (*Parser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	return &Parser{
		scanner:  scanner,
		lineNum:  0,
		filePath: filePath,
	}, nil
}

// Parse parses the TYP file and returns a TYPFile structure
func (p *Parser) Parse() (*TYPFile, error) {
	typFile := &TYPFile{
		FilePath: p.filePath,
		Modified: false,
	}

	for p.scanner.Scan() {
		p.lineNum++
		line := p.scanner.Text()
		line = p.cleanLine(line)

		if line == "" {
			continue
		}

		// Check for section markers
		if strings.HasPrefix(line, "[") {
			section := strings.TrimSpace(strings.Trim(line, "[]"))

			switch section {
			case "_id":
				if err := p.parseHeader(&typFile.Header); err != nil {
					return nil, err
				}
			case "_point":
				point, err := p.parsePoint()
				if err != nil {
					return nil, err
				}
				typFile.Points = append(typFile.Points, *point)
			case "_line":
				line, err := p.parseLine()
				if err != nil {
					return nil, err
				}
				typFile.Lines = append(typFile.Lines, *line)
			case "_polygon":
				polygon, err := p.parsePolygon()
				if err != nil {
					return nil, err
				}
				typFile.Polygons = append(typFile.Polygons, *polygon)
			case "_drawOrder":
				if err := p.parseDrawOrder(&typFile.DrawOrder); err != nil {
					return nil, err
				}
			case "end":
				// Section end marker - ignore
			default:
				// Unknown section - skip it
				if err := p.skipToEnd(); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return typFile, nil
}

// cleanLine removes comments and trims whitespace
func (p *Parser) cleanLine(line string) string {
	// Remove comments (lines starting with ; or #)
	if idx := strings.Index(line, ";"); idx >= 0 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "#"); idx >= 0 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}

// parseHeader parses the [_id] section
func (p *Parser) parseHeader(header *Header) error {
	for p.scanner.Scan() {
		p.lineNum++
		line := p.cleanLine(p.scanner.Text())

		if line == "" {
			continue
		}

		if line == "[end]" {
			return nil
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		var err error
		switch key {
		case "CodePage":
			header.CodePage, err = strconv.Atoi(value)
		case "FID":
			header.FID, err = strconv.Atoi(value)
		case "ProductCode":
			header.ProductCode, err = strconv.Atoi(value)
		case "MapID":
			header.MapID, err = strconv.Atoi(value)
		}

		if err != nil {
			return &ParseError{
				Line:    p.lineNum,
				Message: fmt.Sprintf("invalid value for %s: %s", key, value),
				File:    p.filePath,
			}
		}
	}

	return &ParseError{
		Line:    p.lineNum,
		Message: "unexpected end of file in header section",
		File:    p.filePath,
	}
}

// parsePoint parses a [_point] section
func (p *Parser) parsePoint() (*PointType, error) {
	point := &PointType{
		Labels: make(map[string]string),
	}

	for p.scanner.Scan() {
		p.lineNum++
		line := p.cleanLine(p.scanner.Text())

		if line == "" {
			continue
		}

		if line == "[end]" {
			return point, nil
		}

		if err := p.parsePointProperty(point, line); err != nil {
			return nil, err
		}
	}

	return nil, &ParseError{
		Line:    p.lineNum,
		Message: "unexpected end of file in point section",
		File:    p.filePath,
	}
}

// parsePointProperty parses a single property line in a point definition
func (p *Parser) parsePointProperty(point *PointType, line string) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return nil // Skip malformed lines
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Type":
		point.Type = value
	case "SubType":
		point.SubType = value
	case "String", "String1", "String2", "String3", "String4":
		langCode, label := p.parseString(value)
		if langCode != "" {
			point.Labels[langCode] = label
		}
	case "DayXpm":
		xpm, err := p.parseXPM(value)
		if err != nil {
			return err
		}
		point.DayXpm = xpm
	case "NightXpm":
		xpm, err := p.parseXPM(value)
		if err != nil {
			return err
		}
		point.NightXpm = xpm
	case "FontStyle":
		point.FontStyle = value
	}

	return nil
}

// parseLine parses a [_line] section
func (p *Parser) parseLine() (*LineType, error) {
	line := &LineType{
		Labels: make(map[string]string),
	}

	for p.scanner.Scan() {
		p.lineNum++
		textLine := p.cleanLine(p.scanner.Text())

		if textLine == "" {
			continue
		}

		if textLine == "[end]" {
			return line, nil
		}

		if err := p.parseLineProperty(line, textLine); err != nil {
			return nil, err
		}
	}

	return nil, &ParseError{
		Line:    p.lineNum,
		Message: "unexpected end of file in line section",
		File:    p.filePath,
	}
}

// parseLineProperty parses a single property line in a line definition
func (p *Parser) parseLineProperty(line *LineType, textLine string) error {
	parts := strings.SplitN(textLine, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	var err error
	switch key {
	case "Type":
		line.Type = value
	case "String", "String1", "String2", "String3", "String4":
		langCode, label := p.parseString(value)
		if langCode != "" {
			line.Labels[langCode] = label
		}
	case "LineWidth":
		line.LineWidth, err = strconv.Atoi(value)
	case "BorderWidth":
		line.BorderWidth, err = strconv.Atoi(value)
	case "LineStyle":
		line.LineStyle = value
	case "Xpm":
		line.DayXpm, err = p.parseXPM(value)
	case "UseOrientation":
		line.UseOrientation = strings.ToUpper(value) == "Y" || value == "1"
	}

	if err != nil {
		return &ParseError{
			Line:    p.lineNum,
			Message: fmt.Sprintf("invalid value for %s: %s", key, value),
			File:    p.filePath,
		}
	}

	return nil
}

// parsePolygon parses a [_polygon] section
func (p *Parser) parsePolygon() (*PolygonType, error) {
	polygon := &PolygonType{
		Labels: make(map[string]string),
	}

	for p.scanner.Scan() {
		p.lineNum++
		line := p.cleanLine(p.scanner.Text())

		if line == "" {
			continue
		}

		if line == "[end]" {
			return polygon, nil
		}

		if err := p.parsePolygonProperty(polygon, line); err != nil {
			return nil, err
		}
	}

	return nil, &ParseError{
		Line:    p.lineNum,
		Message: "unexpected end of file in polygon section",
		File:    p.filePath,
	}
}

// parsePolygonProperty parses a single property line in a polygon definition
func (p *Parser) parsePolygonProperty(polygon *PolygonType, line string) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "Type":
		polygon.Type = value
	case "String", "String1", "String2", "String3", "String4":
		langCode, label := p.parseString(value)
		if langCode != "" {
			polygon.Labels[langCode] = label
		}
	case "Xpm":
		xpm, err := p.parseXPM(value)
		if err != nil {
			return err
		}
		polygon.DayXpm = xpm
	case "ExtendedLabels":
		polygon.ExtendedLabels = strings.ToUpper(value) == "Y" || value == "1"
	case "FontStyle":
		polygon.FontStyle = value
	}

	return nil
}

// parseDrawOrder parses the [_drawOrder] section
func (p *Parser) parseDrawOrder(drawOrder *DrawOrder) error {
	for p.scanner.Scan() {
		p.lineNum++
		line := p.cleanLine(p.scanner.Text())

		if line == "" {
			continue
		}

		if line == "[end]" {
			return nil
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Type" {
			continue
		}

		// Draw order types are listed as Type=0x01,1
		// We'll just store the type code for now
		typeCode := strings.Split(strings.TrimSpace(parts[1]), ",")[0]
		// For now, just add to points (we'd need more logic to determine category)
		drawOrder.Points = append(drawOrder.Points, typeCode)
	}

	return nil
}

// parseString parses a string definition like "0x04,Bank"
func (p *Parser) parseString(value string) (langCode, label string) {
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 {
		return "", ""
	}

	langCode = strings.TrimSpace(parts[0])
	label = strings.TrimSpace(parts[1])
	return
}

// parseXPM parses an XPM definition (simplified for now)
func (p *Parser) parseXPM(value string) (*XPMIcon, error) {
	// XPM format: "width height colors chars_per_pixel"
	value = strings.Trim(value, "\"")
	parts := strings.Fields(value)

	if len(parts) < 4 {
		return nil, &ParseError{
			Line:    p.lineNum,
			Message: "invalid XPM format",
			File:    p.filePath,
		}
	}

	xpm := &XPMIcon{
		Palette: make(map[string]Color),
		Data:    make([]string, 0),
	}

	var err error
	xpm.Width, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	xpm.Height, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	xpm.Colors, err = strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	xpm.CharsPerPixel, err = strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	// Parse color palette and data lines
	// We need to read exactly xpm.Colors + xpm.Height lines
	totalLinesToRead := xpm.Colors + xpm.Height
	linesRead := 0

	for linesRead < totalLinesToRead && p.scanner.Scan() {
		p.lineNum++
		line := p.scanner.Text()

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), ";") {
			continue
		}

		// Check if this is a quoted line (XPM data)
		if !strings.HasPrefix(strings.TrimSpace(line), "\"") {
			// Not an XPM line, we're done
			break
		}

		line = strings.Trim(strings.TrimSpace(line), "\"")

		if linesRead < xpm.Colors {
			// This is a color definition line
			// Format: "! c #778899" or "  c none"
			if strings.Contains(line, " c ") {
				parts := strings.SplitN(line, " c ", 2)
				if len(parts) == 2 {
					char := parts[0]
					colorValue := strings.TrimSpace(parts[1])
					xpm.Palette[char] = Color{Hex: colorValue}
				}
			}
		} else {
			// This is pixel data
			xpm.Data = append(xpm.Data, line)
		}

		linesRead++
	}

	return xpm, nil
}

// skipToEnd skips lines until [end] is found
func (p *Parser) skipToEnd() error {
	for p.scanner.Scan() {
		p.lineNum++
		line := p.cleanLine(p.scanner.Text())
		if line == "[end]" {
			return nil
		}
	}
	return &ParseError{
		Line:    p.lineNum,
		Message: "unexpected end of file while skipping section",
		File:    p.filePath,
	}
}

// ParseFile is a convenience function to parse a TYP file
func ParseFile(filePath string) (*TYPFile, error) {
	parser, err := NewParser(filePath)
	if err != nil {
		return nil, err
	}
	return parser.Parse()
}
