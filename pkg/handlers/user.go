package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"my/redditclone/pkg/session"
	"my/redditclone/pkg/user"
	"net/http"
	"text/template"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *log.Logger
	UserRepo user.UserRepo
}
type UserJSON struct {
	Username string
	Password string
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry userHandler.Index: [%s] %s\n",
		r.Method, r.URL.Path)
	err := h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "template error")
		return
	}
	h.Logger.Printf("successfully exit userHandler.Index\n")
}
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry userHandler.Login: [%s] %s\n",
		r.Method, r.URL.Path)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "can not read body")
		return
	}
	userJSON := new(UserJSON)
	err = json.Unmarshal(body, &userJSON)
	if err != nil {
		SendMessage(w, http.StatusBadRequest, "incorrect params")
		return
	}
	u, err := h.UserRepo.Authorize(userJSON.Username, userJSON.Password)
	if err == user.ErrNoUser {
		SendMessage(w, http.StatusBadRequest, "no user")
		return
	}
	if err == user.ErrBadPass {
		SendMessage(w, http.StatusBadRequest, "bad password")
		return
	}
	token, err := session.Create(u.ID, u.Login)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "session error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(map[string]interface{}{"token": token})
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "JSON Marshal error")
		return
	}
	_, err = w.Write(data)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "Write error")
		return
	}
	h.Logger.Printf("successfully exit userHandler.Login\n")
}
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry userHandler.Register: [%s] %s\n",
		r.Method, r.URL.Path)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "can not read body")
		return
	}
	userJSON := new(UserJSON)
	err = json.Unmarshal(body, &userJSON)
	if err != nil {
		SendMessage(w, http.StatusBadRequest, "incorrect params")
		return
	}
	err = h.UserRepo.AddUser(userJSON.Username, userJSON.Password)
	if err == user.ErrUserAlreadyExist {
		SendMessage(w, http.StatusBadRequest, "user already exist")
		return
	}
	u, err := h.UserRepo.Authorize(userJSON.Username, userJSON.Password)
	if err == user.ErrNoUser {
		SendMessage(w, http.StatusBadRequest, "no user")
		return
	}
	if err == user.ErrBadPass {
		SendMessage(w, http.StatusBadRequest, "bad password")
		return
	}
	token, err := session.Create(u.ID, u.Login)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "session error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, err := json.Marshal(map[string]interface{}{"token": token})
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "JSON Marshal error")
		return
	}
	_, err = w.Write(data)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "Write error")
		return
	}
	h.Logger.Printf("successfully exit userHandler.Register\n")
}
