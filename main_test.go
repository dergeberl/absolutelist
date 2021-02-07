package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMainFunc(t *testing.T) {
	t.Run("test main", func(t *testing.T) {

		rIn, wIn, _ := os.Pipe()
		sIn := os.Stdin
		os.Stdin = rIn
		_, _ = wIn.WriteString("testdir/file0")
		_ = wIn.Close()

		rOut, wOut, _ := os.Pipe()
		sOut := os.Stdout
		os.Stdout = wOut

		main()

		_ = wOut.Close()
		os.Stdin = sIn
		os.Stdout = sOut
		output, _ := ioutil.ReadAll(rOut)

		if !strings.HasSuffix(string(output), "testdir/file0\n") {
			t.Errorf("output did not end with 'testdir/file0': %v", string(output))
		}
	})
}

func TestAppendStdin(t *testing.T) {
	t.Run("appendStdin file", func(t *testing.T) {
		var expectFile, gotFile fileList

		r, w, _ := os.Pipe()
		sIn := os.Stdin
		os.Stdin = r
		_, _ = w.WriteString("testdir/file0")
		_ = w.Close()
		err := gotFile.appendStdin()
		if err != nil {
			t.Errorf("got unespected error: %v", err)
		}
		os.Stdin = sIn
		expectFile.files = []file{{givenPath: "testdir/file0"}}

		if !reflect.DeepEqual(expectFile, gotFile) {
			t.Errorf("got %v but expect %v", gotFile, expectFile)
		}
	})
}

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
	t.Run("append testdir", func(t *testing.T) {
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
	})
	t.Run("append not found dir", func(t *testing.T) {
		var gotFile fileList
		err := gotFile.appendDir("notfound")
		if err == nil {
			t.Errorf("expect error but got none")
		}
		if !reflect.DeepEqual(fileList{}, gotFile) {
			t.Errorf("got %v items but expect %v items", len(gotFile.files), 0)
		}
	})
}

func TestCheckSingleDir(t *testing.T) {
	t.Run("test with single dir", func(t *testing.T) {
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
	})
	t.Run("test with more than one dir", func(t *testing.T) {
		var expectFile, gotFile fileList
		expectFile.files = []file{
			{givenPath: "testdir"},
			{givenPath: "testdir2"}}
		_ = gotFile.append([]string{"testdir", "testdir2"})
		err := gotFile.checkSingleDir()
		if err != nil {
			t.Errorf("got unespected error: %v", err)
		}
		if !reflect.DeepEqual(expectFile, gotFile) {
			t.Errorf("got %v items but expect %v items", len(gotFile.files), len(expectFile.files))
		}
	})
	t.Run("with not found dir", func(t *testing.T) {
		var expectFile, gotFile fileList
		_ = gotFile.append([]string{"notfound"})
		_ = expectFile.append([]string{"notfound"})
		err := gotFile.checkSingleDir()
		if err == nil {
			t.Errorf("expect error but got none")
		}
		if !reflect.DeepEqual(expectFile, gotFile) {
			t.Errorf("got %v items but expect %v items", len(gotFile.files), len(expectFile.files))
		}
	})
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

func TestPrint(t *testing.T) {
	t.Run("test with file", func(t *testing.T) {
		var fl fileList
		_ = fl.append([]string{"testdir/file0"})
		_ = fl.getInfos()

		r, w, _ := os.Pipe()
		sOut := os.Stdout
		os.Stdout = w

		_ = fl.print()

		_ = w.Close()
		os.Stdout = sOut
		output, _ := ioutil.ReadAll(r)
		if !strings.HasSuffix(string(output), "testdir/file0\n") {
			t.Errorf("output did not end with 'testdir/file0': %v", string(output))
		}
	})
	t.Run("test with file", func(t *testing.T) {
		var fl fileList
		_ = fl.append([]string{"testdir"})
		_ = fl.getInfos()

		r, w, _ := os.Pipe()
		sOut := os.Stdout
		os.Stdout = w

		_ = fl.print()

		_ = w.Close()
		os.Stdout = sOut
		output, _ := ioutil.ReadAll(r)
		if !strings.HasSuffix(string(output), "testdir/\n") {
			t.Errorf("output did not end with 'testdir/': %v", string(output))
		}
	})
	t.Run("test with hiddenfile", func(t *testing.T) {
		var fl fileList
		_ = fl.append([]string{"testdir/.hiddenfile"})
		_ = fl.getInfos()

		r, w, _ := os.Pipe()
		sOut := os.Stdout
		os.Stdout = w

		_ = fl.print()

		_ = w.Close()
		os.Stdout = sOut
		output, _ := ioutil.ReadAll(r)
		if len(output) != 0 {
			t.Errorf("expect no output")
		}
	})
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
