package js

import (
  "log"
  "github.com/robertkrimen/otto"
  "github.com/syndtr/goleveldb/leveldb"
)

func CreateOtto() *otto.Otto {
    return otto.New();
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