VERSION=`git rev-parse --short HEAD`

gentempl:
	command -v templ &> /dev/null || go install github.com/a-h/templ/cmd/templ@latest
esbuild:
	[ ! -f ./esbuild ] && curl -fsSL https://esbuild.github.io/dl/latest | sh
build:gentempl esbuild
	templ generate
	cat ./style/*.css | ./esbuild --loader=css --minify > ./static/style.css 
	go build -ldflags="-X 'github.com/rramiachraf/dumb/data.Version=$(VERSION)' -s -w"
test:
	go test ./... -v
fmt:
	templ fmt .
