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
	"github.com/spf13/cobra"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branches",
	Short: "Get remote branches",
	Long:  `Get remote branches`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, r := range repos {
			gitService := pkg.RemoteBranch{}
			gitService.GetRemoteBranches(r, filter)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
