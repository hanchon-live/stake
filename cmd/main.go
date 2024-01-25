package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/hanchon-live/stake/src/components"
)

func main() {
	component := components.Body()

	http.Handle("/", templ.Handler(component))

	fs := http.FileServer(http.Dir("./public/assets"))
	http.Handle("/public/assets/", http.StripPrefix("/public/assets/", fs))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
