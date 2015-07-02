package events

import (
  "github.com/robertkrimen/otto"
  "../js"
  "log"
)

func InjectIntoOtto(jsEngine *js.JSEngine, manager *EventManager){
    jsEngine.Run("var events = {};")
    eventsValue,_ := jsEngine.Get("events")
    eventsObj := eventsValue.Object()
    
    eventsObj.Set("on", func(call otto.FunctionCall) otto.Value {
      topic,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      fn := call.Argument(1)
      id := manager.On(topic,func(topic string, payload interface{}){
        t,_ := otto.ToValue(topic)
        p,_ := otto.ToValue(payload)
        _,err := fn.Call(fn,t,p)
        if err!=nil {
            log.Println("Error in JS callback:",err.Error())
        }
      })
      i,_ := otto.ToValue(id)
      return i
    })
    
    eventsObj.Set("once", func(call otto.FunctionCall) otto.Value {
      topic,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      fn := call.Argument(1)
      manager.Once(topic,func(topic string, payload interface{}){
        t,_ := otto.ToValue(topic)
        p,_ := otto.ToValue(payload)
        _,err := fn.Call(fn,t,p)
        if err!=nil {
            log.Println("Error in JS callback:",err.Error())
        }
      })
      return otto.TrueValue()
    })
    
    eventsObj.Set("many", func(call otto.FunctionCall) otto.Value {
      topic,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      count,err := call.Argument(1).ToInteger()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      fn := call.Argument(2)
      manager.Many(topic,int(count),func(topic string, payload interface{}){
        t,_ := otto.ToValue(topic)
        p,_ := otto.ToValue(payload)
        _,err := fn.Call(fn,t,p)
        if err!=nil {
            log.Println("Error in JS callback:",err.Error())
        }
      })
      return otto.TrueValue()
    })
    
    eventsObj.Set("off", func(call otto.FunctionCall) otto.Value {
      id,err := call.Argument(0).ToInteger()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      res := manager.Off(id)
      r,_ := otto.ToValue(res)
      return r
    })

    eventsObj.Set("emit", func(call otto.FunctionCall) otto.Value {
      topic,err := call.Argument(0).ToString()
      if err!=nil {
        log.Println("Error:",err.Error())
        return otto.FalseValue()
      }
      payload := call.Argument(1)
      manager.Emit(topic,payload)
      return otto.TrueValue()
    })
}