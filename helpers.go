package templiconoir

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/a-h/templ"
)

// Reserved attributes for SVG tags that should not be overwritten.
var reservedSVGAttributes = map[string]struct{}{
	"xmlns":        {},
	"viewBox":      {},
	"width":        {},
	"height":       {},
	"stroke-width": {},
	"stroke":       {},
	"fill":         {},
}

// sanitizeAttribute ensures that attribute keys and values are safe for inclusion in the SVG tag.
func sanitizeAttribute(key, value string) (string, string, bool) {
	// Define allowlist for event attributes
	allowedEventAttributes := map[string]struct{}{
		"onclick":  {},
		"onchange": {},
		"onhover":  {},
	}

	// Check for unsafe attributes
	if _, isEvent := allowedEventAttributes[key]; isEvent {
		// For event attributes, only allow simple JS functions (no <script> tags, eval, etc.)
		if strings.Contains(strings.ToLower(value), "<script>") || strings.Contains(strings.ToLower(value), "javascript:") {
			return "", "", false // Unsafe value
		}
	}

	// Escape any unsafe characters for all attributes
	escapedKey := html.EscapeString(key)
	escapedValue := html.EscapeString(value)

	return escapedKey, escapedValue, true // Safe attribute
}

// addAttributesToSVG adds templ.Attributes to the SVG tag, placing them at the end of the <svg> opening tag.
// Reserved attributes are skipped to avoid overwriting critical SVG settings.
// Attributes are sanitized to prevent XSS or injection attacks.
func addAttributesToSVG(builder *strings.Builder, attrs templ.Attributes) {
	if len(attrs) == 0 {
		return
	}

	// Extract keys and sort them for deterministic order
	keys := make([]string, 0, len(attrs))
	for key := range attrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Process attributes in sorted order
	for _, key := range keys {
		value, ok := attrs[key].(string) // Ensure value is a string
		if !ok {
			// Skip attributes with non-string values
			continue
		}

		// Skip reserved attributes
		if _, isReserved := reservedSVGAttributes[key]; isReserved {
			continue
		}

		// Sanitize the attribute
		sanitizedKey, sanitizedValue, ok := sanitizeAttribute(key, value)
		if !ok {
			// Skip attributes that are not safe
			continue
		}

		// Add the sanitized attribute to the SVG tag
		fmt.Fprintf(builder, ` %s="%v"`, sanitizedKey, sanitizedValue)
	}
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func errorSVGComment(err error) string {
	return fmt.Sprintf("<!-- Error: %s -->", err)
}
