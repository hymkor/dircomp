package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetFileHash(fname string) (string, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	sum := md5.Sum(data)
	return string(sum[:]), nil
}

func GetDirHash(dir string) (map[string]string, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	files, err := fd.Readdir(-1)
	if err != nil && err != io.EOF {
		return nil, err
	}
	dict := map[string]string{}
	for _, file1 := range files {
		if !file1.IsDir() {
			path1 := filepath.Join(dir, file1.Name())
			hash1, err := GetFileHash(path1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", path1, err.Error())
			} else {
				dict[file1.Name()] = hash1
			}
		}
	}
	return dict, nil
}

func main1() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: dirdiff DIR1 DIR2")
	}
	dirhash1, err := GetDirHash(os.Args[1])
	if err != nil {
		return fmt.Errorf("%s: %s", os.Args[1], err.Error())
	}
	dirhash2, err := GetDirHash(os.Args[2])
	if err != nil {
		return fmt.Errorf("%s: %s\n", os.Args[2], err.Error())
	}

	for name, hash1 := range dirhash1 {
		if hash2, ok := dirhash2[name]; ok {
			if hash1 != hash2 {
				fmt.Printf("%s: differs\n", name)
			}
		} else {
			fmt.Printf("%s: not found in %s\n", name, os.Args[2])
		}
	}
	for name, _ := range dirhash2 {
		if _, ok := dirhash1[name]; !ok {
			fmt.Printf("%s: not found in %s\n", name, os.Args[1])
		}
	}
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
