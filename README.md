### Installation

- to install dependencies, run `go mod tidy`.
- compile it using `go build .`
- run the executable

### Flags

```
-file string
    markdown file to preview
-retain
    delete the converted file
-skip
    directly preview in browser
```

`-file`: specify the markdown file.  
`-retain`: keep the generated file. by default, it is deleted.  
`-skip`: to skip the preview

All the outputs will be generated in `output` folder. The file names are randomized.
