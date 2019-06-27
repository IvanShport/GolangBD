package dbqueries

import (
	"database/sql"
	"log"
	"strings"

	"github.com/lib/pq"

	"Forum/models"
)

func AddForum(forum *models.Forum) {
	_, err := db.Exec(
		`INSERT INTO forum (forum_slug, forum_title, forum_user)
		VALUES ($1, $2, $3)`,
		forum.Forum_slug, forum.Forum_title, forum.Forum_user)
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

func FindThreadsOfForum(slugOfForum string, desc string, limit string, since string) *[]models.Thread {

	foundThreads := &[]models.Thread{}

	query := strings.Builder{}

	query.WriteString("SELECT * from thread " +
		"WHERE thread_forum = $1 ")

	if since != "" {

		if desc != "" {

			query.WriteString("AND thread_created <= $2 ")

		} else {

			query.WriteString("AND thread_created >= $2 ")
		}
	}

	if desc != "" {

		query.WriteString("ORDER BY thread_created DESC ")
	} else {

		query.WriteString("ORDER BY thread_created ")
	}

	if limit != "" {
		query.WriteString("LIMIT " + limit)
	}

	log.Println(query.String())

	var err error

	if since != "" {

		err = db.Select(
			foundThreads,
			query.String(),
			slugOfForum, since)

	} else {

		err = db.Select(
			foundThreads,
			query.String(),
			slugOfForum)

	}

	if err != nil {

		if err == sql.ErrNoRows {
			return nil
		}

		log.Printf("FindThreadsOfForum: %s", err)
	}

	return foundThreads
}

func FindForumBySlug(slug string) *models.Forum {
	findForum := &models.Forum{}

	err := db.Get(
		findForum,
		`SELECT posts, forum_slug, threads, forum_title, forum_user 
		FROM forum 
		WHERE forum_slug = $1`,
		slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Printf("Error of find forum by slug: %s", err)
		return nil
	}

	return findForum
}
