package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

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
	filename := flag.String("file", "", "markdown file to preview")
	preview := flag.Bool("preview", false, "directly preview in browser")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := run(*filename, *preview)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(file string, skipPreview bool) error {
	input, err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	htmlData := parseContent(input)

	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err = os.Mkdir("output", os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

	}

	temp, err := os.CreateTemp("./output", "markdown*.html")
	if err != nil {
		log.Fatal(err)
	}

	temp.Close()
	outName := temp.Name()

	fmt.Println(outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outName)
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

	return exec.Command(cmd, args...).Run()
}
