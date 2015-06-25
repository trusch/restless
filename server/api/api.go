package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robertkrimen/otto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func buildGetOneHandler(modelName string, jsEngine *otto.Otto) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		code := fmt.Sprintf(`
      var instance = new %v();
      instance.initFromUID(%v);
      JSON.stringify(instance.__data);
    `, modelName, id)
		val, err := jsEngine.Run(code)
		if err != nil {
			log.Println("Error in js:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error in js: %v", err.Error())))
			return
		}
		str := val.String()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(str))
	}
}

func buildPutOneHandler(modelName string, jsEngine *otto.Otto) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("in put one")
		vars := mux.Vars(r)
		id := vars["id"]
		data, _ := ioutil.ReadAll(r.Body)
		code := fmt.Sprintf(`
      var instance = new %v();
      instance.initFromData(JSON.parse('%v'));
      instance.__uid = %v;
      instance.commit();
    `, modelName, string(data), id)
		val, err := jsEngine.Run(code)
		if err != nil {
			log.Println("Error in js:", err.Error())
			log.Println(code)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error in js: %v", err.Error())))
			return
		}
		str := val.String()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(str))
	}
}

func buildDeleteOneHandler(modelName string, jsEngine *otto.Otto) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("in delete one")
		vars := mux.Vars(r)
		id := vars["id"]
		code := fmt.Sprintf(`
      var instance = new %v();
      instance.initFromUID(%v);
      if(!instance.__data){
        false;
      }else{
        instance.remove();
      }
    `, modelName, id)
		val, err := jsEngine.Run(code)
		if err != nil {
			log.Println("Error in js:", err.Error())
			log.Println(code)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error in js: %v", err.Error())))
			return
		}
		str := val.String()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(str))
	}
}

func buildGetAllHandler(modelName string, db *leveldb.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("in get all")
		keys := []string{}
		prefix := strings.ToLower(modelName) + ":"
		iter := db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
		for iter.Next() {
			keys = append(keys, strings.TrimPrefix(string(iter.Key()), prefix))
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		bs, err := json.Marshal(keys)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error in db: %v", err.Error())))
			return
		}
		w.Write(bs)
	}
}

func BuildEndpoint(route *mux.Route, modelName string, jsEngine *otto.Otto, db *leveldb.DB) {
	router := route.Subrouter()
	router.HandleFunc("/{id:[0-9]+}", buildGetOneHandler(modelName, jsEngine)).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", buildPutOneHandler(modelName, jsEngine)).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}", buildDeleteOneHandler(modelName, jsEngine)).Methods("DELETE")
	router.HandleFunc("/", buildGetAllHandler(modelName, db)).Methods("GET")
}
