package main

import (
	_ "embed"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

//go:embed web/app.html
var appHTML string

//go:embed web/vendor/tailwind.min.js
var tailwindJS string

//go:embed web/vendor/htmx.min.js
var htmxJS string

func p(s string) string {
	return "<p>" + s + "</p>"
}

func uploadHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, p("File is missing"))
	}

	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, p("Unable to read file"))
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, p("Unable to create file on the server"))
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.String(http.StatusInternalServerError, p("Unable to write content"))
	}

	log.Printf("received file: %s\n", file.Filename)

	return c.String(http.StatusOK, p("Successfully uploaded "+file.Filename))
}

// GetOutboundIP Ping Cloudflare to get local IP address
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func StartServer(addr string) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, appHTML)
	})
	e.POST("/upload", uploadHandler)
	e.GET("/tailwind.min.js", func(c echo.Context) error {
		return c.HTML(http.StatusOK, tailwindJS)
	})
	e.GET("/htmx.min.js", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmxJS)
	})

	e.HideBanner = true
	e.HidePort = true
	e.Logger.Fatal(e.Start(addr))
}

func main() {
	port := flag.Int("port", 9125, "Port number to listen on")
	qr := flag.Bool("qr", true, "Generate a QR code to connect devices")
	flag.Parse()

	ip := GetOutboundIP()
	addr := ":" + strconv.Itoa(*port)

	log.Println("starting web interface on " + ip + addr)
	if *qr {
		err := RenderQR("http://" + ip + addr)
		if err != nil {
			log.Fatal("unable to generate qr code")
		}
	}

	StartServer(addr)
}
