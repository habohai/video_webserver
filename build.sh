#! /bin/bash

# Build web UI
cd ~/goprojects/src/github.com/haibeichina/video_webserver/web
go install
cp ~/goprojects/bin/web ~/goprojects/bin/video_server_web_ui/web
cp -R ~/work/src/github.com/haibeichina/video_webserver/templates ~/goprojects/bin/video_server_web_ui/