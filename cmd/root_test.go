package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getReposFromFile(t *testing.T) {
	repo1 := "git@github.com:fhopfensperger/amqp-sb-client.git"
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
	repo1 := "git@github.com:fhopfensperger/amqp-sb-client.git"
	repo2 := ""
	repo3 := "git@github.com:fhopfensperger/json-log-to-human-readable.git"
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
