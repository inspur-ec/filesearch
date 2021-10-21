package utils

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

func Open(path string) {
	path = strings.ReplaceAll(path, "/", "\\")
	cmd := `/select,` + path
	exec.Command(`explorer`, cmd).Run()
}

//var files1 map[string][]string
func ListAllFileByName(dirPth string, outpaths []string, suffix ...string) map[string][]string {
	files1 := make(map[string][]string)
	files, _ := ioutil.ReadDir(dirPth)

	for _, onefile := range files {
		if onefile.IsDir() {
			path := dirPth + onefile.Name() + "/"
			oBool := false
			for _, o := range outpaths {
				if path == o {
					oBool = true
					break
				}
			}
			if oBool {
				continue
			}
			fs := ListAllFileByName(path, outpaths, suffix...)
			for suf, filenames := range fs {
				files1[suf] = append(files1[suf], filenames...)
			}
		} else {
			if !strings.HasPrefix(onefile.Name(), "~") {
				for j, _ := range suffix {
					if strings.HasSuffix(strings.ToLower(onefile.Name()), suffix[j]) {
						files1[suffix[j]] = append(files1[suffix[j]], dirPth+onefile.Name())
					}
				}
			}
		}
	}
	return files1

}
