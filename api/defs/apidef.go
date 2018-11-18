package defs

// requests

// UserCredential 用户验证
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AuthorID int    `json:"author_id"`
	Name     string `json:"name"`
}

// response

type SignedUp struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
}

type UserSession struct {
	Username  string `json:"user_name"`
	SessionID string `json:"session_id"`
}

type UserInfo struct {
	ID int `json:"id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"vidoes"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}

// Date model

type User struct {
	ID        int
	LoginName string
	Pwd       string
}

// VideoInfo 视频信息
type VideoInfo struct {
	ID           string `json:"id"`
	AuthorID     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

// Comment 	评论
type Comment struct {
	ID      string `json:"id"`
	VideoID string `json:"video_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

// SimpleSession session用于服务端用户鉴权
type SimpleSession struct {
	Username string // login name
	TTL      int64
}
