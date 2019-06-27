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

func ProcessUser(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	pathes := strings.Split(r.URL.Path, "/")[1:]
	if len(pathes) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nickname := pathes[1]
	//r.Context().Value("nickname") = nickname

	switch pathes[2] {
	case "create":
		if r.Method == http.MethodPost {
			AddUser(w, r, nickname)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "profile":

		if r.Method == http.MethodGet {
			GetInfoAboutUser(w, r, nickname)
		} else if r.Method == http.MethodPost {
			EditUser(w, r, nickname)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func AddUser(w http.ResponseWriter, r *http.Request, nickname string) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u := &models.User{
		Nickname: nickname,
	}

	if err := json.Unmarshal(body, u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("user is %+v", u)

	if isExist, oldUsers := isExistUser(u.Nickname, u.Email); isExist {

		resp, err := json.Marshal(oldUsers)

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

	dbqueries.AddUser(u)

	j, err := json.Marshal(u)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(j); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func GetInfoAboutUser(w http.ResponseWriter, r *http.Request, nickname string) {

	if isExist, foundUsers := isExistUser(nickname, ""); isExist {

		resp, err := json.Marshal(foundUsers[0])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resp); err != nil {
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

func EditUser(w http.ResponseWriter, r *http.Request, nickname string) {

	if isExist, oldUsers := isExistUser(nickname, ""); isExist {

		oldUser := oldUsers[0]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newUser := &models.User{}

		if err := json.Unmarshal(body, newUser); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newUser.About != "" {
			oldUser.About = newUser.About
		}

		if newUser.Fullname != "" {
			oldUser.Fullname = newUser.Fullname
		}

		if newUser.Email != "" {
			oldUser.Email = newUser.Email
		}

		if err := dbqueries.EditUser(oldUser); err {

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

		j, err := json.Marshal(oldUser)

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

func isExistUser(nickname string, email string) (bool, []*models.User) {

	var foundUserByNickname *models.User
	var foundUserByEmail *models.User

	resultUsers := make([]*models.User, 0, 5)

	if nickname != "" {

		foundUserByNickname = dbqueries.FindUserByNickname(nickname)

		if foundUserByNickname != nil {

			resultUsers = append(resultUsers, foundUserByNickname)
		}
	}

	if email != "" {

		foundUserByEmail = dbqueries.FindUserByEmail(email)

		if foundUserByEmail != nil {

			if len(resultUsers) > 0 {

				if *resultUsers[0] != *foundUserByEmail {

					resultUsers = append(resultUsers, foundUserByEmail)
				}
			} else {

				resultUsers = append(resultUsers, foundUserByEmail)
			}
		}
	}

	if len(resultUsers) != 0 {

		return true, resultUsers
	}

	return false, nil
}
