package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiangtao/ccconfig/pkg/config"
	"github.com/jiangtao/ccconfig/pkg/git"
	"github.com/jiangtao/ccconfig/pkg/i18n"
	"github.com/jiangtao/ccconfig/pkg/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize config repository",
	RunE:  runInit,
}

var (
	initRepo   string
	initGitURL string
	initNoGit  bool
)

func init() {
	GetRootCommand().AddCommand(initCmd)

	initCmd.Flags().StringP("repo", "r", "", "Repository path")
	initCmd.Flags().String("git-url", "", "Git remote URL")
	initCmd.Flags().Bool("no-git", false, "Don't initialize git")
}

func runInit(cmd *cobra.Command, args []string) error {
	ui.Title(i18n.T("init.title", nil))

	reader := bufio.NewReader(os.Stdin)

	// Get repository path
	repoPath := initRepo
	if repoPath == "" {
		fmt.Print(i18n.T("init.prompts.repo_path", nil) + " [~/cc-config]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			repoPath = "~/cc-config"
		} else {
			repoPath = input
		}
	}

	expandedRepo, err := config.ExpandPath(repoPath)
	if err != nil {
		return err
	}

	// Check if already exists
	if _, err := os.Stat(expandedRepo); err == nil {
		ui.Warning("Directory already exists: %s", expandedRepo)
	} else {
		// Create directory
		if err := os.MkdirAll(expandedRepo, 0755); err != nil {
			return fmt.Errorf("failed to create repository: %w", err)
		}
		ui.Success(i18n.T("init.messages.created", nil))
	}

	// Initialize git
	if !initNoGit {
		shouldInitGit := false

		if initGitURL != "" {
			shouldInitGit = true
		} else {
			fmt.Print(i18n.T("init.prompts.init_git", nil) + " ")
			input, _ := reader.ReadString('\n')
			input = strings.ToLower(strings.TrimSpace(input))
			shouldInitGit = input == "y" || input == "yes"
		}

		if shouldInitGit {
			// Check if already a git repo
			isGit, _ := git.IsGitRepo(expandedRepo)
			if !isGit {
				if err := git.Init(expandedRepo); err != nil {
					return fmt.Errorf("git init failed: %w", err)
				}
				ui.Success(i18n.T("init.messages.git_initialized", nil))
			}

			// Add remote if specified
			if initGitURL != "" {
				if err := git.AddRemote(expandedRepo, "origin", initGitURL); err != nil {
					ui.Warning("Failed to add git remote: %v", err)
				} else {
					ui.Success("Git remote 'origin' set to: %s", initGitURL)
				}
			}

			// Create .gitignore
			gitignorePath := filepath.Join(expandedRepo, ".gitignore")
			gitignoreContent := "cache/*.tar.gz\n"
			if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err == nil {
				ui.Success("Created .gitignore")
			}
		}
	}

	// Create config directory structure
	configDir := filepath.Join(expandedRepo, "config")
	for _, dir := range []string{
		filepath.Join(configDir, "commands"),
		filepath.Join(configDir, "skills"),
		filepath.Join(configDir, "project-configs"),
		filepath.Join(expandedRepo, "cache"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			ui.Warning("Failed to create directory: %s", dir)
		}
	}

	ui.Success("Repository structure created at: %s", expandedRepo)

	return nil
}
