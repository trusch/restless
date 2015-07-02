package js

import (
	"../auth"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type JSEngine struct {
	mutex sync.Mutex
	otto  *otto.Otto
}

func (ptr *JSEngine) Run(code string) (otto.Value, error) {
	ptr.mutex.Lock()
	val, err := ptr.otto.Run(code)
	ptr.mutex.Unlock()
	return val, err
}

func (ptr *JSEngine) Set(name string, value interface{}) error {
	ptr.mutex.Lock()
	err := ptr.otto.Set(name, value)
	ptr.mutex.Unlock()
	return err
}

func (ptr *JSEngine) Get(name string) (otto.Value, error) {
	ptr.mutex.Lock()
	val, err := ptr.otto.Get(name)
	ptr.mutex.Unlock()
	return val, err
}

func (ptr *JSEngine) Lock() {
	ptr.mutex.Lock()
}

func (ptr *JSEngine) Unlock() {
	ptr.mutex.Unlock()
}

func CreateOtto() *JSEngine {
	engine := new(JSEngine)
	engine.otto = otto.New()
	return engine
}

func headerToJSON(headers http.Header) string {
	b, err := json.Marshal(headers)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(b)
}

func InjectRequestDetails(jsEngine *JSEngine, w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(w, r)
	bodydata, _ := ioutil.ReadAll(r.Body)

	code := fmt.Sprintf(`
      var URL = '%v';
      var METHOD = '%v'
      var HEADERS = JSON.parse('%v');
      var USERNAME = '%v';
      var ROLE = '%v'
      var BODY = '%v'
    `, r.URL.Path, r.Method, headerToJSON(r.Header), user.Username, user.Role, string(bodydata))

	jsEngine.Run(code)
}

func InjectLevelDB(jsEngine *JSEngine, db *leveldb.DB) {
	jsEngine.Run("var db = {};")
	dbValue, _ := jsEngine.Get("db")
	dbObj := dbValue.Object()
	dbObj.Set("put", func(call otto.FunctionCall) otto.Value {
		key, err := call.Argument(0).ToString()
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		value, err := call.Argument(1).ToString()
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		err = db.Put([]byte(key), []byte(value), nil)
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		return otto.TrueValue()
	})
	dbObj.Set("get", func(call otto.FunctionCall) otto.Value {
		key, err := call.Argument(0).ToString()
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		data, err := db.Get([]byte(key), nil)
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		v, _ := otto.ToValue(string(data))
		return v
	})
	dbObj.Set("remove", func(call otto.FunctionCall) otto.Value {
		key, err := call.Argument(0).ToString()
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		err = db.Delete([]byte(key), nil)
		if err != nil {
			log.Println("Error:", err.Error())
			return otto.FalseValue()
		}
		return otto.TrueValue()
	})
}
