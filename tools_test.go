package main

import (
	"fmt"
	"testing"
)

func TestListAllFileByName(t *testing.T) {
	//files1=make(map[string][]string)
	files1 := ListAllFileByName("e:/", "e:/tmp,e:/outlook/", "xlsx", "docx")
	fmt.Printf("files[%v]", files1)
}
