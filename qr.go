package main

// Adapted from:
// https://github.com/piglig/go-qr/blob/main/tools/generator/cmd/cmd.go

import (
	"bytes"
	"fmt"

	go_qr "github.com/piglig/go-qr"
)

const (
	blackBlock = "\033[40m  \033[0m"
	whiteBlock = "\033[47m  \033[0m"
)

// RenderQR Generate a QR code and print it to stdout
func RenderQR(text string) error {
	qr, err := go_qr.EncodeText(text, go_qr.High)
	if err != nil {
		return err
	}

	fmt.Println(toString(qr))
	return nil
}

func toString(qr *go_qr.QrCode) string {
	buf := bytes.Buffer{}
	border := 4

	for y := -border; y < qr.GetSize()+border; y++ {
		for x := -border; x < qr.GetSize()+border; x++ {
			if !qr.GetModule(x, y) {
				buf.WriteString(blackBlock)
			} else {
				buf.WriteString(whiteBlock)
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}
