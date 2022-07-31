package main

import (
	"time"

	"github.com/go-git/go-git/plumbing/object"
	"gopkg.in/src-d/go-git.v4"
)

const (
	daysInLastSixMonths = 183
	outOfRange          = 99999
)

//given an email will return the commits made in the last 6 months
func processRepositories(email string) map[int]int {
	filepath := getDotFilePath()
	repos := parseFileLinesToSlice(filepath)
	daysInMap := daysInLastSixMonths

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits
}

//given a repository get the commits and put them in map
func fillCommits(email, path string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	//get the head of the repo
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	//git the commit history starting from the head
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	//iterate the commits
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset
		if c.Author.Email != email {
			return nil
		}
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return commits
}

func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}
func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

// determines how many days are missing to fill in
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()
	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1

	}
	return offset
}

//generate a nice graph
func stats(email string) {

	commits := processRepositories(email)
	printCommitsStats(commits)
}
