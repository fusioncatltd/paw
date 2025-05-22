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

func InitDefaultSettingsFileAction(ctx context.Context, cmd *cli.Command) error {
	if _, err := os.Stat("fcsettings.yaml"); err == nil {
		return cli.Exit("File 'fcsettings.yaml' already exists in current directory", 1)
	}

	// Default values
	config := contracts.SettingYAMLFile{
		SyntaxVersion: 1,
		Server:        "https://api.fusioncat.dev",
		CodeGeneration: contracts.CodeGeneration{
			OutputFolder: "./fuscioncat",
			Language:     "typescript",
			ClassSuffix:  "FusionCatModel",
		},
	}

	// Get values from CLI flags if provided
	server := cmd.String("server")
	outputFolder := cmd.String("output-folder")
	language := cmd.String("language")
	classSuffix := cmd.String("class-suffix")

	// Prepare questions for missing values
	var questions []*survey.Question

	// Server
	if server == "" {
		questions = append(questions, &survey.Question{
			Name: "server",
			Prompt: &survey.Input{
				Message: "Enter the server URL:",
				Default: config.Server,
			},
		})
	} else {
		config.Server = server
	}

	// Output Folder
	if outputFolder == "" {
		questions = append(questions, &survey.Question{
			Name: "outputFolder",
			Prompt: &survey.Input{
				Message: "Enter the output folder for generated models:",
				Default: config.CodeGeneration.OutputFolder,
			},
		})
	} else {
		config.CodeGeneration.OutputFolder = outputFolder
	}

	// Language
	if language == "" {
		questions = append(questions, &survey.Question{
			Name: "language",
			Prompt: &survey.Select{
				Message: "Select the target language:",
				Options: []string{"typescript", "python", "java", "go"},
				Default: config.CodeGeneration.Language,
			},
		})
	} else {
		// Validate language
		validLanguages := map[string]bool{
			"typescript": true,
			"python":     true,
			"java":       true,
			"go":         true,
		}
		if !validLanguages[language] {
			return cli.Exit(fmt.Sprintf("Invalid language: %s. Must be one of: typescript, python, java, go", language), 1)
		}
		config.CodeGeneration.Language = language
	}

	// Class Suffix
	if classSuffix == "" {
		questions = append(questions, &survey.Question{
			Name: "classSuffix",
			Prompt: &survey.Input{
				Message: "Enter the suffix for generated classes:",
				Default: config.CodeGeneration.ClassSuffix,
			},
		})
	} else {
		config.CodeGeneration.ClassSuffix = classSuffix
	}

	// If there are any missing values, ask for them interactively
	if len(questions) > 0 {
		answers := struct {
			Server       string
			OutputFolder string
			Language     string
			ClassSuffix  string
		}{}

		if err := survey.Ask(questions, &answers); err != nil {
			return fmt.Errorf("error during survey: %w", err)
		}

		// Update config with interactive answers
		if server == "" {
			config.Server = answers.Server
		}
		if outputFolder == "" {
			config.CodeGeneration.OutputFolder = answers.OutputFolder
		}
		if language == "" {
			config.CodeGeneration.Language = answers.Language
		}
		if classSuffix == "" {
			config.CodeGeneration.ClassSuffix = answers.ClassSuffix
		}
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

	fmt.Println("Configuration file 'fcsettings.yaml' has been created successfully!")
	return nil
}
