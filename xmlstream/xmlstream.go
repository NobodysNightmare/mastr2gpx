// Originally taken from https://github.com/srom/xmlstream and then modified
// to work properly with the UTF-16 input files created by the Marktstammdatenregister
// Original license follows (MIT):

// Copyright (C) 2014-2015 Romain Strock
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software
// and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
// LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN
// NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package xmlstream implements a lightweight XML scanner on top of encoding/xml.
// It keeps the flexibility of xml.Unmarshal while allowing the parsing of huge XML files.
package xmlstream

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
)

// Scanner provides a way to read a stream of XML data. It uses an xml.Decoder internally to step
// through the XML elements of the stream.
type Scanner struct {
	decoder    *xml.Decoder
	element    interface{}
	nameToType map[string]reflect.Type // map xml local name to element's type
	err        error
}

// NewScanner returns a new Scanner to read from r.
// Tags must be struct objects or pointer to struct objects, as defined by encoding/xml:
// http://golang.org/pkg/encoding/xml/#Unmarshal
func NewScanner(r io.Reader, tags ...interface{}) *Scanner {
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		return r, nil
	}

	s := Scanner{
		decoder:    decoder,
		nameToType: make(map[string]reflect.Type, len(tags)),
	}

	// Map the xml local name of an element to its underlying type.
	for _, tag := range tags {
		v := reflect.ValueOf(tag)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		t := v.Type()
		name := elementName(v)
		s.nameToType[name] = t
	}
	return &s
}

func elementName(v reflect.Value) string {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("Tags must be of kind Struct but got %s", t.Kind()))
	}
	name := t.Name()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "XMLName" || field.Type.String() == "xml.Name" {
			if field.Tag.Get("xml") != "" {
				name = field.Tag.Get("xml")
			}
		}
	}
	return name
}

// Scan advances the Scanner to the next XML element matching one of the struct passed to NewReader.
// This element will then be available through the Element method.
// It returns false when the scan stops, either by reaching the end of the input or an error.
// After Scan returns false, the Err method will return any error that occurred
// during scanning, except that if it was io.EOF, Err will return nil.
func (s *Scanner) Scan() bool {
	if (*s).err != nil {
		return false
	}
	for {
		// Read next token.
		token, err := (*s).decoder.Token()
		if err != nil {
			(*s).element = nil
			(*s).err = err
			return false
		}
		// Inspect the type of the token.
		switch el := token.(type) {
		case xml.StartElement:
			// Read the element name and compare with the XML element.
			if elementType, ok := (*s).nameToType[el.Name.Local]; ok {
				// create a new element
				element := reflect.New(elementType).Interface()
				// Decode a whole chunk of following XML.
				err := (*s).decoder.DecodeElement(element, &el)
				(*s).element = element
				(*s).err = err
				return err == nil
			}
		}
	}
	return false
}

// Element returns a pointer to the most recent struct object generated by a call to Scan.
// The type of this struct matches the type of one of the custom struct passed to NewReader.
func (s *Scanner) Element() interface{} {
	return (*s).element
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if (*s).err != nil && (*s).err != io.EOF {
		return (*s).err
	}
	return nil
}
