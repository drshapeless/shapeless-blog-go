package data

import "database/sql"

type Models struct {
	Posts PostModel
	Tags  TagModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Posts: PostModel{DB: db},
		Tags:  TagModel{DB: db},
	}
}

// This function is quite useless.
// Migration should be done in client.
func Migrate(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS posts(
id INTEGER PRIMARY KEY,
title TEXT NOT NULL,
filename TEXT NOT NULL,
created TEXT NOT NULL,
updated TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tags(
post_id INTEGER NOT NULL,
tag TEXT NOT NULL
);
`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
