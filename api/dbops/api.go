package dbops

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/haibeichina/video_webserver/api/defs"
	"github.com/haibeichina/video_webserver/api/utils"
)

// AddUserCredential 添加用户认证信息
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	if _, err = stmtIns.Exec(loginName, pwd); err != nil {
		return err
	}

	return nil
}

// GetUserCredential 获取用户认证信息
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	defer stmtOut.Close()

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return pwd, nil
}

// DeleteUser 删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=?")
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}

	defer stmtDel.Close()

	if _, err = stmtDel.Exec(loginName, pwd); err != nil {
		return err
	}

	return nil
}

// GetUser 根据用户名获取用户信息
func GetUser(loginName string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("SELETE id, pwd FROM users WHERE login_name = ?")
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}

	defer stmtOut.Close()

	var (
		id  int
		pwd string
	)

	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err != sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.User{
		ID:        id,
		LoginName: loginName,
		Pwd:       pwd,
	}

	return res, nil
}

// AddNewVideo 添加视频
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info(id, author_id, name, 
		dispaly_ctime) VALUES(?,?,?,?)`)
	if err != nil {
		return nil, err
	}

	defer stmtIns.Close()

	if _, err = stmtIns.Exec(vid, aid, name, ctime); err != nil {
		return nil, err
	}

	res := &defs.VideoInfo{
		ID:           vid,
		AuthorID:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}

	return res, nil
}

// GetVideInfo 获取视频信息
func GetVideInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELETE author_id, name, display_ctime 
		FROM video_info WHERE id = ?`)

	var (
		aid  int
		dct  string
		name string
	)

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{
		ID:           vid,
		AuthorID:     aid,
		Name:         name,
		DisplayCtime: dct,
	}

	return res, nil
}

// ListVideoInfo 获取某段时间内用户的视频
func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELETE video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info 
		INNER JOIN users ON video_info.author_id = users.id 
		WHERE users.login_name = ? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time <= FROM_UNIXTIME(?) 
		ORDER BY video_info.create_time DESC`)

	var res []*defs.VideoInfo

	if err != nil {
		return res, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var (
			id, name, ctime string
			aid             int
		)

		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &defs.VideoInfo{
			ID:           id,
			AuthorID:     aid,
			Name:         name,
			DisplayCtime: ctime,
		}

		res = append(res, vi)
	}

	return res, nil
}

// DeleteVideoInfo 删除视频
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}

	defer stmtDel.Close()

	if _, err = stmtDel.Exec(vid); err != nil {
		return err
	}

	return nil
}

// AddNewComments 添加评论
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare("INSERT INFO comments (id, video_id, author_id, content) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	return nil
}

// ListComments 获取指定时间段的评论
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`SELETE comments.id, users.login_name, comments.content FROM comments 
		INNER JOIN users ON comments.author_id = user.id 
		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time << FROM_UNIXTIME(?) 
		ORDER BY comments.time DESC`)

	var res []*defs.Comment

	if err != nil {
		return res, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{
			ID:      id,
			VideoID: vid,
			Author:  name,
			Content: content,
		}

		res = append(res, c)
	}

	return res, nil
}
