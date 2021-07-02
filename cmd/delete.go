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

var excludes []string
var dryRun bool

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete old branches, keeps every latest hotfix version",
	Long:  `Delete old branches, keeps every latest hotfix version`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRepos()
		auth := http.BasicAuth{
			Username: "123", // Using a PAT this can be anything except an empty string
			Password: pat,
		}
		excludes = viper.GetStringSlice("exclude")
		dryRun = viper.GetBool("dry-run")
		for _, r := range repos {
			gitService := pkg.New(nil, &auth)
			branches := gitService.GetRemoteBranches(r, filter, false)
			branches = pkg.FilterBranches(branches)
			gitService.CleanBranches(branches, excludes, dryRun)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	flags := deleteCmd.Flags()
	flags.StringSliceP("exclude", "e", []string{}, "Exclude branches, e.g. v1.0.1")
	_ = viper.BindPFlag("exclude", flags.Lookup("exclude"))

	flags.Bool("dry-run", false, "Perform dry run, do not delete anything")
	_ = viper.BindPFlag("dry-run", flags.Lookup("dry-run"))

}
