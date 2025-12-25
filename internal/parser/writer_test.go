package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndReload(t *testing.T) {
	// Load the test file
	typFile, err := ParseFile("../../testdata/sample/basic.typ")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Modify something
	if len(typFile.Points) > 0 {
		typFile.Points[0].Labels["0x04"] = "Modified Bank"
	}

	// Write to a temp file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.typ")

	err = WriteFile(typFile, tempFile)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Reload and verify
	reloaded, err := ParseFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to reload file: %v", err)
	}

	// Verify the modification persisted
	if len(reloaded.Points) == 0 {
		t.Fatal("No points in reloaded file")
	}

	if reloaded.Points[0].Labels["0x04"] != "Modified Bank" {
		t.Errorf("Expected 'Modified Bank', got '%s'", reloaded.Points[0].Labels["0x04"])
	}

	// Verify other data is intact
	if reloaded.Header.CodePage != typFile.Header.CodePage {
		t.Errorf("Header CodePage mismatch: expected %d, got %d", typFile.Header.CodePage, reloaded.Header.CodePage)
	}

	if len(reloaded.Lines) != len(typFile.Lines) {
		t.Errorf("Line count mismatch: expected %d, got %d", len(typFile.Lines), len(reloaded.Lines))
	}

	if len(reloaded.Polygons) != len(typFile.Polygons) {
		t.Errorf("Polygon count mismatch: expected %d, got %d", len(typFile.Polygons), len(reloaded.Polygons))
	}
}

func TestWriteHeader(t *testing.T) {
	typFile := &TYPFile{
		Header: Header{
			CodePage:    1252,
			FID:         1234,
			ProductCode: 1,
		},
	}

	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "header_test.typ")

	err := WriteFile(typFile, tempFile)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Read the file content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)

	// Check for expected content
	expectedStrings := []string{
		"[_id]",
		"CodePage=1252",
		"FID=1234",
		"ProductCode=1",
		"[end]",
	}

	for _, expected := range expectedStrings {
		if !contains(contentStr, expected) {
			t.Errorf("Expected content to contain '%s'", expected)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func TestWriteLineProperties(t *testing.T) {
	// Load the test file
	typFile, err := ParseFile("../../testdata/sample/basic.typ")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Modify line properties
	if len(typFile.Lines) > 0 {
		typFile.Lines[0].LineWidth = 10
		typFile.Lines[0].BorderWidth = 2
		typFile.Lines[0].LineStyle = "dashed"
		typFile.Lines[0].UseOrientation = true
	}

	// Write to a temp file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "line_test.typ")

	err = WriteFile(typFile, tempFile)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Reload and verify
	reloaded, err := ParseFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to reload file: %v", err)
	}

	if len(reloaded.Lines) == 0 {
		t.Fatal("No lines in reloaded file")
	}

	line := reloaded.Lines[0]
	if line.LineWidth != 10 {
		t.Errorf("Expected LineWidth 10, got %d", line.LineWidth)
	}
	if line.BorderWidth != 2 {
		t.Errorf("Expected BorderWidth 2, got %d", line.BorderWidth)
	}
	if line.LineStyle != "dashed" {
		t.Errorf("Expected LineStyle 'dashed', got '%s'", line.LineStyle)
	}
	if !line.UseOrientation {
		t.Error("Expected UseOrientation to be true")
	}
}

func TestWritePointSubTypeAndFontStyle(t *testing.T) {
	// Load the test file
	typFile, err := ParseFile("../../testdata/sample/basic.typ")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Modify point properties
	if len(typFile.Points) > 0 {
		typFile.Points[0].SubType = "0x01"
		typFile.Points[0].FontStyle = "LargeFont"
	}

	// Write to a temp file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "point_test.typ")

	err = WriteFile(typFile, tempFile)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Reload and verify
	reloaded, err := ParseFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to reload file: %v", err)
	}

	if len(reloaded.Points) == 0 {
		t.Fatal("No points in reloaded file")
	}

	point := reloaded.Points[0]
	if point.SubType != "0x01" {
		t.Errorf("Expected SubType '0x01', got '%s'", point.SubType)
	}
	if point.FontStyle != "LargeFont" {
		t.Errorf("Expected FontStyle 'LargeFont', got '%s'", point.FontStyle)
	}
}

func TestWritePolygonExtendedLabels(t *testing.T) {
	// Load the test file
	typFile, err := ParseFile("../../testdata/sample/basic.typ")
	if err != nil {
		t.Fatalf("Failed to parse file: %v", err)
	}

	// Modify polygon properties
	if len(typFile.Polygons) > 0 {
		typFile.Polygons[0].ExtendedLabels = true
		typFile.Polygons[0].FontStyle = "SmallFont"
	}

	// Write to a temp file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "polygon_test.typ")

	err = WriteFile(typFile, tempFile)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Read the file content to verify ExtendedLabels=Y is written
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "ExtendedLabels=Y") {
		t.Error("Expected file to contain 'ExtendedLabels=Y'")
	}

	// Reload and verify
	reloaded, err := ParseFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to reload file: %v", err)
	}

	if len(reloaded.Polygons) == 0 {
		t.Fatal("No polygons in reloaded file")
	}

	polygon := reloaded.Polygons[0]
	if !polygon.ExtendedLabels {
		t.Error("Expected ExtendedLabels to be true")
	}
	if polygon.FontStyle != "SmallFont" {
		t.Errorf("Expected FontStyle 'SmallFont', got '%s'", polygon.FontStyle)
	}
}
