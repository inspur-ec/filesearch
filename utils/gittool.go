package utils

import (
	"bytes"
	"compress/zlib"
	"flag"
	log "github.com/golang/glog"
	"io"
	"io/ioutil"
)

func Decompression(packPath string) {
	flag.Parse()
	defer log.Flush()
	//packPath := "/Users/usr/test_pack_no_header"
	b, err := ioutil.ReadFile(packPath)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(b)

	i := 0
	for {
		decoded, err := readObjectRaw(bytes.NewBuffer(buf.Bytes()))
		if err == nil {
			log.Errorf("i %v, decoded is %s", i, decoded)
			//return
		}
		if _, err = buf.ReadByte(); err == io.EOF {
			log.Errorf("EOF")
			return
		}
		i++
	}

}

func readObjectRaw(reader io.Reader) ([]byte, error) {
	r, err := zlib.NewReader(reader)
	if err != nil {
		//log.Error(err)
		return nil, err
	}
	defer r.Close()
	def := bytes.NewBuffer(nil)
	_, err = io.Copy(def, r)
	if err != nil {
		//log.Errorf("%v bytes copied, err %v", i, err)
		return nil, err
	}
	return def.Bytes(), nil
}
