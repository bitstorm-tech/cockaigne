install/templ:
	go install github.com/a-h/templ/cmd/templ@latest

generate/templ:
	templ generate

build: install/templ generate/templ
	go build -tags netgo -ldflags '-s -w' -o app cmd/main.go

build/admin: install/templ generate/templ
	go build -tags netgo -ldflags '-s -w' -o app cmd/admin/main.go

dev: dev/kill
	air& bunx tailwindcss --watch -m -i ./tailwind.css -o ./static/app.css

dev/kill:
	pkill -f cockaigne/tmp/main || echo "Server was not running ..."