package contracts

type SettingYAMLFile struct {
	SyntaxVersion  int            `yaml:"syntax_version"`
	Server         string         `yaml:"server"`
	CodeGeneration CodeGeneration `yaml:"code_generation"`
}

type CodeGeneration struct {
	OutputFolder string `yaml:"output_folder"`
	Language     string `yaml:"language"`
	ClassSuffix  string `yaml:"class_suffix"`
}
