package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"my/redditclone/pkg/post"
	"my/redditclone/pkg/session"
	"net/http"
)

type VoteHandler struct {
	Logger   *log.Logger
	PostRepo post.PostsRepo
}

func (h *VoteHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry voteHandler.Upvote: [%s] %s\n",
		r.Method, r.URL.Path)
	vars := mux.Vars(r)
	id, _ := session.SessionFromContext(r.Context())
	post, err := h.PostRepo.Vote(1, id, vars["id"])
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
	h.Logger.Printf("successfully exit voteHandler.Upvote\n")
}
func (h *VoteHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry voteHandler.Downvote: [%s] %s\n",
		r.Method, r.URL.Path)
	vars := mux.Vars(r)
	id, _ := session.SessionFromContext(r.Context())
	post, err := h.PostRepo.Vote(-1, id, vars["id"])
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
	h.Logger.Printf("successfully exit voteHandler.Downvote\n")
}
func (h *VoteHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("successfully entry voteHandler.Unvote: [%s] %s\n",
		r.Method, r.URL.Path)
	vars := mux.Vars(r)
	id, _ := session.SessionFromContext(r.Context())
	post, err := h.PostRepo.UnVote(id, vars["id"])
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
	h.Logger.Printf("successfully exit voteHandler.Unvote\n")
}
