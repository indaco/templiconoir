package templiconoir

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// 1. Core Tests for Icon Methods
// These tests cover methods like `String`, `SetSize`, and `SetAttrs`.

func TestIcon_default(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Icon
		expected string
	}{
		{
			name: "Outline icon with default attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap",
					Size: "24",
					Type: "Outline",
				}
				icon.body = `<path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><path d="M4.26 10.147a60 60 0 0 0-.491 6.347A48.6 48.6 0 0 1 12 20.904a48.6 48.6 0 0 1 8.232-4.41a61 61 0 0 0-.491-6.347z"/></svg>`,
		},
		{
			name: "Solid icon with default attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap-solid",
					Size: "24",
					Type: "Solid",
				}
				icon.body = `<path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><path d="M12 20a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Mini icon with attributes",
			setup: func() *Icon {
				icon := &Icon{
					Name: "academic-cap-mini",
					Size: "20",
					Type: "Mini",
					Attrs: templ.Attributes{
						"focusable": "false",
					},
				}
				icon.body = `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`,
		},
		{
			name: "Micro icon with stroke-width and color attributes",
			setup: func() *Icon {
				return &Icon{
					Name:        "micro-icon",
					Type:        "Micro",
					Size:        "16",
					StrokeWidth: "2",
					Color:       "#2dd4bf",
					Attrs: templ.Attributes{
						"aria-hidden": "true",
						"class":       "icon-micro",
					},
					body: `<path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/>`,
				}
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="2" color="#2dd4bf" aria-hidden="true" class="icon-micro"><path d="M8 16a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"/></svg>`,
		},
		{
			name: "Fallback case",
			setup: func() *Icon {
				icon := &Icon{
					Name: "unknown-icon",
					Size: "24",
					Type: "Unknown",
				}
				icon.body = `<circle cx="12" cy="12" r="10"/>`
				return icon
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><circle cx="12" cy="12" r="10"/></svg>`,
		},
		{
			name: "SetSize modifies size",
			setup: func() *Icon {
				originalIcon := &Icon{
					Name: "resizable-icon",
					Size: "24",
					Type: "Outline",
				}
				originalIcon.body = `<circle cx="12" cy="12" r="10"/>`
				// Capture the returned icon after setting size
				return ConfigureIcon(originalIcon).SetSize(32).Build()
			},
			expected: `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><circle cx="12" cy="12" r="10"/></svg>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := tt.setup()
			result := strings.TrimSpace(makeSVGTag(icon))
			expected := strings.TrimSpace(tt.expected)

			if result != expected {
				t.Errorf("String() = %q, want %q", result, expected)
			}
		})
	}
}

func TestIcon_makeSVGTag(t *testing.T) {
	// Save the original implementation
	originalGetIconBody := getIconBody

	// Defer the restoration of the original function
	defer func() { getIconBody = originalGetIconBody }()

	// Mock `getIconBody` to return different responses
	mockGetIconBody := func(name string) (string, error) {
		switch name {
		case "existing-icon":
			return `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`, nil
		case "error-icon":
			return "", fmt.Errorf("icon '%s' not found", name)
		default:
			return "", fmt.Errorf("icon '%s' not found", name)
		}
	}
	getIconBody = mockGetIconBody

	tests := []struct {
		name           string
		icon           *Icon
		expectedOutput string
	}{
		{
			name: "Body already set, should not call getIconBody",
			icon: &Icon{
				Name: "existing-icon",
				Size: "24",
				Type: "Outline",
				body: `<path d="M12 2a10 10 0 1 0 0 20a10 10 0 0 0 0-20z"/>`,
			},
			expectedOutput: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><path d="M12 2a10 10 0 1 0 0 20a10 10 0 0 0 0-20z"/></svg>`,
		},
		{
			name: "Body not set, getIconBody returns successfully",
			icon: &Icon{
				Name: "existing-icon",
				Size: "24",
				Type: "Outline",
			},
			expectedOutput: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`,
		},
		{
			name: "Body not set, getIconBody returns an error",
			icon: &Icon{
				Name: "error-icon",
				Size: "24",
				Type: "Outline",
			},
			expectedOutput: `<!-- Error: icon 'error-icon' not found -->`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := makeSVGTag(tt.icon)
			if result != tt.expectedOutput {
				t.Errorf("makeSVGTag() = %q, want %q", result, tt.expectedOutput)
			}
		})
	}
}

