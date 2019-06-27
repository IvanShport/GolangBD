package dbqueries

import "Forum/models"

func ClearData() error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("TRUNCATE TABLE forum_users CASCADE")
	if err != nil {
		return err
	}

	_, err = tx.Exec("TRUNCATE TABLE vote CASCADE")
	if err != nil {
		return err
	}

	_, err = tx.Exec("TRUNCATE TABLE post CASCADE")
	if err != nil {
		return err
	}

	_, err = tx.Exec("TRUNCATE TABLE thread CASCADE")
	if err != nil {
		return err
	}

	_, err = tx.Exec("TRUNCATE TABLE forum CASCADE")
	if err != nil {
		return err
	}

	_, err = tx.Exec("TRUNCATE TABLE user_profile CASCADE")
	if err != nil {
		return err
	}

	return tx.Commit()

}

func GetInfoAboutBD() (*models.Status, error) {

	allData := &models.Status{}

	err := db.Get(allData,
		`SELECT "user", "forum", "thread", "post"
			FROM (SELECT COUNT(*) AS "user" FROM user_profile) u
			, (SELECT COUNT(*) AS "forum" FROM forum) f
			, (SELECT COUNT(*) AS "thread" FROM thread) t
			, (SELECT COUNT(*) AS "post" FROM post) p;`)

	if err != nil {
		return nil, err
	}

	return allData, nil
}
