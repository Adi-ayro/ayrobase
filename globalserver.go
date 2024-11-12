package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

func startGlobalServer() {
	app := fiber.New()

	app.Post("/api/forms/insert/:formname", FormsInsert)
	app.Static("/images", "./Images")

	app.Use(func(c *fiber.Ctx) error {
		db, err := sql.Open("sqlite3", "./analytics.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		qry := `CREATE TABLE IF NOT EXISTS Analytics (
			ip TEXT,
			timestamp DATETIME,
			referer TEXT
		)`
		_, err = db.Exec(qry)
		if err != nil {
			log.Fatalf("Error inserting into Analytics: %s", err)
		}

		insertAnalytics(db, c.IP(), time.Now(), c.Get("Referer", "unknown"), c.Path())
		outputAllAnalytics(db)
		return c.Next()
	})

	app.Static("/", "./Hosting")

	// Listen on all interfaces (0.0.0.0)
	if err := app.Listen("127.0.0.1:4000"); err != nil {
		log.Fatalf("Global server error: %v", err)
	}
}

func FormsInsert(c *fiber.Ctx) error {
	formname := c.Params("formname")
	formData := make(map[string]string)

	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		formData[string(key)] = string(value)
	})

	columns := ""
	placeholders := ""
	values := make([]interface{}, 0, len(formData))

	for key, value := range formData {
		columns += key + ", "
		placeholders += "?, "
		values = append(values, value)
	}
	columns = columns[:len(columns)-2]
	placeholders = placeholders[:len(placeholders)-2]

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO "+formname+" (%s) VALUES (%s)", columns, placeholders)
	_, err = db.Exec(query, values...)
	if err != nil {
		panic(err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func insertAnalytics(db *sql.DB, ip string, timestamp time.Time, referer string, path string) {
	insertStmt := `INSERT INTO Analytics (ip, timestamp, referer, path) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertStmt, ip, timestamp, referer, path)
	if err != nil {
		log.Fatalf("Error inserting into Analytics: %s", err)
	}
}

func outputAllAnalytics(db *sql.DB) {
	rows, err := db.Query("SELECT ip, timestamp, referer, path FROM Analytics")
	if err != nil {
		log.Fatalf("Error querying Analytics table: %s", err)
	}
	defer rows.Close()

	fmt.Println("\nAnalytics Table Data:")
	for rows.Next() {
		var ip string
		var timestamp time.Time
		var referer string
		var path string
		err := rows.Scan(&ip, &timestamp, &referer, &path)
		if err != nil {
			log.Fatalf("Error scanning row: %s", err)
		}
		fmt.Printf("IP: %s, Timestamp: %s, Referer: %s, Path: %s\n", ip, timestamp, referer, path)
	}
}
