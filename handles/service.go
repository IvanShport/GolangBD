package handles

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Forum/dbqueries"
)

func ProcessService(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	pathes := strings.Split(r.URL.Path, "/")[1:]

	if len(pathes) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch pathes[1] {

	case "clear":

		if r.Method == http.MethodPost {
			ClearData(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "status":

		if r.Method == http.MethodGet {
			GetInfoAboutBD(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
func ClearData(w http.ResponseWriter, r *http.Request) {
	err := dbqueries.ClearData()

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetInfoAboutBD(w http.ResponseWriter, r *http.Request) {

	db, err := dbqueries.GetInfoAboutBD()

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("GetInfoAboutBD")
}
