package main

import (
	"math"
)

func lcsProcess(runeStr1, runeStr2 []string) [][]int {
	// 2D Array that will contain str1 and str2 LCS
	lcsMatrix := make([][]int, len(runeStr1)+1)
	for i := 0; i <= len(runeStr1); i++ {
		lcsMatrix[i] = make([]int, len(runeStr2)+1)
		for j := 0; j <= len(runeStr2); j++ {
			lcsMatrix[i][j] = 0
		}
	}

	for i := 1; i <= len(runeStr1); i++ {
		for j := 1; j <= len(runeStr2); j++ {
			if runeStr1[i-1] == runeStr2[j-1] {
				lcsMatrix[i][j] = lcsMatrix[i-1][j-1] + 1
			} else {
				lcsMatrix[i][j] = int(math.Max(float64(lcsMatrix[i][j-1]), float64(lcsMatrix[i-1][j])))
			}
		}
	}

	return lcsMatrix
}

func LCS(str1, str2 []string) []string {
	if len(str1) == 0 || len(str2) == 0 {
		panic("Can't process and backtrack any LCS with empty string")
	} else if AreStringsArraysEqual(str1, str2) {
		return str1
	}

	return processLCSBacktrack(str1, str2, lcsProcess(str1, str2), len(str1), len(str2))
}

func processLCSBacktrack(str1, str2 []string, lcsMatrix [][]int, m, n int) []string {
	// Convert strings to rune array to handle no-ASCII characters
	if m == 0 || n == 0 {
		return nil
	} else if str1[m-1] == str2[n-1] {
		return append(processLCSBacktrack(str1, str2, lcsMatrix, m-1, n-1), string(str1[m-1]))
	} else if lcsMatrix[m][n-1] > lcsMatrix[m-1][n] {
		return processLCSBacktrack(str1, str2, lcsMatrix, m, n-1)
	}

	return processLCSBacktrack(str1, str2, lcsMatrix, m-1, n)
}
