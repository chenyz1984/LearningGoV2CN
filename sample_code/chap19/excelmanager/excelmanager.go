package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/xuri/excelize/v2"
)

// ExcelManager Excel 文件管理器
type ExcelManager struct {
	file *excelize.File
}

// NewExcelManager 创建一个新的 ExcelManager 实例
func NewExcelManager() *ExcelManager {
	return &ExcelManager{
		file: excelize.NewFile(),
	}
}

// CreateExcelFile 函数1：创建 Excel 文档
func (em *ExcelManager) CreateExcelFile(filename, sheetName string) error {
	// 确保至少有一个工作表
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// 创建指定名称的工作表
	_, err := em.file.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("创建工作表失败: %v", err)
	}

	// 设置为活动工作表
	em.file.SetActiveSheet(0)

	// 保存文件
	if err := em.file.SaveAs(filename); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// WriteData 函数2：向 Excel 文档写入数据
func (em *ExcelManager) WriteData(filename, sheetName string, data [][]interface{}) error {
	// 打开文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 检查工作表是否存在
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		// 工作表不存在，创建新工作表
		index, err = f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建工作表失败: %v", err)
		}
	}

	// 写入数据
	for rowIdx, row := range data {
		for colIdx, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				return fmt.Errorf("生成单元格名称失败: %v", err)
			}
			if err := f.SetCellValue(sheetName, cellName, cellValue); err != nil {
				return fmt.Errorf("设置单元格值失败: %v", err)
			}
		}
	}

	// 保存文件
	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// ReadData 函数3：读取 Excel 文档内容
func (em *ExcelManager) ReadData(filename, sheetName string) ([][]string, error) {
	// 打开文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 检查工作表是否存在
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		return nil, fmt.Errorf("工作表 %s 不存在", sheetName)
	}

	// 获取所有行数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %v", err)
	}

	return rows, nil
}

// InsertImage 函数4：向 Excel 文档插入图片
func (em *ExcelManager) InsertImage(filename, sheetName, cell, imagePath string, options *excelize.GraphicOptions) error {
	// 打开文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 检查工作表是否存在
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		// 工作表不存在，创建新工作表
		index, err = f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建工作表失败: %v", err)
		}
	}

	// 插入图片
	if err := f.AddPicture(sheetName, cell, imagePath, options); err != nil {
		return fmt.Errorf("插入图片失败: %v", err)
	}

	// 保存文件
	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// InsertChart 函数5：向 Excel 文档插入图表
func (em *ExcelManager) InsertChart(filename, sheetName, cell string, chart *excelize.Chart) error {
	// 打开文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer f.Close()

	// 检查工作表是否存在
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		// 工作表不存在，创建新工作表
		index, err = f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建工作表失败: %v", err)
		}
	}

	// 插入图表
	if err := f.AddChart(sheetName, cell, chart); err != nil {
		return fmt.Errorf("插入图表失败: %v", err)
	}

	// 保存文件
	if err := f.Save(); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

func main() {
	// 创建 ExcelManager 实例
	em := NewExcelManager()

	// 1. 创建 Excel 文档
	err := em.CreateExcelFile("Book1.xlsx", "Sheet1")
	if err != nil {
		fmt.Println("创建 Excel 文档失败:", err)
		return
	}
	fmt.Println("Excel 文档创建成功")

	// 2. 向 Excel 文档写入数据
	data := [][]interface{}{
		{nil, "Apple", "Orange", "Pear"},
		{"Small", 2, 3, 3},
		{"Normal", 5, 2, 4},
		{"Large", 6, 7, 8},
	}

	err = em.WriteData("Book1.xlsx", "Sheet1", data)
	if err != nil {
		fmt.Println("写入数据失败:", err)
		return
	}
	fmt.Println("数据写入成功")

	// 3. 读取 Excel 文档内容
	rows, err := em.ReadData("Book1.xlsx", "Sheet1")
	if err != nil {
		fmt.Println("读取数据失败:", err)
		return
	}

	fmt.Println("读取到的数据:")
	for _, row := range rows {
		for _, col := range row {
			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}

	// 4. 向 Excel 文档插入图表
	chart := &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: "Fruit 3D Clustered Column Chart",
			},
		},
	}

	err = em.InsertChart("Book1.xlsx", "Sheet1", "F1", chart)
	if err != nil {
		fmt.Println("插入图表失败:", err)
		return
	}
	fmt.Println("图表插入成功")

	// 5. 向 Excel 文档插入图片 (需要实际有图片文件)

	err = em.InsertImage("Book1.xlsx", "Sheet1", "A10", "./pics/image.png", nil)
	if err != nil {
		fmt.Println("插入图片失败:", err)
		return
	}
	fmt.Println("图片插入成功")
}
