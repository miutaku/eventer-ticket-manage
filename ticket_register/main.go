package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// JSONリクエストの構造体
type RequestData struct {
	TicketService    string    `json:"ticketService"`
	EventName        string    `json:"eventName"`
	EventDate        time.Time `json:"eventDate"`
	EventPlace       string    `json:"eventPlace"`
	TicketRegistDate time.Time `json:"ticketRegistDate"`
	TicketCount      int       `json:"ticketCount"`
	IsReserve        bool      `json:"isReserve"`
	PayLimitDate     time.Time `json:"payLimitDate"`
	UserId           string    `json:"userId"`
}

// ユーザーごとの詳細情報を表す構造体
type UserDetail struct {
	UserId           string    `json:"userId"`
	TicketRegistDate time.Time `json:"ticketRegistDate"`
	TicketCount      int       `json:"ticketCount"`
	IsReserve        bool      `json:"isReserve"`
	PayLimitDate     time.Time `json:"payLimitDate"`
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
		ticketTableName := "tickets"
		sql := fmt.Sprintf("INSERT INTO %s (ticketService, eventName, eventDate, eventPlace) VALUES (?, ?, ?, ?)", ticketTableName)
		ticketStmt, err := tx.Prepare(sql)
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL文の準備に失敗しました: %s, SQL: %s", err, sql), http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", sql)
		defer ticketStmt.Close()

		// ユーザー詳細情報の挿入
		userTableName := "user_details"
		userDetailSQL := fmt.Sprintf("INSERT INTO %s (userId, ticketRegistDate, ticketCount, isReserve, payLimitDate) VALUES (?, ?, ?, ?, ?)", userTableName)
		userDetailStmt, err := tx.Prepare(userDetailSQL)
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL文の準備に失敗しました: %s, SQL: %s", err, sql), http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", userDetailSQL)
		defer ticketStmt.Close()

		// ユーザーの所有しているチケット情報の挿入
		userTicketsTableName := "user_tickets"
		userTicketsSQL := fmt.Sprintf("INSERT INTO %s (userID, ticketID) VALUES (?, ?)", userTicketsTableName)
		userTicketsStmt, err := tx.Prepare(userTicketsSQL)
		if err != nil {
			http.Error(w, fmt.Sprintf("SQL文の準備に失敗しました: %s, SQL: %s", err, sql), http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", userTicketsSQL)
		defer ticketStmt.Close()

		// SQLの実行
		_, err = ticketStmt.Exec(
			reqData.TicketService,
			reqData.EventName,
			reqData.EventDate,
			reqData.EventPlace,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
			return
		}
		_, err = userDetailStmt.Exec(
			reqData.UserId,
			reqData.TicketRegistDate,
			reqData.TicketCount,
			reqData.IsReserve,
			reqData.PayLimitDate,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
			return
		}

		_, err = userTicketsStmt.Exec(
			reqData.UserId,
			reqData.TicketRegistDate,
			reqData.TicketCount,
			reqData.IsReserve,
			reqData.PayLimitDate,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Data inserted successfully")
	})

	// サーバーの起動
	http.ListenAndServe(":8080", nil)
}
