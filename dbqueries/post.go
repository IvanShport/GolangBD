package dbqueries

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	"Forum/models"
)

func AddPosts(posts *[]models.Post, slugOrId string) (*[]models.Post, error) {
	newPosts := &[]models.Post{}

	curThread := &models.Thread{}

	if id, err := strconv.Atoi(slugOrId); err != nil {

		curThread = FindThreadBySlug(slugOrId)

	} else {

		curThread = FindThreadById(id)
	}

	if curThread == nil {
		return nil, &NotFound{"Thread", slugOrId}
	}

	if len(*posts) == 0 {
		return newPosts, nil
	}

	now := time.Now()

	stmtForumUser, err := db.Prepare("INSERT INTO forum_users (nickname, forum) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (nickname, forum) DO NOTHING")
	if err != nil {

		return nil, err
	}
	defer stmtForumUser.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO post (post_author, post_created, post_forum, post_id, post_message, parent, post_thread, path, founder) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING post_id, post_created")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	stmtUser, err := tx.Prepare("SELECT nickname FROM user_profile WHERE nickname = $1")
	if err != nil {
		return nil, err
	}
	defer stmtUser.Close()

	stmtId, err := tx.Prepare("SELECT nextval(pg_get_serial_sequence('post', 'post_id'))")

	if err != nil {
		return nil, err
	}
	defer stmtId.Close()

	stmtPost, err := tx.Prepare("SELECT post_thread, path FROM post WHERE post_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmtPost.Close()

	for key, value := range *posts {

		err = stmtUser.QueryRow(value.Post_author).Scan(&(*posts)[key].Post_author)
		if err != nil {

			if err == sql.ErrNoRows {
				return nil, &NotFound{"User", value.Post_author}
			}

			return nil, err
		}
		(*posts)[key].Post_thread = curThread.Thread_id
		(*posts)[key].Post_forum = curThread.Thread_forum

		if value.Parent != 0 {
			parent := models.Post{}

			err = stmtPost.QueryRow(value.Parent).Scan(&parent.Post_thread, pq.Array(&parent.Path))

			if err != nil {

				if err == sql.ErrNoRows {
					return nil, ErrParentPost
				}

				log.Printf("Error of find parent of post: %s", err)
				return nil, err
			}

			if parent.Post_thread != curThread.Thread_id {
				return nil, ErrParentPost
			}

			(*posts)[key].Path = parent.Path
		}

		err := stmtId.QueryRow().Scan(&(*posts)[key].Post_id)

		if err != nil {

			return nil, err
		}

		(*posts)[key].Path = append((*posts)[key].Path, int64((*posts)[key].Post_id))

		(*posts)[key].Founder = int((*posts)[key].Path[0])

		err = stmt.QueryRow(value.Post_author, now, curThread.Thread_forum, (*posts)[key].Post_id, value.Post_message, value.Parent, curThread.Thread_id, pq.Array((*posts)[key].Path), (*posts)[key].Founder).Scan(&(*posts)[key].Post_id, &(*posts)[key].Post_created)

		if err != nil {

			log.Printf("AddPosts: %s", err)
			return nil, err
		}

		_, err = stmtForumUser.Exec(value.Post_author, curThread.Thread_forum)
		if err != nil {

			return nil, err
		}

	}

	_, err = tx.Exec(`UPDATE forum SET posts = posts + $1 WHERE forum_slug = $2`,
		len(*posts), curThread.Thread_forum)

	if err != nil {

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil

}

func EditPost(post *models.Post) bool {

	err := db.Get(&post.IsEdited,
		`UPDATE post 
			   SET post_message = $1,
			   isEdited = CASE WHEN $1 <> (SELECT post_message FROM post WHERE post_id = $2)
			   THEN TRUE
			   ELSE FALSE 
			   END 
			   WHERE post_id = $2
			   RETURNING isEdited`,
		post.Post_message, post.Post_id)

	if err != nil {

		//switch err.(*pq.Error).Code {
		//case "23505":
		//	log.Println("23505")
		//	//
		//default:
		//
		log.Printf("EditPost: %s", err)
		//}
		return true
	}

	return false
}

func FindPostOfThread(thread *models.Thread, desc string, limit string, since string, sort string) *[]models.Post {

	foundPosts := &[]models.Post{}

	query := strings.Builder{}

	query.WriteString("SELECT p.post_author, p.post_created, p.post_forum, p.post_id, p.isEdited, p.post_message, p.parent, p.post_thread " +
		"from post p ")

	if sort == "tree" {

		if since != "" {

			query.WriteString("JOIN post sp ON sp.post_id = $2 WHERE p.path ")

			if desc != "" {

				query.WriteString("< sp.path ")

			} else {

				query.WriteString("> sp.path ")
			}

			query.WriteString("AND p.post_thread = $1")

		} else {

			query.WriteString("WHERE p.post_thread = $1")
		}

		query.WriteString(" ORDER BY p.path")

		if desc != "" {

			query.WriteString(" DESC")

		}
		if limit != "" {

			query.WriteString(" LIMIT " + limit)
		}

	} else if sort == "parent_tree" {

		query.WriteString("WHERE p.founder IN (SELECT p.post_id FROM post p ")

		if since != "" {

			query.WriteString("JOIN post sp ON sp.post_id = $2 ")

			if desc != "" {

				query.WriteString("WHERE p.founder < sp.founder AND p.post_thread = $1 AND p.parent = 0")

			} else {

				query.WriteString("WHERE p.path > sp.path AND p.post_thread = $1 AND p.parent = 0")
			}
		} else {

			query.WriteString("WHERE p.post_thread = $1 AND p.parent = 0 ")
		}

		query.WriteString("ORDER BY p.founder")

		if desc != "" {

			query.WriteString(" DESC")
		}
		if limit != "" {

			query.WriteString(" LIMIT " + limit)
		}

		query.WriteString(`) ORDER BY p.founder`)

		if desc != "" {

			query.WriteString(" DESC")
		}

		query.WriteString(", p.path")
	} else {

		query.WriteString("WHERE p.post_thread = $1 ")

		if since != "" {

			if desc != "" {

				query.WriteString("AND p.post_id < $2 ")

			} else {

				query.WriteString("AND p.post_id > $2 ")
			}
		}

		if desc != "" {

			query.WriteString("ORDER BY p.post_created DESC, p.post_id DESC ")
		} else {

			query.WriteString("ORDER BY p.post_created, p.post_id ")
		}

		if limit != "" {
			query.WriteString("LIMIT " + limit)
		}
	}

	log.Println(query.String())

	var err error

	if since != "" {

		err = db.Select(
			foundPosts,
			query.String(),
			thread.Thread_id, since)

	} else {

		err = db.Select(
			foundPosts,
			query.String(),
			thread.Thread_id)

	}

	if err != nil {

		if err == sql.ErrNoRows {
			return nil
		}

		log.Printf("FindPostOfThread: %s", err)
	}

	return foundPosts
}

func FindPostById(id string, args map[string]bool) (*models.PostInfo, error) {

	findPost := &models.PostInfo{}

	query := strings.Builder{}
	query.WriteString("SELECT p.post_author, p.post_created, p.post_forum, p.post_id, p.isEdited, p.post_message, p.parent, p.post_thread")

	if _, isExist := args["user"]; isExist {
		query.WriteString(", u.email, u.fullname, u.about")
	}

	if _, isExist := args["thread"]; isExist {
		query.WriteString(", t.thread_author, t.thread_created, t.thread_message, t.thread_slug, t.thread_title, t.votes")
	}

	if _, isExist := args["forum"]; isExist {
		query.WriteString(", f.posts, f.threads, f.forum_title, f.forum_user")
	}

	query.WriteString(" FROM post p")

	if _, isExist := args["user"]; isExist {
		query.WriteString(" JOIN user_profile u ON p.post_author = u.nickname")
	}

	if _, isExist := args["thread"]; isExist {
		query.WriteString(" JOIN thread t ON p.post_thread = t.thread_id")
	}

	if _, isExist := args["forum"]; isExist {
		query.WriteString(" JOIN forum f ON p.post_forum = f.forum_slug")
	}

	query.WriteString(" WHERE p.post_id = $1")

	container := &models.PostInfoContainer{}

	log.Println(query.String())
	err := db.QueryRowx(
		query.String(),
		id).StructScan(container)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		log.Printf("Error of find post by id: %s", err)
		return nil, err
	}

	log.Println("Container")
	log.Println(container.Post)
	log.Println(container.User)
	log.Println(container.Thread)
	log.Println(container.Forum)

	findPost.Post = &models.Post{
		Post_author:  container.Post.Post_author,
		Post_created: container.Post.Post_created,
		Post_forum:   container.Post.Post_forum,
		Post_id:      container.Post.Post_id,
		IsEdited:     container.IsEdited,
		Post_message: container.Post.Post_message,
		Parent:       container.Parent,
		Post_thread:  container.Post.Post_thread,
	}

	if _, isExist := args["user"]; isExist {
		findPost.Author = &models.User{
			Nickname: container.Post.Post_author,
			Email:    container.Email,
			Fullname: container.Fullname,
			About:    container.About,
		}
	}

	if _, isExist := args["thread"]; isExist {
		findPost.Thread = &models.Thread{
			Thread_author:  container.Thread.Thread_author,
			Thread_created: container.Thread.Thread_created,
			Thread_forum:   container.Post.Post_forum,
			Thread_id:      container.Post.Post_thread,
			Thread_message: container.Thread.Thread_message,
			Thread_slug:    container.Thread.Thread_slug,
			Thread_title:   container.Thread.Thread_title,
			Votes:          container.Votes,
		}
	}

	if _, isExist := args["forum"]; isExist {
		findPost.Forum = &models.Forum{
			Posts:       container.Posts,
			Forum_slug:  container.Post.Post_forum,
			Threads:     container.Threads,
			Forum_title: container.Forum.Forum_title,
			Forum_user:  container.Forum_user,
		}
	}

	return findPost, nil
}

func FindPost(id string) *models.Post {
	findPost := &models.Post{}

	err := db.Get(
		findPost,
		`SELECT post_author, post_created, post_forum, post_id, isEdited, post_message, parent, post_thread 
		FROM post 
		WHERE post_id = $1`,
		id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Println("Error of find post by id")
		return nil
	}

	return findPost
}
