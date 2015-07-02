package main

import (
	"github.com/gorilla/mux"
	"github.com/robertkrimen/otto"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"log"
	"net/http"

	"./api"
	"./auth"
	"./config"
	"./js"
	"./dynamic"
	"./fs"
)

var jsEngine *otto.Otto = nil
var db *leveldb.DB = nil

func SetUpLevelDB(path string) {
	db_, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Fatal("Leveldb: " + err.Error())
	}
	db = db_
}

func SetUpOtto(codeFile string) {
	jsEngine = js.CreateOtto()
	js.InjectLevelDB(jsEngine, db)
	fs.InjectIntoOtto(jsEngine)
	backendCode, e := ioutil.ReadFile(codeFile)
	if e != nil {
		log.Fatal("Need " + codeFile)
	}
	_,e= jsEngine.Run(backendCode)
	if e != nil {
		log.Fatal("Error in backendjs:" + e.Error())
	}
}

func SetUpAPI(endpoints []config.Endpoint) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	for _, endpoint := range endpoints {
		r := router.PathPrefix("/api" + endpoint.Url)
		api.BuildEndpoint(r, endpoint.Model, jsEngine, db)
	}
	return router
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cfg := config.Load("./config.json")

	SetUpLevelDB(cfg.DB)
	SetUpOtto(cfg.JSBase)

	authRouter := auth.SetupAuth("passwords.gob")
	err := auth.AddUser("root", "root@localhost", "toor", "admin")
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Handle("/auth/", authRouter)

	apiRouter := SetUpAPI(cfg.Endpoints)
	http.Handle("/api/", apiRouter)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(cfg.Assets))))
	http.HandleFunc("/",dynamic.BuildHandler(jsEngine))

	if(cfg.CertFile != "" && cfg.KeyFile != ""){
		log.Println("starting TLS secured server on", cfg.Address)
		log.Fatal(http.ListenAndServeTLS(cfg.Address, cfg.CertFile, cfg.KeyFile, nil))	
	}else{
		log.Println("starting server on", cfg.Address)
		log.Fatal(http.ListenAndServe(cfg.Address, nil))
	}
}
