package libgodone

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Returns the current time and date in the format of
// Day, Nth Month, Year (eg Tuesday, 11th September, 2001).
func CurrentTimeAndDate() string {
	currentTime := time.Now()

	// everything converted to a string for the sake of consistency
	day := fmt.Sprint(currentTime.Day())
	weekday := currentTime.Weekday().String()
	month := currentTime.Month().String()
	year := fmt.Sprint(currentTime.Year())

	// 1st, 2nd, 3rd or Nth suffix of the day
	switch day {
	case "1":
		day += "st"

	case "2":
		day += "nd"

	case "3":
		day += "rd"

	default:
		day += "th"
	}

	return fmt.Sprintf("%s, %s %s, %s", weekday, day, month, year)
}

func Contains(list []string, element string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

func FullPathToBuildTarget(target Target) string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		Die(RUNTIME, "Unable to reach the CWD!")
	}
	return fmt.Sprintf("%s/build/%s/", workingDirectory, TargetAsString(target))
}

// Removes the path of the working directory and only keeps the last folder.
// Used to get the name of the current folder when running godone with 'init'.
func GetwdBaseName() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		Die(RUNTIME, "Unable to reach the CWD!")
	}
	return filepath.Base(workingDirectory)
}

// Recursively collects the names of all files in a given directory. No exceptions made.
func CollectAllFilesFromDir(directory string) []string {
	dir, err := os.Open("src/")
	if err != nil {
		Die(RUNTIME, "Unable to reach src directory!")
	}
	defer dir.Close()

	dirContents, err := dir.ReadDir(0)
	if err != nil {
		Die(RUNTIME, "Unable to reach src directory!")
	}

	var files []string

	for _, v := range dirContents {
		if v.IsDir() {
			files = append(files, CollectAllFilesFromDir(v.Name())...)
		} else {
			files = append(files, v.Name())
		}
	}

	return files
}

// Returns all of the files listed in the build unit of the target makefile.
func CollectAllFilesFromMakefile(target Target) []string {
	file, err := os.Open(FullPathToBuildTarget(target) + "/Makefile")
	if err != nil {
		Die(RUNTIME, "Unable to update target Makefile!")
	}
	defer file.Close()

	files := []string{"main"} // contains main by default
	substring := "build: main"
	scanner := bufio.NewScanner(file)
	foundSubstring := false

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substring) {
			foundSubstring = true
			break
		}
	}

	if !foundSubstring {
		Die(RUNTIME, "Unable to update target Makefile!")
	}

	return append(files, strings.Split(scanner.Text(), " ")[2:]...)
}
