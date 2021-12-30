package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func return400OnError(w http.ResponseWriter, err error, message string) bool {

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, message)
		return true
	}

	return false

}

func interfaceIsString(variable interface{}) bool {

	return reflect.TypeOf(variable).Kind() == reflect.String

}

func interfaceIsInt(variable interface{}) bool {

	return reflect.TypeOf(variable).Kind() == reflect.Int

}

func interfaceIsFloat64(variable interface{}) bool {

	return reflect.TypeOf(variable).Kind() == reflect.Float64

}

func mapValueNotNil(mapObject map[string]interface{}, mapKey string) (interface{}, bool) {

	if mapObject[mapKey] == nil {

		return "", false

	} else {

		return mapObject[mapKey], true

	}

}

func parseMapToJson(mapObject map[string]interface{}) (string, bool) {

	jsonString, err := json.Marshal(mapObject)

	if err == nil {

		return string(jsonString), true

	} else {

		return "{}", false

	}

}

func parseJsonToMap(bytes []byte, mapObject *[]map[string]interface{}) bool {

	if json.Valid(bytes) {

		err := json.Unmarshal(bytes, mapObject)
		return err == nil

	} else {

		return false

	}

}
