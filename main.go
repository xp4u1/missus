package main

import (
	_ "embed"
	"flag"
	"log"
	"net"
	"strconv"
)

// GetOutboundIP Ping Cloudflare to get local IP address
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func main() {
	port := flag.Int("port", 9125, "Port number to listen on")
	qr := flag.Bool("qr", true, "Generate a QR code to connect devices")
	flag.Parse()

	ip, err := GetOutboundIP()
	if err != nil {
		log.Println("unable to determine local address")
		ip = "localhost"
		*qr = false
	}

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
