/*
Copyright © 2020 Florian Hopfensperger <f.hopfensperger@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var repos []string
var filter string
var fileName string
var pat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-remote-cleanup",
	Short: "Simple command line utility to get and delete branches from a remote git repo",
	Long:  `Simple command line utility to get and delete branches from a remote git repo`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pf := rootCmd.PersistentFlags()
	pf.StringSliceP("repos", "r", []string{}, "Git Repo urls e.g. git@github.com:fhopfensperger/my-repo.git")
	_ = viper.BindPFlag("repos", pf.Lookup("repos"))

	pf.StringP("filter", "b", "", "Which branches should be filtered e.g. release")
	_ = viper.BindPFlag("filter", pf.Lookup("filter"))
	_ = cobra.MarkFlagRequired(pf, "branch-filter")

	pf.StringP("file", "f", "", "Uses repos from file (one repo per line)")
	_ = viper.BindPFlag("file", pf.Lookup("file"))
	pf.StringP("pat", "p", "", `Use a Git Personal Access Token instead of the default private certificate! You could also set a environment variable. "export PAT=123456789" `)
	_ = viper.BindPFlag("pat", pf.Lookup("pat"))

	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	repos = viper.GetStringSlice("repos")
	filter = viper.GetString("filter")
	fileName = viper.GetString("file")
	pat = viper.GetString("pat")
}

func getReposFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msgf("Could not open file %s", fileName)
		return nil
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && string(line[0]) != "#" {
			lines = append(lines, line)
		}
	}
	return lines
}

func checkRepos() {
	if fileName != "" {
		repos = getReposFromFile(fileName)
	}
	if len(repos) == 0 && fileName == "" {
		fmt.Println("Either -f (file) or -r (repos) must be set")
		os.Exit(1)
	}
}
