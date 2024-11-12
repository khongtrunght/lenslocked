package templates

import "embed"

// FS is the embedded filesystem for the templates.
//
//go:embed *
var FS embed.FS
