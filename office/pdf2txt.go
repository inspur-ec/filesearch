package office

import (
	"bytes"
	"github.com/ledongthuc/pdf"
)

func ReadPdf(path string) (res string, err error) {

	defer func() (string, error) {
		if err := recover(); err != nil {
			return "", err.(error)
		}
		return "", nil
	}()
	_, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	var textBuilder bytes.Buffer
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		t, _ := p.GetPlainText(nil)
		textBuilder.WriteString(t)
	}
	return textBuilder.String(), nil
}
