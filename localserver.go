package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

func startLocalServer() {
	app := fiber.New()

	app.Static("/", "./Admin")

	app.Post("/host/:filename", host)
	app.Post("/images/upload/:filename", ImgUpload)
	app.Get("/images/get", ImgFetch)
	app.Delete("/images/delete/:filename", ImgDelete)
	app.Post("/forms/create/", FormsCreate)
	app.Post("/forms/insert/:formname", FormsInsert)

	// Listen on localhost only
	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Local server error: %v", err)
	}
}

func host(c *fiber.Ctx) error {
	filename := c.Params("filename")

	// Open the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Failed to get uploaded file")
	}

	// Define path to save the file
	savePath := filepath.Join("./hosting", filename)

	// Save the uploaded file
	inFile, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Failed to open uploaded file")
	}
	defer inFile.Close()

	outFile, err := os.Create(savePath)
	if err != nil {
		return c.Status(500).SendString("Failed to create file on server")
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, inFile); err != nil {
		return c.Status(500).SendString("Failed to save file")
	}

	return c.SendString("File uploaded successfully!")
}

func ImgUpload(c *fiber.Ctx) error {
	filename := c.Params("filename")
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("Failed to get uploaded file")
	}

	savePath := filepath.Join("./Images", filename)
	inFile, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Failed to open uploaded file")
	}
	defer inFile.Close()

	outFile, err := os.Create(savePath)
	if err != nil {
		return c.Status(500).SendString("Failed to create file on server")
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, inFile); err != nil {
		return c.Status(500).SendString("Failed to save file")
	}

	return c.SendString("File uploaded successfully!")
}

func ImgFetch(c *fiber.Ctx) error {
	root := "./Images"
	files, err := os.ReadDir(root)
	if err != nil {
		return c.Status(500).SendString("Failed to read directory")
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() { // Ensure it's a file, not a directory
			fileNames = append(fileNames, file.Name())
		}
	}

	return c.JSON(fileNames)
}

func ImgDelete(c *fiber.Ctx) error {
	root := "./Images"
	filename := c.Params("filename")
	filePath := filepath.Join(root, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).SendString("File not found")
	}

	if err := os.Remove(filePath); err != nil {
		return c.Status(500).SendString("Failed to delete file")
	}

	return c.SendString("File deleted successfully")
}

func FormsCreate(c *fiber.Ctx) error {
	formData := make(map[string][]string)

	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		formData[string(key)] = append(formData[string(key)], string(value))
	})

	var qry string = "Create Table If Not Exists " + formData["tablename"][0] + "("
	var table string = ""
	var idx int = len(formData["column_names[]"])

	for i := 0; i < idx-1; i++ {
		table = table + formData["column_names[]"][i] + " " + formData["column_types[]"][i] + ","
	}
	table = table + formData["column_names[]"][idx-1] + " " + formData["column_types[]"][idx-1]

	qry = qry + table + ");"

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(qry)
	if err != nil {
		panic(err)
	}

	return c.SendStatus(fiber.StatusOK)
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
