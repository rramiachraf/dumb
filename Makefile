VERSION=`git rev-parse --short HEAD`

gentempl:
	@command -v templ &> /dev/null || go install github.com/a-h/templ/cmd/templ@latest
build:gentempl
	templ generate && go build -ldflags="-X 'github.com/rramiachraf/dumb/data.Version=$(VERSION)' -s -w"
