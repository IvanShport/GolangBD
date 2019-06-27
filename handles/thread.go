package handles

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Forum/dbqueries"
	"Forum/models"
)

func ProcessThread(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.Path)
	pathes := strings.Split(r.URL.Path, "/")[1:]

	if len(pathes) < 3 {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slugOrId := pathes[1]
	//r.Context().Value("slugOrId") = slugOrId

	switch pathes[2] {

	case "create":

		if r.Method == http.MethodPost {
			AddPost(w, r, slugOrId)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "details":

		if r.Method == http.MethodGet {
			GetInfoAboutThread(w, r, slugOrId)
		} else if r.Method == http.MethodPost {
			EditThread(w, r, slugOrId)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "posts":

		if r.Method == http.MethodGet {
			GetPostsOfThread(w, r, slugOrId)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "vote":

		if r.Method == http.MethodPost {
			VoteForThread(w, r, slugOrId)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func AddPost(w http.ResponseWriter, r *http.Request, slugOrId string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	posts := &[]models.Post{}

	if err := json.Unmarshal(body, posts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newPosts, err := dbqueries.AddPosts(posts, slugOrId)

	if err != nil {

		if err == dbqueries.ErrParentPost {

			error := &models.Error{
				Message: "Parent post was created in another thread",
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

		switch err.(type) {
		case *dbqueries.NotFound:

			error := &models.Error{
				Message: err.Error(),
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
			return
		}

	}

	j, err := json.Marshal(newPosts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(j); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func GetInfoAboutThread(w http.ResponseWriter, r *http.Request, slugOrId string) {

	if err, foundThread := isExistThread(slugOrId); err {

		resp, err := json.Marshal(foundThread)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {

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
}

func EditThread(w http.ResponseWriter, r *http.Request, slugOrId string) {

	if isExist, oldThread := isExistThread(slugOrId); isExist {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newThread := &models.Thread{}

		if err := json.Unmarshal(body, newThread); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newThread.Thread_slug != nil {
			oldThread.Thread_slug = newThread.Thread_slug
		}

		if newThread.Thread_message != "" {
			oldThread.Thread_message = newThread.Thread_message
		}

		if newThread.Thread_title != "" {
			oldThread.Thread_title = newThread.Thread_title
		}

		if err := dbqueries.EditThread(oldThread); err {

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

		j, err := json.Marshal(oldThread)

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

func GetPostsOfThread(w http.ResponseWriter, r *http.Request, slugOrId string) {

	if err, thread := isExistThread(slugOrId); err {

		desc := ""

		if r.URL.Query().Get("desc") == "true" {
			desc = "DESC"
		}

		limit := ""

		if r.URL.Query().Get("limit") != "" {
			limit = r.URL.Query().Get("limit")
		}

		since := ""

		if r.URL.Query().Get("since") != "" {
			since = r.URL.Query().Get("since")
		}

		sort := ""

		if r.URL.Query().Get("sort") != "" {
			sort = r.URL.Query().Get("sort")
		}

		foundPost := dbqueries.FindPostOfThread(thread, desc, limit, since, sort)

		resp, err := json.Marshal(foundPost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
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
}

func VoteForThread(w http.ResponseWriter, r *http.Request, slugOrId string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vote := &models.Vote{}

	if err := json.Unmarshal(body, vote); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newThread, err := dbqueries.VoteForThread(vote, slugOrId)

	if err != nil {
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
		return

	}

	resp, err := json.Marshal(newThread)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func isExistThread(slugOrId string) (bool, *models.Thread) {

	var foundThread *models.Thread

	if slugOrId != "" {

		var number int
		var err error

		if number, err = strconv.Atoi(slugOrId); err != nil {

			foundThread = dbqueries.FindThreadBySlug(slugOrId)

			if foundThread != nil {

				return true, foundThread
			}
		}

		foundThread = dbqueries.FindThreadById(number)

		if foundThread != nil {

			return true, foundThread
		}
	}

	return false, nil
}
