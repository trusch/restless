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
  "./js"
)

var jsEngine *otto.Otto = nil
var db *leveldb.DB = nil

func SetUpLevelDB(path string){
  db_,err := leveldb.OpenFile(path,nil);
  if err!=nil {
    log.Fatal("Leveldb: "+err.Error())
  }
  db = db_
}

func SetUpOtto(codeFile string){
  jsEngine = js.CreateOtto();
  js.InjectLevelDB(jsEngine,db);
  backendCode, e := ioutil.ReadFile(codeFile)
  if e!=nil {
    log.Fatal("Need "+codeFile)
  }
  jsEngine.Run(backendCode)
}

func SetUpAPI(endpoints []config.Endpoint) *mux.Router {
  router := mux.NewRouter()
  router.StrictSlash(true)
  for _,endpoint := range endpoints {
    r := router.PathPrefix("/api"+endpoint.Url)
    api.BuildEndpoint(r, endpoint.Model, jsEngine, db)
  }
  return router  
}

func main(){
  log.SetFlags(log.Lshortfile | log.LstdFlags)

  cfg := config.Load("./config.json")

  SetUpLevelDB(cfg.DB);
  SetUpOtto(cfg.JSBase);

  router := SetUpAPI(cfg.Endpoints);

  http.Handle("/api/", router)
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(cfg.Assets))))
  http.Handle("/", http.RedirectHandler("/assets/index.html",301))
  
  log.Println("starting server on",cfg.Address)  
  log.Fatal(http.ListenAndServe(cfg.Address, nil))
}
