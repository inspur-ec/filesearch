package office

import (
	"fmt"
	"github.com/nguyenthenguyen/docx"
)

func Modify() {
	// Read from docx file
	r, err := docx.ReadDocxFile("./template.docx")
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	if err != nil {
		panic(err)
	}
	docx1 := r.Editable()
	fmt.Printf("content:[%s]", docx1.GetContent())
	// Replace like https://golang.org/pkg/strings/#Replace
	docx1.Replace("old_1_1", "new_1_1", -1)
	docx1.Replace("old_1_2", "new_1_2", -1)
	docx1.ReplaceHeader("out with the old", "in with the new")
	docx1.ReplaceFooter("Change This Footer", "new footer")
	docx1.WriteToFile("./new_result_1.docx")

	docx2 := r.Editable()
	docx2.Replace("old_2_1", "new_2_1", -1)
	docx2.Replace("old_2_2", "new_2_2", -1)
	docx2.WriteToFile("./new_result_2.docx")

	// Or write to ioWriter
	// docx2.Write(ioWriter io.Writer)

	r.Close()
}

func ReadWord(filename string) string {
	// Read from docx file
	r, err := docx.ReadDocxFile(filename)
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	if err != nil {
		fmt.Printf("error is cuur [%v]\n", err)
		return ""
		//panic(err)
	}
	docx1 := r.Editable()
	return docx1.GetContent()

}
