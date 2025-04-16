package contracts

type SettingYAMLFile struct {
	SyntaxVersion int    `yaml:"syntax_version"`
	Server        string `yaml:"server"`
}
