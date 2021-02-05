package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fileList struct {
	files []file
}
type file struct {
	givenPath   string
	absolutPath string
	isDir       bool
	isHidden    bool
}

var onlyDir, onlyFiles, allFiles, showVersion bool
var version = "dev"
var date = "none"
var commit = "none"

func main() {
	//init arguments
	flag.BoolVar(&onlyDir, "d", false, "print only directories")
	flag.BoolVar(&onlyFiles, "f", false, "print only files")
	flag.BoolVar(&allFiles, "a", false, "print all files and directories including hidden files")
	flag.BoolVar(&showVersion, "v", false, "print version")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options] [files]...\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	//return error if onlyFiles and onlyDir is active
	if onlyDir && onlyFiles {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -d and -f cannot be activated at the same time\n\n")
		flag.Usage()
		os.Exit(1)
	}

	//print version if -v flag is set
	if showVersion {
		fmt.Printf("%s\nVersion: %s\nDate: %v\nCommit: %v\n", os.Args[0], version, date, commit)
		os.Exit(0)
	}

	var fl fileList
	//read stdin for pipes
	err := fl.appendStdin()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	//append files from arguments
	err = fl.append(flag.Args())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	//if only one dir is set use this dir as root
	err = fl.checkSingleDir()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	//if no file from stdin and argument is specified append currentDir
	if len(fl.files) == 0 {
		err = fl.appendDir(".")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	//append file infos if needed
	err = fl.getInfos()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	//print result
	err = fl.print()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//append appends the list to fileList
func (fl *fileList) append(list []string) error {
	for _, listItem := range list {
		fl.files = append(fl.files, file{
			givenPath: listItem,
		})
	}
	return nil
}

//appendStdin appends the list from stdin to fileList
func (fl *fileList) appendStdin() error {
	stdin, _ := os.Stdin.Stat()
	if stdin.Mode()&os.ModeCharDevice == 0 {
		reader := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for reader.Scan() {
			fl.files = append(fl.files, file{
				givenPath: reader.Text(),
			})
		}
	}
	return nil
}

//appendDir appends the content from a dir to fileList
func (fl *fileList) appendDir(dir string) error {
	dirPath, _ := filepath.Abs(dir)
	dirFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, dirFile := range dirFiles {
		fl.files = append(fl.files, file{givenPath: dirPath + "/" + dirFile.Name(), absolutPath: dirPath + "/" + dirFile.Name(), isDir: dirFile.IsDir(), isHidden: strings.HasPrefix(dirFile.Name(), ".")})
	}
	return nil
}

//checkSingleDir check if only a single dir is set and appends the dir content to fileList
func (fl *fileList) checkSingleDir() error {
	if len(fl.files) != 1 {
		return nil
	}
	f, err := os.Stat(fl.files[0].givenPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		fl.appendDir(fl.files[0].givenPath)
	}
	return nil
}

//getInfos get all needed infos for all files in fileList
func (fl *fileList) getInfos() error {
	for i, fileItem := range fl.files {
		if fileItem.absolutPath != "" {
			continue
		}
		f, err := os.Stat(fileItem.givenPath)
		if err != nil {
			return err
		}
		fl.files[i].isDir = f.IsDir()
		fl.files[i].isHidden = strings.HasPrefix(f.Name(), ".")
		fl.files[i].absolutPath, err = filepath.Abs(fileItem.givenPath)
		if err != nil {
			return err
		}
	}
	return nil
}

//print prints the absolut path for all files and checks if only dirs and files should be printed
func (fl *fileList) print() error {
	for _, f := range fl.files {
		if f.isHidden && !allFiles {
			continue
		}
		if f.isDir && !onlyFiles && f.absolutPath != "" {
			fmt.Printf("%v/\n", f.absolutPath)
		} else if !f.isDir && !onlyDir && f.absolutPath != "" {
			fmt.Printf("%v\n", f.absolutPath)
		}
	}
	return nil
}
