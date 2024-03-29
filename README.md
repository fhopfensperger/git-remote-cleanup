# git-remote-cleanup
![Go](https://github.com/fhopfensperger/git-remote-cleanup/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/git-remote-cleanup)](https://goreportcard.com/report/github.com/fhopfensperger/git-remote-cleanup)
[![Coverage Status](https://coveralls.io/repos/github/fhopfensperger/git-remote-cleanup/badge.svg?branch=master)](https://coveralls.io/github/fhopfensperger/git-remote-cleanup?branch=master)
[![Release](https://img.shields.io/github/release/fhopfensperger/git-remote-cleanup.svg?style=flat-square)](https://github.com//fhopfensperger/git-remote-cleanup/releases/latest)


Get and delete no longer needed release branches from a remote repository.

# Usage

## All commands and flags

```bash
Available Commands:
  branches    Get remote branches
  delete      Delete old branches, keeps every latest patch version
  help        Help about any command

Flags:
  -f, --file string     Uses repos from file (one repo per line)
  -b, --filter string   Which branches should be filtered e.g. release
  -h, --help            help for git-remote-cleanup
  -p, --pat string      Use a Git Personal Access Token instead of the default private certificate! You could also set a environment variable. "export PAT=123456789" 
  -r, --repos strings   Git Repo urls e.g. git@github.com:fhopfensperger/my-repo.git
  -v, --version         version for git-remote-cleanup
```

Note: All flags can be set using environment variables, for example:
```bash
export REPOS=git@github.com:fhopfensperger/my-repo.git
export PAT=1234567890abcdef
...
```

# Installation

## Homebrew

```bash
brew install fhopfensperger/tap/git-remote-cleanup
```

## Script

```bash
curl https://raw.githubusercontent.com/fhopfensperger/git-remote-cleanup/master/get.sh | bash
```

## Manually

Go to [Releases](https://github.com/fhopfensperger/git-remote-cleanup/releases) download the latest release according to your processor architecture and operating system, then unarchive and copy it to the right location

```bash
tar xvfz git-remote-cleanup_x.x.x_darwin_amd64.tar.gz
cd git-remote-cleanup_x.x.x_darwin_amd64
chmod +x git-remote-cleanup
sudo mv git-remote-cleanup /usr/local/bin/
```

## Run as container

Besides installing the binary on the local computer, you have the option to run the program as a container
```bash
# Using a single repo
docker run -it --rm ghcr.io/fhopfensperger/git-remote-cleanup branches -r https://github.com/fhopfensperger/git-remote-cleanup.git -b master -p 123

# Using a file to define multiple repos
docker run -it --rm -v $(pwd)/repos_http.txt:/app/repos_http.txt ghcr.io/fhopfensperger/git-remote-cleanup branches -f repos_http.txt -b master -p 123
```
