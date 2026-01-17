package restore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// SettingsRestore handles settings.json restore
type SettingsRestore struct{}

// NewSettingsRestore creates a new settings restore handler
func NewSettingsRestore() *SettingsRestore {
	return &SettingsRestore{}
}

// Restore restores settings.json
func (sr *SettingsRestore) Restore(configDir, claudeDir string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("restore.steps.settings", nil))

	srcPath := filepath.Join(configDir, "settings.json")
	dstPath := filepath.Join(claudeDir, "settings.json")

	// Read source file
	data, err := os.ReadFile(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			//nolint:govet // msg is from i18n.T() which returns runtime strings
			ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
				"Item": "settings.json",
			}))
			return nil
		}
		return err
	}

	// Parse JSON
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	// Check if destination exists and has API token
	existingToken := ""
	if _, err := os.Stat(dstPath); err == nil {
		existingData, _ := os.ReadFile(dstPath)
		var existingSettings map[string]interface{}
		if json.Unmarshal(existingData, &existingSettings) == nil {
			if env, ok := existingSettings["env"].(map[string]interface{}); ok {
				if token, ok := env["ANTHROPIC_AUTH_TOKEN"].(string); ok && token != "" {
					existingToken = token
					//nolint:govet // msg is from i18n.T() which returns runtime strings
					ui.Println(ui.Yellow, "  "+i18n.T("restore.messages.existing_token", nil))
				}
			}
		}
	}

	// Handle API token
	if existingToken == "" {
		// Check environment variable
		if envToken := os.Getenv("ANTHROPIC_AUTH_TOKEN"); envToken != "" {
			existingToken = envToken
			ui.Println(ui.Yellow, "  Using ANTHROPIC_AUTH_TOKEN from environment")
		}
	}

	if existingToken == "" {
		// Prompt for token
		fmt.Print("  " + i18n.T("restore.prompts.api_token", nil) + " ")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)
		if token != "" {
			existingToken = token
		}
	}

	// Merge token into settings
	if existingToken != "" {
		if settings["env"] == nil {
			settings["env"] = make(map[string]interface{})
		}
		if env, ok := settings["env"].(map[string]interface{}); ok {
			env["ANTHROPIC_AUTH_TOKEN"] = existingToken
		}
	}

	// Write settings
	result, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return err
	}

	ui.Success("settings.json restored")
	return os.WriteFile(dstPath, result, 0644)
}
