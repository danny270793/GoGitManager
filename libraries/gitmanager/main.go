package gitmanager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"danny270793.github.com/gitmanager/libraries/shell"
)

type GitRepo struct {
	Path   string
	Status string
}

type Gitmanager struct {
	Shell shell.Shell
}

const (
	NO_GIT_REPO       string = "NO_GIT_REPO"
	PENDING_ADD       string = "PENDING_ADD"
	PENDING_COMMIT    string = "PENDING_COMMIT"
	PENDING_PUSH      string = "PENDING_PUSH"
	UP_TO_DATE        string = "UP_TO_DATE"
	CREATED_NEWS      string = "CREATED_NEWS"
	NOTHING_TO_COMMIT string = "NOTHING_TO_COMMIT"
	UNKNOW_STATE      string = "UNKNOW_STATE"
)

func (g *Gitmanager) Status(path string) (string, error) {
	stdout, stderr, _ := g.Shell.Execute(fmt.Sprintf("cd %s && git status --porcelain", path))
	response := strings.ToLower(stdout + stderr)

	if strings.Contains(response, "fatal: not a git repository") {
		return NO_GIT_REPO, nil
	}

	pendingToCommit := false
	responseLines := strings.Split(response, "\n")
	for _, line := range responseLines {
		if line == "" {
			continue
		}
		lineWithoutSpaces := strings.TrimSpace(response)
		lineItems := strings.Split(lineWithoutSpaces, " ")
		if lineItems[0] == "??" {
			return PENDING_ADD, nil
		}
		options := []string{"a", "d", "m"}
		position := sort.SearchStrings(options, lineItems[0])
		if position < len(options) && options[position] == lineItems[0] {
			pendingToCommit = true
		}
	}
	if pendingToCommit {
		return PENDING_COMMIT, nil
	}

	stdout, stderr, _ = g.Shell.Execute(fmt.Sprintf("cd %s && git status", path))
	response = strings.ToLower(stdout + stderr)
	if strings.Contains(response, "your branch is ahead of") {
		return PENDING_PUSH, nil
	} else if strings.Contains(response, "your branch is up") {
		return NOTHING_TO_COMMIT, nil
	} else if strings.Contains(response, "nothing to commit") {
		return NOTHING_TO_COMMIT, nil
	} else {
		return "", errors.New(fmt.Sprintf("status of the repo %s could not be obtained\n%s\n", path, response))
	}
}

func (g *Gitmanager) IsRepo(path string) (bool, error) {
	status, err := g.Status(path)
	if err != nil {
		return false, err
	}

	return status != NO_GIT_REPO, nil
}

func (g *Gitmanager) GetRepos(basePath string) ([]GitRepo, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []GitRepo{}, err
	}

	allRepos := []GitRepo{}

	for _, f := range files {
		fullPath := path.Join(basePath, f.Name())
		if f.IsDir() {
			status, err2 := g.Status(fullPath)
			if err2 != nil {
				return []GitRepo{}, err2
			}

			if status == NO_GIT_REPO {
				repos, err3 := g.GetRepos(fullPath)
				if err3 != nil {
					return []GitRepo{}, err3
				}
				allRepos = append(allRepos, repos...)
			} else {
				repo := GitRepo{
					Path:   fullPath,
					Status: status,
				}
				allRepos = append(allRepos, repo)
			}
		}
	}

	return allRepos, nil
}

func (g *Gitmanager) GetGitVersion() (string, error) {
	stdout, stderr, err := g.Shell.Execute("git --version")
	if err != nil {
		return "", errors.New(stderr)
	}

	tokens := strings.Split(stdout, " ")
	version := tokens[len(tokens)-1]
	return version, nil
}
