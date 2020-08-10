package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/go-git/go-git/v5/config"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func Test_contains(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"test-contains", args{
			s: []string{"v1.1.2", "v.1.1.9"},
			e: "/head/release/v1.1.2",
		}, "v1.1.2", true},
		{"test-not-contains", args{
			s: []string{"v1.1.2", "v.1.1.9"},
			e: "head/release/v1.1.8",
		}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualString, actualBool := contains(tt.args.s, tt.args.e)
			assert.Equal(tt.want, actualString)
			assert.Equal(tt.want1, actualBool)

		})
	}
}

func TestFilterBranches(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		branches []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "filter-return-empty",
			args: args{branches: []string{"head/release/v1.2.8", "head/release/v1.1.9"}},
			want: []string{},
		},
		{
			name: "filter-return-one",
			args: args{branches: []string{"head/release/v1.2.8", "head/release/v1.1.9", "head/release/v1.1.10"}},
			want: []string{"head/release/v1.1.9"},
		},
		{
			name: "filter-return-multiple",
			args: args{branches: []string{"head/release/v1.1.0", "head/release/v1.1.8", "head/release/v1.1.9", "head/release/v1.1.10", "head/release/v1.1.11",
				"head/release/v1.2.0", "head/release/v1.2.1"}},
			want: []string{"head/release/v1.1.0", "head/release/v1.1.8", "head/release/v1.1.9", "head/release/v1.1.10", "head/release/v1.2.0"},
		},
		{
			name: "filter-return-multiple",
			args: args{branches: []string{"head/release/v10.1.0", "head/release/v10.1.8", "head/release/v1.1.9", "head/release/v1.1.10", "head/release/v1.1.11",
				"head/release/v1.2.0", "head/release/v1.2.1"}},
			want: []string{"head/release/v10.1.0", "head/release/v1.1.9", "head/release/v1.1.10", "head/release/v1.2.0"},
		},
		{
			name: "filter-return-multiple-major-versions",
			args: args{branches: []string{"head/release/v2.1.8", "head/release/v2.1.9", "head/release/v1.1.10", "head/release/v1.1.11",
				"head/release/v1.2.0", "head/release/v1.2.1"}},
			want: []string{"head/release/v1.1.10", "head/release/v1.2.0", "head/release/v2.1.8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, FilterBranches(tt.args.branches))
		})
	}
}

type remoteBranchMock struct {
	mock.Mock
}

func (m *remoteBranchMock) Push(options *git.PushOptions) error {
	fmt.Println("Mocked Push function")
	return nil
}

func (m *remoteBranchMock) Config() *config.RemoteConfig {
	fmt.Println("Mocked Config function")
	args := m.Called()
	return args.Get(0).(*config.RemoteConfig)
}

func (m *remoteBranchMock) List(l *git.ListOptions) ([]*plumbing.Reference, error) {
	fmt.Println("Mocked List function")
	args := m.Called(l)
	return args.Get(0).([]*plumbing.Reference), nil
}

func TestGetRemoteBranches(t *testing.T) {
	remote := new(remoteBranchMock)
	ref := plumbing.NewHashReference("refs/heads/release/v2.2.2", plumbing.Hash{})
	ref2 := plumbing.NewHashReference(("refs/heads/release/v2.2.0"), plumbing.Hash{})
	ref3 := plumbing.NewHashReference(("refs/heads/master"), plumbing.Hash{})
	ref4 := plumbing.NewHashReference(("refs/heads/release/v1.0.0"), plumbing.Hash{})
	ref5 := plumbing.NewHashReference(("refs/heads/release/v11.0.0"), plumbing.Hash{})

	mockRemoteBranch := New(remote)

	remote.On("List", &git.ListOptions{}).Return([]*plumbing.Reference{ref, ref2, ref3, ref4, ref5}, nil)

	foundBranches := mockRemoteBranch.GetRemoteBranches("https://github.com/fhopfensperger/amqp-sb-client.git", "release", false)
	remote.AssertExpectations(t)

	assert.Equal(t, "refs/heads/release/v1.0.0", foundBranches[0])
	assert.Equal(t, "refs/heads/release/v2.2.0", foundBranches[1])
	assert.Equal(t, "refs/heads/release/v2.2.2", foundBranches[2])
	assert.Equal(t, "refs/heads/release/v11.0.0", foundBranches[3])

}

func TestGetRemoteBranches_latest(t *testing.T) {
	remote := new(remoteBranchMock)
	ref1 := plumbing.NewHashReference(("refs/heads/master"), plumbing.Hash{})
	ref2 := plumbing.NewHashReference(("refs/heads/release/v1.0.0"), plumbing.Hash{})
	ref3 := plumbing.NewHashReference(("refs/heads/release/v1.1.99"), plumbing.Hash{})
	ref4 := plumbing.NewHashReference(("refs/heads/release/v11.0.1"), plumbing.Hash{})
	ref5 := plumbing.NewHashReference(("refs/heads/release/v11.0.0"), plumbing.Hash{})

	mockRemoteBranch := New(remote)

	remote.On("List", &git.ListOptions{}).Return([]*plumbing.Reference{ref1, ref2, ref3, ref4, ref5}, nil)

	foundBranches := mockRemoteBranch.GetRemoteBranches("https://github.com/fhopfensperger/amqp-sb-client.git", "release", true)
	remote.AssertExpectations(t)

	assert.Equal(t, "refs/heads/release/v11.0.1", foundBranches[0])
}

