package main

import (
	"sort"

	"github.com/hbollon/go-edlib"
)

func SortCommits(commits []Commit, query string, alg edlib.Algorithm) []Commit {
	sort.Slice(commits, func(i, j int) bool {
		x, err := edlib.StringsSimilarity(commits[i].Message, query, alg)
		if err != nil {
			panic(err)
		}
		y, err := edlib.StringsSimilarity(commits[j].Message, query, alg)
		if err != nil {
			panic(err)
		}
		return x > y
	})

	// select where similarity is greater than 0.5

	last := sort.Search(len(commits), func(i int) bool {
		x, _ := edlib.StringsSimilarity(commits[i].Message, query, edlib.Levenshtein)
		return x < 0.5
	})

	return commits[:last]
}

func SortCommitsByLevenshteinSim(commits []Commit, query string) []Commit {
	return SortCommits(commits, query, edlib.Levenshtein)
}

func SortCommitsByCosineSim(commits []Commit, query string) []Commit {
	return SortCommits(commits, query, edlib.Cosine)
}

func SortCommitsByJaccardSim(commits []Commit, query string) []Commit {
	return SortCommits(commits, query, edlib.Jaccard)
}

func SortCommitsByLCSSim(commits []Commit, query string) []Commit {
	return SortCommits(commits, query, edlib.Lcs)
}
