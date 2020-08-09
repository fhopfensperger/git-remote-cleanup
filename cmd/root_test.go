package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getReposFromFile(t *testing.T) {
	repo1 := "https://github.com/fhopfensperger/amqp-sb-client.git"
	repo2 := "git@github.com:fhopfensperger/json-log-to-human-readable.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))
	f.WriteString(fmt.Sprintln(repo2))
	repos := getReposFromFile(fileName)
	assert.Equal(t, []string{repo1, repo2}, repos)

	os.Remove(fileName)
}

func Test_getReposFromFile_ignore_empty_and_hashtag_lines(t *testing.T) {
	repo1 := "https://github.com/fhopfensperger/amqp-sb-client.git"
	repo2 := ""
	repo3 := "https://github.com/fhopfensperger/json-log-to-human-readable.git"
	repo4 := "#git@github.com:fhopfensperger/json-log-to-human-readable.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))
	f.WriteString(fmt.Sprintln(repo2))
	f.WriteString(fmt.Sprintln(repo3))
	f.WriteString(fmt.Sprintln(repo4))
	repos := getReposFromFile(fileName)
	assert.Equal(t, []string{repo1, repo3}, repos)

	os.Remove(fileName)
}

func Test_getReposFromFile_file_not_found(t *testing.T) {
	fileName := "test.txt"
	os.Create(fileName)
	repos := getReposFromFile(fileName + "unknown")
	assert.Equal(t, []string(nil), repos)

	os.Remove(fileName)
}

func TestExecute_version(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--version"})
	Execute("0.0.0")
	out, _ := ioutil.ReadAll(b)
	assert.Equal(t, "v0.0.0\n", string(out))
}

func TestExecute_repos_from_args(t *testing.T) {
	cmd := rootCmd
	testRepos := []string{"git@github.com:fhopfensperger/my-repo.git"}
	cmd.SetArgs([]string{"branches", "-b", "release", "-r", testRepos[0]})
	Execute("0.0.0")

	assert.Equal(t, repos, testRepos)
	assert.Equal(t, filter, "release")
}

func TestExecute_repos_from_file(t *testing.T) {
	repo1 := "git@github.com:fhopfensperger/my-repo.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))

	cmd := rootCmd
	cmd.SetArgs([]string{"branches", "-b", "release", "-f", fileName})
	Execute("0.0.0")

	assert.Equal(t, repos, []string{repo1})
	assert.Equal(t, filter, "release")
	os.Remove(fileName)
}

func TestExecute_repos_or_file_must_be_defined(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		cmd := rootCmd
		cmd.SetArgs([]string{"branches", "-b", "release"})
		Execute("0.0.0")
		return
	}
	cmdtest := exec.Command(os.Args[0], "-test.run=TestExecute_repos_or_file_must_be_defined")
	cmdtest.Env = append(os.Environ(), "FLAG=1")
	err := cmdtest.Run()
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestExecute_delete_exclude_dry_run(t *testing.T) {
	repo1 := "git@github.com:fhopfensperger/my-repo.git"
	fileName := "test.txt"
	f, _ := os.Create(fileName)
	f.WriteString(fmt.Sprintln(repo1))

	excludeList := []string{"release/v1", "release/v2"}

	cmd := rootCmd
	cmd.SetArgs([]string{"delete", "-b", "release", "-f",
		fileName, "--dry-run", "-e", excludeList[0] + "," + excludeList[1]})
	Execute("0.0.0")

	assert.Equal(t, excludeList, excludes)
	assert.Equal(t, true, dryRun)
	os.Remove(fileName)
}
