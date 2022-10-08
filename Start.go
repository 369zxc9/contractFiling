package main

import (
	"contractFiling/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {

	readFile, err := ioutil.ReadFile("./config.yaml")

	if err != nil {
		fmt.Println(err)
		return
	}

	c := config{}
	err = yaml.Unmarshal(readFile, &c)

	if err != nil {
		fmt.Println(err)
		return
	}

	searchPath := c.FileSearchPath

	row := c.Row

	name := c.Name

	if !(len(searchPath) == len(row) && len(searchPath) == len(name)) {
		fmt.Println("请正确配置参数")
		return
	}

	fmt.Println("确认配置的参数是否正确，\n如果有问题，手动关闭cmd窗口：\n", c, " \n三秒后开始执行文件整理程序……")

	time.Sleep(3 * time.Second)

	var searchSave [][]string
	for _, s := range searchPath {
		var save []string
		searchSave = append(searchSave, append(save, utils.GetAllFile(s)...))
	}

	if searchSave == nil {
		fmt.Println("未找到文件")
		return
	}

	//path := c.ContractSearchPath
	//
	//var contractSearchFiles []string
	//for _, searchPath := range path {
	//	contractSearchFiles = append(contractSearchFiles, utils.GetAllFile(searchPath)...)
	//}
	//
	//anyPath := c.AnyThingSearchPath
	//var anyThingSearchFiles []string
	//for _, searchPath := range anyPath {
	//	anyThingSearchFiles = append(anyThingSearchFiles, utils.GetAllFile(searchPath)...)
	//}

	// companyMap := make(map[string]string)

	openFile, err := excelize.OpenFile(c.ExcelReadPath)
	if err != nil {
		fmt.Println("读取excel失败：", err)
		return
	}
	excelReadSheetName := c.ExcelReadSheetName
	rows, err := openFile.GetRows(excelReadSheetName)
	for i := range rows {
		if c.ExcelReadStartRow > i {
			continue
		}
		s := strconv.Itoa(i)
		companyRow := c.CompanyNameRow + s
		contractRow := c.ContractRow + s
		categoryRow := c.CategoryRow + s

		companyName, err := openFile.GetCellValue(excelReadSheetName, companyRow)
		categoryName, err := openFile.GetCellValue(excelReadSheetName, categoryRow)

		contractNo, err := openFile.GetCellValue(excelReadSheetName, contractRow)

		var companySavePath string

		//if companyMap[companyName] == "" {
		//	companySavePath = c.BaseSavePath + "\\" + companyName
		//	err := os.MkdirAll(companySavePath, os.ModePerm)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//	companyMap[companyName] = companySavePath
		//} else {
		//	companySavePath = companyMap[companyName]
		//}

		companySavePath = c.BaseSavePath + "\\" + categoryName + "\\" + companyName
		err = os.MkdirAll(companySavePath, os.ModePerm)
		contractSavePath := companySavePath + "\\" + contractNo + "\\"

		err = os.MkdirAll(contractSavePath, os.ModePerm)
		for index, saveDirName := range name {
			savePath := contractSavePath + saveDirName

			fmt.Println("搜索的文件名：", searchSave[index])
			rowValue, err := openFile.GetCellValue(excelReadSheetName, row[index]+s)
			fmt.Println(err)
			getFileByName := utils.GetFileByName(searchSave[index], rowValue)

			for _, file := range getFileByName {
				err := utils.CopyDir(file, savePath)
				fmt.Println(err)
			}
		}

		fmt.Println(err)
	}
}

//defer func() {
//	if err := file.Close(); err != nil {
//		fmt.Println(err)
//	}
//}()
//strings.Replace("", ".xlsx", "--副本.xlsx", -1)
//// 获取 Sheet1 上所有单元格
//rows, err := file.GetRows("测试")
//if err != nil {
//	fmt.Println(err)
//	return
//}
//err = file.SetCellValue("测试", "A11", "Hello world.")
//
//err = file.SaveAs("C:\\Users\\admin\\Desktop\\新建 XLSX 工作表--副本.xlsx")
//if err != nil {
//	return
//}
//if err != nil {
//	fmt.Println(err)
//}
//for _, row := range rows {
//	for _, cellValue := range row {
//		fmt.Print(cellValue, "\t")
//	}
//	fmt.Println()
//}

// value, err := file.GetCellValue("Sheet1", "B2")

type config struct {
	// 合同文件搜索目录
	ContractSearchPath []string `yaml:"contractSearchPath,flow"`
	// 整理后保存目录
	BaseSavePath string `yaml:"baseSavePath"`
	// 文件搜索目录
	FileSearchPath []string `yaml:"fileSearchPath,flow"`
	// excel文件所在目录
	ExcelReadPath string `yaml:"excelReadPath"`
	// excel 读取的sheet名称
	ExcelReadSheetName string `yaml:"excelReadSheetName"`
	// excel开始读取的行数
	ExcelReadStartRow int `yaml:"excelReadStartRow"`
	// 公司名称所在行
	CompanyNameRow string `yaml:"companyNameRow"`

	ContractRow string `yaml:"contractRow"`
	// 搜索条件所在列
	Row []string `yaml:"row,flow"`
	// 搜索之后保存的文件夹名称
	Name []string `yaml:"name,flow"`
	// 类别
	CategoryRow string `yaml:"categoryRow"`
}
