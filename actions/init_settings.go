package actions

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func InitDefaultSettingsFileAction(context.Context, *cli.Command) error {
	if _, err := os.Stat("fcsettings.yaml"); err == nil {
		return cli.Exit("File 'fcsettings.yaml' already exists in current directory", 1)
	}

	// Default values
	config := contracts.SettingYAMLFile{
		SyntaxVersion: 1,
		Server:        "https://api.fusioncat.dev",
		CodeGeneration: contracts.CodeGeneration{
			OutputFolder: "./fuscioncat/models",
			Language:     "typescript",
			ClassSuffix:  "FusionCatModel",
		},
	}

	// Interactive prompts
	questions := []*survey.Question{
		{
			Name: "server",
			Prompt: &survey.Input{
				Message: "Enter the server URL:",
				Default: config.Server,
			},
		},
		{
			Name: "outputFolder",
			Prompt: &survey.Input{
				Message: "Enter the output folder for generated models:",
				Default: config.CodeGeneration.OutputFolder,
			},
		},
		{
			Name: "language",
			Prompt: &survey.Select{
				Message: "Select the target language:",
				Options: []string{"typescript", "python", "java", "go"},
				Default: config.CodeGeneration.Language,
			},
		},
		{
			Name: "classSuffix",
			Prompt: &survey.Input{
				Message: "Enter the suffix for generated classes:",
				Default: config.CodeGeneration.ClassSuffix,
			},
		},
	}

	answers := struct {
		Server       string
		OutputFolder string
		Language     string
		ClassSuffix  string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return fmt.Errorf("error during survey: %w", err)
	}

	// Update config with user answers
	config.Server = answers.Server
	config.CodeGeneration.OutputFolder = answers.OutputFolder
	config.CodeGeneration.Language = answers.Language
	config.CodeGeneration.ClassSuffix = answers.ClassSuffix

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

	fmt.Println("Configuration file 'fcsettings.yaml' has been created successfully!")
	return nil
}
