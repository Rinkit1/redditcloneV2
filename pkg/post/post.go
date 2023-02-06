package post

type Post struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Created  string `json:"created"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Text     string `json:"text"`
	URL      string `json:"url"`

	Vote             []*Vote `json:"votes"`
	Score            int     `json:"score"`
	UpvotePercentage int     `json:"upvotePercentage"`

	Comment       []*Comment `json:"comments"`
	lastCommentID int        `json:"-"`

	Views int `json:"views"`

	Author Author `json:"author"`
}

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type PostsRepo interface {
	GetAll() []*Post
	AddPost(postJSON *Post, id, login string)
	OpenPost(id string) (*Post, error)
	Vote(vote int, authorID string, postID string) (*Post, error)
	UnVote(authorID string, postID string) (*Post, error)
	Delete(postID string, authorID string) (err error)
	Category(name string) []*Post
	User(name string) []*Post
	AddComment(postID, body, authorID, login string) (*Post, error)
	DeleteComment(postID, commentID, authorID string) (*Post, error)
}
