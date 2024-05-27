package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	bid, err := getCotacao()
	if err != nil {
		log.Fatal(err)
	}
	err = writeToFile(bid)
	if err != nil {
		log.Fatal(err)
	}
}

func getCotacao() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return "", err
	}
	return c.Bid, nil
}

func writeToFile(bid string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
	if err != nil {
		return err
	}
	return nil
}
