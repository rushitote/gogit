package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func Initialize() {
	os.MkdirAll(".gogit", 0755)
	os.MkdirAll(".gogit/objects", 0755)
	os.MkdirAll(".gogit/commits", 0755)
	os.MkdirAll(".gogit/branches", 0755)
	os.OpenFile(".gogit/config", os.O_RDONLY|os.O_CREATE, 0666)
}

func GetAllCommits() *[]Commit {
	configFile, err := ioutil.ReadFile(".gogit/config")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(configFile), "\n")

	var commits []Commit
	for _, line := range lines {
		commit := GetCommit(line)
		if commit != nil {
			commits = append(commits, *commit)
		}
	}

	return &commits
}

func SaveConfig(commits *[]Commit) {
	var config string
	for _, commit := range *commits {
		config += commit.Hash + "\n"
	}
	ioutil.WriteFile(".gogit/config", []byte(config), 0666)
}

func SaveHead(commit *Commit) {
	ioutil.WriteFile(".gogit/HEAD", []byte(commit.Hash), 0666)
}
