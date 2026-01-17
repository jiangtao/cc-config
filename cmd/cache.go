package cmd

import (
	"os"

	"github.com/jiangtao/ccconfig/pkg/cache"
	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage plugin cache",
}

var (
	cacheRepo string
)

func init() {
	GetRootCommand().AddCommand(cacheCmd)

	cacheCmd.PersistentFlags().StringP("repo", "r", "", "Config repository path")

	// Add subcommands
	cacheCmd.AddCommand(cacheBackupCmd)
	cacheCmd.AddCommand(cacheRestoreCmd)
	cacheCmd.AddCommand(cacheCleanCmd)
}

var cacheBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup plugin cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := getCacheRepo()
		if err := cache.Backup(repoPath); err != nil {
			if os.IsNotExist(err) {
				ui.Warning("Plugin cache not found")
				return nil
			}
			return err
		}
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Success(i18n.T("cache.backup.done", nil))
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Println(ui.Yellow, i18n.T("cache.backup.hint", nil))
		return nil
	},
}

var cacheRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore plugin cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := getCacheRepo()
		if err := cache.Restore(repoPath); err != nil {
			if os.IsNotExist(err) {
				ui.Warning("Plugin cache backup not found")
				ui.Println(ui.Yellow, "  Run 'ccconfig cache backup' first")
				return nil
			}
			return err
		}
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Success(i18n.T("cache.restore.done", nil))
		return nil
	},
}

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean plugin cache backups",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := getCacheRepo()
		return cache.Clean(repoPath)
	},
}

func getCacheRepo() string {
	repoPath := cacheRepo
	if repoPath == "" {
		repoPath = config.Get().Repo
	}
	if repoPath == "" {
		repoPath = "~/cc-config"
	}
	expanded, err := config.ExpandPath(repoPath)
	if err != nil {
		return repoPath
	}
	return expanded
}
