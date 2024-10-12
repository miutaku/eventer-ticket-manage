package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true", user, password, host, database, charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQLへの接続に失敗しました:", err)
	}
	log.Println("MySQLへの接続に成功しました")
	defer db.Close()

	// チケット情報更新API
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// リクエストボディの読み込み
		type UpdateData struct {
			UserId     string    `json:"userId"`
			EventName  string    `json:"eventName"`
			EventDate  time.Time `json:"eventDate"`
			EventPlace string    `json:"eventPlace"`
			IsPaid     bool      `json:"isPaid"`
		}
		var reqData UpdateData
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, "JSONデコードに失敗しました:", http.StatusBadRequest)
			return
		}

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

		// ticketsテーブルからticketIdを取得
		PaidSelectSQL := "SELECT ticketId FROM tickets WHERE eventName = ? AND eventDate = ? AND eventPlace = ?"
		var ticketId int64
		err = tx.QueryRow(PaidSelectSQL, reqData.EventName, reqData.EventDate, reqData.EventPlace).Scan(&ticketId)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "指定されたイベントが見つかりません:", http.StatusNotFound)
			} else {
				handleError(w, err, http.StatusInternalServerError)
			}
			return
		}

		// user_ticketsテーブルを更新
		paidUpdateSQL := "UPDATE user_tickets SET isPaid = ? WHERE userId = ? AND ticketId = ?"
		updateStmt, err := tx.Prepare(paidUpdateSQL)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		defer updateStmt.Close()

		_, err = updateStmt.Exec(reqData.IsPaid, reqData.UserId, ticketId)
		if err != nil {
			http.Error(w, fmt.Sprintf("データの更新に失敗しました: %s", err), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "Data updated successfully")
		}
	})

	// チケット新規登録API
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		type RequestData struct {
			TicketService    string    `json:"ticketService"`
			EventName        string    `json:"eventName"`
			EventDate        time.Time `json:"eventDate"`
			EventPlace       string    `json:"eventPlace"`
			TicketRegistDate time.Time `json:"ticketRegistDate"`
			TicketCount      int       `json:"ticketCount"`
			IsReserve        bool      `json:"isReserve"`
			PayLimitDate     time.Time `json:"payLimitDate"`
			IsPaid           bool      `json:"isPaid"`
			UserId           string    `json:"userId"`
		}

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

		// 同日のチケットが既に存在するか確認（丸め処理）
		selectSQL := "SELECT ticketId FROM tickets WHERE eventName = ? AND DATE(eventDate) = DATE(?) AND eventPlace = ?"
		var existingTicketId int64
		err = tx.QueryRow(selectSQL, reqData.EventName, reqData.EventDate, reqData.EventPlace).Scan(&existingTicketId)

		// 新規チケット挿入
		var ticketId int64
		if err == sql.ErrNoRows {
			log.Println("このチケットはDBにまだ存在しません。新規チケットを挿入します。")
			// 新規チケットを挿入
			ticketSQL := "INSERT INTO tickets (ticketService, ticketRegistDate, eventName, eventDate, eventPlace) VALUES (?, ?, ?, ?, ?)"
			result, err := tx.Exec(ticketSQL, reqData.TicketService, reqData.TicketRegistDate, reqData.EventName, reqData.EventDate, reqData.EventPlace)
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
		} else if err != nil {
			http.Error(w, fmt.Sprintf("チケット重複確認中にエラーが発生しました: %s", err), http.StatusInternalServerError)
			return
		} else {
			log.Printf("既存のチケットが見つかりました。ticketId: %d", existingTicketId)
			ticketId = existingTicketId
		}

		// user_ticketsテーブルにおける重複確認
		var isDuplicate bool
		var duplicateTicketIds []int64 // 重複チケットのIDを格納するスライス
		checkDuplicateSQL := `
			SELECT ut.ticketId 
			FROM user_tickets ut 
			JOIN tickets t ON ut.ticketId = t.ticketId 
			WHERE ut.userId = ? AND DATE(t.eventDate) = DATE(?)`

		rows, err := tx.Query(checkDuplicateSQL, reqData.UserId, reqData.EventDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("重複確認中にエラーが発生しました: %s", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var duplicateTicketId int64
			if err := rows.Scan(&duplicateTicketId); err != nil {
				http.Error(w, fmt.Sprintf("重複チケットのID取得中にエラーが発生しました: %s", err), http.StatusInternalServerError)
				return
			}
			duplicateTicketIds = append(duplicateTicketIds, duplicateTicketId)
		}

		if len(duplicateTicketIds) > 0 {
			isDuplicate = true
			log.Printf("重複チケットが見つかりました。duplicateTicketIds: %v", duplicateTicketIds)
		} else {
			isDuplicate = false
		}

		// duplicateTicketId をカンマ区切りの文字列に変換
		var duplicateTicketIdStr string
		if isDuplicate {
			for i, id := range duplicateTicketIds {
				if i > 0 {
					duplicateTicketIdStr += ","
				}
				duplicateTicketIdStr += fmt.Sprintf("%d", id)
			}
		}

		// user_tickets挿入
		userTicketSQL := "INSERT INTO user_tickets (userId, ticketId, ticketCount, isReserve, payLimitDate, isPaid, duplicateTicketId, isDuplicate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = tx.Exec(userTicketSQL, reqData.UserId, ticketId, reqData.TicketCount, reqData.IsReserve, reqData.PayLimitDate, reqData.IsPaid, duplicateTicketIdStr, isDuplicate)
		if err != nil {
			http.Error(w, fmt.Sprintf("データの挿入に失敗しました: %s", err), http.StatusInternalServerError)
			return
		}

		log.Printf("ユーザーのチケット情報が挿入されました。userId: %s, ticketId: %d, isDuplicate: %t", reqData.UserId, ticketId, isDuplicate)

		// 新規チケットが重複している場合、duplicateTicketIdStrを使って既存の重複チケットを更新
		if isDuplicate {
			// duplicateTicketIdsをスライスに変換
			duplicateIds := strings.Split(duplicateTicketIdStr, ",")

			// プレースホルダーを作成
			placeholders := make([]string, len(duplicateIds))
			for i := range duplicateIds {
				placeholders[i] = "?"
			}
			placeholderStr := strings.Join(placeholders, ",")

			updateSQL := fmt.Sprintf(`
				UPDATE user_tickets
				SET duplicateTicketId = ?, isDuplicate = ?
				WHERE userId = ? AND ticketId IN (%s)
			`, placeholderStr)

			// 更新するチケットIDを引数に渡す
			args := []any{duplicateTicketIdStr, true, reqData.UserId}
			for _, id := range duplicateIds {
				args = append(args, id)
			}

			_, err = tx.Exec(updateSQL, args...)
			if err != nil {
				http.Error(w, fmt.Sprintf("重複チケットの更新に失敗しました: %s", err), http.StatusInternalServerError)
				return
			}
		}

		fmt.Fprintf(w, "Data inserted successfully")
	})

	// ユーザーのチケット情報取得API
	http.HandleFunc("/fetchUserTickets", func(w http.ResponseWriter, r *http.Request) {
		// リクエストからユーザーIDを抽出
		userId := r.URL.Query().Get("userId")
		if userId == "" {
			http.Error(w, "Missing userId parameter", http.StatusBadRequest)
			return
		}

		// SQL文の準備
		query := `
        SELECT
            u.userId,
            t.eventName,
            t.eventDate,
            u.ticketCount
        FROM
            user_tickets u
        INNER JOIN tickets t ON u.ticketId = t.ticketId
        WHERE
            u.userId = ?
            AND
            (
                u.isPaid = 1
                OR
                (u.isPaid = 0 AND u.payLimitDate > CURDATE())
            );
    `

		// SQLの実行
		rows, err := db.Query(query, userId)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 結果を格納するスライス
		var results []struct {
			UserId      string    `json:"userId"`
			EventName   string    `json:"eventName"`
			EventDate   time.Time `json:"eventDate"`
			TicketCount int       `json:"ticketCount"`
		}

		for rows.Next() {
			var r struct {
				UserId      string    `json:"userId"`
				EventName   string    `json:"eventName"`
				EventDate   time.Time `json:"eventDate"`
				TicketCount int       `json:"ticketCount"`
			}
			if err := rows.Scan(&r.UserId, &r.EventName, &r.EventDate, &r.TicketCount); err != nil {
				handleError(w, err, http.StatusInternalServerError)
				return
			}
			results = append(results, r)
		}

		// JSONエンコードしてレスポンスを返す
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			handleError(w, err, http.StatusInternalServerError)
		}
	})

	// サーバーの起動
	http.ListenAndServe(":8080", nil)
}
