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

		if strings.Contains(file, ".") {
			// 说明是文件
			if strings.Contains(file, "\\"+name+".") {
				result = append(result, file)
			}
		} else {
			if strings.HasSuffix(file, "\\"+name) {
				result = append(result, file)
			}
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
		if !strings.Contains(destPath, ".") {
			err := os.MkdirAll(destPath, stat.Mode())
			if err != nil {
				fmt.Println(err)
			}
		}

		if stat.IsDir() {

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
			if err != nil {
				fmt.Println(err)
			}
			defer srcFile.Close()

			split := strings.Split(srcPath, "\\")

			s := split[len(split)-1]
			fmt.Println("保存的文件名", destPath+"\\"+s)
			destFile, err := os.OpenFile(destPath+"\\"+s, os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("创建文件失败")
				fmt.Println(err)
			}
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				fmt.Println("拷贝文件失败")
				fmt.Println(err)
			}

			//input, err := ioutil.ReadFile(srcPath)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//split := strings.Split(srcPath, "\\")
			//
			//s := split[len(split)-1]
			//err = ioutil.WriteFile(destPath + "\\" + s, input, 0644)
			//if err != nil {
			//	fmt.Println("Error creating", destPath + "\\" + s)
			//	fmt.Println(err)
			//}
			return nil
		}
	}
}
