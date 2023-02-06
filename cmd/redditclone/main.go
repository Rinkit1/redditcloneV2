package main

import (
	"github.com/gorilla/mux"
	"log"
	"my/redditclone/pkg/handlers"
	"my/redditclone/pkg/middleware"
	"my/redditclone/pkg/post"
	"my/redditclone/pkg/user"
	"net/http"
	"os"
	"text/template"
)

func main() {
	templates := template.Must(template.ParseGlob("./redditclone/static/html/index.html"))
	logger := log.New(os.Stdout, "STD ", log.Ltime|log.Lshortfile)

	userRepo := user.NewMemoryRepo()
	postRepo := post.NewMemoryRepo()

	postHandlers := &handlers.PostHandler{
		Logger:   logger,
		PostRepo: postRepo,
	}
	userHandlers := &handlers.UserHandler{
		Tmpl:     templates,
		Logger:   logger,
		UserRepo: userRepo,
	}
	voteHandlers := &handlers.VoteHandler{
		Logger:   logger,
		PostRepo: postRepo,
	}
	commentHandlers := &handlers.CommentHandler{
		Logger:   logger,
		PostRepo: postRepo,
	}

	r := mux.NewRouter()

	r.Handle("/static/js/{file}", http.FileServer(http.Dir("./redditclone")))
	r.Handle("/static/css/{file}", http.FileServer(http.Dir("./redditclone")))

	r.NotFoundHandler = http.HandlerFunc(userHandlers.Index)
	r.HandleFunc("/api/login", userHandlers.Login).Methods("POST")
	r.HandleFunc("/api/register", userHandlers.Register).Methods("POST")

	r.HandleFunc("/api/posts/", postHandlers.List).Methods("GET")
	r.HandleFunc("/api/posts", postHandlers.Create).Methods("POST")
	r.HandleFunc("/api/post/{id}", postHandlers.Open).Methods("GET")
	r.HandleFunc("/api/post/{id}", postHandlers.Delete).Methods("DELETE")
	r.HandleFunc("/api/posts/{category}", postHandlers.Category).Methods("GET")
	r.HandleFunc("/api/user/{user}", postHandlers.User).Methods("GET")

	r.HandleFunc("/api/post/{id}/upvote", voteHandlers.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/downvote", voteHandlers.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/unvote", voteHandlers.Unvote).Methods("GET")

	r.HandleFunc("/api/post/{id}", commentHandlers.Create).Methods("POST")
	r.HandleFunc("/api/post/{postID}/{commentID}", commentHandlers.Delete).Methods("DELETE")

	mux := middleware.Auth(r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Println("listen and serve error")
	}
}
