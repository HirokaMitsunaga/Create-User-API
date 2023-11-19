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

func escapeSQL(value string) string {
    // 単純なエスケープ例（実際の実装ではより堅牢な方法を検討すること）
    return strings.Replace(value, "'", "''", -1)
}

func main() {
    db, err := sql.Open("mysql", "root:IhVmfFZo8R?@tcp(13.112.0.29:3306)/my_schema4")
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

        _, err = db.Exec(fmt.Sprintf("CREATE USER '%s' IDENTIFIED BY '%s'", 
        escapeSQL(u.Username), escapeSQL(u.Password)))
        if err != nil {
            http.Error(w, "Failed to create user", http.StatusInternalServerError)
            return
        }
    
        privileges := strings.ToUpper(u.Privileges)
        query := fmt.Sprintf("GRANT %s ON `%s`.* TO '%s'", privileges, escapeSQL(u.Schema), escapeSQL(u.Username))

        _, err = db.Exec(query)
        if err != nil {
            http.Error(w, "Failed to grant privileges", http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "User created with privileges: %s on schema %s", privileges, u.Schema)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
