package dbops

import (
	"log"

	// 初始化数据库
	_ "github.com/go-sql-driver/mysql"
)

// AddVideoDeletionRecord 添加视频删除记录
func AddVideoDeletionRecord(vid string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO video_del_rec(video_id) VALUES(?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}
