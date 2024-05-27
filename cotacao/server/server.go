package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// JSON of the whole request to the API
type CotacaoRequest struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

// JSON that will be returned to the client
type CotacaoResponse struct {
	Bid string `json:"bid"`
}

// SQL statement to create the table on the DB
const createTableStmt string = `CREATE TABLE IF NOT EXISTS cotacao(
	timestamp INTEGER NOT NULL PRIMARY KEY,
	bid TEXT NOT NULL)`

// SQL statemnent to insert data into the table
const insertCotacaoStmt string = "INSERT INTO cotacao VALUES (?, ?);"

func main() {
	http.HandleFunc("/cotacao", handleCotacao)
	http.ListenAndServe(":8080", nil)
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	// make the request to the API
	cotacao, err := getCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR: ", err)
		return
	}

	// write the response to the client
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CotacaoResponse{cotacao.USDBRL.Bid})

	// write the values into the DB
	db, err := openDB()
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
	defer db.Close()
	err = insertData(db, cotacao)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}
}

func getCotacao() (*CotacaoRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cr CotacaoRequest
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createTableStmt)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertData(db *sql.DB, cr *CotacaoRequest) error {

	stmt, err := db.Prepare(insertCotacaoStmt)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err = stmt.ExecContext(ctx, cr.USDBRL.Timestamp, cr.USDBRL.Bid)
	return err
}
