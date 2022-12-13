package main

import (
	"fmt"
	"strings"
)

type Object struct {
	Hash         string
	RelativePath string
}

func (o *Object) Serialize() string {
	return fmt.Sprintf("%s|%s", o.Hash, o.RelativePath)
}

func DeserializeObject(s string) *Object {
	var hash, relativePath string
	strs := strings.Split(s, "|")
	hash = strs[0]
	relativePath = strs[1]
	return &Object{hash, relativePath}
}
