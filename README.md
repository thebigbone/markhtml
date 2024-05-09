### Installation

- to install dependencies, run `go mod tidy`.
- compile it using `go build .`
- run the executable

### Flags

```
-delete
    delete the converted html file
-file string
    markdown file
-pdf
    export to pdf file
-preview
    directly preview in browser
```

`-file`: specify the markdown file.  
`-delete`: delete the converted html file.
`-pdf`: export it as pdf.
`-preview`: directly preview in browser.

All the outputs will be generated in `output` folder. The file names are randomized.
