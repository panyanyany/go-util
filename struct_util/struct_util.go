package struct_util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func Map2Struct(input interface{}, output interface{}) (err error) {
	var bs []byte
	bs, err = json.Marshal(input)
	if err != nil {
		err = fmt.Errorf("Map2Struct-json.Marshal: %w", err)
		return
	}
	err = json.Unmarshal(bs, output)
	if err != nil {
		err = fmt.Errorf("Map2Struct-json.Unmarshal: %w", err)
		return
	}

	return err
}

func XmlToMap(xmlData []byte) map[string]interface{} {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	m := make(map[string]interface{})
	var token xml.Token
	var err error
	var k string
	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {
		if v, ok := token.(xml.StartElement); ok {
			k = v.Name.Local
			continue
		}
		if v, ok := token.(xml.CharData); ok {
			data := string(v.Copy())
			if strings.TrimSpace(data) == "" {
				continue
			}
			m[k] = data
		}
	}

	if err != nil && err != io.EOF {
		panic(err)
	}
	return m
}
