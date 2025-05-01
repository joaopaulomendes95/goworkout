package web

import "embed"

// efs = embed file server

//go:embed "assets"
var Files embed.FS
