package backup

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// SkillsBackup handles custom skills backup
type SkillsBackup struct {
	skipSkills bool
}

// NewSkillsBackup creates a new skills backup handler
func NewSkillsBackup(skip bool) *SkillsBackup {
	return &SkillsBackup{skipSkills: skip}
}

// Backup backs up custom skills
func (skb *SkillsBackup) Backup(claudeDir, configDir string) error {
	if skb.skipSkills {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Skipped(i18n.T("backup.steps.skills", nil))
		return nil
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("backup.steps.skills", nil))

	srcDir := filepath.Join(claudeDir, "skills")
	dstDir := filepath.Join(configDir, "skills")

	// Check if source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
			"Item": "skills directory",
		}))
		return nil
	}

	// Create destination directory
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Copy all skill files
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

		if err := copySkillFile(srcPath, dstPath); err != nil {
			return err
		}
		count++
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Success(i18n.T("backup.messages.skills_count", map[string]interface{}{
		"Count": count,
	}))

	return nil
}

// copySkillFile copies a skill file from src to dst
func copySkillFile(src, dst string) (err error) {
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
