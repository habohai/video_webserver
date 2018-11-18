package dbops

import (
	"log"

	// 初始化数据库
	_ "github.com/go-sql-driver/mysql"
)

// ReadVideoDeletionRecord 读取要删除的视频列表
func ReadVideoDeletionRecord(count int) ([]string, error) {
	var ids []string
	stmtOut, err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	if err != nil {
		return ids, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	//for _, id:=range ids {
	//	log.Printf("SELECT video_id FROM video_del_rec includes %s\n", id)
	//}

	return ids, nil
}

// DelVideoDeletionRecord 删除成功删除记录后的的视频记录
func DelVideoDeletionRecord(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id=?")
	if err != nil {
		return err
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}

	return nil
}
