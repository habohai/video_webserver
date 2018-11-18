package main

import (
	"net/http"

	"github.com/haibeichina/video_webserver/api/defs"
	"github.com/haibeichina/video_webserver/api/session"
)

var (
	// HeadweFieldSession 自定义header
	HeadweFieldSession = "X-Session-Id"
	// HeadweFieldUname 自定义header
	HeadweFieldUname = "X-User-Name"
)

// validateUserSession Check if the current user has the permission
// Use session id to do the check
func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HeadweFieldSession)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HeadweFieldUname, uname)
	return true
}

// ValidateUser 验证用户
func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HeadweFieldUname)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}
