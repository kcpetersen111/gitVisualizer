package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

//scan path searching for git repos
func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repos := recursiveScanFolder(folder)
	filepath := getDotFilePath()
	addNewSliceElementsToFile(filepath, repos)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

//returns a list of subfolders of folder ending with .git
//returns the base folder of the repo, parent of .git
//recursively searches in te subfolder by passing an existing folder slice
func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}

//recursive scan folder starts the search
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

//returns the dot file for the repos list.
// Creates it and the enclosing folder if it does not exist
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotFile := usr.HomeDir + "/.gogitlocalstats"

	return dotFile
}

func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// reads evert line of a file and will put it into a slice of strings
func parseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}

// opens the file located at filepath and will create it if it does not exist
func openFile(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		if os.IsNotExist(err) {
			//if the file does not exist
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			//any error that has nothing to do with the file not existing
			panic(err)
		}
	}
	return f
}

//join slices joins two slices
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

//true if slice contains value
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

//will over write all of the old repos in the file with the new list of repos
func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	ioutil.WriteFile(filePath, []byte(content), 0755)

}
