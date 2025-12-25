package parser

import (
	"fmt"
	"os"
	"strings"
)

// WriteFile writes a TYPFile to disk in TYP text format
func WriteFile(typFile *TYPFile, filePath string) error {
	var b strings.Builder

	// Write header section
	if err := writeHeader(&b, typFile.Header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write point types
	for _, point := range typFile.Points {
		if err := writePointType(&b, point); err != nil {
			return fmt.Errorf("failed to write point type %s: %w", point.Type, err)
		}
	}

	// Write line types
	for _, line := range typFile.Lines {
		if err := writeLineType(&b, line); err != nil {
			return fmt.Errorf("failed to write line type %s: %w", line.Type, err)
		}
	}

	// Write polygon types
	for _, polygon := range typFile.Polygons {
		if err := writePolygonType(&b, polygon); err != nil {
			return fmt.Errorf("failed to write polygon type %s: %w", polygon.Type, err)
		}
	}

	// Write to file
	content := b.String()
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// writeHeader writes the [_id] header section
func writeHeader(b *strings.Builder, header Header) error {
	b.WriteString("[_id]\n")

	if header.CodePage > 0 {
		b.WriteString(fmt.Sprintf("CodePage=%d\n", header.CodePage))
	}
	if header.FID > 0 {
		b.WriteString(fmt.Sprintf("FID=%d\n", header.FID))
	}
	if header.ProductCode > 0 {
		b.WriteString(fmt.Sprintf("ProductCode=%d\n", header.ProductCode))
	}
	if header.MapID > 0 {
		b.WriteString(fmt.Sprintf("MapID=%d\n", header.MapID))
	}

	b.WriteString("[end]\n\n")
	return nil
}

// writePointType writes a [_point] section
func writePointType(b *strings.Builder, point PointType) error {
	b.WriteString("[_point]\n")

	// Type (required)
	b.WriteString(fmt.Sprintf("Type=%s\n", point.Type))

	// SubType (optional)
	if point.SubType != "" {
		b.WriteString(fmt.Sprintf("SubType=%s\n", point.SubType))
	}

	// Labels
	if err := writeLabels(b, point.Labels); err != nil {
		return err
	}

	// DayXpm
	if point.DayXpm != nil {
		if err := writeXPM(b, "DayXpm", point.DayXpm); err != nil {
			return err
		}
	}

	// NightXpm
	if point.NightXpm != nil {
		if err := writeXPM(b, "NightXpm", point.NightXpm); err != nil {
			return err
		}
	}

	// FontStyle
	if point.FontStyle != "" {
		b.WriteString(fmt.Sprintf("FontStyle=%s\n", point.FontStyle))
	}

	b.WriteString("[end]\n\n")
	return nil
}

// writeLineType writes a [_line] section
func writeLineType(b *strings.Builder, line LineType) error {
	b.WriteString("[_line]\n")

	// Type (required)
	b.WriteString(fmt.Sprintf("Type=%s\n", line.Type))

	// Labels
	if err := writeLabels(b, line.Labels); err != nil {
		return err
	}

	// LineWidth
	if line.LineWidth > 0 {
		b.WriteString(fmt.Sprintf("LineWidth=%d\n", line.LineWidth))
	}

	// BorderWidth
	if line.BorderWidth > 0 {
		b.WriteString(fmt.Sprintf("BorderWidth=%d\n", line.BorderWidth))
	}

	// LineStyle
	if line.LineStyle != "" {
		b.WriteString(fmt.Sprintf("LineStyle=%s\n", line.LineStyle))
	}

	// UseOrientation
	if line.UseOrientation {
		b.WriteString("UseOrientation=Y\n")
	}

	// DayXpm
	if line.DayXpm != nil {
		if err := writeXPM(b, "Xpm", line.DayXpm); err != nil {
			return err
		}
	}

	// NightXpm
	if line.NightXpm != nil {
		if err := writeXPM(b, "NightXpm", line.NightXpm); err != nil {
			return err
		}
	}

	b.WriteString("[end]\n\n")
	return nil
}

// writePolygonType writes a [_polygon] section
func writePolygonType(b *strings.Builder, polygon PolygonType) error {
	b.WriteString("[_polygon]\n")

	// Type (required)
	b.WriteString(fmt.Sprintf("Type=%s\n", polygon.Type))

	// Labels
	if err := writeLabels(b, polygon.Labels); err != nil {
		return err
	}

	// ExtendedLabels
	if polygon.ExtendedLabels {
		b.WriteString("ExtendedLabels=Y\n")
	}

	// FontStyle
	if polygon.FontStyle != "" {
		b.WriteString(fmt.Sprintf("FontStyle=%s\n", polygon.FontStyle))
	}

	// DayXpm
	if polygon.DayXpm != nil {
		if err := writeXPM(b, "Xpm", polygon.DayXpm); err != nil {
			return err
		}
	}

	// NightXpm
	if polygon.NightXpm != nil {
		if err := writeXPM(b, "NightXpm", polygon.NightXpm); err != nil {
			return err
		}
	}

	b.WriteString("[end]\n\n")
	return nil
}

// writeLabels writes String= lines for all labels
func writeLabels(b *strings.Builder, labels map[string]string) error {
	if len(labels) == 0 {
		return nil
	}

	// Write labels in a consistent order (English first, then sorted)
	if label, ok := labels["0x04"]; ok {
		b.WriteString(fmt.Sprintf("String=0x04,%s\n", label))
	}

	// Write other labels
	for code, label := range labels {
		if code != "0x04" {
			b.WriteString(fmt.Sprintf("String=%s,%s\n", code, label))
		}
	}

	return nil
}

// writeXPM writes an XPM icon/pattern
func writeXPM(b *strings.Builder, fieldName string, xpm *XPMIcon) error {
	// Write XPM header: "width height numColors charsPerPixel"
	b.WriteString(fmt.Sprintf("%s=\"%d %d %d %d\"\n",
		fieldName, xpm.Width, xpm.Height, xpm.Colors, xpm.CharsPerPixel))

	// Write color palette
	for char, color := range xpm.Palette {
		if color.Hex == "none" || color.Hex == "transparent" {
			b.WriteString(fmt.Sprintf("\"%s c none\"\n", char))
		} else {
			b.WriteString(fmt.Sprintf("\"%s c %s\"\n", char, color.Hex))
		}
	}

	// Write pixel data
	for _, line := range xpm.Data {
		b.WriteString(fmt.Sprintf("\"%s\"\n", line))
	}

	return nil
}
