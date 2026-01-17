package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/git"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/restore"
	"github.com/jiangtao/ccconfig/pkg/ui"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Claude Code configuration",
	RunE:  runRestore,
}

var (
	restoreRepo   string
	restoreNoPull bool
	restoreCache  bool
	restoreDryRun bool
)

func init() {
	GetRootCommand().AddCommand(restoreCmd)

	restoreCmd.Flags().StringP("repo", "r", "", "Config repository path")
	restoreCmd.Flags().Bool("pull", true, "Auto git pull")
	restoreCmd.Flags().Bool("restore-cache", false, "Restore plugin cache")
	restoreCmd.Flags().Bool("dry-run", false, "Preview mode, don't write")
}

func runRestore(cmd *cobra.Command, args []string) error {
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Title(i18n.T("restore.title", nil))

	if restoreDryRun {
		ui.Warning("Dry run mode - no changes will be made")
	}

	// Get paths
	repoPath := restoreRepo
	if repoPath == "" {
		repoPath = config.Get().Repo
	}
	if repoPath == "" {
		repoPath = "~/cc-config"
	}
	expandedRepo, err := config.ExpandPath(repoPath)
	if err != nil {
		return err
	}

	claudeDir, err := config.GetClaudeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(expandedRepo, "config")

	// Check if repo exists
	if _, err := os.Stat(expandedRepo); os.IsNotExist(err) {
		return fmt.Errorf("config repository not found: %s", expandedRepo)
	}

	// Step 1: Git pull
	if restoreNoPull {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Skipped(i18n.T("restore.steps.update", nil))
	} else {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Println(ui.Cyan, i18n.T("restore.steps.update", nil))
		isGit, _ := git.IsGitRepo(expandedRepo)
		if !isGit {
			ui.Warning("Not a git repository, skipping pull")
		} else if !restoreDryRun {
			if err := git.Pull(expandedRepo); err != nil {
				ui.Warning("Git pull failed: %v (continuing)", err)
			} else {
				ui.Success("Repository updated")
			}
		}
	}

	// Step 2: Create Claude directories
	//nolint:govet // msg is from i18n.T() which returns runtime strings
	ui.Println(ui.Cyan, i18n.T("restore.steps.create_dirs", nil))
	if !restoreDryRun {
		for _, dir := range []string{"commands", "skills", "plugins"} {
			if err := os.MkdirAll(filepath.Join(claudeDir, dir), 0755); err != nil {
				return err
			}
		}
	}
	ui.Success("Directories created")

	// Step 3: Restore settings
	settingsRestore := restore.NewSettingsRestore()
	if !restoreDryRun {
		if err := settingsRestore.Restore(configDir, claudeDir); err != nil {
			ui.Error("Settings restore failed: %v", err)
		}
	}

	// Step 4: Restore commands
	commandsRestore := restore.NewCommandsRestore()
	if !restoreDryRun {
		if err := commandsRestore.Restore(configDir, claudeDir); err != nil {
			ui.Error("Commands restore failed: %v", err)
		}
	}

	// Step 5: Restore skills
	skillsRestore := restore.NewSkillsRestore()
	if !restoreDryRun {
		if err := skillsRestore.Restore(configDir, claudeDir); err != nil {
			ui.Error("Skills restore failed: %v", err)
		}
	}

	// Step 6: Restore cache
	if restoreCache {
		//nolint:govet // msg is from i18n.T() which returns runtime strings
		ui.Println(ui.Cyan, i18n.T("restore.steps.cache", nil))
		cacheFile := filepath.Join(expandedRepo, "cache", "plugins-cache.tar.gz")
		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			ui.Warning("Plugin cache not found")
		} else {
			ui.Warning("Plugin cache restore not yet implemented")
		}
	}

	ui.Println(ui.BoldBlue, "\n=== Restore Complete ===\n")

	return nil
}
