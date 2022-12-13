package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogit",
	Short: "A simple VCS",
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes to the repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commits := *Initialize()
		var newCommit *Commit
		message := strings.Join(args[0:], " ")
		if len(commits) == 0 {
			newCommit = CreateCommit("user", message, []*Commit{})
		} else {
			head := GetHead()
			headCommit := GetCommit(head)
			newCommit = CreateCommit("user", message, []*Commit{headCommit})
		}
		if newCommit != nil {
			commits = append(commits, *newCommit)
			SaveHead(newCommit)
		}
		SaveConfig(&commits)
	},
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout a commit",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitHash := args[0]
		commit := GetCommit(commitHash)
		if commit != nil {
			ApplyCommit(commit)
			SaveHead(commit)
		} else {
			fmt.Println("Commit not found")
		}
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Run: func(cmd *cobra.Command, args []string) {
		commitHash := GetHead()
		queue := []string{commitHash}
		for len(queue) != 0 {
			commitHash := queue[0]
			queue = queue[1:]
			commit := GetCommit(commitHash)
			if commit == nil {
				continue
			}
			commit.LogCommit()
			for _, prevCommit := range commit.PrevCommits {
				queue = append(queue, prevCommit.Hash)
			}
		}
	},
}

var searchCommitCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a commit",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args[0:], " ")
		commitHash := GetHead()
		queue := []string{commitHash}
		for len(queue) != 0 {
			commitHash := queue[0]
			queue = queue[1:]
			commit := GetCommit(commitHash)
			if commit == nil {
				continue
			}
			if strings.Contains(commit.Message, query) {
				commit.LogCommit()
			}
			for _, prevCommit := range commit.PrevCommits {
				queue = append(queue, prevCommit.Hash)
			}
		}
	},
}

var mergeCommitsCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merges two commits",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		commit1 := GetCommit(args[0])
		commit2 := GetCommit(args[1])
		if commit1 == nil || commit2 == nil {
			fmt.Println("Commit not found")
			return
		}
		Merge("user", commit1, commit2, true)
	},
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Create/Delete/Rename a branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		Initialize()
		flag := args[0]
		var branchName string	
		if len(args) > 1 {
			branchName = args[1]
		}
		if flag == "create" {
			var commitHash string
			if len(args) != 3 {
				commitHash = GetHead()
			} else {
				commitHash = args[2]
			}
			CreateBranch(branchName, commitHash)
		} else if flag == "delete" {
			DeleteBranch(branchName)
		} else if flag == "rename" {
			RenameBranch(branchName, args[2])
		} else if flag == "all" {
			LogAllBranches()
		} else {
			fmt.Println("Invalid flag")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(searchCommitCmd)
	rootCmd.AddCommand(mergeCommitsCmd)
	rootCmd.AddCommand(branchCmd)
}
