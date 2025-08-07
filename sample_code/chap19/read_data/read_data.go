package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 打开文件
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 检查工作表是否存在
	sheetName := "Sheet2"
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		fmt.Errorf("工作表 %s 不存在", sheetName)
		return
	}

	// 获取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Errorf("读取工作表数据失败: %v", err)
		return
	}

	fmt.Println("读取到的数据:")
	// 遍历并打印所有行和列数据
	for _, row := range rows {
		for _, col := range row {
			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
}
