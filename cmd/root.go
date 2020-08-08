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

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var repos []string
var filter string
var fileName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-remote-cleanup",
	Short: "Simple command line utility to get and delete branches from a git hub repo",
	Long:  `Simple command line utility to get and delete branches from a git hub repo`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-remote-cleanup.yaml)")
	pf := rootCmd.PersistentFlags()
	pf.StringSliceP("repos", "r", []string{}, "Git Repo urls e.g. git@github.com:fhopfensperger/my-repo.git")
	viper.BindPFlag("repos", pf.Lookup("repos"))
	//cobra.MarkFlagRequired(pf, "repos")
	pf.StringP("branch-filter", "b", "", "Which branches should be filtered e.g. release")
	viper.BindPFlag("branch-filter", pf.Lookup("branch-filter"))
	cobra.MarkFlagRequired(pf, "branch-filter")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	pf.StringP("file", "f", "", "Uses repos from file (one repo per line)")
	viper.BindPFlag("file", pf.Lookup("file"))
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".git-remote-cleanup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".git-remote-cleanup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	repos = viper.GetStringSlice("repos")
	filter = viper.GetString("branch-filter")
	fileName = viper.GetString("file")

	if fileName != "" {
		repos = getReposFromFile()
	}
	if len(repos) == 0 && fileName == "" {
		fmt.Println("Either -f (file) or -r (repos) must be set")
		os.Exit(1)
	}
}

func getReposFromFile() []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msgf("Could not open file %s", fileName)
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
