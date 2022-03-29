// Package properties is used to read or write or modify the properties document.
package properties

import (
	"fmt"
	"io"
	"strings"
)

type element struct {
	//  #   注释行
	//  !   注释行
	//  ' ' 空白行或者空行
	//  =   等号分隔的属性行
	//  :   冒号分隔的属性行
	typo  byte   //  行类型
	value string //  行的内容,如果是注释注释引导符也包含在内
	key   string //  如果是属性行这里表示属性的key
}

// Document The properties document in memory.
type Document struct {
	props map[string]*element
}

// New is used to create a new and empty properties document.
//
// It's used to generate a new document.
func New() *Document {
	doc := new(Document)
	doc.props = make(map[string]*element)
	return doc
}

// Save is used to save the doc to file or stream.
func Save(doc *Document, writer io.Writer) error {
	var err error
	doc.Accept(func(typo byte, value string, key string) bool {
		switch typo {
		case '=', ':':
			_, err = fmt.Fprintf(writer, "%s%c%s\n", escapeKey(key), typo, escapeValue(value))
		}

		return nil == err
	})

	return err
}

// Get Retrieve the value from Document.
//
// If the item is not exist, the exist is false.
func (p *Document) Get(key string) (value string, exist bool) {
	e, ok := p.props[key]
	if !ok {
		return "", ok
	}

	return e.value, ok
}

// Set Update the value of the item of the key.
//
// Create a new item if the item of the key is not exist.
func (p *Document) Set(key string, value string) {
	e, ok := p.props[key]
	if !ok {
		p.props[key] = &element{typo: '=', key: key, value: value}
		return
	}

	e.value = value
	return
}

// Accept Traverse every element of the document, include comment.
//
// The typo parameter special the element type.
// If typo is '#' or '!' means current element is a comment.
// If typo is ' ' means current element is a empty or a space line.
// If typo is '=' or ':' means current element is a key-value pair.
// The traverse will be terminated if f return false.
func (p *Document) Accept(f func(typo byte, value string, key string) bool) {
	for _, elem := range p.props {
		continues := f(elem.typo, elem.value, elem.key)
		if !continues {
			return
		}
	}
}

func escapeKey(value string) string {
	replacer := strings.NewReplacer("=", "\\=", ":", "\\:", " ", "\\ ", "\\", "\\\\")
	return replacer.Replace(value)
}

func escapeValue(value string) string {
	replacer := strings.NewReplacer(" ", "\\ ", "\\", "\\\\", "\n", "\\n")
	return replacer.Replace(value)
}
