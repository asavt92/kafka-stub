package json

import (
	"github.com/xeipuuv/gojsonschema"
)

type JsonConfig struct {
	
}

type JsonService struct {
	JsonConfig JsonConfig
	ValidationSchema gojsonschema.JSONLoader
	ResponseExamples []gojsonschema.JSONLoader
}



func NewJsonService(c *JsonConfig) *JsonService {


}