package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 正規表現で許可するテーブル名文字列
var validTableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// JSONリクエストの構造体
type RequestData struct {
	TicketService string    `json:"ticketService"`
	RegistDate    time.Time `json:"registDate"`
	EventDate     time.Time `json:"eventDate"`
	EventPlace    string    `json:"eventPlace"`
	EventName     string    `json:"eventName"`
	TicketCount   int       `json:"ticketCount"`
	IsReserve     bool      `json:"isReserve"`
	PayLimitDate  time.Time `json:"payLimitDate"`
}

func main() {
	// 環境変数から接続情報を取得
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	database := os.Getenv("MYSQL_DATABASE")
	charset := os.Getenv("MYSQL_CHARSET")

	// 接続情報が不足している場合、エラーを返す
	if user == "" || password == "" || host == "" || database == "" || charset == "" {
		log.Fatal("MySQL接続に必要な環境変数が不足しています")
	}

	// MySQLへの接続
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", user, password, host, database, charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQLへの接続に失敗しました:", err)
	}
	log.Println("MySQLへの接続に成功しました")
	defer db.Close()

	// HTTPハンドラー
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		// リクエストボディの読み込み
		var reqData RequestData
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, "JSONデコードに失敗しました:", http.StatusBadRequest)
			return
		}
		log.Printf("Request Data: %+v\n", reqData)

		// テーブル名を動的に設定し、正規表現で検証
		tableName := reqData.TicketService
		if !validTableNameRegex.MatchString(tableName) {
			http.Error(w, "不正なテーブル名です:", http.StatusBadRequest)
			return
		}
		log.Print("Table Name: ", tableName)

		// トランザクション開始
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "トランザクション開始に失敗しました:", http.StatusInternalServerError)
			return
		}
		log.Println("Transaction started")
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()

		// SQL文の準備
		sql := fmt.Sprintf("INSERT INTO %s (ticketService, registDate, eventDate, eventPlace, eventName, ticketCount, isReserve, payLimitDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", tableName)
		stmt, err := tx.Prepare(sql)
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL文の準備に失敗しました: %s, SQL: %s", err, sql), http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", sql)
		defer stmt.Close()

		// SQLの実行
		_, err = stmt.Exec(
			reqData.TicketService,
			reqData.RegistDate,
			reqData.EventDate,
			reqData.EventPlace,
			reqData.EventName,
			reqData.TicketCount,
			reqData.IsReserve,
			reqData.PayLimitDate,
		)
		if err != nil {
			http.Error(w, "データの挿入に失敗しました:", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Data inserted successfully")
	})

	// サーバーの起動
	http.ListenAndServe(":8080", nil)
}
