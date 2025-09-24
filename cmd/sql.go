/*
Copyright © 2025 Marcel Zapf
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "Wait for SQL connection",
	Long: `wait for SQL connection until available
	
Simple example:
  waitfor sql -u root -p mysecretpassword -s mariadb.mydatabase.cluster.local -d mydb

Example, if you have a non-standard port, set it with -P, default is 3306:
  waitfor sql -u root -p mysecretpassword -s mariadb.mydatabase.cluster.local -P 3307 -d mydb
`,
	Run: func(cmd *cobra.Command, args []string) {

		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		server, _ := cmd.Flags().GetString("server")
		port, _ := cmd.Flags().GetString("port")
		database, _ := cmd.Flags().GetString("database")
		retries, _ := cmd.Flags().GetInt("retries")

		// check if required flags are set
		if user == "" || password == "" || server == "" || database == "" {
			log.Fatalf("Error: user, password, server and database are required")
		}

		// build dsn
		dsn := user + ":" + password + "@tcp(" + server + ":" + port + ")/" + database

		// Try to connect with db
		db, err := waitForDB(dsn, retries, time.Second*time.Duration(timer))
		if err != nil {
			log.Fatalf("DB connection failed: %v", err)
		}
		defer db.Close()

		log.Println("Connection established!")
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)
	sqlCmd.Flags().StringP("user", "u", "", "Database user")
	sqlCmd.Flags().StringP("password", "p", "", "Database password")
	sqlCmd.Flags().StringP("server", "s", "", "Database server")
	sqlCmd.Flags().StringP("port", "P", "3306", "Database port")
	sqlCmd.Flags().StringP("database", "d", "", "Database name")
	sqlCmd.Flags().IntP("retries", "r", 10, "Number of retries")
}

func waitForDB(dsn string, retries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < retries; i++ {
		// Open does not establish a connection immediately, it just prepares the handle
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Error opening DB handle: %v", err)
		} else {
			// Ping actually checks if the DB is reachable
			err = db.Ping()
			if err == nil {
				return db, nil // success
			}
			log.Printf("DB not ready yet: %v", err)
		}

		// Wait before trying again
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("could not connect after %d retries: %w", retries, err)
}
