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
	"github.com/fhopfensperger/git-remote-cleanup/pkg"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var latest bool

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branches",
	Short: "Get remote branches",
	Long:  `Get remote branches`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRepos()
		auth := http.BasicAuth{
			Username: "123", // Using a PAT this can be anything except an empty string
			Password: pat,
		}
		for _, r := range repos {
			latest = viper.GetBool("latest")
			gitService := pkg.New(nil, &auth)
			gitService.GetRemoteBranches(r, filter, latest)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)

	flags := branchCmd.Flags()
	flags.BoolP("latest", "l", false, "Print latest remote branch for filter")
	_ = viper.BindPFlag("latest", flags.Lookup("latest"))
}
