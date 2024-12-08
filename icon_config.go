package templiconoir

import (
	"strconv"

	"github.com/a-h/templ"
)

type IconBuilder struct {
	icon Icon
}

// ConfigureIcon creates a new builder from an existing icon.
func ConfigureIcon(icon *Icon) *IconBuilder {
	return &IconBuilder{
		icon: *icon, // Copy the icon (not reference)
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

// SetColor sets the color of the icon.
func (b *IconBuilder) SetColor(value string) *IconBuilder {
	b.icon.Color = value
	return b
}

// SetAttrs sets the attributes for the SVG tag.
func (b *IconBuilder) SetAttrs(attrs templ.Attributes) *IconBuilder {
	b.icon.Attrs = attrs
	return b
}

// Build returns the final configured Icon.
func (b *IconBuilder) Build() *Icon {
	return &b.icon
}
