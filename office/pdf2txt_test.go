package office

import (
	"fmt"
	"testing"
)

func TestWritepdf(t *testing.T) {
	content, err := ReadPdf("e:\\云计算发展白皮书2019 - 副本.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
