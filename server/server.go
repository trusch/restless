package main

import (
  "log"
  "github.com/gorilla/mux"
  "net/http"
  "github.com/robertkrimen/otto"
  "github.com/syndtr/goleveldb/leveldb"
  "io/ioutil"

  "./config"
  "./api"
)

var jsEngine *otto.Otto = nil
var db *leveldb.DB = nil

func SetUpLeveldb(path string){
  db_,err := leveldb.OpenFile(path,nil);
  if err!=nil {
    log.Fatal("Leveldb: "+err.Error())
  }
  db = db_
}

func SetUpOtto(){
  jsEngine = otto.New();
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

func main(){
  log.SetFlags(log.Lshortfile | log.LstdFlags)

  cfg := config.Load("./config.json")

  SetUpLeveldb(cfg.DB);
  SetUpOtto();

  backendCode, e := ioutil.ReadFile(cfg.JSBase)
  if e!=nil {
    log.Fatal("Need "+cfg.JSBase)
  }

  jsEngine.Run(backendCode)

  router := mux.NewRouter()
  router.StrictSlash(true)
  for _,endpoint := range cfg.Endpoints {
    r := router.PathPrefix("/api"+endpoint.Url)
    api.BuildEndpoint(r, endpoint.Model, jsEngine, db)
  }

  http.Handle("/api/", router)
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(cfg.Assets))))
  http.Handle("/", http.RedirectHandler("/assets/index.html",301))
  
  log.Println("starting server on",cfg.Address)  
  log.Fatal(http.ListenAndServe(cfg.Address, nil))
}
