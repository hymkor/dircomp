package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var option_i = flag.String("i", "*", "compared files pattern")

func AreSameFiles(fname1, fname2 string) (bool, error) {
	fd1, err := os.Open(fname1)
	if err != nil {
		return false, err
	}
	defer fd1.Close()

	fd2, err := os.Open(fname2)
	if err != nil {
		return false, err
	}
	defer fd2.Close()

	r1 := bufio.NewReader(fd1)
	r2 := bufio.NewReader(fd2)
	for {
		b1, err1 := r1.ReadByte()
		b2, err2 := r2.ReadByte()

		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		}
		if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}
		if err1 != nil {
			return false, err1
		}
		if err2 != nil {
			return false, err2
		}
		if b1 != b2 {
			return false, nil
		}
	}

}

type FileT struct {
	Size int64
	Path string
}

func (file1 FileT) Equals(file2 FileT) (bool, error) {
	if file1.Size != file2.Size {
		return false, nil
	}
	return AreSameFiles(file1.Path, file2.Path)
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
				fmt.Printf("M\t%s\n", filepath.Base(info1.Path))
			}
		} else {
			fmt.Printf("D\t%s\n", filepath.Base(info1.Path))
		}
	}
	for name, info2 := range dirlist2 {
		if _, ok := dirlist1[name]; !ok {
			fmt.Printf("A\t%s\n", filepath.Base(info2.Path))
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
