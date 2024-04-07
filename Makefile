install_templ:
	go install github.com/a-h/templ/cmd/templ@latest

generate_templ:
	templ generate

build: install_templ generate_templ
	go build -tags netgo -ldflags '-s -w' -o app cmd/main.go

build_admin: install_templ generate_templ
	go build -tags netgo -ldflags '-s -w' -o app cmd/admin/main.go

dev: dev_kill
	air& bunx tailwindcss --watch -m -i ./tailwind.css -o ./static/app.css

dev_kill:
	pkill -f cockaigne/tmp/main || echo "Server was not running ..."

stripe_dev:
	stripe listen --forward-to localhost:3000/api/stripe/webhook
