# git-remote-cleanup
![Go](https://github.com/fhopfensperger/git-remote-cleanup/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/git-remote-cleanup)](https://goreportcard.com/report/github.com/fhopfensperger/git-remote-cleanup)
[![Coverage Status](https://coveralls.io/repos/github/fhopfensperger/git-remote-cleanup/badge.svg?branch=master)](https://coveralls.io/github/fhopfensperger/git-remote-cleanup?branch=master)

Get and delete no longer needed branches from a remote repository.

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/git-remote-cleanup/master/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/git-remote-cleanup/releases) download the latest release according to your processor architecture and operating system, then unarchive and copy it to the right location

```bash
tar xvfz git-remote-cleanup_x.x.x_darwin_amd64.tar.gz
cd git-remote-cleanup_x.x.x_darwin_amd64
chmod +x git-remote-cleanup
sudo mv git-remote-cleanup /usr/local/bin/
```