func TestIcon_SetSize(t *testing.T) {
	tests := []struct {
		name     string
		initial  Size
		newSize  int
		expected Size
	}{
		{
			name:     "Set size to 16",
			initial:  "24",
			newSize:  16,
			expected: "16",
		},
		{
			name:     "Set size to 32",
			initial:  "24",
			newSize:  32,
			expected: "32",
		},
		{
			name:     "Set size to 48",
			initial:  "24",
			newSize:  48,
			expected: "48",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the original icon with the initial size
			originalIcon := &Icon{
				Size: tt.initial,
			}

			// Use the ConfigureIcon builder to modify the size
			modifiedIcon := ConfigureIcon(originalIcon).SetSize(tt.newSize).Build()

			// Check that the modified icon has the expected size
			if modifiedIcon.Size != tt.expected {
				t.Errorf("SetSize() = %q, want %q", modifiedIcon.Size, tt.expected)
			}

			// Ensure the original icon is unchanged
			if originalIcon.Size != tt.initial {
				t.Errorf("Original icon size modified: got %q, want %q", originalIcon.Size, tt.initial)
			}
		})
	}
}

func TestIcon_Setters(t *testing.T) {
	originalIcon := &Icon{
		Name: "test-icon",
		Size: "24",
		Type: "Outline",
	}

	// Chain the setters and capture the returned icon
	finalIcon := ConfigureIcon(originalIcon).SetColor("#FF0000").SetSize(32).Build()

	// Validate the individual fields on the returned icon
	if finalIcon.Color != "#FF0000" {
		t.Errorf("SeColor failed: expected #FF0000, got %s", finalIcon.Color)
	}
	if finalIcon.Size.String() != "32" {
		t.Errorf("SetSize failed: expected 32, got %s", finalIcon.Size.String())
	}

	// Ensure the original icon remains unchanged
	if originalIcon.Size.String() != "24" {
		t.Errorf("Original icon size modified: expected 24, got %s", originalIcon.Size.String())
	}
	if originalIcon.Color != "" {
		t.Errorf("Original icon color modified: expected empty, got %s", originalIcon.Color)
	}

}

func TestIcon_SetAttrs(t *testing.T) {
	t.Parallel() // Run test in parallel.

	originalIcon := &Icon{
		Name: "test-icon",
		Size: "24",
		Type: "Outline",
	}
	originalIcon.body = `<path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/>`

	attrs := templ.Attributes{
		"aria-hidden": "true",
		"custom-attr": "custom-val",
		"focusable":   "false",
	}

	// Capture the returned icon after setting attributes
	finalIcon := ConfigureIcon(originalIcon).SetAttrs(attrs).Build()

	if len(finalIcon.Attrs) != len(attrs) {
		t.Errorf("expected %d attributes, got %d", len(attrs), len(finalIcon.Attrs))
	}

	for key, expectedValue := range attrs {
		if value, exists := finalIcon.Attrs[key]; !exists || value != expectedValue {
			t.Errorf("expected attribute %s=%s, got %s", key, expectedValue, value)
		}
	}

	expectedSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000" aria-hidden="true" custom-attr="custom-val" focusable="false"><path d="M10 20a10 10 0 1 0 0-20 10 10 0 0 0 0 20z"/></svg>`
	if svg := makeSVGTag(finalIcon); svg != expectedSVG {
		t.Errorf("String() = %s, want %s", svg, expectedSVG)
	}

	// Ensure the original icon remains unchanged
	if len(originalIcon.Attrs) != 0 {
		t.Errorf("Original icon attributes modified: expected 0, got %d", len(originalIcon.Attrs))
	}
}

// 2. Tests for JSON-Based Functionality
// These tests cover JSON parsing, caching, and error handling.

func TestGetIconBody_RealData(t *testing.T) {
	tests := []struct {
		name           string
		iconName       string
		expectedBody   string
		expectingError bool
	}{
		{
			name:           "Retrieve existing icon",
			iconName:       "accessibility",
			expectedBody:   `<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2S2 6.477 2 12s4.477 10 10 10M7 9l5 1m5-1l-5 1m0 0v3m0 0l-2 5m2-5l2 5"/><path fill="currentColor" d="M12 7a.5.5 0 1 1 0-1a.5.5 0 0 1 0 1"/></g>`,
			expectingError: false,
		},
		{
			name:           "Retrieve another existing icon",
			iconName:       "page-left",
			expectedBody:   `<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"><path d="M13 8.5L9.5 12l3.5 3.5"/><path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2S2 6.477 2 12s4.477 10 10 10"/></g>`,
			expectingError: false,
		},
		{
			name:           "Icon not found",
			iconName:       "non-existing-icon",
			expectedBody:   "",
			expectingError: true,
		},
		{
			name:           "Empty icon name",
			iconName:       "",
			expectedBody:   "",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := getIconBody(tt.iconName)

			if tt.expectingError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if body != tt.expectedBody {
					t.Errorf("getIconBody() = %q, want %q", body, tt.expectedBody)
				}
			}
		})
	}
}

func TestGetIconBody_OnceWithRealData(t *testing.T) {
	// First call should initialize the data
	_, err := getIconBody("check-circle")
	if err != nil {
		t.Fatalf("unexpected error during first call: %v", err)
	}

	// Ensure no error on subsequent calls for valid icons
	_, err = getIconBody("chromecast")
	if err != nil {
		t.Fatalf("unexpected error during subsequent call: %v", err)
	}
}

