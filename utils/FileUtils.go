package utils

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func getFileAndDirectory(fileBasePath string, myFile []fs.FileInfo) []fs.FileInfo {
	dir, err := ioutil.ReadDir(fileBasePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		myFile = append(myFile, file)
		if file.IsDir() {
			path := fileBasePath + "\\" + file.Name()
			getFileAndDirectory(path, myFile)
		}
	}

	return myFile
}

func GetFileAndDirectory(fileBasePath string) []fs.FileInfo {
	var myFile = make([]fs.FileInfo, 100)
	myFile = getFileAndDirectory(fileBasePath, myFile)
	return myFile
}

func GetAllFile(basePath string) []string {
	var result []string
	dir, err := ioutil.ReadDir(basePath)

	if err != nil {
		log.Fatal(err)
	}
	for _, info := range dir {
		fullName := basePath + "\\" + info.Name()
		result = append(result, fullName)
		if info.IsDir() {
			file := GetAllFile(fullName)
			result = append(result, file...)
		}
	}
	return result
}

func GetFileByName(files []string, name string) []string {
	var result []string
	for _, file := range files {
		//fmt.Println("根据[", name, "]，查找文件，匹配文件[", file, "]")
		if strings.Contains(file, name) {
			result = append(result, file)
		}
	}
	if len(result) == 0 {
		fmt.Println("根据[", name, "]，查找文件，未能找到匹配文件")
	}
	return result
}

func CopyDir(srcPath string, destPath string) error {

	if stat, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if stat.IsDir() {
			err := os.MkdirAll(destPath, stat.Mode())
			if err != nil {
				return err
			}
			files, err := ioutil.ReadDir(srcPath)
			if err != nil {
				return err
			}
			for _, file := range files {
				srcFilePath := srcPath + "\\" + file.Name()
				destFilePath := destPath + "\\" + file.Name()
				CopyDir(srcFilePath, destFilePath)
			}
			return nil
		} else {
			srcFile, err := os.Open(srcPath)
			fmt.Println("复制的文件名", srcPath)
			fmt.Println("组装之后的文件名", destPath+"\\"+srcFile.Name())
			if err != nil {
				return err
			}
			defer srcFile.Close()
			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
