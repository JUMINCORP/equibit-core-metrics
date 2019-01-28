package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

type responseIface interface {
	append(map[string]interface{})
	export()
	setCommandName(string)
}

type response struct {
	command string
	data    []map[string]interface{}
}

func (r *response) setCommandName(command string) {
	r.command = command
}

func (r *response) append(obj map[string]interface{}) {
	//log.Printf(":: OBJ :: %v\n", obj)
	r.data = append(r.data, obj)
}

func sanitizeName(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "%", "_Percent", -1)
	return name
}

func (r *response) metricName(data string) string {
	return strings.Join([]string{"bfx", r.command, sanitizeName(data)}, "_")
}

func newResponse(command string, responseBytes []byte) responseIface {

	var responseData = make(map[string]interface{})
	err := json.Unmarshal(bytes.Trim(responseBytes, "\x00"), &responseData)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("JSON0: %v\n\n", data)

	delete(responseData, "STATUS")
	delete(responseData, "id")

	log.Printf("JSON: %v\n\n", responseData)

	// By now we should have only one element
	if len(responseData) != 1 {
		log.Printf("len(responseData) != 1 : %v\n", len(responseData))
		return nil
	}

	var r responseIface

	switch command {
	case "devs":
		r = newDevsResponse()
	}

	r.setCommandName(command)

	var respList []interface{}
	for _, wrapped := range responseData {
		respList = wrapped.([]interface{})
	}
	//log.Printf("respList: %v\n\n", respList)

	for _, respElement := range respList {
		r.append(respElement.(map[string]interface{}))
	}

	return r
}

func (r *response) export() {
	for _, element := range r.data {
		log.Printf("---")
		for name, val := range element {
			metricName := r.metricName(name)
			switch casted := val.(type) {
			case int64:
				log.Printf("i! %s %d\n", name, casted)
			case float64:
				log.Printf("f  %s >> %s %f\n", name, metricName, casted)
			case string:
				log.Printf("s! %s %s\n", name, casted)
			default:
				log.Printf("?! %s %s\n", name, casted)
			}
		}
	}
}

type devsResponse struct {
	response
}

func newDevsResponse() *devsResponse {
	r := new(devsResponse)
	//r.response.data = make([]map[string]interface{}, 10)

	return r
}
