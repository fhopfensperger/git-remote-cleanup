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

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/rs/zerolog/log"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

type GitInterface interface {
	// all of the functions we use from the third party client
	List(*git.ListOptions) ([]*plumbing.Reference, error)
	Config() *config.RemoteConfig
	Push(*git.PushOptions) error
}

type RemoteBranch struct {
	gitClient GitInterface
	repo      *git.Repository
}

func New(client GitInterface) RemoteBranch {
	return RemoteBranch{client, nil}
}

func (m *RemoteBranch) AddRepo(repo *git.Repository) {
	m.repo = repo
}

func (m *RemoteBranch) GetRemoteBranches(repoUrl string, branchFilter string) []string {
	if branchFilter == "" {
		log.Warn().Msg("No branchfilter defined")
		os.Exit(1)
	}
	if m.gitClient == nil {
		m.gitClient = git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
			Name: "origin",
			URLs: []string{repoUrl},
		})
	}

	// We can then use every Remote functions to retrieve wanted information
	refs, err := m.gitClient.List(&git.ListOptions{})
	if err != nil {
		log.Err(err).Msg("")
	}

	// Filters the references list and only branches which apply to the filter
	var branches []string
	for _, ref := range refs {
		if ref.Name().IsBranch() && strings.Contains(ref.Name().Short(), branchFilter) {
			branches = append(branches, ref.Name().String())
		}
	}
	sort.SliceStable(branches, func(i, j int) bool {
		return branches[i] < branches[j]
	})
	log.Info().Msgf("Remote branches found: %v for repo %s and filter %s", branches, repoUrl, branchFilter)
	return branches
}

func FilterBranches(branches []string) []string {
	filteredBranches := branches[:0]
	// TODO make it configuable
	regex := `v[0-9].[0-9]`
	//var minorVersionRegex = regexp.MustCompile(`v[0-9]`)
	var minorVersionRegex = regexp.MustCompile(regex)

	// Sort branches by regex in ascending order
	sort.SliceStable(branches, func(i, j int) bool {
		return minorVersionRegex.FindString(branches[i]) < minorVersionRegex.FindString(branches[j])
	})

	minorBranches := make(map[string][]string)
	tempBranches := branches[:0]
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

	// As iteration over maps doesnt guarantee order, we need to do this workaround
	// https://blog.golang.org/maps#TOC_7.
	keys := make([]string, 0)
	for k, _ := range minorBranches {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for i, branch := range minorBranches[k] {
			// Keep last version for minor
			if i == len(minorBranches[k])-1 {
				continue
			}
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}

func (m *RemoteBranch) CleanBranches(branchesToDelete []string, exclusionList []string, dryRun bool) (deletedBranches []string) {

	repoUrl := m.gitClient.Config().URLs[0]
	if len(branchesToDelete) == 0 {
		log.Info().Msgf("Nothing to delete for repo %s", repoUrl)
		return nil
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

	if len(branchesToDelete) == 0 {
		log.Info().Msgf("Nothing to delete, all branches are excluded")
		return nil
	}

	log.Info().Msgf("Going to delete branches: %v from repo %s", branchesToDelete, repoUrl)

	// Clone repo temp
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		log.Err(err).Msg("")
		os.Exit(1)
	}
	if m.repo == nil {
		m.AddRepo(r)
	}

	var refspecs []config.RefSpec
	// Add branches to Delete into refspecs
	for _, b := range branchesToDelete {
		refspecs = append(refspecs, config.RefSpec(b+":"+b))
	}

	log.Info().Msg("Deleting...")
	// push to delete branches which are matches the refspecs
	if !dryRun {
		err = m.gitClient.Push(&git.PushOptions{
			Prune:    true,
			RefSpecs: refspecs,
		})
		if err != nil {
			log.Err(err).Msg("")
		}
		return branchesToDelete
		log.Info().Msg("Branches deleted")
	} else {
		log.Info().Msg("Dry run! Nothing deleted")
	}
	return nil
}

func contains(s []string, e string) (string, bool) {
	for _, a := range s {
		if strings.Contains(e, a) {
			return a, true
		}
	}
	return "", false
}
