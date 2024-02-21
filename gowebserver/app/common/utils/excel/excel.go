package excel

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/tealeg/xlsx"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//CreateFilePath
//@summary 创建路径,路径不存在创建路径
func CreateFilePath(filePath string) error {
	path, _ := filepath.Split(filePath) //获取路径
	_, err := os.Stat(path)             //检查路径状态，不存在创建
	if err != nil || os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	return err
}

//DownloadExcel
//@summary gdb数据查询结果导出到Excel表
//@param1  heads []string     "Excel表头"
//@param2  key   []string     "Excel表头对应字段"
//@param3  data  gdb.Result   "gdb数据查询结果集"
//@param4  filename ...string "可选,指定导出Excel文件名"
//@return1 string "导出Excel文件名"
//@return2 error  "报错信息"
func DownloadExcel(heads, key []string, data gdb.Result, filename ...string) (string, error) {
	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		row   *xlsx.Row
		cell  *xlsx.Cell
	)

	//导出Excel文件名
	var fileName string
	if len(filename) > 0 && filename[0] != "" {
		fileName = filename[0]
	} else {
		fileName = strconv.FormatInt(time.Now().UnixNano(), 10) + ".xls"
	}

	//创建路径
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filePath := curDir + g.Cfg().GetString("download.downPath") + "/" + fileName
	err = CreateFilePath(filePath)
	if err != nil {
		log.Printf("%s", err.Error())
		return "", err
	}

	//创建Excel文件
	file = xlsx.NewFile()
	//添加新工作表
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	//设置样式
	//sheet.SetColWidth(5, 5, 60) //设置单元格宽度 0-A 1-B 2-C

	//向工作表中添加表头信息
	row = sheet.AddRow()
	for _, head := range heads {
		cell = row.AddCell()
		cell.Value = head
	}

	//主体写入数据
	for _, record := range data {
		row = sheet.AddRow()
		for _, v := range key {
			row.AddCell().Value = record[v].String()
		}
	}

	//将数据保存到xlsx文件
	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

//DownloadExcelMap
//@summary map数据集保存为Excel表
//@param1  heads []string     "Excel表头"
//@param2  key   []string     "Excel表头对应字段"
//@param3  data  []map[string]interface{}   "Map数据集"
//@param4  filename ...string "可选,指定导出Excel文件名"
//@return1 string "导出Excel文件名"
//@return2 error  "报错信息"
func DownloadExcelMap(heads, key []string, data []map[string]interface{}, filename ...string) (string, error) {
	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		row   *xlsx.Row
		cell  *xlsx.Cell
	)

	//导出Excel文件名
	var fileName string
	if len(filename) > 0 && filename[0] != "" {
		fileName = filename[0]
	} else {
		fileName = strconv.FormatInt(time.Now().UnixNano(), 10) + ".xls"
	}

	//创建路径
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filePath := curDir + g.Cfg().GetString("download.downPath") + "/" + fileName
	err = CreateFilePath(filePath)
	if err != nil {
		log.Printf("%s", err.Error())
		return "", err
	}

	//创建Excel文件
	file = xlsx.NewFile()
	//添加新工作表
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	//设置样式
	//sheet.SetColWidth(5, 5, 60) //设置单元格宽度 0-A 1-B 2-C

	//向工作表中添加表头信息
	row = sheet.AddRow()
	for _, head := range heads {
		cell = row.AddCell()
		cell.Value = head
	}

	//主体写入数据
	for _, record := range data {
		row = sheet.AddRow()
		for _, v := range key {
			row.AddCell().Value = gconv.String(record[v])
		}
	}

	//将数据保存到xlsx文件
	err = file.Save(filePath)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
