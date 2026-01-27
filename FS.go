package FS

import "embed"

//go:embed front-end/templates
var TemplatesFS embed.FS
var StaticFS    embed.FS
