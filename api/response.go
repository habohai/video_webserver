package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/haibeichina/video_webserver/api/defs"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	w.WriteHeader(errResp.HTTPSC)

	resStr, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(resStr))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
