package backup

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// CopyFile copies a file - exported for use in other backup packages
func CopyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := srcFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Make executable
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	return CopyFile(src, dst)
}

// CommandsBackup handles custom commands backup
type CommandsBackup struct {
	skipCommands bool
}

// NewCommandsBackup creates a new commands backup handler
func NewCommandsBackup(skip bool) *CommandsBackup {
	return &CommandsBackup{skipCommands: skip}
}

// Backup backs up custom commands
func (cb *CommandsBackup) Backup(claudeDir, configDir string) error {
	if cb.skipCommands {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Skipped(i18n.T("backup.steps.commands", nil))
		return nil
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("backup.steps.commands", nil))

	srcDir := filepath.Join(claudeDir, "commands")
	dstDir := filepath.Join(configDir, "commands")

	// Check if source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
			"Item": "commands directory",
		}))
		return nil
	}

	// Create destination directory
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Copy all command files
	count := 0
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
		count++
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Success(i18n.T("backup.messages.commands_count", map[string]interface{}{
		"Count": count,
	}))

	return nil
}
