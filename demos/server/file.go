package main

import (
	"net/http"
	"fmt"
	"html"
	"log"
)

func main() {
	// To serve a directory on disk (/tmp) under an alternate URL
	// path (/tmpfiles/), use StripPrefix to modify the request
	// URL's path before the FileServer sees it:
	http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("./runtime/log"))))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		r.ParseMultipartForm(16 * 1024 * 1024)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path), r)
		fmt.Fprintf(w, fmt.Sprintln("Request method:", r.Method, " header:", r.Header, " Params:", r.Form, " PostParams:", r.PostForm, " Host:", r.Host, " Remote:", r.RemoteAddr))

		if f, ok := r.Form["f"]; ok {
			for key, value := range f {
				fmt.Fprintln(w, " param f value:", key, value)
			}
		}else {
			fmt.Fprintln(w, " param f needed.")
		}


	})

	log.Println("Server started...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}




