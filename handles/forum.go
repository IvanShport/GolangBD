package handles

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"Forum/dbqueries"
	"Forum/models"
)

func ProcessForum(w http.ResponseWriter, r *http.Request) {

	pathes := strings.Split(r.URL.Path, "/")[2:]

	if len(pathes) == 2 {

		if pathes[1] == "create" {

			if r.Method == http.MethodPost {
				AddForum(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		} else {

			w.WriteHeader(http.StatusNotFound)
		}

		return

	}

	if len(pathes) < 3 {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slug := pathes[1]
	//r.Context().Value("slug") = slug

	switch pathes[2] {

	case "create":

		if r.Method == http.MethodPost {
			AddThread(w, r, slug)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "details":

		if r.Method == http.MethodGet {
			GetInfoAboutForum(w, r, slug)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "threads":

		if r.Method == http.MethodGet {
			GetThreadsOfForum(w, r, slug)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	case "users":

		if r.Method == http.MethodGet {
			GetUserOfForum(w, r, slug)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func AddForum(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forum := &models.Forum{}

	if err := json.Unmarshal(body, forum); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if statusUser, foundUsers := isExistUser(forum.Forum_user, ""); statusUser {

		if statusForum, foundForum := isExistForum(forum.Forum_slug); statusForum {

			resp, err := json.Marshal(foundForum)

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

		forum.Forum_user = foundUsers[0].Nickname
		dbqueries.AddForum(forum)

		j, err := json.Marshal(forum)

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

func AddThread(w http.ResponseWriter, r *http.Request, slug string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thread := &models.Thread{}

	if err := json.Unmarshal(body, thread); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if statusUser, foundUsers := isExistUser(thread.Thread_author, ""); statusUser {

		if statusForum, foundForum := isExistForum(slug); statusForum {

			thread.Thread_author = foundUsers[0].Nickname
			thread.Thread_forum = foundForum.Forum_slug
			newThread, err := dbqueries.AddThread(thread)

			if err != nil {

				resp, err := json.Marshal(newThread)

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

			j, err := json.Marshal(newThread)

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

func GetInfoAboutForum(w http.ResponseWriter, r *http.Request, slug string) {

	if err, foundForum := isExistForum(slug); err {

		resp, err := json.Marshal(foundForum)
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

func GetThreadsOfForum(w http.ResponseWriter, r *http.Request, slug string) {

	if err, _ := isExistForum(slug); err {

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

		foundThreads := dbqueries.FindThreadsOfForum(slug, desc, limit, since)

		resp, err := json.Marshal(foundThreads)
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

func GetUserOfForum(w http.ResponseWriter, r *http.Request, slug string) {

	if err, _ := isExistForum(slug); err {

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

		foundUsers, err := dbqueries.FindUsersOfForum(slug, desc, limit, since)

		resp, err := json.Marshal(foundUsers)
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

func isExistForum(slug string) (bool, *models.Forum) {

	var foundForumBySlug *models.Forum

	if slug != "" {

		foundForumBySlug = dbqueries.FindForumBySlug(slug)

		if foundForumBySlug != nil {

			return true, foundForumBySlug
		}
	}

	return false, nil
}
