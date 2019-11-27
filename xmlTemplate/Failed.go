package xmlTemplate

import "encoding/xml"

type Failed struct {
	Root xml.Name `xml:"Root"`

	ResultFlag string `xml:"resultFlag"`
	FailCode   string `xml:"failCode"`
	FailInfo   string `xml:"failInfo"`
	RetData    string `xml:"retData"`
}
