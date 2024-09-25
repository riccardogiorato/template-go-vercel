package handler

import (
	"fmt"
	"net/http"
)

func HtmlRendering(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Go HTML Rendering</title>
	</head>
	<body>
		<h1>Hello, this is a HTML document rendered with Go!</h1>
	</body>
	</html>
	`
	fmt.Fprint(w, htmlContent)
}
