package gui

import (
	"filesearch/es"
	"filesearch/utils"
	"fmt"
	"log"
	"math/rand"
	"strings"

	//"sort"
	//"strings"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Foo struct {
	Index    int
	Filename string
	Filetype string
	Filepath string
	Content  string
	//	checked  bool
}

type FooModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Foo
}

const PIXERS = 4
const HIGH_LIGHT_COLOR="green"
func NewFooModel() *FooModel {
	m := new(FooModel)
	//m.Search()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *FooModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *FooModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Filename
	case 2:
		return item.Filepath

	case 3:
		return item.Filetype

	case 4:
		return item.Content
	}

	panic("unexpected col")
}

/*
// Called by the TableView to retrieve if a given row is checked.
func (m *FooModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *FooModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

*/
func (m *FooModel) Search() {
	// Create some random data.
	//now := time.Now()
	c := strings.Trim(myMainWindow.le.Text(), " ")
	c = strings.ReplaceAll(c, " ", "*")
	/**
	if myMainWindow.checkdoc.Checked(){
		filetype="doc"
	}
	if myMainWindow.checkpdf.Checked(){
		filetype="pdf"
	}
	if myMainWindow.checkexcel.Checked(){
		filetype="xls"
	}*/
	filetype:=""
	ss:=""
	//ss := es.SearchALL("*"+c+"*",filetype, HIGH_LIGHT_COLOR)
	if (myMainWindow.checkdoc.Checked()&&myMainWindow.checkpdf.Checked()&&myMainWindow.checkexcel.Checked()){
		ss=es.Search("*"+c+"*", HIGH_LIGHT_COLOR)
	}else{
		if myMainWindow.checkdoc.Checked(){
			filetype="doc"
		}else if myMainWindow.checkpdf.Checked(){
			filetype="pdf"
		}else if myMainWindow.checkexcel.Checked() {
			filetype = "xls"
		}
		if (filetype!=""){
			ss = es.SearchALL("*"+c+"*",filetype, HIGH_LIGHT_COLOR)
		}
	}
	s := es.FormatResult(ss)
	m.items = make([]*Foo, len(s.Hits.Hits))
	for i, r := range s.Hits.Hits {
		files := r.Source.Filename
		files = strings.ReplaceAll(files, "\\", "/")
		n := strings.LastIndex(files, "/")
		filename1 := files
		filepath1 := ""
		if n > 0 {
			filename1 = files[n+1:]
			filepath1 = files[:n]
		}
		m.items[i] = &Foo{
			Index:    i,
			Filename: filename1,
			Filepath: filepath1,
			Filetype: r.Source.Filetype,
			Content:  r.Highlight.Content[0],
		}

	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

func (m *FooModel) InitIndex() {
	i := walk.MsgBox(myMainWindow, "提示", "该操作初始化索引库，是否确认重新初始化？", walk.MsgBoxOKCancel)
	if i != 1 {
		return
	}
	t := time.Now()
	es.Delete("localindex")
	success, errors, errs := es.InitIndex()
	msg := fmt.Sprintf("初始化索引完成,files success count[%d],errors count [%d] time consume [%d]s\n files error is [%s]", success, errors, time.Now().Sub(t)/1e9, errs)
	walk.MsgBox(myMainWindow, "提示", msg, walk.MsgBoxIconInformation)
}

func OnMouseUp() {
	v := myMainWindow.tv
	i := v.CurrentIndex()
	if i < 0 {
		return
	}
	m := v.Model().(*FooModel)

	utils.Open(m.items[i].Filepath + "/" + m.items[i].Filename)
}

func OnMouseMove() {
}

type TableMainWindow struct {
	*walk.MainWindow
	checkdoc        *walk.CheckBox
	checkpdf        *walk.CheckBox
	checkexcel        *walk.CheckBox
	le        *walk.LineEdit
	wv        *walk.WebView
	tv        *walk.TableView
	query     *walk.PushButton
	initindex *walk.PushButton
}

var myMainWindow *TableMainWindow

func OpenSearchWindows() {
	rand.Seed(time.Now().UnixNano())
	//boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)

	//boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	//goodIcon, _ := walk.Resources.Icon("../img/check.ico")
	//badIcon, _ := walk.Resources.Icon("../img/stop.ico")

	barBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas, err := walk.NewCanvasFromImage(barBitmap)
	if err != nil {
		panic(err)
	}
	defer barBitmap.Dispose()

	canvas.GradientFillRectangle(walk.RGB(255, 0, 0), walk.RGB(0, 255, 0), walk.Horizontal, walk.Rectangle{0, 0, 100, 1})

	canvas.Dispose()

	myMainWindow = &TableMainWindow{}
	model := NewFooModel()
	//var tv *walk.TableView
	//var le *walk.LineEdit
	MainWindow{
		AssignTo:   &myMainWindow.MainWindow,
		Title:      "文件查询",
		Size:       Size{1366, 768},
		Layout:     VBox{MarginsZero: true},
		Persistent: true,
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					CheckBox{
						AssignTo: &myMainWindow.checkdoc,
						Name: "word",
						Text: "word",
						Checked:true,
					},
					CheckBox{
						AssignTo: &myMainWindow.checkpdf,
						Name: "pdf",
						Text: "pdf",
						Checked:true,
					},
					CheckBox{
						AssignTo: &myMainWindow.checkexcel,
						Name: "excel",
						Text: "excel",
						Checked:true,
					},
					LineEdit{
						AssignTo: &myMainWindow.le,
						Text:     "",
						OnKeyPress: func(key walk.Key) {
							if key == 13 {
								model.Search()
							}
						},
					},
					PushButton{
						AssignTo:  &myMainWindow.query,
						Text:      "查询",
						OnClicked: model.Search,
					},
					PushButton{
						AssignTo:  &myMainWindow.initindex,
						Text:      "初始化",
						OnClicked: model.InitIndex,
					},
				},
			},
			TableView{
				AssignTo:         &myMainWindow.tv,
				AlternatingRowBG: true,
				CheckBoxes:       false,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Columns: []TableViewColumn{
					{Title: "序号", Width: 40},
					{Title: "文件名", Width: 300},
					{Title: "路径", Width: 400},
					{Title: "类型", Width: 40},
					{Title: "内容", Width: 600},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]

					//if item.checked {
					if style.Row()%2 == 0 {
						style.BackgroundColor = walk.RGB(225, 255, 255)
					} else {
						style.BackgroundColor = walk.RGB(255, 255, 255)
					}
					//}
					switch style.Col() {
					case 1:
						/**
						if canvas := style.Canvas(); canvas != nil {
							bounds := style.Bounds()
							bounds.X += 2
							bounds.Y += 2
							bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Filename)))
							bounds.Height -= 4
							canvas.DrawBitmapPartWithOpacity(barBitmap, bounds, walk.Rectangle{0, 0, 100 / 5 * len(item.Filename), 1}, 127)

							bounds.X += 4
							bounds.Y += 2
							canvas.DrawText(item.Filename, myMainWindow.tv.Font(), 0, bounds, walk.TextLeft)
						}
						*/
					case 2:
						//if item.Baz >= 900.0 {
						//	style.TextColor = walk.RGB(0, 191, 0)
						//	style.Image = goodIcon
						//} else if item.Baz < 100.0 {
						//	style.TextColor = walk.RGB(255, 0, 0)
						//	style.Image = badIcon
						//}

					case 3:
						//if item.Quux.After(time.Now().Add(-365 * 24 * time.Hour)) {
						//	style.Font = boldFont
						//}
					case 4:
						{
							if canvas := style.Canvas(); canvas != nil {
								bounds := style.Bounds()
								s := strings.Split(item.Content, "</font>")
								for _, t := range s {
									t1 := strings.Split(t, "<font color='"+HIGH_LIGHT_COLOR+"'>")
									bounds.Width = len(t1[0])*PIXERS
									log.Printf("bounds.X=%d;t1[0]=%s", bounds.X, t1[0])
									canvas.DrawTextPixels(t1[0], myMainWindow.tv.Font(), 0, bounds, walk.TextLeft)
									if (len(t1) > 1) {
										bounds.X = bounds.X + len(t1[0])*PIXERS
										bounds.Width = len(t1[1])*PIXERS
										log.Printf("bounds.X=%d;t1[1]=%s", bounds.X, t1[1])
										canvas.DrawTextPixels(t1[1],  myMainWindow.tv.Font(), walk.RGB(255, 0, 0), bounds, walk.TextLeft)
										bounds.X = bounds.X + len(t1[1])*PIXERS
									}
								}
							}
						}
					}
				},
				Model: model,
				OnSelectedIndexesChanged: func() {
					log.Printf("SelectedIndexes: %v\n", myMainWindow.tv.SelectedIndexes())
				},
				OnMouseUp: func(x, y int, button walk.MouseButton) {
					OnMouseUp()
				},
				OnMouseMove: func(x, y int, button walk.MouseButton) {
					OnMouseMove()
				},
			},
		},
	}.Run()
}
