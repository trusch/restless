package js

import (
  "log"
  "github.com/robertkrimen/otto"
  "github.com/syndtr/goleveldb/leveldb"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "fmt"
  "../auth"
)

func CreateOtto() *otto.Otto {
    return otto.New();
}

func headerToJSON(headers http.Header) string {
  b,err := json.Marshal(headers);
  if err!=nil {
    log.Println(err.Error())
    return "";
  }
  return string(b)
}

func InjectRequestDetails(jsEngine *otto.Otto, w http.ResponseWriter, r *http.Request){
    user := auth.GetUser(w,r)
    bodydata,_ := ioutil.ReadAll(r.Body)

    code := fmt.Sprintf(`
      var URL = '%v';
      var METHOD = '%v'
      var HEADERS = JSON.parse('%v');
      var USERNAME = '%v';
      var ROLE = '%v'
      var BODY = '%v'
    `,r.URL.Path,r.Method,headerToJSON(r.Header),user.Username,user.Role,string(bodydata))

    jsEngine.Run(code)
}

func InjectLevelDB(jsEngine *otto.Otto, db *leveldb.DB){
  jsEngine.Run("var db = {};")
  dbValue,_ := jsEngine.Get("db")
  dbObj := dbValue.Object()
  dbObj.Set("put", func(call otto.FunctionCall) otto.Value {
      key,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      value,err := call.Argument(1).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      err = db.Put([]byte(key),[]byte(value),nil)
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      return otto.TrueValue()
  })
  dbObj.Set("get", func(call otto.FunctionCall) otto.Value {
      key,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      data, err := db.Get([]byte(key),nil)
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      v,_ := otto.ToValue(string(data))
      return v
  })
  dbObj.Set("remove", func(call otto.FunctionCall) otto.Value {
      key,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      err = db.Delete([]byte(key),nil)
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      return otto.TrueValue()
  })
}