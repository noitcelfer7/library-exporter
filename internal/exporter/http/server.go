package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"library_exporter/internal/exporter/config"
	postgresql "library_exporter/internal/exporter/database"
	"log"
	"net"
	"net/http"
	"time"
)

type IssuesByDateResponse struct {
	Dates map[string]int `json:"dates"`
}

type SpreadsheetRequest struct {
	SpreadsheetID string          `json:"spreadsheetId"`
	Data          [][]interface{} `json:"data"`
}

func Serve(config *config.Config, db *postgresql.Database) {
	http.HandleFunc("/issues", issuesByDateHandler(db))

	addr := net.JoinHostPort(config.Http.Server.Host, config.Http.Server.Port)

	http.ListenAndServe(addr, nil)
}

func issuesByDateHandler(db *postgresql.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Выполняем SQL-запрос
		rows, err := db.QueryContext(ctx,
			`SELECT 
							issue_date, 
							COUNT(uuid) 
					 FROM issues 
					 GROUP BY issue_date 
					 ORDER BY issue_date DESC`)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		result := make(map[string]int)
		var dataArray [][]interface{}
		for rows.Next() {
			var date time.Time
			var count int

			if err := rows.Scan(&date, &count); err != nil {
				http.Error(w, fmt.Sprintf("Row scan error: %v", err), http.StatusInternalServerError)
				return
			}

			result[date.Format("2006-01-02")] = count

			dataArray = append(dataArray, []interface{}{count, date})
		}

		// Проверяем ошибки итерации
		if err = rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Rows iteration error: %v", err), http.StatusInternalServerError)
			return
		}

		requestBody := SpreadsheetRequest{
			SpreadsheetID: "1cg996SvJJwQAU0C7DaFOvnySnH3PsTcDen1NegkFOdk",
			Data:          dataArray,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalf("Error marshaling request: %v", err)
		}

		url := "https://script.google.com/macros/s/AKfycbzZj9blTvFTq371iOzMvTdlA_tgqTCMJYdKEKGrqu_B25kIZYR-_24TGLJfUYgfXIKp/exec"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IssuesByDateResponse{Dates: result})
	}
}
