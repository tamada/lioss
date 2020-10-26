package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

/*
1: not found
2: dir
3: regular file
*/
func existFlag(path string) int {
	stat, err := os.Stat(path)
	if err != nil {
		return 1
	}
	if stat.Mode().IsRegular() {
		return 2
	}
	if stat.IsDir() {
		return 3
	}
	return 4
}

func findParentPath(path string) (string, error) {
	if path != "." {
		return filepath.Dir(path), nil
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Dir(absPath), nil
}

func findParent(path string) (string, error) {
	parent, err := findParentPath(path)
	if parent == path {
		return "", fmt.Errorf("not found")
	}
	givenPath, err := findGitRepository(parent)
	if err != nil {
		err = fmt.Errorf("%s: not found", path)
	}
	return givenPath, nil
}

func readString(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	data, _ := ioutil.ReadAll(reader)
	strArray := strings.TrimSpace(string(data))
	for _, str := range strings.Split(strArray, "\n") {
		if !strings.HasPrefix(str, "gitdir: ") {
			continue
		}
		return strings.TrimPrefix(str, "gitdir: "), nil
	}
	return strArray, nil
}

func readLink(path string) (string, error) {
	relPath, err := readString(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(filepath.Dir(path), relPath)), nil
}

func findGitRepository(path string) (string, error) {
	target := filepath.Join(path, ".git")
	existFlag := existFlag(target)
	if existFlag == 1 { // not found
		parentPath, err := findParent(path)
		if err != nil {
			return "", fmt.Errorf("%s: not found", path)
		}
		return parentPath, nil
	}
	if existFlag == 2 { // regular file
		return readLink(target)
	}
	if existFlag == 3 {
		return target, nil
	}
	return "", fmt.Errorf("%s: not found", path)
}

func readCommitID(dir string) (string, error) {
	repoPath, err := findGitRepository(dir)
	if err != nil {
		return "", err
	}
	// fmt.Printf("findGitRepository(%s): %s\n", dir, repoPath)
	repository, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	head, err := repository.Reference("refs/heads/master", false)
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}
