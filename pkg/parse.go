package pkg

import (
	"archive/zip"
	"fmt"
	"github.com/secsy/goftp"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// matching items from folder by mask pattern

func MatchingFolder(list []string, pattern string) []string {
	var outputList []string

	for _, item := range list {
		if strings.ToLower(item[0:len(pattern)-1]) == strings.ToLower(pattern[:len(pattern)-1]) {
			outputList = append(outputList, item)
		}
	}

	//if len(outputList) == 0 {
	//	fmt.Println("0 files was matched")
	//}
	return outputList
}

func GetNextDay(s string) string {
	now, _ := time.Parse("2006010200", s)
	tomorrow := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, time.UTC).
		AddDate(0, 0, 1).
		Format("2006010200")
	return tomorrow
}

func ExtractZip(list []string, path, dest string) {
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
func DownloadZips(c *goftp.Client, list []string, path, output string) {
	if len(list) == 0 {
		return
	}
	for _, item := range list {
		out, err := os.Create(output + item)
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		_ = c.Retrieve(path+item, out)
	}

}
