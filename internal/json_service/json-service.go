package json_service

import (
	"encoding/json"
	"errors"
	"github.com/asavt92/kafka-stub/internal/configs"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"math/rand"
	"path/filepath"
)

const (
	ALL_ALLOW_JSON_SCHEMA = "{}"
)

type JsonService struct {
	JsonConfig       configs.JsonConfig
	ValidationSchema gojsonschema.JSONLoader
	ResponseExamples []map[string]interface{}
}

func NewJsonService(c *configs.JsonConfig) *JsonService {

	validationSchema := loadJsonSchema(c.InJsonSchemaPath)
	responseExamples := loadJsonFiles(filepath.Join(c.OutJsonExamplesDirPath, "*.json"))
	return &JsonService{
		JsonConfig:       *c,
		ResponseExamples: responseExamples,
		ValidationSchema: validationSchema}
}

func loadJsonSchema(file string) gojsonschema.JSONLoader {
	if file == "" {
		return gojsonschema.NewStringLoader(ALL_ALLOW_JSON_SCHEMA)
	}

	s, err := readJsonFile(file)
	if err != nil {
		log.Error("Error loading json schema", err)
		return gojsonschema.NewStringLoader(ALL_ALLOW_JSON_SCHEMA)
	}
	return gojsonschema.NewStringLoader(s)
}

func readJsonFile(file string) (string, error) {
	j, err := ioutil.ReadFile(file)
	return string(j), err
}

func loadJsonFiles(pattern string) []map[string]interface{} {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	var list []map[string]interface{}

	if len(files) == 0 {
		log.Fatal("No one json files match pattern: ", pattern)
	}

	for _, file := range files {
		j, err := readJsonFile(file)
		if err != nil {
			log.Warnf("error while loading file %s", file)
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(j), &raw); err != nil {
			panic(err)
		}

		list = append(list, raw)
		log.Infof("loaded json file: %s", file)
		log.Debug(raw)
	}

	return list
}

func (service *JsonService) IsJsonStringValid(jsonString string) bool {
	result, err := gojsonschema.Validate(service.ValidationSchema, gojsonschema.NewStringLoader(jsonString))
	if err != nil {
		log.Error("Error validating", err, jsonString)
	}
	if result.Valid() {
		return true
	} else {
		log.Warnf("The document %s is not valid. see errors ", jsonString)
		for _, desc := range result.Errors() {
			log.Warnf("- %s", desc)
		}
		return false
	}
}

func (service *JsonService) GetRandomJsonResponseString() map[string]interface{} {
	return service.ResponseExamples[rand.Intn(len(service.ResponseExamples))]
}

func (service *JsonService) GetJsonResponseStringByConfiguredFieldMapping(value string) (map[string]interface{}, error) {
	//keys := reflect.ValueOf(service.JsonConfig.IdFieldMapping).MapKeys()
	//key := service.JsonConfig.IdFieldMapping[keys[0].String()]
	key := service.JsonConfig.IdFieldMapping["to"]
	for _, res := range service.ResponseExamples {
		if res[key] == value {
			return res, nil
		}
	}
	return nil, errors.New("There is no one response example with this mapping " + key + " = " + value)
}
