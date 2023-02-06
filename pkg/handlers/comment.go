package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"my/redditclone/pkg/post"
	"my/redditclone/pkg/session"
	"net/http"
)

type CommentHandler struct {
	PostRepo post.PostsRepo
	Logger   *log.Logger
}

type CommentInfo struct {
	Comment string `json:"comment"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry commentHandler.Create: [%s] %s\n",
		r.Method, r.URL.Path)
	vars := mux.Vars(r)
	ID, login := session.SessionFromContext(r.Context())
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "can not read body")
		return
	}
	info := new(CommentInfo)
	err = json.Unmarshal(body, &info)
	if err != nil {
		SendMessage(w, http.StatusBadRequest, "incorrect params")
		return
	}
	post, err := h.PostRepo.AddComment(vars["id"], info.Comment, ID, login)
	if err != nil {
		SendMessage(w, http.StatusBadRequest, "no post")
		return
	}
	data, err := json.Marshal(post)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "JSON Marshal error")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "Write error")
		return
	}
	h.Logger.Printf("successfully exit commentHandler.Create\n")
}
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry commentHandler.Delete: [%s] %s\n",
		r.Method, r.URL.Path)
	vars := mux.Vars(r)
	ID, _ := session.SessionFromContext(r.Context())
	findPost, err := h.PostRepo.DeleteComment(vars["postID"], vars["commentID"], ID)
	if err != nil {
		text := "unknown error"
		switch {
		case err == post.ErrNoPost:
			text = "no post"
		case err == post.ErrNoAuthor:
			text = "you are not author"
		case err == post.ErrNoComment:
			text = "no comment"
		}
		SendMessage(w, http.StatusBadRequest, text)
		return
	}
	data, err := json.Marshal(findPost)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "JSON Marshal error")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		SendMessage(w, http.StatusInternalServerError, "Write error")
		return
	}
	h.Logger.Printf("successfully exit commentHandler.Delete\n")
}
