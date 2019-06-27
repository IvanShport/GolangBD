package dbqueries

import (
	"database/sql"
	"log"
	"strings"

	"github.com/lib/pq"

	"Forum/models"
)

func AddUser(u *models.User) {
	//q := "INSERT INTO user_profile VALUES ($1, $2, $3)"
	//
	_, err := db.Exec(`INSERT INTO user_profile VALUES ($1, $2, $3, $4)`,
		u.Nickname, u.Email, u.Fullname, u.About)
	if err != nil {

		switch err.(*pq.Error).Code {
		case "23505":
			log.Println("23505")
			//
		default:

			log.Println("Error!!!")
		}
	}

	log.Println("Ok!")

}

func EditUser(u *models.User) bool {

	_, err := db.Exec(
		`UPDATE user_profile 
			   SET email = $1, fullname = $2, about = $3
			   WHERE nickname = $4`,
		u.Email, u.Fullname, u.About, u.Nickname)

	if err != nil {

		//switch err.(*pq.Error).Code {
		//case "23505":
		//	log.Println("23505")
		//	//
		//default:
		//
		log.Println("Error!!!")
		//}
		return true
	}

	return false
}

func FindUsersOfForum(slugOfForum string, desc string, limit string, since string) (*[]models.User, error) {

	foundUsers := &[]models.User{}

	query := strings.Builder{}

	query.WriteString("SELECT u.nickname, u.email, u.fullname, u.about from user_profile u " +
		"JOIN forum_users fu ON u.nickname = fu.nickname " +
		"WHERE fu.forum = $1 ")

	if since != "" {

		if desc != "" {

			query.WriteString("AND fu.nickname < $2 COLLATE \"POSIX\"")

		} else {

			query.WriteString("AND fu.nickname > $2 COLLATE \"POSIX\"")
		}
	}

	if desc != "" {

		query.WriteString("ORDER BY fu.nickname COLLATE \"POSIX\" DESC ")
	} else {

		query.WriteString("ORDER BY fu.nickname COLLATE \"POSIX\" ")
	}

	if limit != "" {
		query.WriteString("LIMIT " + limit)
	}

	log.Println(query.String())

	var err error

	if since != "" {

		err = db.Select(
			foundUsers,
			query.String(),
			slugOfForum, since)

	} else {

		err = db.Select(
			foundUsers,
			query.String(),
			slugOfForum)

	}

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Printf("FindUsersOfForum: %s", err)
		return nil, err
	}

	return foundUsers, nil
}

func FindUserByNickname(nickname string) *models.User {
	findUser := &models.User{}

	err := db.Get(
		findUser,
		`SELECT nickname, email, fullname, about 
			   FROM user_profile 
			   WHERE nickname = $1`,
		nickname)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Println("Error of find user by nickname")
		return nil
	}

	return findUser
}

func FindUserByEmail(email string) *models.User {
	findUser := &models.User{}

	err := db.Get(
		findUser,
		`SELECT nickname, email, fullname, about 
			   FROM user_profile 
			   WHERE email = $1`,
		email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Println("Error of find user by email")
		return nil
	}

	return findUser
}
