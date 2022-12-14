package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
	Merge Strategy when a same file is modified in two commits (provides no choice to the user):

	- Find the closest LCA of the two commits (let it be 'lca')
	- Find the LCS of 'lca' and 'x' (let it be 'lx')
	- Find the LCS of 'lca' and 'y' (let it be 'ly')
	- The base string now is the LCS of 'lx' and 'ly' (let it be 'lxy')
	- Give the base string line numbers and map the lines of 'lx' and 'ly' to it
	- Find the additions of 'lx' and 'x' and apply it to 'lxy'
	- Find the additions of 'ly' and 'y' and apply it to 'lxy'

*/

func Merge(user string, x, y *Commit, force bool) {
	lca := LCA(x, y)
	if lca == nil {
		panic("LCA not found")
	}

	baseObjects := map[string]string{} // relative path -> hash
	for _, object := range lca.Objects {
		baseObjects[object.RelativePath] = object.Hash
	}

	xObjects := map[string]string{}
	for _, object := range x.Objects {
		xObjects[object.RelativePath] = object.Hash
	}

	yObjects := map[string]string{}
	for _, object := range y.Objects {
		yObjects[object.RelativePath] = object.Hash
	}

	finalObjects := map[string]string{}
	filesWithConflicts := []string{}

	hasConflicts := false

	for relativePath, hash := range xObjects {
		if baseObjects[relativePath] == hash && yObjects[relativePath] == baseObjects[relativePath] {
			finalObjects[relativePath] = hash
		} else if baseObjects[relativePath] == hash && yObjects[relativePath] != baseObjects[relativePath] {
			// x is the same as lca, y is different
			finalObjects[relativePath] = yObjects[relativePath]
		} else if baseObjects[relativePath] != hash && yObjects[relativePath] == baseObjects[relativePath] {
			// y is the same as lca, x is different
			finalObjects[relativePath] = xObjects[relativePath]
		} else if baseObjects[relativePath] != hash && yObjects[relativePath] != baseObjects[relativePath] {
			// x and y are different from lca
			fmt.Println("Conflicts found in ", relativePath)
			hasConflicts = true
			if force {
				filesWithConflicts = append(filesWithConflicts, relativePath)
			}
		}
	}

	for relativePath, hash := range yObjects {
		if baseObjects[relativePath] == hash && xObjects[relativePath] == baseObjects[relativePath] {
			finalObjects[relativePath] = hash
		} else if baseObjects[relativePath] == hash && xObjects[relativePath] != baseObjects[relativePath] {
			// y is the same as lca, x is different
			finalObjects[relativePath] = xObjects[relativePath]
		} else if baseObjects[relativePath] != hash && xObjects[relativePath] == baseObjects[relativePath] {
			// x is the same as lca, y is different
			finalObjects[relativePath] = yObjects[relativePath]
		}
	}


	if !hasConflicts {
		for relativePath, hash := range finalObjects {
			CopyFile(".gogit/objects/"+hash, relativePath)
		}
		ApplyMerge(finalObjects)
		return
	}

	if !force {
		fmt.Println("Merge failed")
		return
	}

	for _, relativePath := range filesWithConflicts {
		fmt.Println("Resolving conflicts in ", relativePath)
		newContent := ResolveConflicts(baseObjects[relativePath], xObjects[relativePath], yObjects[relativePath])
		ioutil.WriteFile(relativePath, []byte(newContent), 0644)
		f, err := os.OpenFile(relativePath, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		hash := HashFile(f)
		CopyFile(relativePath, ".gogit/objects/"+hash)
		finalObjects[relativePath] = hash
	}
	ApplyMerge(finalObjects)
}

func ApplyMerge(objects map[string]string) {
	// delete all files except .gogit
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() != ".gogit" {
			os.RemoveAll(file.Name())
		}
	}

	// copy all files from objects to current directory
	for relativePath, hash := range objects {
		if hash != "" {
			CopyFile(".gogit/objects/"+hash, relativePath)
		}
	}
}

