package dbqueries

import (
	"log"
	"strconv"

	"Forum/models"
)

func VoteForThread(vote *models.Vote, slugOrId string) (*models.Thread, error) {

	newThread := &models.Thread{}

	var id int
	var err error

	if id, err = strconv.Atoi(slugOrId); err != nil {

		err = db.Get(
			&id,
			`SELECT thread_id FROM thread
			WHERE thread_slug = $1`,
			slugOrId)

		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	_, err = db.Exec(
		`INSERT INTO vote VALUES 
		((SELECT nickname FROM user_profile where nickname = $1), $2, $3)
		ON CONFLICT (nickname, thread) DO UPDATE SET voice = $3`,
		vote.Nickname, id, vote.Voice)

	if err != nil {

		log.Printf("VoteForThread: %s", err)
		return nil, err
	}

	newThread = FindThreadById(id)

	return newThread, nil
}
