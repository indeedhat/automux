package configs

import _ "embed"

//go:embed template.automux.tmpl
var IclTemplate string

//go:embed template.automux.json.tmpl
var JsonTemplate string

//go:embed template.automux.yml.tmpl
var YamlTemplate string
