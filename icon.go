package templiconoir

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/tidwall/gjson"
)

var (
	iconBodyCache = map[string]string{}
	cacheMutex    sync.Mutex
)

// Size represents the size of UI components (e.g., small, medium, large).
type Size string

// String returns the string representation of a Size.
func (s Size) String() string {
	return string(s)
}

// Icon represents a single icon with its attributes.
type Icon struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        Size   `json:"size"`
	StrokeWidth string
	Color       string
	Attrs       templ.Attributes
	body        string // Cached Body
}

func (i *Icon) Render() templ.Component {
	return templ.Raw(makeSVGTag(i))
}

func (i *Icon) fetchBody() error {
	if i.body != "" {
		return nil // Body is already cached
	}

	body, err := getIconBody(i.Name)
	if err != nil {
		return err
	}

	i.body = body
	return nil
}

func makeSVGTag(icon *Icon) string {
	// Set default values for stroke width and color
	strokeWidth := defaultIfEmpty(icon.StrokeWidth, "1.5")
	color := defaultIfEmpty(icon.Color, "#000000")

	// Ensure the icon body is fetched and cached
	if err := icon.fetchBody(); err != nil {
		return errorSVGComment(err)
	}

	svgTag := fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 24 24" fill="none" stroke-width="%s" color="%s"`,
		icon.Size.String(),
		icon.Size.String(),
		strokeWidth,
		color,
	)

	// Add user-defined attributes and close the opening <svg> tag
	var builder strings.Builder
	builder.WriteString(svgTag)
	addAttributesToSVG(&builder, icon.Attrs)
	builder.WriteString(">")

	// Add the icon body and close the </svg> tag
	builder.WriteString(icon.body)
	builder.WriteString(`</svg>`)

	return builder.String()
}

// getIconBody retrieves the body of an icon by its name, with thread-safe caching.
var getIconBody = func(name string) (string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the body is already cached.
	if body, found := iconBodyCache[name]; found {
		return body, nil
	}

	// Read and parse the JSON data.
	jsonFilename := "data/iconoir_cache.json"
	iconoirData, _ := iconoirJSONSource.Open(jsonFilename)
	defer iconoirData.Close()

	data, _ := io.ReadAll(iconoirData)

	// Check for valid JSON (parsing)
	if !gjson.ValidBytes(data) {
		return "", fmt.Errorf("failed to parse iconoir JSON")
	}

	// Extract the icons key
	iconsResult := gjson.GetBytes(data, "icons")

	// If the icons key exists, populate the cache
	if iconsResult.Exists() {
		iconsResult.ForEach(func(key, value gjson.Result) bool {
			iconBody := value.Get("body").String()
			iconBodyCache[key.String()] = iconBody
			return true
		})
	}

	// Return the requested icon body
	body, exists := iconBodyCache[name]
	if !exists {
		return "", fmt.Errorf("icon '%s' not found", name)
	}
	return body, nil
}
