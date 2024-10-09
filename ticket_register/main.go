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

func handleError(w http.ResponseWriter, err error, status int) {
	log.Printf("Error: %s", err)
	http.Error(w, fmt.Sprintf("An error occurred: %s", err), status)
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
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				err = tx.Rollback()
			} else if err != nil {
				err = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()

		// SQL文の準備
		// tickets
		ticketSQL := "INSERT INTO tickets (ticketService, ticketRegistDate, eventName, eventDate, eventPlace) VALUES (?, ?, ?, ?, ?)"
		ticketStmt, err := tx.Prepare(ticketSQL)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", ticketSQL)
		defer ticketStmt.Close()

		// user_tickets
		userTicketSQL := "INSERT INTO user_tickets (userId, ticketId,  ticketCount, isReserve, payLimitDate) VALUES (?, ?, ?, ?, ?)"
		userTicketStmt, err := tx.Prepare(userTicketSQL)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		log.Printf("SQL statement prepared: %s", userTicketSQL)
		defer userTicketStmt.Close()

		// SQLの実行
		// tickets
		selectSQL := "SELECT ticketId FROM tickets WHERE eventName = ? AND eventDate = ? AND eventPlace = ?" // 重複チェック
		var ticketId int64
		err = tx.QueryRow(selectSQL, reqData.EventName, reqData.EventDate, reqData.EventPlace).Scan(&ticketId)

		if err != nil && err != sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("チケット重複確認中にエラーが発生しました: %s", err), http.StatusInternalServerError)
			return
		}

		if err == nil {
			// チケットが既に存在する場合
			log.Printf("既存のチケットが見つかりました。ticketId: %d", ticketId)
		} else {
			// チケットが存在しない場合、新規に挿入
			result, err := ticketStmt.Exec(
				reqData.TicketService,
				reqData.TicketRegistDate,
				reqData.EventName,
				reqData.EventDate,
				reqData.EventPlace,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
				return
			}

			ticketId, err = result.LastInsertId()
			if err != nil {
				http.Error(w, fmt.Sprintf("ticketIdの確認ができません: %s", err), http.StatusInternalServerError)
				return
			}
			log.Printf("新しいチケットが挿入されました。ticketId: %d", ticketId)
		}

		// user_tickets
		_, err = userTicketStmt.Exec(
			reqData.UserId,
			ticketId,
			reqData.TicketCount,
			reqData.IsReserve,
			reqData.PayLimitDate,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "Data inserted successfully")
	})

	// サーバーの起動
	http.ListenAndServe(":8080", nil)
}
