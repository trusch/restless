package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apexskier/httpauth"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	backend     httpauth.GobFileAuthBackend
	aaa         httpauth.Authorizer
	roles       map[string]httpauth.Role
	backendfile = "auth.gob"
)

func SetupAuth(pwfile string) *mux.Router {

	backendfile = pwfile
	os.Create(backendfile)
	backend_, err := httpauth.NewGobFileAuthBackend(backendfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	backend = backend_

	roles = make(map[string]httpauth.Role)
	roles["user"] = 30
	roles["admin"] = 80

	aaa, err = httpauth.NewAuthorizer(backend, []byte("cookie-encryption-key"), "user", roles)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := mux.NewRouter()
	route := router.PathPrefix("/auth")
	r := route.Subrouter().StrictSlash(true)
	r.HandleFunc("/login", handleLogin).Methods("POST")
	r.HandleFunc("/logout", handleLogout)

	return router
}

func AddUser(name, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := httpauth.UserData{Username: name, Email: email, Hash: hash, Role: role}
	err = backend.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(rw http.ResponseWriter, req *http.Request) (user httpauth.UserData, e error) {
	return aaa.CurrentUser(rw, req)
}

func DeleteUser(name string) error {
	return aaa.DeleteUser(name)
}

func handleLogin(rw http.ResponseWriter, req *http.Request) {
	log.Println("login!")
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	if err := aaa.Login(rw, req, username, password, "/"); err != nil && err.Error() == "already authenticated" {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	} else if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(err.Error()))
	}
}

func handleLogout(rw http.ResponseWriter, req *http.Request) {
	log.Println("logout!")
	if err := aaa.Logout(rw, req); err != nil {
		fmt.Println(err)
		// this shouldn't happen
		return
	}
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}
