package config

import (
	"embed"
)

//go:embed boot.yaml
var Boot []byte

//go:embed swagger
var SwaggerFS embed.FS
