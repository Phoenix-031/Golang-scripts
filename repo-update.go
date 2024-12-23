package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func findGitRepos(root string) []string {
	var repos []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
		}
		return nil
	})

	if err != nil {
		fmt.Printf("[%s] Error walking the file tree: %s\n", time.Now(), err)
	}

	return repos
}

func pullFromRepo(repoPath string) {
	fmt.Printf("[%s] Pulling updates in %s...\n", time.Now(), repoPath)

	if err := os.Chdir(repoPath); err != nil {
		fmt.Printf("[%s] Failed to change directory to %s: %s\n", time.Now(), repoPath, err)
		return
	}

	cmd := exec.Command("git", "branch", "--show-current")
	var branchOut bytes.Buffer
	cmd.Stdout = &branchOut
	if err := cmd.Run(); err != nil {
		fmt.Printf("[%s] Failed to get branch name in %s: %s\n", time.Now(), repoPath, err)
		return
	}
	branch := strings.TrimSpace(branchOut.String())

	if branch == "" {
		fmt.Printf("[%s] No active branch in %s, skipping.\n", time.Now(), repoPath)
		return
	}

	// Pull the latest changes
	cmd = exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[%s] Failed to pull updates in %s: %s\n", time.Now(), repoPath, err)
	}
}

func main() {
	root := "/home/debayan/MY/"
	fmt.Printf("[%s] Searching for git repositories in %s...\n", time.Now(), root)

	repos := findGitRepos(root)
	fmt.Printf("[%s] Found %d git repositories.\n", time.Now(), len(repos))

	// LOGGING ALL GIT REPO NAMES

	// for _, repo := range repos {
	// 	fmt.Printf("- %s\n", repo)
	// }

	for _, repo := range repos {
		pullFromRepo(repo)
	}
}