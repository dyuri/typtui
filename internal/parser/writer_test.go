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
