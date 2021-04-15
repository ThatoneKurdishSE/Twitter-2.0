package routes

import(
	"github.com/gorilla/mux"
	"net/http"
	"../models"
	"../sessions"
	"../tools"
	"../middleware"
)

func NewRouter() *mux.Router{
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	// r.HandleFunc("/{username}",
	// 	middleware.AuthRequired(userGetHandler)).Methods("GET")
	return r
}
func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		session, _ := sessions.Store.Get(r, "session")
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func indexGetHandler(w http.ResponseWriter, r *http.Request)  {
	updates, err := models.GetAllUpdates()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	tools.ExecuteTemplate(w, "update.html", updates)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request)  {
	
	session, _ := sessions.Store.Get(r, "session")
	untypeduserId := session.Values["user_id"]
	userId, ok := untypeduserId.(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	http.Redirect(w, r, "/", 302)
}

// func userGetHandler(w http.ResponseWriter, r *http.Request)  {
// 	vars := mux.Vars(r)
// 	username := vars["username"]
// 	user, err := models.GetUserByUsername(username)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Internal Server Error"))
// 		return
// 	}
// 	userId, err := user.GetId()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Internal Server Error"))
// 		return
// 	}
// 	updates, err := models.GetUpdates(userId)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte("Internal Server Error"))
// 		return
// 	}

func loginGetHandler(w http.ResponseWriter, r *http.Request){
	tools.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err{
		case models.ErrUserNotFound:
			tools.ExecuteTemplate(w, "login.html", "Unknown User")
		case models.ErrInvalidLogin:
			tools.ExecuteTemplate(w, "login.html", "Invalid Login!")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
		return
	}
	userId, err := user.GetId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request){
	tools.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	http.Redirect(w, r, "/login", 302)
}