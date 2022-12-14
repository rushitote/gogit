package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CreateBranch(branch, commit string) {
	f, err := os.Create(".gogit/branches/" + branch)
	if err != nil {
		panic(err)
	}
	f.WriteString(commit)
}

func GetBranchCommit(branch string) string {
	f, err := os.Open(".gogit/branches/" + branch)
	if err != nil {
		panic(err)
	}
	return FileToString(f)
}

func DeleteBranch(branch string) {
	os.Remove(".gogit/branches/" + branch)
}

func RenameBranch(oldBranch, newBranch string) {
	os.Rename(".gogit/branches/"+oldBranch, ".gogit/branches/"+newBranch)
}

func LogAllBranches() {
	files, err := os.ReadDir(".gogit/branches")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Print(file.Name() + " ")
	}
	fmt.Print("\n")
}

func GetAllBranchHeads() *[]string {
	files, err := os.ReadDir(".gogit/branches")
	if err != nil {
		panic(err)
	}
	var commitHashes []string
	for _, file := range files {
		commitHashes = append(commitHashes, GetBranchCommit(file.Name()))
	}
	return &commitHashes
}

func GetHeadBranch() string {
	f, err := os.Open(".gogit/HEAD_BRANCH")
	if err != nil {
		panic(err)
	}
	return FileToString(f)
}

func SaveHeadBranch(branch string) {
	ioutil.WriteFile(".gogit/HEAD_BRANCH", []byte(branch), 0666)
}

func UpdateHeadBranch(commitHash string) {
	branch := GetHeadBranch()
	CreateBranch(branch, commitHash)
}
