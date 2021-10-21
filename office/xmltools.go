package office

/*

使用 etree 解析复杂结构的 xml 文件
https://godoc.org/github.com/beevik/etree
https://pkg.go.dev/github.com/beevik/etree?tab=doc
https://github.com/beevik/etree
*/

import (
	"errors"
	"github.com/beevik/etree" // go get github.com/beevik/etree
	"strings"
)

func ReadXml(xmls string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes([]byte(xmls)); err != nil {
		return "", err
	}

	root := doc.SelectElement("w:document")
	//fmt.Println("ROOT element:", root.Tag)
	if root == nil {
		return "", errors.New(" xml 文件 root不存在")
	}
	body := root.SelectElement("w:body")
	txt := readElement(body)
	//fmt.Printf("content[%s]",txt)
	return txt, nil
}

func readElement(body *etree.Element) (result string) {

	if len(body.ChildElements()) == 0 {
		if (body.Tag == "rFonts") || (body.Tag == "proofErr") { //这两个标签是word中标记不符合规范的和字体的，不加空格分离
			result = ""
			return
		} else if body.Tag == "t" {
			//fmt.Printf("content[%s:%s]\n",body.Tag,body.Text())
			result = body.Text()
			return
		} else {
			return " "
		}
	} else {
		for _, body1 := range body.ChildElements() {
			r := readElement(body1)
			result = result + r
			if strings.HasSuffix(result, " ") {
				result = strings.Trim(result, " ") + " "
			}
		}
	}
	return
}