// Test exit status 1 if no branchFilter is defined
func TestGetRemoteBranchesNoBranchFilter(t *testing.T) {
	if os.Getenv("FLAG") == "1" {
		remote := new(remoteBranchMock)
		mockRemoteBranch := New(remote)
		mockRemoteBranch.GetRemoteBranches("https://github.com/fhopfensperger/amqp-sb-client.git", "", false)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestGetRemoteBranchesNoBranchFilter")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestRemoteBranch_CleanBranches(t *testing.T) {
	remote := new(remoteBranchMock)
	// Actually we need to use a real Repo, because inside CleanBranches we need a real repo
	// however, the deletion of branches is mocked.
	remoteConfing := config.RemoteConfig{
		Name:  "amqp-sb-client.git",
		URLs:  []string{"https://github.com/fhopfensperger/amqp-sb-client.git"},
		Fetch: nil,
	}
	repo := git.Repository{}

	mockRemoteBranch := New(remote)
	mockRemoteBranch.AddRepo(&repo)

	var refspecs []config.RefSpec
	refspecs = append(refspecs, "refs/heads/release/v2.2.2:refs/heads/release/v2.2.2")
	pushOptions := &git.PushOptions{
		Prune:    true,
		RefSpecs: refspecs,
	}

	remote.On("Config").Return(&remoteConfing)
	remote.On("Push", &pushOptions).Return(nil)
	deletedBranches := mockRemoteBranch.CleanBranches([]string{"refs/heads/release/v2.2.2"}, []string{"v2.2.1"}, false)
	assert.Equal(t, []string{"refs/heads/release/v2.2.2"}, deletedBranches)
}

func TestRemoteBranch_CleanBranches_All_Excluded(t *testing.T) {
	remote := new(remoteBranchMock)
	// Actually we need to use a real Repo, because inside CleanBranches we need a real repo
	// however, the deletion of branches is mocked.
	remoteConfing := config.RemoteConfig{
		Name:  "amqp-sb-client.git",
		URLs:  []string{"https://github.com/fhopfensperger/amqp-sb-client.git"},
		Fetch: nil,
	}
	repo := git.Repository{}

	mockRemoteBranch := New(remote)
	mockRemoteBranch.AddRepo(&repo)

	var refspecs []config.RefSpec
	refspecs = append(refspecs, "refs/heads/release/v2.2.2:refs/heads/release/v2.2.2")
	refspecs = append(refspecs, "refs/heads/release/v2.2.1:refs/heads/release/v2.2.1")
	pushOptions := &git.PushOptions{
		Prune:    true,
		RefSpecs: refspecs,
	}

	remote.On("Config").Return(&remoteConfing)
	remote.On("Push", &pushOptions).Return(nil)
	deletedBranches := mockRemoteBranch.CleanBranches([]string{"refs/heads/release/v2.2.2", "refs/heads/release/v2.2.1"}, []string{"v2.2.1", "v2.2.2"}, false)
	assert.Empty(t, deletedBranches)
}

func TestRemoteBranch_CleanBranches_Some_Excluded(t *testing.T) {
	remote := new(remoteBranchMock)
	// Actually we need to use a real Repo, because inside CleanBranches we need a real repo
	// however, the deletion of branches is mocked.
	remoteConfing := config.RemoteConfig{
		Name:  "amqp-sb-client.git",
		URLs:  []string{"https://github.com/fhopfensperger/amqp-sb-client.git"},
		Fetch: nil,
	}
	repo := git.Repository{}

	mockRemoteBranch := New(remote)
	mockRemoteBranch.AddRepo(&repo)

	var refspecs []config.RefSpec
	refspecs = append(refspecs, "refs/heads/release/v2.2.2:refs/heads/release/v2.2.2")
	refspecs = append(refspecs, "refs/heads/release/v2.2.1:refs/heads/release/v2.2.1")
	pushOptions := &git.PushOptions{
		Prune:    true,
		RefSpecs: refspecs,
	}

	remote.On("Config").Return(&remoteConfing)
	remote.On("Push", &pushOptions).Return(nil)
	deletedBranches := mockRemoteBranch.CleanBranches([]string{"refs/heads/release/v2.2.3", "refs/heads/release/v2.2.2", "refs/heads/release/v2.2.1"}, []string{"v2.2.2"}, false)
	assert.Equal(t, []string{"refs/heads/release/v2.2.3", "refs/heads/release/v2.2.1"}, deletedBranches)
}

func TestRemoteBranch_CleanBranches_branches_to_delete_empty(t *testing.T) {
	remote := new(remoteBranchMock)
	// Actually we need to use a real Repo, because inside CleanBranches we need a real repo
	// however, the deletion of branches is mocked.
	remoteConfing := config.RemoteConfig{
		Name:  "amqp-sb-client.git",
		URLs:  []string{"https://github.com/fhopfensperger/amqp-sb-client.git"},
		Fetch: nil,
	}
	repo := git.Repository{}

	mockRemoteBranch := New(remote)
	mockRemoteBranch.AddRepo(&repo)

	var refspecs []config.RefSpec
	refspecs = append(refspecs, "refs/heads/release/v2.2.2:refs/heads/release/v2.2.2")
	refspecs = append(refspecs, "refs/heads/release/v2.2.1:refs/heads/release/v2.2.1")
	pushOptions := &git.PushOptions{
		Prune:    true,
		RefSpecs: refspecs,
	}

	remote.On("Config").Return(&remoteConfing)
	remote.On("Push", &pushOptions).Return(nil)
	deletedBranches := mockRemoteBranch.CleanBranches([]string{}, []string{"v2.2.2"}, false)
	assert.Empty(t, deletedBranches)
}
