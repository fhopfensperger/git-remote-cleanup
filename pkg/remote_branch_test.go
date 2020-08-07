package pkg

import (
	"fmt"
	"testing"

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
			want: []string(nil),
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

func (m *remoteBranchMock) List() bool {
	fmt.Println("Mocked charge notification function")
	//fmt.Printf("Value passed in: %s %s\n", s, c)
	// this records that the method was called and passes in the value
	// it was called with
	args := m.Called("nase")
	// it then returns whatever we tell it to return
	// in this case true to simulate an SMS Service Notification
	// sent out
	//remoteConfig := config.RemoteConfig{
	//	Name:  "nase",
	//	URLs:  []string{"nase.com"},
	//	Fetch: nil,
	//}

	return args.Bool(0)
}

func TestGetRemoteBranches(t *testing.T) {
	remote := new(remoteBranchMock)
	nase := "nase"
	ref := plumbing.NewHashReference(plumbing.ReferenceName(nase), plumbing.Hash{})
	rfs := []*plumbing.Reference{ref}

	//rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
	//	Name: "origin",
	//	URLs: []string{nase},
	//})

	remote.On("rem.List").Return(&rfs, nil)

	GetRemoteBranches(nase, nase)
}
