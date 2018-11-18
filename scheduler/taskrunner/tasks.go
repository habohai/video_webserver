package taskrunner

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/haibeichina/video_webserver/scheduler/dbops"
)

// deleteVideo 删除实际的视频文件
func deleteVideo(vid string) error {
	path, _ := filepath.Abs(VIDEO_PATH + vid)
	log.Println(path)
	err := os.Remove(VIDEO_PATH + vid)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}
	return nil
}

// VideoClearDispatcher 删除视频的任务分发方法
func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3) // 3 作为参数传进来
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
	}
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}

// VideoClearExecutor 执行删除视频的方法
func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) { // 这里待优化 使用 sync.WaitGroup
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}
