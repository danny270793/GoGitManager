package main

import (
	"flag"
	"fmt"
	"os"

	"danny270793.github.com/gitmanager/libraries/directory"
	"danny270793.github.com/gitmanager/libraries/gitmanager"
	"danny270793.github.com/gitmanager/libraries/shell"
)

var VERSION string = "1.0.1"

func main() {
	// parse arguments
	versionFlag := flag.Bool("version", false, "Show the version of gitmanager")
	helpFlag := flag.Bool("help", false, "Shows the current help")
	flag.Parse()

	if *versionFlag {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if *helpFlag {
		fmt.Printf("analyse paths to search for git repos and get the status of them\n")
		fmt.Printf("usage: gitmanager [options...] paths...\n")
		fmt.Printf("\noptions:\n")
		flag.PrintDefaults()
		fmt.Printf("\nexample:\n\tgitmanager /home/user/gitfolders\n")
		os.Exit(0)
	}

	g := gitmanager.Gitmanager{
		Shell: shell.NewBashShell(),
	}

	// check if git is installed
	_, err := g.GetGitVersion()
	if err != nil {
		fmt.Printf("%s\ncheck if git is available on the path\n", err)
		return
	}

	// check if the paths are valid
	hasErrors := false
	for _, basePath := range os.Args[1:] {
		isDirectory, err := directory.IsDirectory(basePath)
		if err != nil {
			hasErrors = true
			fmt.Printf("%s\n", err)
		} else if !isDirectory {
			hasErrors = true
			fmt.Printf("%s is not a directory\n", basePath)
		}
	}
	if hasErrors {
		return
	}

	// search the repositories
	var totalSize int64 = 0
	fmt.Printf("%-20s\t%-10s\t%s\n", "Status", "Size", "Path")
	for _, basePath := range os.Args[1:] {
		repos, err := g.GetRepos(basePath)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		for _, repo := range repos {
			size, err := directory.GetSize(repo.Path)
			if err != nil {
				panic(fmt.Sprintf("Error getting folder size: %s", err))
			}
			totalSize += size

			sizeAsString := directory.SizeToString(size)
			fmt.Printf("%-20s\t%10s\t%s\n", repo.Status, sizeAsString, repo.Path)
		}
	}
	totalSizeAsString := directory.SizeToString(totalSize)
	fmt.Printf("\ntotal size: %s\n", totalSizeAsString)
}
