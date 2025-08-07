package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 创建一个新的工作簿
	f := excelize.NewFile()
	// 确保工作簿在关闭时释放资源
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建一个工作表，返回工作表的索引和错误信息
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 定义数据
	data := [][]interface{}{
		{nil, "Apple", "Orange", "Pear"},
		{"Small", 2, 3, 3},
		{"Normal", 5, 2, 4},
		{"Large", 6, 7, 8},
	}

	// 写入数据
	for rowIdx, row := range data {
		for colIdx, cellValue := range row {
			// 将单元格的行列号转换为单元格名称，例如：1, 2 -> B3
			cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				fmt.Printf("生成单元格名称失败: %v\n", err)
				return
			}
			// 设置单元格的值
			if err := f.SetCellValue("Sheet2", cellName, cellValue); err != nil {
				fmt.Errorf("设置单元格值失败: %v", err)
				return
			}
		}
	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
