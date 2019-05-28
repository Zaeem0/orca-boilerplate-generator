package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createZip(templateName string) error {

	//get files from standard folder
	os.Rename("./tmp", fmt.Sprintf("./%s", templateName))
	files, err := listFiles(fmt.Sprintf("./%s", templateName))
	if err != nil {
		return err
	}
	os.Mkdir("./downloads", 0700)
	err = zipMe(files, fmt.Sprintf("./downloads/%s.zip", templateName))
	if err != nil {
		return err
	}
	for _, f := range files {
		fmt.Println(f)
	}
	//delete folder now that zip is created
	os.RemoveAll(fmt.Sprintf("./%s", templateName))
	fmt.Println("Done!")

	return nil
}

func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil

}

func zipMe(filepaths []string, target string) error {

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(target, flags, 0644)

	if err != nil {
		return fmt.Errorf("Failed to open zip for writing: %s", err)
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range filepaths {
		if err := addFileToZip(filename, zipw); err != nil {
			return fmt.Errorf("Failed to add file %s to zip: %s", filename, err)
		}
	}
	return nil

}

func addFileToZip(filename string, zipw *zip.Writer) error {
	file, err := os.Open(filename)

	if err != nil {
		return fmt.Errorf("Error opening file %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {

		return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Error writing %s to zip: %s", filename, err)
	}

	return nil
}
