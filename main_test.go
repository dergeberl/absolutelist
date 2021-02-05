package main

import (
	"os"
	"reflect"
	"testing"
)

func TestAppend(t *testing.T) {
	var expectFile, gotFile fileList
	expectFile.files = []file{{givenPath: "testdir/file0"}}
	err := gotFile.append([]string{"testdir/file0"})
	if err != nil {
		t.Errorf("got unespected error: %v", err)
	}
	if !reflect.DeepEqual(expectFile, gotFile) {
		t.Errorf("got %v but expect %v", gotFile, expectFile)
	}
}

func TestAppendDir(t *testing.T) {
	var expectFile, gotFile fileList
	workdir, _ := os.Getwd()
	expectFile.files = []file{
		{givenPath: workdir + "/testdir/.hiddenfile", absolutPath: workdir + "/testdir/.hiddenfile", isHidden: true},
		{givenPath: workdir + "/testdir/file0", absolutPath: workdir + "/testdir/file0"},
		{givenPath: workdir + "/testdir/file1", absolutPath: workdir + "/testdir/file1"},
		{givenPath: workdir + "/testdir/file2", absolutPath: workdir + "/testdir/file2"},
		{givenPath: workdir + "/testdir/file3", absolutPath: workdir + "/testdir/file3"},
		{givenPath: workdir + "/testdir/file4", absolutPath: workdir + "/testdir/file4"},
		{givenPath: workdir + "/testdir/file5", absolutPath: workdir + "/testdir/file5"},
		{givenPath: workdir + "/testdir/file6", absolutPath: workdir + "/testdir/file6"},
		{givenPath: workdir + "/testdir/file7", absolutPath: workdir + "/testdir/file7"},
		{givenPath: workdir + "/testdir/file8", absolutPath: workdir + "/testdir/file8"},
		{givenPath: workdir + "/testdir/file9", absolutPath: workdir + "/testdir/file9"}}
	err := gotFile.appendDir("testdir")
	if err != nil {
		t.Errorf("got unespected error: %v", err)
	}
	if !reflect.DeepEqual(expectFile, gotFile) {
		t.Errorf("got %v items but expect %v items", len(gotFile.files), len(expectFile.files))
	}
}

func TestCheckSingleDir(t *testing.T) {
	var expectFile, gotFile fileList
	workdir, _ := os.Getwd()
	expectFile.files = []file{
		{givenPath: "testdir"},
		{givenPath: workdir + "/testdir/.hiddenfile", absolutPath: workdir + "/testdir/.hiddenfile", isHidden: true},
		{givenPath: workdir + "/testdir/file0", absolutPath: workdir + "/testdir/file0"},
		{givenPath: workdir + "/testdir/file1", absolutPath: workdir + "/testdir/file1"},
		{givenPath: workdir + "/testdir/file2", absolutPath: workdir + "/testdir/file2"},
		{givenPath: workdir + "/testdir/file3", absolutPath: workdir + "/testdir/file3"},
		{givenPath: workdir + "/testdir/file4", absolutPath: workdir + "/testdir/file4"},
		{givenPath: workdir + "/testdir/file5", absolutPath: workdir + "/testdir/file5"},
		{givenPath: workdir + "/testdir/file6", absolutPath: workdir + "/testdir/file6"},
		{givenPath: workdir + "/testdir/file7", absolutPath: workdir + "/testdir/file7"},
		{givenPath: workdir + "/testdir/file8", absolutPath: workdir + "/testdir/file8"},
		{givenPath: workdir + "/testdir/file9", absolutPath: workdir + "/testdir/file9"}}
	_ = gotFile.append([]string{"testdir"})
	err := gotFile.checkSingleDir()
	if err != nil {
		t.Errorf("got unespected error: %v", err)
	}
	if !reflect.DeepEqual(expectFile, gotFile) {
		t.Errorf("got %v items but expect %v items", len(gotFile.files), len(expectFile.files))
	}
}

func TestGetInfo(t *testing.T) {
	var expectFile, gotFile fileList
	workdir, _ := os.Getwd()
	expectFile.files = []file{
		{givenPath: "testdir", absolutPath: workdir + "/testdir", isDir: true},
		{givenPath: workdir + "/testdir/.hiddenfile", absolutPath: workdir + "/testdir/.hiddenfile", isHidden: true},
		{givenPath: workdir + "/testdir/file0", absolutPath: workdir + "/testdir/file0"},
		{givenPath: workdir + "/testdir/file1", absolutPath: workdir + "/testdir/file1"},
		{givenPath: workdir + "/testdir/file2", absolutPath: workdir + "/testdir/file2"},
		{givenPath: workdir + "/testdir/file3", absolutPath: workdir + "/testdir/file3"},
		{givenPath: workdir + "/testdir/file4", absolutPath: workdir + "/testdir/file4"},
		{givenPath: workdir + "/testdir/file5", absolutPath: workdir + "/testdir/file5"},
		{givenPath: workdir + "/testdir/file6", absolutPath: workdir + "/testdir/file6"},
		{givenPath: workdir + "/testdir/file7", absolutPath: workdir + "/testdir/file7"},
		{givenPath: workdir + "/testdir/file8", absolutPath: workdir + "/testdir/file8"},
		{givenPath: workdir + "/testdir/file9", absolutPath: workdir + "/testdir/file9"}}
	_ = gotFile.append([]string{"testdir"})
	_ = gotFile.appendDir("testdir")
	err := gotFile.getInfos()
	if err != nil {
		t.Errorf("got unespected error: %v", err)
	}
	if !reflect.DeepEqual(expectFile, gotFile) {
		t.Errorf("got %v items but expect %v items", len(gotFile.files), len(expectFile.files))
	}
}

func BenchmarkAppend(b *testing.B) {
	b.ResetTimer()
	appendList := []string{"testdir/file0"}
	for i := 0; i < b.N; i++ {
		var fl fileList
		fl.append(appendList)
	}
}

func BenchmarkAppendDir(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var fl fileList
		fl.appendDir("testdir")
	}
}

func BenchmarkCheckSingleDir(b *testing.B) {
	b.ResetTimer()
	var fl fileList
	fl.append([]string{"testdir/file0"})
	for i := 0; i < b.N; i++ {
		fl.checkSingleDir()
	}
}

func BenchmarkGetInfos(b *testing.B) {
	b.ResetTimer()
	var fl fileList
	fl.appendDir("testdir")
	for i := 0; i < b.N; i++ {
		fl.getInfos()
	}
}

func BenchmarkPrint(b *testing.B) {
	devNull, _ := os.Open(os.DevNull)
	sOut := os.Stdout
	os.Stdout = devNull
	var fl fileList
	fl.appendDir("testdir")
	for i := 0; i < b.N; i++ {
		fl.print()
	}
	_ = devNull.Close()
	os.Stdout = sOut
}
