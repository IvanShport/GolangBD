package main

import (
	"net/http"

	"Forum/dbqueries"
	"Forum/handles"
)

func AppJsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func main() {
	dbqueries.InitDB()

	http.Handle("/forum/", AppJsonMiddleware(http.HandlerFunc(handles.ProcessForum)))

	http.Handle("/post/", AppJsonMiddleware(http.HandlerFunc(handles.ProcessPost)))

	http.Handle("/service/", AppJsonMiddleware(http.HandlerFunc(handles.ProcessService)))

	http.Handle("/thread/", AppJsonMiddleware(http.HandlerFunc(handles.ProcessThread)))

	http.Handle("/user/", AppJsonMiddleware(http.HandlerFunc(handles.ProcessUser)))

	http.ListenAndServe(":5000", nil)
}
