package backup

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// ProjectsBackup handles project configs backup
type ProjectsBackup struct {
	projects []string
	scanAll  bool
}

// NewProjectsBackup creates a new projects backup handler
func NewProjectsBackup(projects []string, scanAll bool) *ProjectsBackup {
	return &ProjectsBackup{
		projects: projects,
		scanAll:  scanAll,
	}
}

// Backup backs up project configurations
func (pb *ProjectsBackup) Backup(configDir string) error {
	if len(pb.projects) == 0 && !pb.scanAll {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Skipped(i18n.T("backup.steps.projects", nil))
		return nil
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("backup.steps.projects", nil))

	// Collect directories to scan
	scanDirs := pb.projects
	if pb.scanAll {
		cfg := config.Get()
		scanDirs = cfg.ScanDirs
	}

	count := 0
	for _, baseDir := range scanDirs {
		expanded, err := config.ExpandPath(baseDir)
		if err != nil {
			continue
		}

		found, err := pb.scanDirectory(expanded, configDir)
		if err != nil {
			ui.Warning("Error scanning %s: %v", baseDir, err)
			continue
		}
		count += found
	}

	if count == 0 {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Warning(i18n.T("backup.warnings.not_found", map[string]interface{}{
			"Item": "project configurations",
		}))
	} else {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Success(i18n.T("backup.messages.projects_count", map[string]interface{}{
			"Count": count,
		}))
	}

	return nil
}

// scanDirectory scans a directory for .claude folders
func (pb *ProjectsBackup) scanDirectory(baseDir, configDir string) (int, error) {
	count := 0

	// Check if base directory exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return 0, nil
	}

	// Walk through directory looking for .claude folders
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on error
		}

		if !info.IsDir() {
			return nil
		}

		// Check if this is a .claude directory
		if info.Name() == ".claude" {
			projectPath := filepath.Dir(path)
			projectName := filepath.Base(projectPath)

			if err := pb.backupProject(path, projectName, configDir); err == nil {
				count++
				//nolint:govet // format string contains non-constant parts
				ui.Println(ui.Gray, "  ✓ "+projectName)
			}
			// Don't descend into .claude directory
			return filepath.SkipDir
		}

		return nil
	})

	return count, err
}

// backupProject backs up a single project's .claude config
func (pb *ProjectsBackup) backupProject(claudeDir, projectName, configDir string) error {
	dstDir := filepath.Join(configDir, "project-configs", projectName)

	// Create destination directory
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Copy all files from .claude directory
	entries, err := os.ReadDir(claudeDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(claudeDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// Copy directory recursively
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFileCopy(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFileCopy(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFileCopy copies a file from src to dst (renamed to avoid conflict)
func copyFileCopy(src, dst string) (err error) {
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

	// Preserve file mode
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
