package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

func HashString(s string, timestamp int64) string {
	h := sha512.New()
	io.WriteString(h, s)
	io.WriteString(h, fmt.Sprintf("%d", timestamp))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashFile(f *os.File) string {
	h := sha512.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}