package dbops

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/haibeichina/video_webserver/api/defs"
)

// InsertSession 插入session
func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (session_id, TTL, login_name) VALUES (?,?,?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveSession 根据sid获取session的信息
func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare("SELETE TTL, login_name FROM sessions WHERE session_id=?")
	if err != nil {
		return nil, err
	}

	defer stmtOut.Close()

	var (
		ttl   string
		uname string
	)

	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
		return nil, err
	}

	return ss, nil
}

// RetrieveAllSessions 获取所有session信息，用于程序初始化
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELETE * FROM sessions")
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()
	if err != nil {
		fmt.Printf("")
	}

	for rows.Next() {
		var (
			id        string
			ttlstr    string
			loginName string
		)
		if er := rows.Scan(&id, &ttlstr, &loginName); er != nil {
			fmt.Printf("retrive sessions error: %s", er)
			break
		}

		if ttl, er := strconv.ParseInt(ttlstr, 10, 64); er == nil {
			ss := &defs.SimpleSession{
				Username: loginName,
				TTL:      ttl,
			}
			m.Store(id, ss)
			fmt.Printf("session id: %s, ttl: %d", id, ss.TTL)
		} else {
			return nil, err
		}
	}

	return m, nil
}

// DeleteSession 删除session信息
func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	return nil
}
