gentempl:
	@command -v templ &> /dev/null || go install github.com/a-h/templ/cmd/templ@latest
build:gentempl
	templ generate && go build -o dumb
