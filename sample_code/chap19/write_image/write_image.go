package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 新建 Excel 文件
	f := excelize.NewFile()

	// 确保文件在退出时关闭
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 插入图片
	if err := f.AddPicture("Sheet1", "A2", "./pics/Golang.png", nil); err != nil {
		fmt.Println(err)
		return
	}
	// 在工作表中插入图片，并设置图片的缩放比例
	if err := f.AddPicture("Sheet1", "E2", "./pics/Linux.jpg",
		&excelize.GraphicOptions{ScaleX: 0.4, ScaleY: 0.4}); err != nil {
		fmt.Println(err)
		return
	}
	// 在工作表中插入图片，并设置图片的打印属性
	enable, disable := true, false
	if err := f.AddPicture("Sheet1", "I2", "./pics/Java.gif",
		&excelize.GraphicOptions{
			PrintObject:     &enable,
			LockAspectRatio: false,
			OffsetX:         15,
			OffsetY:         10,
			Locked:          &disable,
			ScaleX:          0.8,
			ScaleY:          0.8,
		}); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
