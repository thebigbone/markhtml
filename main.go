package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>Markdown Preview Tool</title>
	</head>
	<body>
	`
	footer = `
	</body>
	</html>
	`
)

func main() {
	filename := flag.String("file", "", "markdown file")
	preview := flag.Bool("preview", false, "directly preview in browser")
	delete := flag.Bool("delete", false, "delete the converted html file")
	pdf := flag.Bool("pdf", false, "export to pdf file")

	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := run(*filename, *preview, *delete, *pdf)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(file string, previewFile bool, deleteFile bool, exportPDF bool) error {
	input, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	htmlData := parseContent(input)

	if _, err := os.Stat("output_html"); os.IsNotExist(err) {
		err = os.Mkdir("output_html", os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

	}

	temp, err := os.CreateTemp("./output_html", "markdown*.html")
	if err != nil {
		log.Fatal(err)
	}

	temp.Close()
	outName := temp.Name()

	fmt.Println(outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if exportPDF {
		exportPdf(outName)
	}

	if previewFile {
		preview(outName)
	}

	if deleteFile {
		defer os.Remove(outName)
	}
	return nil
}

func parseContent(input []byte) []byte {
	// use blackfridy to convert to html
	output := blackfriday.Run(input)

	// santize the output
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}

func preview(file string) error {
	cmd := ""
	args := []string{}

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "cmd.exe"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// so the string would be "/c start file"
	args = append(args, file)

	err := exec.Command(cmd, args...).Run()
	time.Sleep(2 * time.Second)

	return err
}

// https://github.com/SebastiaanKlippert/go-wkhtmltopdf#usage
func exportPdf(file string) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat("output_pdf"); os.IsNotExist(err) {
		err = os.Mkdir("output_pdf", os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

	}
	pdfg.Dpi.Set(300)
	pdfg.Grayscale.Set(true)

	html_file := wkhtmltopdf.NewPage(file)

	html_file.FooterRight.Set("[page]")
	html_file.FooterFontSize.Set(10)
	html_file.Zoom.Set(0.95)

	pdfg.AddPage(html_file)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	temp, err := os.CreateTemp("./output_pdf", "pdf*.pdf")
	if err != nil {
		log.Fatal(err)
	}

	temp.Close()
	outName := temp.Name()
	err = pdfg.WriteFile(outName)
	if err != nil {
		log.Fatal(err)
	}
}
