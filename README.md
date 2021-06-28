## Start

Start http server:

`python3 -m http.server`

## Edit

### How to compile `.go` file to `wasm` module:

Make your changes in `main.go` file, then type in console

```GOOS=js GOARCH=wasm go build -o file_name.wasm```

after that you will get the file `file_name.wasm` you have to to import into `index.html` file

### Sources

https://github.com/anthonynsimon/bild