package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jiangtao/ccconfig/pkg/backup"
	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/git"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Claude Code configuration",
	RunE:  runBackup,
}

var (
	backupRepo       string
	backupProjects   []string
	backupAll        bool
	backupNoSettings bool
	backupNoCommands bool
	backupNoSkills   bool
	backupNoCommit   bool
)

func init() {
	GetRootCommand().AddCommand(backupCmd)

	backupCmd.Flags().StringP("repo", "r", "", "Config repository path")
	backupCmd.Flags().StringSliceP("projects", "p", nil, "Project directories to backup")
	backupCmd.Flags().BoolP("all-projects", "A", false, "Auto-scan common directories")
	backupCmd.Flags().Bool("no-settings", false, "Skip settings.json")
	backupCmd.Flags().Bool("no-commands", false, "Skip custom commands")
	backupCmd.Flags().Bool("no-skills", false, "Skip custom skills")
	backupCmd.Flags().Bool("no-commit", false, "Skip git commit")
}

func runBackup(cmd *cobra.Command, args []string) error {
	ui.Title(i18n.T("backup.title", nil))

	// Get paths
	repoPath := backupRepo
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

	// Ensure repo exists
	if err := os.MkdirAll(expandedRepo, 0755); err != nil {
		return fmt.Errorf("failed to create repo directory: %w", err)
	}

	// Step 1: Backup settings
	settingsBackup := backup.NewSettingsBackup(backupNoSettings)
	if err := settingsBackup.Backup(claudeDir, configDir); err != nil {
		ui.Error("Settings backup failed: %v", err)
	}

	// Step 2: Backup commands
	commandsBackup := backup.NewCommandsBackup(backupNoCommands)
	if err := commandsBackup.Backup(claudeDir, configDir); err != nil {
		ui.Error("Commands backup failed: %v", err)
	}

	// Step 3: Backup skills
	skillsBackup := backup.NewSkillsBackup(backupNoSkills)
	if err := skillsBackup.Backup(claudeDir, configDir); err != nil {
		ui.Error("Skills backup failed: %v", err)
	}

	// Step 4: Backup projects
	projectsBackup := backup.NewProjectsBackup(backupProjects, backupAll)
	if err := projectsBackup.Backup(configDir); err != nil {
		ui.Error("Projects backup failed: %v", err)
	}

	// Step 5: Git commit
	if !backupNoCommit {
		ui.Println(ui.Cyan, i18n.T("backup.steps.commit", nil))

		// Check if it's a git repo
		isGit, _ := git.IsGitRepo(expandedRepo)
		if !isGit {
			ui.Warning("Not a git repository, skipping commit")
		} else {
			// Check for changes
			hasChanges, _ := git.HasChanges(expandedRepo)
			if !hasChanges {
				ui.Warning(i18n.T("backup.messages.no_changes", nil))
			} else {
				// Add and commit
				if err := git.Add(expandedRepo, "config/"); err != nil {
					ui.Error("git add failed: %v", err)
				} else {
					commitMsg := fmt.Sprintf("chore: backup Claude Code configuration\n\n%s", time.Now().Format("2006-01-02 15:04:05"))
					if err := git.Commit(expandedRepo, commitMsg); err != nil {
						ui.Error("git commit failed: %v", err)
					} else {
						ui.Success(i18n.T("backup.messages.committed", nil))
						ui.Println(ui.Yellow, "  Tip: Run 'git push' to push to remote")
					}
				}
			}
		}
	}

	ui.Println(ui.BoldBlue, "\n=== Backup Complete ===\n")

	return nil
}
