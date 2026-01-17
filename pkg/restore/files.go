package restore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// CommandsRestore handles custom commands restore
type CommandsRestore struct{}

// NewCommandsRestore creates a new commands restore handler
func NewCommandsRestore() *CommandsRestore {
	return &CommandsRestore{}
}

// Restore restores custom commands
func (cr *CommandsRestore) Restore(configDir, claudeDir string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("restore.steps.commands", nil))

	srcDir := filepath.Join(configDir, "commands")
	dstDir := filepath.Join(claudeDir, "commands")

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
			"Item": "commands directory",
		}))
		return nil
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

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
	ui.Success(i18n.T("restore.messages.restored", map[string]interface{}{
		"Item": fmt.Sprintf("%d commands", count),
	}))
	return nil
}

// SkillsRestore handles custom skills restore
type SkillsRestore struct{}

// NewSkillsRestore creates a new skills restore handler
func NewSkillsRestore() *SkillsRestore {
	return &SkillsRestore{}
}

// Restore restores custom skills
func (sr *SkillsRestore) Restore(configDir, claudeDir string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("restore.steps.skills", nil))

	srcDir := filepath.Join(configDir, "skills")
	dstDir := filepath.Join(claudeDir, "skills")

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
			"Item": "skills directory",
		}))
		return nil
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

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
	ui.Success(i18n.T("restore.messages.restored", map[string]interface{}{
		"Item": fmt.Sprintf("%d skills", count),
	}))
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) (err error) {
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

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
