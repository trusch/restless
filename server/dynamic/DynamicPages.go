package dynamic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"log"
	"net/http"

	"../js"
)

func headerToJSON(headers http.Header) string {
	b, err := json.Marshal(headers)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(b)
}

func objToHeader(data *otto.Value, header http.Header) {
	dataMap, _ := data.Export()
	for key, val := range dataMap.(map[string]interface{}) {
		switch val.(type) {
		case string:
			{
				header.Add(key, val.(string))
				break
			}
		case []interface{}:
			{
				for _, headerVal := range val.([]interface{}) {
					header.Add(key, headerVal.(string))
				}
				break
			}
		}
	}
}

func BuildHandler(jsEngine *js.JSEngine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errHandler := func(err error) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error in js: %v", err.Error())))
		}

		js.InjectRequestDetails(jsEngine, w, r)
		code := "restless.renderPage();"
		val, err := jsEngine.Run(code)
		if err != nil {
			errHandler(err)
			return
		}
		obj := val.Object()
		if obj == nil {
			errHandler(errors.New("restless.renderPage() did not return an object"))
			return
		}
		bodyVal, err := obj.Get("body")
		if err != nil {
			errHandler(err)
			return
		}
		body, err := bodyVal.ToString()
		if err != nil {
			errHandler(err)
			return
		}
		statusVal, err := obj.Get("status")
		if err != nil {
			errHandler(err)
			return
		}
		status, err := statusVal.ToInteger()
		if err != nil {
			errHandler(err)
			return
		}
		headerVal, err := obj.Get("header")
		if err != nil {
			errHandler(err)
			return
		}
		objToHeader(&headerVal, w.Header())

		w.WriteHeader(int(status))
		w.Write([]byte(body))
	}
}
