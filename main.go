package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var option_i = flag.String("i", "*", "compared files pattern")

func GetFileHash(fname string) (string, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	sum := md5.Sum(data)
	return string(sum[:]), nil
}

type FileT struct {
	Size int64
	Path string
}

func (file1 FileT) Equals(file2 FileT) (bool, error) {
	if file1.Size != file2.Size {
		return false, nil
	}
	hash1, err := GetFileHash(file1.Path)
	if err != nil {
		return false, err
	}
	hash2, err := GetFileHash(file2.Path)
	if err != nil {
		return false, err
	}
	return hash1 == hash2, nil
}

func GetDirList(dir string) (map[string]FileT, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	files, err := fd.Readdir(-1)
	if err != nil && err != io.EOF {
		return nil, err
	}
	dict := map[string]FileT{}
	for _, file1 := range files {
		if file1.IsDir() {
			continue
		}
		matched, err := filepath.Match(
			strings.ToUpper(*option_i),
			strings.ToUpper(file1.Name()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if !matched {
			continue
		}
		path1 := filepath.Join(dir, file1.Name())
		stat1, err := os.Stat(path1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path1, err.Error())
		} else {
			dict[strings.ToUpper(file1.Name())] = FileT{
				Size: stat1.Size(),
				Path: path1,
			}
		}
	}
	return dict, nil
}

func main1(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: %s DIR1 DIR2", os.Args[0])
	}
	dirlist1, err := GetDirList(args[0])
	if err != nil {
		return fmt.Errorf("%s: %s", args[0], err.Error())
	}
	dirlist2, err := GetDirList(args[1])
	if err != nil {
		return fmt.Errorf("%s: %s\n", args[1], err.Error())
	}

	for name, info1 := range dirlist1 {
		if info2, ok := dirlist2[name]; ok {
			same, err := info1.Equals(info2)
			if err != nil {
				return err
			}
			if !same {
				fmt.Printf("M\t%s\n", name)
			}
		} else {
			fmt.Printf("D\t%s\n", name)
		}
	}
	for name, _ := range dirlist2 {
		if _, ok := dirlist1[name]; !ok {
			fmt.Printf("A\t%s\n", name)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := main1(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
