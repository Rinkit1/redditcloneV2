package post

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// WARNING! completly unsafe in multi-goroutine use, need add mutexes

type PostsMemoryRepository struct {
	data   []*Post
	lastID int
	mu     *sync.RWMutex
}

var (
	ErrNoPost    = errors.New("no post")
	ErrNoAuthor  = errors.New("you are not author")
	ErrNoComment = errors.New("no comment")
)

func NewMemoryRepo() *PostsMemoryRepository {
	return &PostsMemoryRepository{
		data:   make([]*Post, 0),
		lastID: 0,
		mu:     &sync.RWMutex{},
	}
}

func (repo *PostsMemoryRepository) GetAll() []*Post {
	return repo.data
}
func (repo *PostsMemoryRepository) AddPost(postJSON *Post, id, login string) {
	repo.lastID++
	postJSON.NewVote(1, id)
	postJSON.ID = strconv.Itoa(repo.lastID)
	postJSON.Created = time.Now().Format(time.RFC3339)
	postJSON.Views = 0
	postJSON.Author = Author{
		ID:       id,
		Username: login,
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.data = append(repo.data, postJSON)
}
func (repo *PostsMemoryRepository) OpenPost(id string) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.ID == id {
			val.Views++
			return val, nil
		}
	}
	return nil, ErrNoPost
}
func (repo *PostsMemoryRepository) Category(name string) []*Post {
	posts := make([]*Post, 0)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.Category == name {
			posts = append(posts, val)
		}
	}
	return posts
}
func (repo *PostsMemoryRepository) User(name string) []*Post {
	posts := make([]*Post, 0)
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.Author.Username == name {
			posts = append(posts, val)
		}
	}
	return posts
}
func (repo *PostsMemoryRepository) Delete(postID string, authorID string) (err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for ind, val := range repo.data {
		if val.ID == postID {
			if val.Author.ID == authorID {
				repo.data = append(repo.data[:ind], repo.data[ind+1:]...)
				return nil
			} else {
				return ErrNoAuthor
			}
		}
	}
	return ErrNoPost
}
func (repo *PostsMemoryRepository) Vote(vote int, authorID string, postID string) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.ID == postID {
			val.NewVote(vote, authorID)
			return val, nil
		}
	}
	return nil, ErrNoPost
}
func (repo *PostsMemoryRepository) UnVote(authorID string, postID string) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.ID == postID {
			val.DeleteVote(authorID)
			return val, nil
		}
	}
	return nil, ErrNoPost
}
func (repo *PostsMemoryRepository) AddComment(postID, body, authorID, login string) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.ID == postID {
			val.NewComment(body, authorID, login)
			return val, nil
		}
	}
	return nil, ErrNoPost
}
func (repo *PostsMemoryRepository) DeleteComment(postID, commentID, authorID string) (*Post, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, val := range repo.data {
		if val.ID == postID {
			err := val.DeleteComment(authorID, commentID)
			return val, err
		}
	}
	return nil, ErrNoPost
}
