package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
    Username   string `json:"username"`
    Password   string `json:"password"`
    Schema     string `json:"schema"`
    Privileges string `json:"privileges"`
}

func main() {
    db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        var u User
        err := json.NewDecoder(r.Body).Decode(&u)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err = db.Exec("CREATE USER ? IDENTIFIED BY ?", u.Username, u.Password)
        if err != nil {
            http.Error(w, "Failed to create user", http.StatusInternalServerError)
            return
        }

        // 権限リストの検証（例: "SELECT, INSERT"）
        // ...

        privileges := strings.ToUpper(u.Privileges)
        query := fmt.Sprintf("GRANT %s ON %s.* TO ?", privileges, u.Schema)

        _, err = db.Exec(query, u.Username)
        if err != nil {
            http.Error(w, "Failed to grant privileges", http.StatusInternalServerError)
            return
        }
		

        fmt.Fprintf(w, "User created with privileges: %s on schema %s", privileges, u.Schema)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
