package main

import (
	"net/http"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"errors"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-gorp/gorp"
	log "github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ygrei/taskr/validators"
)

var dbMap *gorp.DbMap

const (
	salt = "here be dragons"
)

var (
	port = flag.String("port", "8080", "Port to listen on")
	dbPath           = flag.String("db_path", "/tmp/taskr.db", "Database path.")
)

var (
	errNoSuchUser = errors.New("no such user")
)

var (
	usernameRE = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")
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
		errors["email"] = "Email already taken"
	}
	if _, err := getUserByHandle(username); err != nil && err != errNoSuchUser {
		return err
	} else if err != errNoSuchUser {
		errors["username"] = "Username already taken"
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

func getUserByHandle(handle string) (*User, error) {
	var userID *int64
	if strings.Contains(handle, "@") {
		err := dbMap.SelectOne(&userID, "SELECT ID FROM users WHERE email=?", handle)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("unexpected DB error in getUserByHandle(%v): %v", handle, err)
		}
	} else {
		err := dbMap.SelectOne(&userID, "SELECT ID FROM users WHERE username=?", handle)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("unexpected DB error in getUser(%v): %v", handle, err)
		}
	}
	if userID == nil {
		return nil, errNoSuchUser
	}
	return getUser(*userID)
}

func getUser(id int64) (*User, error) {
	var user *User
	err := dbMap.SelectOne(&user, "SELECT * FROM users WHERE id=?", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("unexpected DB error in getUser(%v): %v", id, err)
	}
	if user == nil {
		return nil, errNoSuchUser
	}
	user.Expand()
	return user, nil
}

func (u *User) Expand() {
	if u == nil || u.TechsRaw == "" {
		return
	}
	techsSlice := strings.Split(u.TechsRaw, ",")
	for i, _ := range techsSlice {
		techsSlice[i] = strings.TrimSpace(techsSlice[i])
	}
	u.Techs = techsSlice
	return
}

type User struct {
	ID             int64
	Username       string
	RealName       string
	GitHubUsername string
	GitHubToken    string
	Location       string
	Techs          []string `db:"-"` // gorp does not support slices :(, see https://github.com/coopernurse/gorp/issues/5
	TechsRaw       string
	Email          string
	PasswordHash   string
	Bio            string
	PastProjects   string

	// Disbursement options
	DisburseStyle string
	PayPal        string
	Bitcoin       string
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

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func hashPassword(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
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
