package actions

import (
	"context"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
	"os"
)

func InitDefaultSettingsFileAction(context.Context, *cli.Command) error {
	if _, err := os.Stat("fcsettings.yaml"); err == nil {
		return cli.Exit("File 'fcsettings.yaml' already exists in current directory", 1)
	}

	config := contracts.SettingYAMLFile{
		SyntaxVersion: 1,
		Server:        "https://api.fusioncatalyst.io",
	}

	file, err := os.Create("fcsettings.yaml")
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
