package templiconoir

import (
	"embed"
	"io/fs"
)

//go:embed data/iconoir_cache.json
var iconoirJSON embed.FS

var iconoirJSONSource fs.FS = iconoirJSON // Default to the embedded FS
