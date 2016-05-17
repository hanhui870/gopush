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

	http.Handle("/form", &uploadFormHandle{})

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

type uploadFormHandle struct {

}

func (*uploadFormHandle)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("server", "gopush")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#375EAB">
<title>Go push http form test</title>
</head>
<body>
	<form name="upload" method="post" enctype="multipart/form-data" action="/form">
		<label>点击选择文件: </label>
		<input type="file" name="file">
		<br><label>文件名称: </label><input type="text" name="name">
		<input type="submit" value="提交">
	</form>
</body>
</html>
		`)
	}else {
		w.Header().Add("server", "gopush")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		r.ParseForm()
		r.ParseMultipartForm(16 * 1024 * 1024)
		fmt.Fprintln(w, "Hello world.</br>")
		if f, ok := r.Form["name"]; ok {
			for key, value := range f {
				fmt.Fprintln(w, " param name value:", key, value, "</br>")
			}
		}else {
			fmt.Fprintln(w, " param name needed.")
		}

		//multi forms
		if _, ok := r.Form["name"]; ok {
			for key, value := range r.MultipartForm.Value {
				fmt.Fprintln(w, " param MultipartForm value:", key, value, "</br>")
			}

			for key, value := range r.MultipartForm.File {
				for iter, file := range value {
					fmt.Fprintln(w, " param MultipartForm File:", key, iter, file.Header, file.Filename, "</br>")
				}
			}
		}else {
			fmt.Fprintln(w, " param name needed.")
		}

		fmt.Fprintln(w, " request form:", r.Form)
	}

}