// 3. Tests for Mocked Data
// These tests cover cases where mocked FS and invalid JSON are used.

func TestIcon_String_FetchBody(t *testing.T) {
	resetTestState()

	// Mock the embedded JSON with valid data
	validJSON := `{
        "icons": {
            "voice-xmark": { "body": "<path fill=\"none\" stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"1.5\" d=\"M12 3v16M8 8v6m12-5v4M4 9v4m12-7v8m.121 7.364l2.122-2.121m0 0l2.121-2.122m-2.121 2.122L16.12 17.12m2.122 2.122l2.121 2.121\"/>7" },
			"meter-arrow-down-right": { "body": "<g fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M2.5 3.5L7 8m0 0V4m0 4H3m12 8l-3.5-3.5\"/><path d=\"M14.5 9C10.358 9 7 12.283 7 16.333a7.2 7.2 0 0 0 .733 3.165a.93.93 0 0 0 .84.502h11.853a.93.93 0 0 0 .841-.502A7.2 7.2 0 0 0 22 16.333C22 12.283 18.642 9 14.5 9Z\"/></g>" }
        }
    }`
	iconoirJSONSource = mockInvalidJSONFS(validJSON)
	defer func() {
		iconoirJSONSource = iconoirJSON // Restore original embedded JSON
	}()

	t.Run("Fetches and caches body", func(t *testing.T) {
		icon := &Icon{
			Name: "meter-arrow-down-right",
			Size: "24",
			Type: "Outline",
		}

		// Call String() for the first time to trigger the body fetch
		result := makeSVGTag(icon) // Pass a pointer

		// Validate the resulting SVG
		expected := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke-width="1.5" color="#000000"><g fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M2.5 3.5L7 8m0 0V4m0 4H3m12 8l-3.5-3.5"/><path d="M14.5 9C10.358 9 7 12.283 7 16.333a7.2 7.2 0 0 0 .733 3.165a.93.93 0 0 0 .84.502h11.853a.93.93 0 0 0 .841-.502A7.2 7.2 0 0 0 22 16.333C22 12.283 18.642 9 14.5 9Z"/></g></svg>`
		if result != expected {
			t.Errorf("String() = %q, want %q", result, expected)
		}
	})
}

func TestGetIconBody_JSONParsing(t *testing.T) {
	tests := []struct {
		name           string
		mockJSON       string
		iconName       string
		expectedError  string
		expectedResult string
	}{
		{
			name:          "Invalid JSON format",
			mockJSON:      `{"icons": "invalid"`, // Invalid JSON structure
			iconName:      "academic-cap",
			expectedError: "failed to parse iconoir JSON",
		},
		{
			name:          "Missing icons field",
			mockJSON:      `{"missingIcons": {}}`, // No `icons` key
			iconName:      "academic",
			expectedError: "icon 'academic' not found",
		},
		{
			name:           "Valid JSON",
			mockJSON:       `{"icons": {"academic-cap": {"body": "<path d='...'/>"}}}`,
			iconName:       "academic-cap",
			expectedError:  "",
			expectedResult: "<path d='...'/>",
		},
		{
			name:          "Icon not found",
			mockJSON:      `{"icons": {"academic-cap": {"body": "<path d='...'/>"}}}`,
			iconName:      "non-existent-icon",
			expectedError: "icon 'non-existent-icon' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetTestState()

			// Replace iconoirJSONSource with a mocked FS
			iconoirJSONSource = mockInvalidJSONFS(tt.mockJSON)
			defer func() {
				iconoirJSONSource = iconoirJSON // Restore original embedded FS
			}()

			result, err := getIconBody(tt.iconName)

			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error %q, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if result != tt.expectedResult {
				t.Errorf("Expected result %q, got %q", tt.expectedResult, result)
			}
		})
	}
}

// 4. Utility Functions for Testing
// These utilities mock data and manage state resets.

type mockFS struct {
	data map[string]string
}

func mockInvalidJSONFS(data string) fs.FS {
	return &mockFS{
		data: map[string]string{"data/iconoir_cache.json": data},
	}
}

func (m *mockFS) Open(name string) (fs.File, error) {
	content, exists := m.data[name]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", name)
	}
	return &mockFile{content: strings.NewReader(content)}, nil
}

type mockFile struct {
	content io.Reader
}

func (f *mockFile) Read(p []byte) (int, error) {
	return f.content.Read(p)
}

func (f *mockFile) Close() error {
	return nil
}

func (f *mockFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func resetTestState() {
	iconBodyCache = map[string]string{}
}

func TestMockFS(t *testing.T) {
	data := `{"icons": invalid}`
	mockFS := mockInvalidJSONFS(data)
	content, err := fs.ReadFile(mockFS, "data/iconoir_cache.json")
	if err != nil {
		t.Fatalf("Failed to read mock file: %v", err)
	}
	if string(content) != data {
		t.Fatalf("Expected mock content %q, got %q", data, string(content))
	}
}
