package contracts

type SettingYAMLFile struct {
	SyntaxVersion      int            `yaml:"syntaxVersion"`
	Server             string         `yaml:"server"`
	CodeGeneration     CodeGeneration `yaml:"codeGeneration"`
	WorkingWithProject *string        `yaml:"workingWithProject,omitempty"`
}

type CodeGeneration struct {
	Language string `yaml:"language"`
}
