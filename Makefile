VERSION=`git rev-parse --short HEAD`

build:
	go tool templ generate
	cat ./style/*.css | go tool esbuild --loader=css --minify > ./static/style.css
	go build -ldflags="-X 'github.com/rramiachraf/dumb/data.Version=$(VERSION)' -s -w"
test:
	go test ./... -v
fmt:
	go tool templ fmt .
