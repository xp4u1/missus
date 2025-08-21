package main

import (
	_ "embed"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
)

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
