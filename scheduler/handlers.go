package main

import (
	"net/http"

	"github.com/haibeichina/video_webserver/scheduler/dbops"

	"github.com/julienschmidt/httprouter"
)

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	if len(vid) == 0 {
		sendResponse(w, 400, "video id should not be empty")
	}
	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		sendResponse(w, 500, "Internal server error")
		return
	}
	sendResponse(w, 200, "delete video successfully")
	return
}
