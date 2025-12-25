package parser

import (
	"testing"
)

func TestParseMinimalFile(t *testing.T) {
	typFile, err := ParseFile("../../testdata/sample/minimal.typ")
	if err != nil {
		t.Fatalf("Failed to parse minimal.typ: %v", err)
	}

	if typFile.Header.CodePage != 1252 {
		t.Errorf("Expected CodePage 1252, got %d", typFile.Header.CodePage)
	}

	if typFile.Header.FID != 1234 {
		t.Errorf("Expected FID 1234, got %d", typFile.Header.FID)
	}

	if typFile.Header.ProductCode != 1 {
		t.Errorf("Expected ProductCode 1, got %d", typFile.Header.ProductCode)
	}

	if len(typFile.Points) != 1 {
		t.Fatalf("Expected 1 point, got %d", len(typFile.Points))
	}

	point := typFile.Points[0]
	if point.Type != "0x2f06" {
		t.Errorf("Expected point type 0x2f06, got %s", point.Type)
	}

	if label, ok := point.Labels["0x04"]; !ok || label != "Bank" {
		t.Errorf("Expected label 'Bank' for language 0x04, got %s", label)
	}
}

func TestParseBasicFile(t *testing.T) {
	typFile, err := ParseFile("../../testdata/sample/basic.typ")
	if err != nil {
		t.Fatalf("Failed to parse basic.typ: %v", err)
	}

	// Check header
	if typFile.Header.CodePage != 1252 {
		t.Errorf("Expected CodePage 1252, got %d", typFile.Header.CodePage)
	}

	// Check points
	if len(typFile.Points) != 1 {
		t.Fatalf("Expected 1 point, got %d", len(typFile.Points))
	}

	point := typFile.Points[0]
	if point.Type != "0x2f06" {
		t.Errorf("Expected point type 0x2f06, got %s", point.Type)
	}

	if len(point.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(point.Labels))
	}

	if point.DayXpm == nil {
		t.Error("Expected DayXpm to be parsed")
	} else {
		if point.DayXpm.Width != 8 {
			t.Errorf("Expected XPM width 8, got %d", point.DayXpm.Width)
		}
		if point.DayXpm.Height != 8 {
			t.Errorf("Expected XPM height 8, got %d", point.DayXpm.Height)
		}
	}

	// Check lines
	if len(typFile.Lines) != 1 {
		t.Fatalf("Expected 1 line, got %d", len(typFile.Lines))
	}

	line := typFile.Lines[0]
	if line.Type != "0x01" {
		t.Errorf("Expected line type 0x01, got %s", line.Type)
	}

	if line.LineWidth != 4 {
		t.Errorf("Expected line width 4, got %d", line.LineWidth)
	}

	if line.BorderWidth != 1 {
		t.Errorf("Expected border width 1, got %d", line.BorderWidth)
	}

	if line.LineStyle != "solid" {
		t.Errorf("Expected line style 'solid', got '%s'", line.LineStyle)
	}

	// Check polygons
	if len(typFile.Polygons) != 1 {
		t.Fatalf("Expected 1 polygon, got %d", len(typFile.Polygons))
	}

	polygon := typFile.Polygons[0]
	if polygon.Type != "0x13" {
		t.Errorf("Expected polygon type 0x13, got %s", polygon.Type)
	}

	if !polygon.ExtendedLabels {
		t.Error("Expected ExtendedLabels to be true")
	}

	if polygon.DayXpm == nil {
		t.Error("Expected polygon XPM to be parsed")
	} else {
		if polygon.DayXpm.Width != 32 {
			t.Errorf("Expected polygon XPM width 32, got %d", polygon.DayXpm.Width)
		}
	}
}

func TestParseString(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		input    string
		wantLang string
		wantLabel string
	}{
		{"0x04,Bank", "0x04", "Bank"},
		{"0x01,Banque", "0x01", "Banque"},
		{"0x04,City Hall", "0x04", "City Hall"},
		{"invalid", "", ""},
	}

	for _, tt := range tests {
		lang, label := p.parseString(tt.input)
		if lang != tt.wantLang || label != tt.wantLabel {
			t.Errorf("parseString(%q) = (%q, %q), want (%q, %q)",
				tt.input, lang, label, tt.wantLang, tt.wantLabel)
		}
	}
}

func TestCleanLine(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		input string
		want  string
	}{
		{"; This is a comment", ""},
		{"Type=0x2f06 ; with comment", "Type=0x2f06"},
		{"  Type=0x2f06  ", "Type=0x2f06"},
		{"# Another comment", ""},
		{"Type=0x2f06", "Type=0x2f06"},
	}

	for _, tt := range tests {
		got := p.cleanLine(tt.input)
		if got != tt.want {
			t.Errorf("cleanLine(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
