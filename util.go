package main

import (
	"io"
	"os"
	"strconv"
	"strings"
)

func FileToString(f *os.File) string {
	str := ""
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		str += string(buf[:n])
	}
	return str
}

func CopyFile(src string, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	io.Copy(dstFile, srcFile)
}

func GetMessageFromCommitFile(commitHash string) string {
	commitFile, err := os.OpenFile(".gogit/commits/"+commitHash, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	commitString := FileToString(commitFile)
	lines := strings.Split(commitString, "\n")
	return lines[len(lines)-2]
}

func GetTimeFromCommitFile(commitHash string) int64 {
	commitFile, err := os.OpenFile(".gogit/commits/"+commitHash, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	commitString := FileToString(commitFile)
	lines := strings.Split(commitString, "\n")
	time, err := strconv.ParseInt(lines[len(lines)-3], 10, 64)
	if err != nil {
		panic(err)
	}
	return time
}

func AreStringsArraysEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
