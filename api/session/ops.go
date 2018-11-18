package session

import (
	"sync"
	"time"

	"github.com/avenssi/video_server/api/dbops"
	"github.com/haibeichina/video_webserver/api/defs"
	"github.com/haibeichina/video_webserver/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

// LoadSessionsFromDB 从数据库中获取session
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

// GenerateNewSessionID 生成一个session id
func GenerateNewSessionID(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Serverside session valid time: 30 min 改进为可以设置的选项

	ss := &defs.SimpleSession{
		Username: un,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, un)

	return id
}

// IsSessionExpired 判断session是否过期
func IsSessionExpired(sid string) (string, bool) {
	ct := nowInMilli()
	if ss, ok := sessionMap.Load(sid); ok {
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		return ss.(*defs.SimpleSession).Username, false
	}

	ss, err := dbops.RetrieveSession(sid)
	if err != nil || ss == nil {
		return "", true
	}

	if ss.TTL < ct {
		deleteExpiredSession(sid)
		return "", true
	}

	sessionMap.Store(sid, ss)
	return ss.Username, false
}
