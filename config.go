package config

import (
	"embed"
)

//go:embed boot.yaml
var Boot []byte

//go:embed api
var SwaggerFS embed.FS
