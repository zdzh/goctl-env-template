package types

type ConfigField struct {
	Name         string
	EnvVar       string
	Type         string
	IsOptional   bool
	DefaultValue string
	Comment      string
	Group        string
}

type ConfigStruct struct {
	Name    string
	Comment string
	Fields  []ConfigField
}
