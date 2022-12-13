package main

import (
	"fmt"
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
