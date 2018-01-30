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
		hash1, err := GetFileHash(path1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path1, err.Error())
		} else {
			dict[strings.ToUpper(file1.Name())] = hash1
		}
	}
	return dict, nil
}

func main1(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: dirdiff DIR1 DIR2")
	}
	dirhash1, err := GetDirHash(args[0])
	if err != nil {
		return fmt.Errorf("%s: %s", args[0], err.Error())
	}
	dirhash2, err := GetDirHash(args[1])
	if err != nil {
		return fmt.Errorf("%s: %s\n", args[1], err.Error())
	}

	for name, hash1 := range dirhash1 {
		if hash2, ok := dirhash2[name]; ok {
			if hash1 != hash2 {
				fmt.Printf("M\t%s\n", name)
			}
		} else {
			fmt.Printf("D\t%s\n", name)
		}
	}
	for name, _ := range dirhash2 {
		if _, ok := dirhash1[name]; !ok {
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
