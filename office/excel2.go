package office

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func Write() {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
func ReadExcel(filename string) (result string, err error) {
	result = ""
	f, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	m := f.GetSheetMap()
	for j := 0; j < len(m); j++ {
		rows := f.GetRows(m[j-1])
		for _, row := range rows {
			for _, colCell := range row {
				result = result + " " + colCell
			}
		}
	}
	// Get value from cell by given worksheet name and axis.
	//cell:= f.GetCellValue("Sheet1", "B2")
	//fmt.Printf("cell[%s]\n",cell)
	// Get all the rows in the Sheet1.
	//rows := f.GetRows("Sheet1")
	//for i, row := range rows {
	//	for j, colCell := range row {
	//		fmt.Printf("col[%d]row[%d]:%s\n",i,j,colCell)
	//	}
	//}
	return
}

//插入图表
func AddGraph() {
	categories := map[string]string{"A2": "Small", "A3": "Normal", "A4": "Large", "B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	f := excelize.NewFile()
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}
	if err := f.AddChart("Sheet1", "E1", `{"type":"col3DClustered","series":[{"name":"Sheet1!$A$2","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$2:$D$2"},{"name":"Sheet1!$A$3","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$3:$D$3"},{"name":"Sheet1!$A$4","categories":"Sheet1!$B$1:$D$1","values":"Sheet1!$B$4:$D$4"}],"title":{"name":"Fruit 3D Clustered Column Chart"}}`); err != nil {
		println(err.Error())
		return
	}
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		println(err.Error())
	}
}

func AddPic() {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		println(err.Error())
		return
	}
	// 插入图片
	if err := f.AddPicture("Sheet1", "A2", "image.png", ""); err != nil {
		println(err.Error())
	}
	// 在工作表中插入图片，并设置图片的缩放比例
	if err := f.AddPicture("Sheet1", "D2", "image.jpg", `{"x_scale": 0.5, "y_scale": 0.5}`); err != nil {
		println(err.Error())
	}
	// 在工作表中插入图片，并设置图片的打印属性
	if err := f.AddPicture("Sheet1", "H2", "image.gif", `{"x_offset": 15, "y_offset": 10, "print_obj": true, "lock_aspect_ratio": false, "locked": false}`); err != nil {
		println(err.Error())
	}
	// 保存文件
	if err = f.Save(); err != nil {
		println(err.Error())
	}
}
