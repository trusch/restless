package fs

import (
    "github.com/robertkrimen/otto"
    "../js"
    "io/ioutil"
    "log"
)

func InjectIntoOtto(jsEngine *js.JSEngine){
      jsEngine.Run("var fs = {};")
      fsValue,_ := jsEngine.Get("fs")
      fsObj := fsValue.Object()
      fsObj.Set("readFile", func(call otto.FunctionCall) otto.Value {
          filename,err := call.Argument(0).ToString()
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          bs,err := ioutil.ReadFile(filename)
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          val,err := otto.ToValue(string(bs))
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          return val
      })

      fsObj.Set("writeFile", func(call otto.FunctionCall) otto.Value {
          filename,err := call.Argument(0).ToString()
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          data,err := call.Argument(1).ToString()
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          err = ioutil.WriteFile(filename,[]byte(data),0777)
          if err!=nil {
            log.Println("Error:",err.Error())
            return otto.FalseValue()
          }
          return otto.TrueValue()
      })
}