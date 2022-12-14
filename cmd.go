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
		commits := *GetAllCommits()
		var newCommit *Commit
		message := strings.Join(args[0:], " ")
		if len(commits) == 0 {
			newCommit = CreateCommit("user", message, []*Commit{})
			CreateBranch("MASTER", newCommit.Hash)
			SaveHeadBranch("MASTER")
		} else {
			head := GetHead()
			headCommit := GetCommit(head)
			newCommit = CreateCommit("user", message, []*Commit{headCommit})
		}
		if newCommit != nil {
			commits = append(commits, *newCommit)
			SaveHead(newCommit)
			UpdateHeadBranch(newCommit.Hash)
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

var checkoutBranchCmd = &cobra.Command{
	Use:   "cb",
	Short: "Checkout a branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		commitHash := GetBranchCommit(branchName)
		if commitHash != "" {
			commit := GetCommit(commitHash)
			ApplyCommit(commit)
			SaveHead(commit)
			SaveHeadBranch(branchName)
		} else {
			fmt.Println("Branch not found")
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

var mergeBranchesCmd = &cobra.Command{
	Use:   "mb",
	Short: "Merges two branches",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch1 := GetBranchCommit(args[0])
		branch2 := GetBranchCommit(GetHeadBranch())
		if branch1 == "" || branch2 == "" || branch1 == branch2 {
			fmt.Println("Merge not possible")
			return
		}
		commit1 := GetCommit(branch1)
		commit2 := GetCommit(branch2)
		Merge("user", commit1, commit2, true)
		newCommit := CreateCommit("user", "Merge "+args[0]+" into "+GetHeadBranch(), []*Commit{commit1, commit2})
		if newCommit != nil {
			commits := *GetAllCommits()
			commits = append(commits, *newCommit)
			SaveHead(newCommit)
			UpdateHeadBranch(newCommit.Hash)
			SaveConfig(&commits)
		}
	},
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Create/Delete/Rename a branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetAllCommits()
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
			SaveHeadBranch(branchName)
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

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collection",
	Run: func(cmd *cobra.Command, args []string) {
		visited := map[string]bool{}
		commits := *GetAllCommits()
		queue := []string{}
		queue = append(queue, *GetAllBranchHeads()...)
		for len(queue) != 0 {
			commitHash := queue[0]
			queue = queue[1:]
			visited[commitHash] = true
			commit := GetCommit(commitHash)
			if commit == nil {
				continue
			}
			for _, prevCommit := range commit.PrevCommits {
				queue = append(queue, prevCommit.Hash)
			}
		}
		var allVisCommits []Commit
		for _, commit := range commits {
			if _, ok := visited[commit.Hash]; !ok {
				commit.DeleteCommit()
			} else {
				allVisCommits = append(allVisCommits, commit)
			}
		}
		SaveConfig(&allVisCommits)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	Initialize()
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(checkoutBranchCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(searchCommitCmd)
	rootCmd.AddCommand(mergeCommitsCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(mergeBranchesCmd)
	rootCmd.AddCommand(gcCmd)
}
