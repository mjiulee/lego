package utils

import (
	"encoding/xml"
)

// StringMap is a map[string]string.
type StringMap map[string]string

// StringMap marshals into XML.
func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	start.Name.Local = "xml"
	tokens := []xml.Token{start}
	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func Map2Xml(smap StringMap) string {
	if x, err := xml.MarshalIndent(smap, "", " "); err == nil {
		return string(x)
	}
	return ""
}
