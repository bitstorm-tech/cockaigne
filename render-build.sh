go install github.com/a-h/templ/cmd/templ@latest && 
templ generate && 
go build -tags netgo -ldflags '-s -w' -o app cmd/main.go