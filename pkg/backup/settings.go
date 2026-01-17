package backup

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// SettingsBackup handles settings.json backup
type SettingsBackup struct {
	skipSettings bool
}

// NewSettingsBackup creates a new settings backup handler
func NewSettingsBackup(skip bool) *SettingsBackup {
	return &SettingsBackup{skipSettings: skip}
}

// Backup backs up settings.json, removing sensitive information
func (sb *SettingsBackup) Backup(claudeDir, configDir string) error {
	if sb.skipSettings {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Skipped(i18n.T("backup.steps.settings", nil))
		return nil
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("backup.steps.settings", nil))

	srcPath := filepath.Join(claudeDir, "settings.json")
	dstPath := filepath.Join(configDir, "settings.json")

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

	// Remove sensitive information
	if env, ok := settings["env"].(map[string]interface{}); ok {
		delete(env, "ANTHROPIC_AUTH_TOKEN")
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Success(i18n.T("backup.messages.settings_removed", nil))
	}

	// Write cleaned JSON
	cleaned, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	// Create destination directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(dstPath, cleaned, 0644)
}
