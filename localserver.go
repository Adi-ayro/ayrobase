package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

func startLocalServer() {
	engine := html.New("./Admin", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Post("/host", Host)
	app.Post("/newhost", ChangeHost)
	app.Post("/content/upload/", ImgUpload)
	app.Get("/content/get", ImgFetch)
	app.Delete("/content/delete/:filename", ImgDelete)
	app.Post("/forms/create/", FormsCreate)
	app.Post("/forms/insert/:formname", FormsInsert)

	app.Get("/admin/content", AdminContent)
	app.Get("/admin/hosting", AdminHosting)
	app.Get("/admin/analytics", AdminAnalytics)
	app.Static("/admin/forms", "./Admin/forms.html")

	// Listen on localhost only
	if err := app.Listen("127.0.0.1:3000"); err != nil {
		log.Fatalf("Local server error: %v", err)
	}
}

func clearDirectory(path string) error {
	// Read all entries in the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Loop through each entry and remove it
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}
	return nil
}

func ChangeHost(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	fmt.Println("Started Change Host")

	// Open file for reading
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	h := clearDirectory("./Hosting")
	if h != nil {
		log.Println("Error clearing directory:", h)
	} else {
		log.Println("Directory cleared successfully.")
	}

	// Read the zip file into a zip.Reader
	zipReader, err := zip.NewReader(f, file.Size)
	if err != nil {
		return err
	}

	// Create the target directory if it doesn’t exist
	targetDir := "./Hosting"
	os.MkdirAll(targetDir, os.ModePerm)

	// Detect the top-level directory if it exists
	var rootDir string
	for _, zipFile := range zipReader.File {
		parts := strings.SplitN(zipFile.Name, "/", 2)
		if len(parts) > 1 {
			rootDir = parts[0] + "/"
			break
		}
	}

	// Loop through each file and directory in the zip
	for _, zipFile := range zipReader.File {
		// Remove the root directory from the target path if it exists
		targetPath := zipFile.Name
		if rootDir != "" && strings.HasPrefix(zipFile.Name, rootDir) {
			targetPath = strings.TrimPrefix(zipFile.Name, rootDir)
		}
		targetPath = filepath.Join(targetDir, targetPath)

		if zipFile.FileInfo().IsDir() {
			// Create directory if it doesn’t exist
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		// Create all necessary parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// Create and copy file contents
		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return c.Redirect("/admin/hosting")
}

func Host(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Open file for reading
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the zip file into a zip.Reader
	zipReader, err := zip.NewReader(f, file.Size)
	if err != nil {
		return err
	}

	// Create the target directory if it doesn’t exist
	targetDir := "./Hosting"
	os.MkdirAll(targetDir, os.ModePerm)

	// Detect the top-level directory if it exists
	var rootDir string
	for _, zipFile := range zipReader.File {
		parts := strings.SplitN(zipFile.Name, "/", 2)
		if len(parts) > 1 {
			rootDir = parts[0] + "/"
			break
		}
	}

	// Loop through each file and directory in the zip
	for _, zipFile := range zipReader.File {
		// Remove the root directory from the target path if it exists
		targetPath := zipFile.Name
		if rootDir != "" && strings.HasPrefix(zipFile.Name, rootDir) {
			targetPath = strings.TrimPrefix(zipFile.Name, rootDir)
		}
		targetPath = filepath.Join(targetDir, targetPath)

		if zipFile.FileInfo().IsDir() {
			// Create directory if it doesn’t exist
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		// Create all necessary parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// Create and copy file contents
		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return c.Redirect("/admin/hosting")
}

func ImgUpload(c *fiber.Ctx) error {
	filename := c.FormValue("filename")
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

	return c.Redirect("/admin/content")
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

	return c.SendStatus(fiber.StatusOK)
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

func AdminContent(c *fiber.Ctx) error {
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

	return c.Render("content", fiber.Map{
		"files": fileNames,
	})
}

func AdminHosting(c *fiber.Ctx) error {
	f, err := os.Open("./Hosting")
	if err != nil {
		return c.SendStatus(fiber.ErrBadGateway.Code)
	}

	res := false
	_, file := f.ReadDir(1)
	if file == io.EOF {
		res = false // dir is empty
	} else {
		res = true
	}

	return c.Render("host", fiber.Map{
		"host": res,
	})
}

type PageData struct {
	Path  string
	Count int
}

type RefererData struct {
	Referer string
	Count   int
}

type DailyData struct {
	Date  string
	Count int
}

func AdminAnalytics(c *fiber.Ctx) error {
	db, err := sql.Open("sqlite3", "./analytics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pageRows, err := db.Query("SELECT path, COUNT(*) as count FROM Analytics GROUP BY path")
	if err != nil {
		return err
	}
	defer pageRows.Close()

	var pages []PageData
	for pageRows.Next() {
		var page PageData
		if err := pageRows.Scan(&page.Path, &page.Count); err != nil {
			return err
		}
		pages = append(pages, page)
	}

	// Query for referer counts
	refererRows, err := db.Query("SELECT COALESCE(referer, 'None') AS referer, COUNT(*) as count FROM Analytics GROUP BY referer")
	if err != nil {
		return err
	}
	defer refererRows.Close()

	var referers []RefererData
	for refererRows.Next() {
		var ref RefererData
		if err := refererRows.Scan(&ref.Referer, &ref.Count); err != nil {
			return err
		}
		referers = append(referers, ref)
	}

	// Query for daily unique IP count
	rows, err := db.Query(`
		SELECT DATE(timestamp) AS day, COUNT(DISTINCT ip) as unique_ips 
		FROM Analytics 
		GROUP BY day
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Fetch results into a slice
	var dailyData []DailyData
	for rows.Next() {
		var data DailyData
		if err := rows.Scan(&data.Date, &data.Count); err != nil {
			log.Fatal(err)
		}
		dailyData = append(dailyData, data)
	}

	dates := []string{}
	counts := []int{}
	for _, data := range dailyData {
		dates = append(dates, data.Date)
		counts = append(counts, data.Count)
	}

	// Pass both data sets to the HTML template
	return c.Render("analytics", fiber.Map{
		"pages":    pages,
		"referers": referers,
		"Dates":    dates,
		"Counts":   counts,
	})
}