func ResolveConflicts(base, hash1, hash2 string) string {
	b, err := ioutil.ReadFile(".gogit/objects/" + base)
	if err != nil {
		panic(err)
	}
	linesBase := strings.Split(string(b), "\n")

	f1, err := ioutil.ReadFile(".gogit/objects/" + hash1)
	if err != nil {
		panic(err)
	}
	lines1 := strings.Split(string(f1), "\n")

	f2, err := ioutil.ReadFile(".gogit/objects/" + hash2)
	if err != nil {
		panic(err)
	}
	lines2 := strings.Split(string(f2), "\n")

	lx := LCS(linesBase, lines1)
	ly := LCS(linesBase, lines2)

	lxy := LCS(lx, ly)

	lxLineNums := make([]int, len(lx))
	for i := range lxLineNums {
		lxLineNums[i] = -1
	}

	lyLineNums := make([]int, len(ly))
	for i := range lyLineNums {
		lyLineNums[i] = -1
	}

	i1 := 0
	for i, line := range lxy {
		for i1 < len(lx) && lx[i1] != line {
			i1++
		}
		lxLineNums[i1] = i
	}

	i2 := 0
	for i, line := range lxy {
		for i2 < len(ly) && ly[i2] != line {
			i2++
		}
		lyLineNums[i2] = i
	}

	xAdditions := []struct {
		Line    string
		LineNum int
	}{}
	yAdditions := []struct {
		Line    string
		LineNum int
	}{}

	i1 = 0
	for i, line := range lx {
		for i1 < len(lines1) && lines1[i1] != line {
			xAdditions = append(xAdditions, struct {
				Line    string
				LineNum int
			}{Line: lines1[i1], LineNum: i - 1})
			i1++
		}
		i1++
	}
	for i1 < len(lines1) {
		xAdditions = append(xAdditions, struct {
			Line    string
			LineNum int
		}{Line: lines1[i1], LineNum: len(lx) - 1})
		i1++
	}

	i2 = 0
	for i, line := range ly {
		for i2 < len(lines2) && lines2[i2] != line {
			yAdditions = append(yAdditions, struct {
				Line    string
				LineNum int
			}{Line: lines2[i2], LineNum: i - 1})
			i2++
		}
		i2++
	}
	for i2 < len(lines2) {
		yAdditions = append(yAdditions, struct {
			Line    string
			LineNum int
		}{Line: lines2[i2], LineNum: len(ly) - 1})
		i2++
	}

	finalLines := []string{}
	i1 = 0
	i2 = 0
	for i, line := range lxy {
		for i1 < len(xAdditions) && xAdditions[i1].LineNum < i {
			finalLines = append(finalLines, xAdditions[i1].Line)
			i1++
		}
		for i2 < len(yAdditions) && yAdditions[i2].LineNum < i {
			finalLines = append(finalLines, yAdditions[i2].Line)
			i2++
		}
		finalLines = append(finalLines, line)
	}
	for i1 < len(xAdditions) {
		finalLines = append(finalLines, xAdditions[i1].Line)
		i1++
	}
	for i2 < len(yAdditions) {
		finalLines = append(finalLines, yAdditions[i2].Line)
		i2++
	}

	finalStringContent := strings.Join(finalLines, "\n")
	return finalStringContent
}

func LCA(x, y *Commit) *Commit {
	visX := map[string]bool{}
	visY := map[string]bool{}

	queueX := []*Commit{x}
	queueY := []*Commit{y}

	for len(queueX) > 0 || len(queueY) > 0 {
		if len(queueX) > 0 {
			curX := queueX[0]
			queueX = queueX[1:]
			if visY[curX.Hash] {
				return curX
			}
			visX[curX.Hash] = true
			queueX = append(queueX, curX.PrevCommits...)
		}

		if len(queueY) > 0 {
			curY := queueY[0]
			queueY = queueY[1:]
			if visX[curY.Hash] {
				return curY
			}
			visY[curY.Hash] = true
			queueY = append(queueY, curY.PrevCommits...)
		}
	}

	return nil
}
