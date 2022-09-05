package pkg

import (
	"archive/zip"
	"fmt"
	"github.com/secsy/goftp"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetFolders(client *goftp.Client, path string) []string {
	files, err := client.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var list []string
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list
}

func MatchingFolder(list []string, pattern string) []string {
	var outputList []string

	for _, item := range list {
		if strings.ToLower(item[0:len(pattern)-1]) == strings.ToLower(pattern[:len(pattern)-1]) {
			outputList = append(outputList, item)
		}
	}
	return outputList
}
func GetNextDay(s string) string {
	num, _ := strconv.Atoi(s[6:8])
	num++
	if num >= 10 {
		return s[0:6] + strconv.Itoa(num) + s[8:10]
	} else {
		return s[0:6] + "0" + strconv.Itoa(num) + s[8:10]
	}
}

func Extract(list []string, path, dest string) {
	for _, file := range list {

		dst := dest + file[0:len(file)-7]
		fmt.Println(dest)
		archive, err := zip.OpenReader(path + file)
		if err != nil {
			panic(err)
		}
		defer archive.Close()

		for _, f := range archive.File {
			filePath := filepath.Join(dst, f.Name)

			if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return
			}
			if f.FileInfo().IsDir() {
				fmt.Println("creating directory...")
				os.MkdirAll(filePath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				panic(err)
			}

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(err)
			}

			fileInArchive, err := f.Open()
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				panic(err)
			}

			dstFile.Close()
			fileInArchive.Close()
		}
	}
}
