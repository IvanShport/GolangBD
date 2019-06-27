package handles

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"Forum/dbqueries"
	"Forum/models"
)

func ProcessPost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	pathes := strings.Split(r.URL.Path, "/")[2:]

	if len(pathes) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := pathes[1]
	//r.Context().Value("id") = id

	if pathes[2] == "details" {

		if r.Method == http.MethodPost {

			EditPost(w, r, id)

		} else if r.Method == http.MethodGet {

			GetInfoAboutPost(w, r, id)

		} else {

			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {

		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func EditPost(w http.ResponseWriter, r *http.Request, id string) {
	if isExist, oldPost := isExistPost(id); isExist {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newPost := &models.Post{}

		if err := json.Unmarshal(body, newPost); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newPost.Post_message != "" {
			oldPost.Post_message = newPost.Post_message
		}

		if err := dbqueries.EditPost(oldPost); err {

			error := &models.Error{
				Message: "Can't find user with id #42\n",
			}

			resp, err := json.Marshal(error)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusConflict)

			if _, err := w.Write(resp); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			return
		}

		j, err := json.Marshal(oldPost)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(j); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	error := &models.Error{
		Message: "Can't find user with id #42\n",
	}

	resp, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetInfoAboutPost(w http.ResponseWriter, r *http.Request, id string) {

	args := make(map[string]bool, 3)
	if q, ok := r.URL.Query()["related"]; ok {
		params := strings.Split(q[0], ",")

		for _, value := range params {
			args[value] = true
		}
	}

	log.Println(args)
	if foundPost, err := dbqueries.FindPostById(id, args); err != nil {

		error := &models.Error{
			Message: "Can't find user with id #42\n",
		}

		resp, err := json.Marshal(error)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)

		if _, err := w.Write(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {

		resp, err := json.Marshal(foundPost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func isExistPost(id string) (bool, *models.Post) {

	var foundPost *models.Post

	if id != "" {

		foundPost = dbqueries.FindPost(id)

		if foundPost != nil {

			return true, foundPost
		}

	}

	return false, nil
}
