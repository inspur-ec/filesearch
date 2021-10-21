package office

import (
	"github.com/unidoc/unioffice/presentation"
	"log"
)

func ReadPpt() {
	doc, err := presentation.Open("e:/基于CPS的智能工厂管理平台技术交流-0625.pptx")
	if err != nil {
		log.Fatalf("error opening document: %s", err)
	}
	doc.ExtractText()
	/**
	//doc.Paragraphs()得到包含文档所有的段落的切片
	for i, para := range doc.Paragraphs() {
		//run为每个段落相同格式的文字组成的片段
		fmt.Println("-----------第", i, "段-------------")
		for j, run := range para.Runs() {
			fmt.Print("\t-----------第", j, "格式片段-------------")
			fmt.Print(run.Text())

		}
		fmt.Println()
	}
	 */
}
