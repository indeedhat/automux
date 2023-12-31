package configs

import _ "embed"

//go:embed template.automux.hcl.tmpl
var ConfigTemplate string
