package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed html/app.html
var appHTML string

func uploadHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "missing file")
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "unable to read file")
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "unable to save file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, "unable to save file")
	}

	fmt.Printf("received file: %s\n", file.Filename)

	return c.String(http.StatusOK, "successfully uploaded file")
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, appHTML)
	})
	e.POST("/upload", uploadHandler)

	fmt.Println("Starting web interface on :9125") // todo: port flag

	e.HideBanner = true
	e.HidePort = true
	e.Logger.Fatal(e.Start(":9125"))
}
