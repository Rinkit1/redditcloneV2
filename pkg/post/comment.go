package post

import (
	"strconv"
	"time"
)

type Comment struct {
	Created string `json:"created"`
	Author  Author `json:"author"`
	Body    string `json:"body"`
	ID      string `json:"id"`
}

func (p *Post) NewComment(body, authorID, login string) {
	p.lastCommentID++
	p.Comment = append(p.Comment, &Comment{
		Created: time.Now().Format(time.RFC3339),
		Author: Author{
			ID:       authorID,
			Username: login,
		},
		Body: body,
		ID:   strconv.Itoa(p.lastCommentID),
	})
}

func (p *Post) DeleteComment(authorID, commentID string) error {
	for ind, val := range p.Comment {
		if val.ID == commentID {
			if val.Author.ID != authorID {
				return ErrNoAuthor
			}
			p.Comment = append(p.Comment[:ind], p.Comment[ind+1:]...)
			return nil
		}
	}
	return ErrNoComment
}
