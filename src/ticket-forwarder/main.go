package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// リクエストボディを受け取る構造体
type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"` // HTMLの本文
}

// 別のAPIに転送するための関数
func forwardEmail(email EmailRequest) error {
	forwardURL := "https://example.com/forward"

	// 転送するデータ
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %v", err)
	}

	// POSTリクエストを送信
	req, err := http.NewRequest("POST", forwardURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %s, body: %s", resp.Status, body)
	}

	return nil
}

// メールの転送を処理するハンドラー
func handleEmailForward(w http.ResponseWriter, r *http.Request) {
	var emailReq EmailRequest

	// リクエストボディをパース
	if err := json.NewDecoder(r.Body).Decode(&emailReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 転送元と転送先のアドレスをチェック
	if strings.ToLower(emailReq.From) == "lt-mail@l-tike.com" && strings.ToLower(emailReq.To) == "mtakumi.0925@gmail.com" {
		// 転送処理
		if err := forwardEmail(emailReq); err != nil {
			log.Printf("Failed to forward email: %v", err)
			http.Error(w, "Failed to forward email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email forwarded successfully"))
	} else {
		// 転送条件に合わない場合は403エラーを返す
		http.Error(w, "Forbidden: Invalid sender or recipient", http.StatusForbidden)
	}
}

func main() {
	http.HandleFunc("/forward-email", handleEmailForward)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
