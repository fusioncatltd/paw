package actions

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/urfave/cli/v3"
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
			Language: "typescript",
		},
	}

	// Get values from CLI flags if provided
	server := cmd.String("server")
	language := cmd.String("language")

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

	// If there are any missing values, ask for them interactively
	if len(questions) > 0 {
		answers := struct {
			Server   string
			Language string
		}{}

		if err := survey.Ask(questions, &answers); err != nil {
			return fmt.Errorf("error during survey: %w", err)
		}

		// Update config with interactive answers
		if server == "" {
			config.Server = answers.Server
		}
		if language == "" {
			config.CodeGeneration.Language = answers.Language
		}
	}

	// 2. Marshal to YAML bytes
	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling to YAML: %v\n", err)
		os.Exit(1)
	}

	// 3. Write the bytes to a file
	if err := os.WriteFile("fcsettings.yaml", data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}

	//file, err := os.Create("fcsettings.yaml")
	//if err != nil {
	//	return err
	//}
	//
	//defer file.Close()
	//
	//encoder := yaml.NewEncoder(file)
	//defer encoder.Close()
	//
	//if err := encoder.Encode(config); err != nil {
	//	return err
	//}

	fmt.Println("Configuration file 'fcsettings.yaml' has been created successfully!")
	return nil
}
