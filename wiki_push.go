package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

var tmpDir string = "tmpDir"
var workingDir string

func main() {
	var err error
	workingDir, err = prepareTmpDir(tmpDir)
	if err != nil {
		fmt.Printf(err.Error())
		return
	} else {
		fmt.Printf("Temp. dir created: %v\n", workingDir)
	}

	err = getWiki()
	if err != nil {
		fmt.Printf("error preparing the wiki git repo: %v\n", err)
		return
	}

	err = clearWiki()
	if err != nil {
		fmt.Printf("error cleansing git repo: %v\n", err)
		return
	}

	os.Chdir("./wiki")

	err = filepath.WalkDir(".", CopyFilesToTemp)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", tmpDir, err)
		return
	}

	os.Chdir("../")

	err = pushToWiki()
	if err != nil {
		fmt.Printf("Error with the wiki repo\n")
	}
}

func prepareTmpDir(tmpDir string) (string, error) {
	tempPath, err := os.MkdirTemp("", "")

	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	return tempPath, nil
}

func getWiki() error {

	cmd := exec.Command("git", "init")
	cmd.Dir = workingDir
	stdOutStdErr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error initializing git repo: %v\n", err)
	}
	fmt.Printf("%s\n", string(stdOutStdErr))

	cmd = exec.Command("git", "config", "user.name", os.Getenv("GITHUB_ACTOR"))
	cmd.Dir = workingDir
	stdOutStdErr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error configuring git repo: %v\n", err)
	}
	fmt.Printf("%s\n", string(stdOutStdErr))

	cmd = exec.Command("git", "config", "user.email", fmt.Sprintf("%s@wiki-push-plugin.com", os.Getenv("GITHUB_ACTOR")))
	cmd.Dir = workingDir
	stdOutStdErr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error configuring git repo: %v\n", err)
	}
	fmt.Printf("%s\n", string(stdOutStdErr))

	cmd = exec.Command("git", "pull", GetURL())
	cmd.Dir = workingDir
	stdOutStdErr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error configuring git repo: %v\n", err)
	}
	fmt.Printf("%s\n", string(stdOutStdErr))

	return nil
}

func clearWiki() error {
	err := filepath.WalkDir(workingDir, ClearWikiFiles)
	if err != nil {
		return fmt.Errorf("error walking the path %q: %v\n", tmpDir, err)
	}
	return nil
}

func ClearWikiFiles(path string, info fs.DirEntry, err error) error {
	subDirToSkip := ".git"

	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if info.IsDir() && info.Name() == subDirToSkip {
		fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		return filepath.SkipDir
	}
	if info.IsDir() == false {
		os.Remove(path)
	}
	return nil
}

func CopyFilesToTemp(path string, info fs.DirEntry, err error) error {
	subDirToSkip := ""

	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	if info.IsDir() && info.Name() == subDirToSkip {
		fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		return filepath.SkipDir
	}
	joinedPath := filepath.Join(workingDir, path)
	if info.IsDir() {
		fmt.Printf("visited dir: %q\n", path)
		cmd := exec.Command("mkdir", "-p", joinedPath)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}
	} else {

		fmt.Printf("visited file : %q\n", path)
		fmt.Printf("moving file to %q\n", joinedPath)

		if err != nil {
			return fmt.Errorf("Error making directory: %q", joinedPath)
		}
		cmd := exec.Command("cp", path, joinedPath)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Command finished with error: %v\n", err)
		}
	}

	return nil
}

func pushToWiki() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = workingDir
	stdOutStdErr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error staging files: %v\n", err)
	}
	fmt.Printf("%s\n", stdOutStdErr)

	cmd = exec.Command("git", "commit", "-m", "Automatically pushed commit")
	cmd.Dir = workingDir
	stdOutStdErr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error creating commit: %v\n", err)
	}
	fmt.Printf("%s\n", stdOutStdErr)

	cmd = exec.Command("git", "push", "--set-upstream", GetURL(), "master")
	cmd.Dir = workingDir
	stdOutStdErr, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error pushing branch: %v\n", err)
	}
	fmt.Printf("%s\n", stdOutStdErr)

	return nil
}

func GetURL() string {
	return fmt.Sprintf("https://%s@%s/%s.wiki.git",
		os.Getenv("GH_PERSONAL_ACCESS_TOKEN"),
		os.Getenv("GITHUB_SERVER_URL")[8:],
		os.Getenv("GITHUB_REPOSITORY"))
}
