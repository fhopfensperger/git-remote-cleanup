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
package pkg

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

func GetRemoteBranches(repoUrl string, branchFilter string) (*git.Remote, []string) {
	if branchFilter == "" {
		log.Warn().Msg("No branchfilter defined")
		os.Exit(1)
	}
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoUrl},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Err(err)
	}

	// Filters the references list and only branches which apply to the filter
	var branches []string
	for _, ref := range refs {
		if ref.Name().IsBranch() && strings.Contains(ref.Name().Short(), branchFilter) {
			branches = append(branches, ref.Name().String())
		}
	}
	log.Info().Msgf("Remote branches found: %v for repo %s and filter %s", branches, repoUrl, branchFilter)
	return rem, branches
}

func FilterBranches(branches []string) []string {
	var filteredBranches []string

	// Sort branches in ascending order
	sort.Slice(branches, func(i, j int) bool {
		return branches[i] < branches[j]
	})

	var minorVersionRegex = regexp.MustCompile(`v[0-9]`)
	// TODO make it configuable
	//var minorVersionRegex = regexp.MustCompile(`v[0-9].[0-9]`)

	minorBranches := make(map[string][]string)
	var tempBranches []string
	for _, b := range branches {
		minorVersion := minorVersionRegex.FindString(b)
		if len(tempBranches) == 0 {
			// temp branch is empty
			tempBranches = append(tempBranches, b)
		} else if minorVersionRegex.FindString(tempBranches[0]) == minorVersion {
			// branch has the same major version
			tempBranches = append(tempBranches, b)

		} else {
			// branch has other minor version
			tempBranches = nil
			tempBranches = append(tempBranches, b)
		}
		minorBranches[minorVersion] = tempBranches
	}

	for _, b := range minorBranches {
		for i, branch := range b {
			// Keep last version for minor
			if i == len(b)-1 {
				continue
			}
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}

func CleanBranches(remote *git.Remote, branchesToDelete []string, exclusionList []string, dryRun bool) {

	repoUrl := remote.Config().URLs[0]
	if len(branchesToDelete) == 0 {
		log.Info().Msgf("Nothing to delete for repo %s", repoUrl)
		return
	}

	// Exclude branches from deletion
	if len(exclusionList) != 0 {
		tmp := branchesToDelete[:0]
		for _, branch := range branchesToDelete {
			if exclude, ok := contains(exclusionList, branch); !ok {
				tmp = append(tmp, branch)
			} else {
				log.Info().Msgf("Excluding branch %s as it matches the exclusion list %s", branch, exclude)
			}
		}
		branchesToDelete = tmp
	}

	log.Info().Msgf("Going to delete branches: %v from repo %s ", branchesToDelete, repoUrl)

	// Clone repo temp
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		log.Err(err)
		os.Exit(1)
	}

	var refspecs []config.RefSpec
	// Add branches to Delete into refspecs
	for _, b := range branchesToDelete {
		refspecs = append(refspecs, config.RefSpec(b+":"+b))
	}

	log.Info().Msg("Deleting...")
	// push to delete branches which are matches the refspecs
	if !dryRun {
		err = r.Push(&git.PushOptions{
			Prune:    true,
			RefSpecs: refspecs,
		})
		if err != nil {
			log.Err(err)
		}
		log.Info().Msg("Branches deleted")
	} else {
		log.Info().Msg("Dry run! Nothing deleted")
	}

}

func contains(s []string, e string) (string, bool) {
	for _, a := range s {
		if strings.Contains(e, a) {
			return a, true
		}
	}
	return "", false
}
