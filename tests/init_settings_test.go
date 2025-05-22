package tests

import (
	"context"
	"os"
	"testing"

	"github.com/fusioncatalyst/paw/actions"
	"github.com/fusioncatalyst/paw/contracts"
	"github.com/fusioncatalyst/paw/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func TestInitDefaultSettingsFileAction(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	t.Run("create settings file with default values", func(t *testing.T) {
		utils.Cleanup()
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "https://api.fusioncat.dev",
				},
				&cli.StringFlag{
					Name:  "language",
					Value: "typescript",
				},
			},
		}
		cmd.Set("server", "https://api.fusioncat.dev")
		cmd.Set("language", "typescript")

		output, err := utils.CaptureOutputInTests(actions.InitDefaultSettingsFileAction, context.Background(), cmd)
		assert.NoError(t, err)
		assert.Contains(t, output, "Configuration file 'fcsettings.yaml' has been created")

		// Read and verify content
		content, err := os.ReadFile("fcsettings.yaml")
		assert.NoError(t, err)

		var config contracts.SettingYAMLFile
		err = yaml.Unmarshal(content, &config)
		assert.NoError(t, err)
		assert.Equal(t, 1, config.SyntaxVersion)
		assert.Equal(t, "https://api.fusioncat.dev", config.Server)
		assert.Equal(t, "typescript", config.CodeGeneration.Language)
	})

	t.Run("create settings file with custom values", func(t *testing.T) {
		utils.Cleanup()
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "https://custom-api.example.com",
				},
				&cli.StringFlag{
					Name:  "language",
					Value: "python",
				},
			},
		}
		cmd.Set("server", "https://custom-api.example.com")
		cmd.Set("language", "python")

		output, err := utils.CaptureOutputInTests(actions.InitDefaultSettingsFileAction, context.Background(), cmd)
		assert.NoError(t, err)
		assert.Contains(t, output, "Configuration file 'fcsettings.yaml' has been created")

		// Read and verify content
		content, err := os.ReadFile("fcsettings.yaml")
		assert.NoError(t, err)

		var config contracts.SettingYAMLFile
		err = yaml.Unmarshal(content, &config)
		assert.NoError(t, err)
		assert.Equal(t, 1, config.SyntaxVersion)
		assert.Equal(t, "https://custom-api.example.com", config.Server)
		assert.Equal(t, "python", config.CodeGeneration.Language)
	})

	t.Run("file already exists", func(t *testing.T) {
		// Create the file first
		utils.Cleanup()
		file, err := os.Create("fcsettings.yaml")
		assert.NoError(t, err)
		file.Close()

		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "https://api.fusioncat.dev",
				},
				&cli.StringFlag{
					Name:  "language",
					Value: "typescript",
				},
			},
		}
		cmd.Set("server", "https://api.fusioncat.dev")
		cmd.Set("language", "typescript")

		_, err = utils.CaptureOutputInTests(actions.InitDefaultSettingsFileAction, context.Background(), cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "File 'fcsettings.yaml' already exists in current directory")
	})

	t.Run("invalid language value", func(t *testing.T) {
		utils.Cleanup()
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "server",
					Value: "https://api.fusioncat.dev",
				},
				&cli.StringFlag{
					Name:  "language",
					Value: "invalid",
				},
			},
		}
		cmd.Set("server", "https://api.fusioncat.dev")
		cmd.Set("language", "invalid")

		_, err := utils.CaptureOutputInTests(actions.InitDefaultSettingsFileAction, context.Background(), cmd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid language: invalid")
	})
}
