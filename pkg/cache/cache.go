package cache

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
)

// Backup backs up plugin cache to a tar.gz file
func Backup(repoPath string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Title(i18n.T("cache.backup.title", nil))

	claudeDir, err := config.GetClaudeDir()
	if err != nil {
		return err
	}

	pluginsCacheDir := filepath.Join(claudeDir, "plugins", "cache")
	cacheDir := filepath.Join(repoPath, "cache")
	cacheFile := filepath.Join(cacheDir, "plugins-cache.tar.gz")

	// Check if source exists
	if _, err := os.Stat(pluginsCacheDir); os.IsNotExist(err) {
		return os.ErrNotExist
	}

	// Create cache directory
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	// Show detected plugins
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Green, i18n.T("cache.backup.detected", nil))
	entries, _ := os.ReadDir(pluginsCacheDir)
	for _, entry := range entries {
		if entry.IsDir() {
			info, _ := entry.Info()
			//nolint:govet // format string contains non-constant parts
			ui.Println(ui.Gray, "  • "+entry.Name()+" ("+formatSize(info.Size())+")")
		}
	}

	// Create tar.gz file
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("cache.backup.packing", nil))

	outFile, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			// Log but don't fail if close fails during cleanup
			_ = closeErr
		}
	}()

	gzw := gzip.NewWriter(outFile)
	defer func() {
		if closeErr := gzw.Close(); closeErr != nil {
			// Log but don't fail if close fails during cleanup
			_ = closeErr
		}
	}()

	tw := tar.NewWriter(gzw)
	defer func() {
		if closeErr := tw.Close(); closeErr != nil {
			// Log but don't fail if close fails during cleanup
			_ = closeErr
		}
	}()

	return filepath.Walk(pluginsCacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(pluginsCacheDir, path)
		header.Name = filepath.Join("cache", relPath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				// Log but don't fail if close fails during cleanup
				_ = closeErr
			}
		}()

		_, err = io.Copy(tw, file)
		return err
	})
}

// Restore restores plugin cache from tar.gz file
func Restore(repoPath string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Title(i18n.T("cache.restore.title", nil))

	claudeDir, err := config.GetClaudeDir()
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(repoPath, "cache", "plugins-cache.tar.gz")

	// Check if cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return os.ErrNotExist
	}

	// Show file info
	info, _ := os.Stat(cacheFile)
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Green, i18n.T("cache.restore.info", nil))
	//nolint:govet // format string contains non-constant parts
	ui.Println(ui.Gray, "  File: "+cacheFile)
	//nolint:govet // format string contains non-constant parts
	ui.Println(ui.Gray, "  Size: "+formatSize(info.Size()))

	// Extract
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("cache.restore.extracting", nil))

	file, err := os.Open(cacheFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log but don't fail if close fails during cleanup
			_ = closeErr
		}
	}()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := gzr.Close(); closeErr != nil {
			// Log but don't fail if close fails during cleanup
			_ = closeErr
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(claudeDir, "plugins", header.Name)

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				_ = outFile.Close()
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Clean removes plugin cache backup files
func Clean(repoPath string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Title(i18n.T("cache.clean.title", nil))

	cacheDir := filepath.Join(repoPath, "cache")

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Green, i18n.T("cache.clean.files", nil))
	totalSize := int64(0)
	for _, entry := range entries {
		if !entry.IsDir() {
			info, _ := entry.Info()
			//nolint:govet // format string contains non-constant parts
			ui.Println(ui.Gray, "  • "+entry.Name()+" ("+formatSize(info.Size())+")")
			totalSize += info.Size()
		}
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Green, i18n.T("cache.clean.total", map[string]interface{}{
		"Size": totalSize / 1024 / 1024,
	}))

	// Delete files
	for _, entry := range entries {
		if !entry.IsDir() {
			//nolint:errcheck // Cleanup failure is acceptable
			os.Remove(filepath.Join(cacheDir, entry.Name()))
		}
	}

	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Success(i18n.T("cache.clean.done", nil))
	return nil
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1fGB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
