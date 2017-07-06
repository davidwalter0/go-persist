package main

import (
	"fmt"
	"github.com/davidwalter0/go-persist"
	"github.com/davidwalter0/go-persist/pgsql"
	uuid "github.com/davidwalter0/go-persist/uuid"
	"time"
)

var db *persist.Database = &persist.Database{}

func init() {
	db.Configure()
	db.Connect()
	db.Initialize(pgsql.Schema)
}

func main() {
	row := db.QueryRow("INSERT INTO pages (page_guid, page_title, page_content, page_date) VALUES ('" + uuid.GUID().String() + "', 'Hello, World', 'I''m so glad you found this page!  It''s been sitting patiently on the Internet for some time, just waiting for a visitor.', CURRENT_TIMESTAMP)")
	fmt.Printf("%v\n", *row)
	rows := db.Query("select * from pages")
	defer rows.Close()
	for rows.Next() {
		var pageId int
		var pageGuid string
		var pageTitle string
		var pageContent string
		var pageDate time.Time
		rows.Scan(&pageId, &pageGuid, &pageTitle, &pageContent, &pageDate)
		fmt.Println(pageId, pageGuid, pageTitle, pageContent, pageDate)
	}
}
