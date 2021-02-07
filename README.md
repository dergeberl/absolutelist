# Absolutelist

[![Go Report Card](https://goreportcard.com/badge/github.com/dergeberl/absolutelist)](https://goreportcard.com/report/github.com/dergeberl/absolutelist) ![Coveralls github](https://img.shields.io/coveralls/github/dergeberl/absolutelist) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dergeberl/absolutelist) ![GitHub](https://img.shields.io/github/license/dergeberl/absolutelist) ![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/dergeberl/absolutelist?include_prereleases)

Absolutelist is a small tool to get the absolute path of files and folders.  
I know there are a few other tools that can do this already, but it was a nice small project to code in golang and set up a github-action pipeline for it.

## Installation

Download the binary from the github release and mv it to a folder that is in your path (example /usr/local/bin).   
Use `chmod +x /path/to/absolutelist` to make it executable.

## Usage

```
absolutelist [options] [files]...

Options:
  -d    print only directories
  -f    print only files
```

For easy typing you can set an alias like `al`:
```
alias al=absolutelist
```

## Examples:

`absolutelist` prints the absolut path of files and directories from the current directory   
`absolutelist some-dir` prints the absolut path of files and directories in the folder `some-dir`   
`absolutelist file1 some-dir` prints the absolut path from `file1` and `some-dir` (not the content of the dir)   
`absolutelist -d` prints the absolut path of **files** from the current directory   
`absolutelist -f` prints the absolut path of **directories** from the current directory   
`ls *.txt | absolutelist` prints the absolute path of all files and directories piped by `ls`
