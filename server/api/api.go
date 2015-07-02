package api

import (
	"../js"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func buildGetOneHandler(modelName string, jsEngine *js.JSEngine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		js.InjectRequestDetails(jsEngine, w, r)
		code := fmt.Sprintf(`
      var instance = restless.CreateModel('%v');
      instance.getFromUID('%v');
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

func buildPutOneHandler(modelName string, jsEngine *js.JSEngine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		data, _ := ioutil.ReadAll(r.Body)
		js.InjectRequestDetails(jsEngine, w, r)
		code := fmt.Sprintf(`
      var instance = restless.CreateModel('%v');
      instance.initFromData(JSON.parse('%v'));
      instance.__uid = '%v';
      instance.put();
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

func buildPostHandler(modelName string, jsEngine *js.JSEngine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		js.InjectRequestDetails(jsEngine, w, r)
		code := fmt.Sprintf(`
      var instance = restless.CreateModel('%v');
      instance.initFromData(JSON.parse('%v'));
      instance.put();
      instance.__uid;
    `, modelName, string(data))
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

func buildDeleteOneHandler(modelName string, jsEngine *js.JSEngine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		js.InjectRequestDetails(jsEngine, w, r)
		code := fmt.Sprintf(`
      var instance = restless.CreateModel('%v');
      instance.getFromUID('%v');
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
		keys := []string{}
		prefix := modelName + ":"
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

func BuildEndpoint(route *mux.Route, modelName string, jsEngine *js.JSEngine, db *leveldb.DB) {
	router := route.Subrouter()
	router.StrictSlash(true)
	router.HandleFunc("/{id:.+}", buildGetOneHandler(modelName, jsEngine)).Methods("GET")
	router.HandleFunc("/{id:.+}", buildPutOneHandler(modelName, jsEngine)).Methods("PUT")
	router.HandleFunc("/{id:.+}", buildDeleteOneHandler(modelName, jsEngine)).Methods("DELETE")
	router.HandleFunc("/", buildGetAllHandler(modelName, db)).Methods("GET")
	router.HandleFunc("/", buildPostHandler(modelName, jsEngine)).Methods("POST")
}
