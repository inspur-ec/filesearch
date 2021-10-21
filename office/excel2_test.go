package office

import (
	"fmt"
	"testing"
)

func TestReadExcel(t *testing.T) {
	c, _ := ReadExcel("e:/区块链/质量链/阿胶/附件1企业及产品上链-数据模板-东阿阿胶.xlsx")
	fmt.Printf("content is [%s]", c)
}
