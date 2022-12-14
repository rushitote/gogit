package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Commit struct {
	User        string    // User who made the commit
	Hash        string    // Hash of the commit
	Objects     []Object  // Objects (files) that were changed
	PrevCommits []*Commit // Previous commit
	Time        int64     // Time of the commit
	Message     string    // Message of the commit
}

func (c *Commit) Serialize() string {
	var prevCommitHashes []string
	for _, commit := range c.PrevCommits {
		prevCommitHashes = append(prevCommitHashes, commit.Hash)
	}
	prevCommitHashConcat := strings.Join(prevCommitHashes, ",")
	var objectsString string
	for _, object := range c.Objects {
		objectsString += object.Serialize() + ","
	}
	return fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n%d\n%s", c.User, c.Hash, len(c.Objects), objectsString, prevCommitHashConcat, c.Time, c.Message)
}

func DeserializeCommit(s string) *Commit {
	var user, hash, prevCommitHashConcat string
	var objectCount int
	var objectsString string
	var time int64
	fmt.Sscanf(s, "%s\n%s\n%d\n%s\n%s\n%d", &user, &hash, &objectCount, &objectsString, &prevCommitHashConcat, &time)
	message := GetMessageFromCommitFile(hash)
	objects := strings.Split(objectsString, ",")
	var objectsList []Object
	for i := 0; i < objectCount; i++ {
		objectsList = append(objectsList, *DeserializeObject(objects[i]))
	}
	return &Commit{user, hash, objectsList, GetMultipleCommits(prevCommitHashConcat), time, message}
}

func GetMultipleCommits(hashesConcat string) []*Commit {
	var commits []*Commit
	hashes := strings.Split(hashesConcat, ",")
	for _, hash := range hashes {
		commit := GetCommit(hash)
		if commit != nil {
			commits = append(commits, commit)
		}
	}
	return commits
}

func GetCommit(hash string) *Commit {
	if hash == "" {
		return nil
	}
	f, err := os.OpenFile(".gogit/commits/"+hash, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	commitString := FileToString(f)
	return DeserializeCommit(commitString)
}

func SaveCommit(c *Commit) {
	err := os.MkdirAll(".gogit/commits", 0755)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(".gogit/commits/"+c.Hash, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, c.Serialize())
	if err != nil {
		panic(err)
	}
}

func GetSnapshot() []Object {
	var objects []Object
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() == ".gogit" {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				f, err := os.OpenFile(path, os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				hash := HashFile(f)
				CopyFile(path, ".gogit/objects/"+hash)
				objects = append(objects, Object{hash, path})
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return objects
}

func HashObjects(objects []Object, time int64) string {
	var hash string
	for _, object := range objects {
		hash += object.Hash
	}
	return HashString(hash, time)
}

func CreateCommit(user string, message string, parentCommits []*Commit) *Commit {
	objects := GetSnapshot()

	//check if the commit is the same as the previous one
	if len(parentCommits) != 0 {
		same := true
		currSnap := map[string]string{} // relative path -> hash
		for _, object := range objects {
			currSnap[object.RelativePath] = object.Hash
		}
		if len(objects) != len(parentCommits[0].Objects) {
			same = false
		}
		for _, object := range parentCommits[0].Objects {
			if currSnap[object.RelativePath] != object.Hash {
				same = false
			}
		}
		if same {
			fmt.Println("Nothing to commit")
			return nil
		}
	}

	currTime := time.Now().Unix()

	commit := Commit{
		User:        user,
		Hash:        HashObjects(objects, currTime),
		Objects:     objects,
		PrevCommits: parentCommits,
		Time:        currTime,
		Message:     message,
	}
	SaveCommit(&commit)
	return &commit
}

func ApplyCommit(c *Commit) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".gogit" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, object := range c.Objects {
		CopyFile(".gogit/objects/"+object.Hash, object.RelativePath)
	}
}

func (c *Commit) LogCommit() {
	fmt.Printf("Commit: %s\n", c.Hash)
	fmt.Printf("Author: %s\n", c.User)
	fmt.Printf("Date: %s\n", time.Unix(c.Time, 0))
	fmt.Printf("Message: %s\n\n", c.Message)
}

func (c *Commit) DeleteCommit(){
	err := os.Remove(".gogit/commits/"+c.Hash)
	if err != nil {
		panic(err)
	}
	// TODO: delete objects
}

func GetHead() string {
	head, err := os.OpenFile(".gogit/HEAD", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer head.Close()
	commitHash := FileToString(head)
	return commitHash
}
