package dbqueries

import (
	"database/sql"
	"log"

	"github.com/lib/pq"

	"Forum/models"
)

func AddThread(thread *models.Thread) (*models.Thread, error) {
	newThread := &models.Thread{}

	err := db.Get(
		newThread,
		`INSERT INTO thread (thread_author, thread_created, thread_forum, thread_message, thread_slug, thread_title)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`,
		thread.Thread_author, thread.Thread_created, thread.Thread_forum, thread.Thread_message, thread.Thread_slug, thread.Thread_title)
	if err != nil {

		switch err.(*pq.Error).Code {
		case "23505":
			log.Println("23505")
			//
		default:

			log.Println("Error!!!")
		}
		log.Printf("AddThread: %s", err)

		newThread = FindThreadBySlug(*thread.Thread_slug)
		return newThread, err
	}

	_, err = db.Exec(`INSERT INTO forum_users (nickname, forum)
		VALUES ($1, $2)
		ON CONFLICT (nickname, forum) DO NOTHING`,
		thread.Thread_author, thread.Thread_forum)

	if err != nil {

		return newThread, err
	}
	return newThread, nil

}

func EditThread(thread *models.Thread) bool {

	_, err := db.Exec(
		`UPDATE thread 
			   SET thread_title = $1, thread_message = $2, thread_slug = $3
			   WHERE thread_id = $4`,
		thread.Thread_title, thread.Thread_message, thread.Thread_slug, thread.Thread_id)

	if err != nil {

		//switch err.(*pq.Error).Code {
		//case "23505":
		//	log.Println("23505")
		//	//
		//default:
		//
		log.Printf("EditThread: %s", err)
		//}
		return true
	}

	return false
}

func FindThreadBySlug(slug string) *models.Thread {

	findThread := &models.Thread{}

	err := db.Get(
		findThread,
		`SELECT * 
		FROM thread 
		WHERE thread_slug = $1`,
		slug)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Println("Error of find thread by slug")
		return nil
	}

	return findThread
}

func FindThreadById(id int) *models.Thread {

	findThread := &models.Thread{}

	err := db.Get(
		findThread,
		`SELECT * 
		FROM thread 
		WHERE thread_id = $1`,
		id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Println("Error of find thread by id")
		return nil
	}

	return findThread
}

//func FindForumBySlug(slug string) *models.Post_forum {
//	findForum := &models.Post_forum{}
//
//	err := db.Get(
//		findForum,
//		`SELECT posts, slug, threads, title, forum_user
//		FROM forum
//		WHERE slug = $1`,
//		slug)
//
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil
//		}
//
//		log.Println("Error of find forum by slug")
//		return nil
//	}
//
//	return findForum
//}
