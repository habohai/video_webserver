package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/haibeichina/video_webserver/api/utils"

	"github.com/haibeichina/video_webserver/api/dbops"
	"github.com/haibeichina/video_webserver/api/defs"
	"github.com/haibeichina/video_webserver/api/session"
	"github.com/julienschmidt/httprouter"
)

// CreateUser 创建用户
func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
	}

	id := session.GenerateNewSessionID(ubody.Username)
	su := &defs.SignedUp{
		Success:   true,
		SessionID: id,
	}

	if resp, err := json.Marshal(su); err != nil {
		sendNormalResponse(w, string(resp), 201)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

// Login 用户登陆
// func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	uname := p.ByName("user_name")
// 	io.WriteString(w, uname)
// }

// Login 用户登陆
func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s", res)
	ubody := &defs.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	uname := p.ByName("username")
	fmt.Printf("Login url name: %s", uname)
	fmt.Printf("Login body name: %s", ubody.Username)
	if uname != ubody.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	fmt.Printf("%s", ubody.Username)
	pwd, err := dbops.GetUserCredential(ubody.Username)
	fmt.Printf("Login pwd: %s", pwd)
	fmt.Printf("Login body pwd: %s", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	id := session.GenerateNewSessionID(ubody.Username)
	si := &defs.SignedIn{
		Success:   true,
		SessionID: id,
	}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

// GentUserInfo 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		fmt.Printf("Unathorized user \n")
		return
	}

	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		fmt.Printf("Error in GetUserInfo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{
		ID: u.ID,
	}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

// AddNewVideo 添加新视频
func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		fmt.Printf("Unathorized user \n")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		fmt.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorID, nvbody.Name)
	fmt.Printf("Author id: %d, name: %s \n", nvbody.AuthorID, nvbody.Name)
	if err != nil {
		fmt.Printf("Error in AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

// ListAllVideos 获取视频列表
func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		fmt.Printf("Error in ListAllVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{
		Videos: vs,
	}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

// DeleteVideo 删除视频
func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		fmt.Printf("Error in DeleteVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", 204)
}

// PostComment 发送评论
func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body) // 这里的错误要处理一下

	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		fmt.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorID, cbody.Content); err != nil {
		fmt.Printf("Error in PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
	} else {
		sendNormalResponse(w, "ok", 201)
	}
}

// ShowComments 展示所有评论
func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		return
	}

	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		fmt.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{
		Comments: cm,
	}

	if resp, err := json.Marshal(cms); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}
