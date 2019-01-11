package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func replace(filename string, str1 string, str2 string, str3 string, str4 string) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	br := bufio.NewReader(f)
	allStr := []string{}
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		allStr = append(allStr, string(line))
	}
	f.Truncate(0)
	f.Seek(0, 0)

	bw := bufio.NewWriter(f)
	for _, v := range allStr {
		n1 := strings.Index(v, str1)
		if n1 >= 0 {
			if len(str2) == 0 {
				vs := v[n1+len(str1):]
				v = strings.Replace(v, str1+vs, str3+vs+str4, 1)
			} else {
				n2 := strings.Index(v[n1+len(str1):], str2)
				if n2 >= 0 {
					vs := v[n1+len(str1) : n1+len(str1)+n2]
					v = strings.Replace(v, str1+vs+str2, str3+vs+str4, 1)
				}
			}
		}
		bw.WriteString(v)
		bw.WriteString("\n")
	}
	bw.Flush()
	defer f.Close()
}

func getFilelist(path string) []string {
	paths := []string{}
	err := filepath.Walk(path, func(fi string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			if fi != path {
				paths = append(paths, getFilelist(fi)...) //若不递归遍历则注释此行即可
			}
		} else {
			paths = append(paths, fi)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return paths
}

func main() {
	paths := getFilelist("test")
	for _, v := range paths {
		replace(v, "hello", "world", "HELLO", "WORLD")
		replace(v, "HELLO", "", "hello", "world")
	}
}
