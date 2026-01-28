package FS

import "embed"

//go:embed front-end/templates
var TemplatesFS embed.FS

//go:embed front-end/static
var StaticFS    embed.FS
