package templiconoir

import (
	_ "embed"
	"fmt"
	"io"
	"strconv"
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

// IconBuilder is a builder for configuring an Icon.
// It allows method chaining to update the icon's properties.
type IconBuilder struct {
	icon *Icon // Reference to the icon being configured
}

// Config returns an IconBuilder to allow chaining configuration methods on the icon.
func (icon *Icon) Config() *IconBuilder {
	return &IconBuilder{
		icon: icon.clone(), // Clone the icon to ensure immutability
	}
}

// ConfigureIcon creates a new builder from an existing icon.
func ConfigureIcon(icon *Icon) *IconBuilder {
	return &IconBuilder{
		icon: icon.clone(), // Clone the icon to ensure immutability
	}
}

// SetSize sets the size of the icon.
func (b *IconBuilder) SetSize(size int) *IconBuilder {
	b.icon.Size = Size(strconv.Itoa(size))
	return b
}

// SetStrokeWidth sets the stroke-width of the icon.
func (b *IconBuilder) SetStrokeWidth(value string) *IconBuilder {
	b.icon.StrokeWidth = value
	return b
}

// SetColor sets the fill color of the icon.
func (b *IconBuilder) SetColor(value string) *IconBuilder {
	b.icon.Color = value
	return b
}

// SetAttrs sets custom attributes for the SVG tag (e.g., `aria-hidden`, `focusable`).
func (b *IconBuilder) SetAttrs(attrs templ.Attributes) *IconBuilder {
	b.icon.Attrs = attrs
	return b
}

// GetIcon returns the configured icon instance.
func (b *IconBuilder) GetIcon() *Icon {
	return b.icon
}

// Render generates the SVG for the configured icon.
func (b *IconBuilder) Render() templ.Component {
	return b.icon.Render()
}

func (i *Icon) clone() *Icon {
	attrsCopy := make(templ.Attributes, len(i.Attrs))
	for k, v := range i.Attrs {
		attrsCopy[k] = v
	}
	return &Icon{
		Name:        i.Name,
		Type:        i.Type,
		Size:        i.Size,
		StrokeWidth: i.StrokeWidth,
		Color:       i.Color,
		Attrs:       attrsCopy,
		body:        i.body, // The body is shared since it's immutable
	}
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
