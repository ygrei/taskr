package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/golang/glog"
	"github.com/go-gorp/gorp"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ygrei/taskr/validators"
)

var dbMap *gorp.DbMap

var (
	usernameRE = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")
)

var (
	port   = flag.String("port", "8080", "Port to listen on")
	dbPath = flag.String("db_path", "/tmp/taskr.db", "Database path.")
)

func newMuxRouter() *mux.Router {
	r := mux.NewRouter()
	return r
}

type JSONResponse map[string]interface{}

func SignupPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := SignupPostHandlerWrapper(w, r); err != nil {
		log.Errorf("%v", err)
		log.Flush()
	}

}

func SignupPostHandlerWrapper(w http.ResponseWriter, r *http.Request) error {
	defer log.Flush()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	errors := make(map[string]string)
	validators.RequireMatchRegexp("usernameError", username, &errors, "Username can only contain alphanumeric characters, or underscores", usernameRE)
	validators.RequireString("usernameError", username, &errors, "Username is required")
	validators.MaxLength("usernameError", username, &errors, "Max length is 300 charaters", 300)
	validators.RequireSaneEmail("emailError", email, &errors, "Enter a valid email address")
	validators.RequireString("emailError", email, &errors, "Email is required")
	validators.MaxLength("emailError", email, &errors, "Max length is 128 charaters", 128) /*from phab*/
	validators.MaxLength("emailError", email, &errors, "Max length is 500 charaters", 500)
	validators.RequireString("passwordError", password, &errors, "Password is required")

	if _, err := getUserByHandle(email); err != nil && err != errNoSuchUser {
		return err
	} else if err != errNoSuchUser {
		errors["emailError"] = "Email already taken"
	}
	if _, err := getUserByHandle(username); err != nil && err != errNoSuchUser {
		return err
	} else if err != errNoSuchUser {
		errors["usernameError"] = "Username already taken"
	}

	hashed, err := hashPassword([]byte(password + salt))
	if err != nil {
		return err
	}

	if len(errors) > 0 {
		resp := JSONResponse{"status": "invalid", "errors": errors}
		w.Header().Set("Content-Type", "application/json")
		log.Infof("Response is %v ", resp)
		fmt.Fprint(w, resp)
		return nil
	}

	user := User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
	}
	log.Infof("New user: %v", user)
	if err := dbMap.Insert(&user); err != nil {
		log.Errorf("Couldn't insert user %v: %v", user, err)
		return err
	}
	log.Infof("Successfully inserted into DB: %v", user)

	w.Header().Set("Content-Type", "application/json")
	resp := JSONResponse{"status": "success", "errors": errors}
	fmt.Fprint(w, resp)
	return nil
}

func main() {
	flag.Parse()
	initDB(*dbPath)

	http.HandleFunc("/json/signup", SignupPostHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	p := ":" + *port

	log.Infof("Listening on %v", p)
	http.ListenAndServe(p, nil)
}

func initDB(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	if err := dbMap.CreateTablesIfNotExists(); err != nil {
		return err
	}
	return nil
}

func (jr JSONResponse) String() string {
	b, err := json.Marshal(jr)
	if err != nil {
		return ""
	}
	return string(b)
}
