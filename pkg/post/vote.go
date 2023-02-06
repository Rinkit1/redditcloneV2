package post

type Vote struct {
	UserID string `json:"user"`
	Votes  int    `json:"vote"`
}

func (p *Post) NewVote(vote int, id string) {
	p.DeleteVote(id)
	p.Vote = append(p.Vote, &Vote{
		Votes:  vote,
		UserID: id,
	})
	p.Score += vote
	p.UpvotePercentage = 100 * (p.Score + len(p.Vote)) / (2 * len(p.Vote))
}

func (p *Post) DeleteVote(id string) {
	for ind, val := range p.Vote {
		if val.UserID == id {
			p.Score -= val.Votes
			if len(p.Vote) == 1 {
				p.UpvotePercentage = 0
			} else {
				p.UpvotePercentage = (p.Score + len(p.Vote) - 1) / 2 / (len(p.Vote) - 1) * 100
			}
			p.Vote = append(p.Vote[:ind], p.Vote[ind+1:]...)
			break
		}
	}
}
