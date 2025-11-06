package types

type Post struct {
	ID         string `json:"id"`
	Author     string `json:"author"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"createdAt"`
	LikesCount int    `json:"likesCount"`
}
