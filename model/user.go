package model

import (
	"database/sql"
)

type User struct {
	_id      *int64 `db:"id"`
	Username string `db:"username"`

	State         int    `db:"state"`
	Enabled       bool   `db:"enabled"`
	PreviousEmoji string `db:"previous_emoji"`
	SpotifyEmoji  string `db:"spotify_emoji"`

	SlackAccessToken    string `db:"slack_access_token"`
	SpotifyRefreshToken string `db:"spotify_refresh_token"`
}

func (u *User) id() *int64 {
	return u._id
}

func (u *User) setId(id int64) {
	u._id = &id
}

func (u User) table() string {
	return "user"
}

// Creates a new user in the DB or updates an existing one
func (u *User) Save(db *sql.DB) error {
	return save(u, db)
}

func LoadUser(id int64, db *sql.DB) (*User, error) {
	u := &User{_id: &id}
	err := load(u, db)
	return u, err
}

func LoadAllUsers(db *sql.DB) ([]*User, error) {
	ifaces, err := loadAll(&User{}, db)
	if err != nil {
		return nil, err
	}
	users := make([]*User, len(ifaces))
	for i, iface := range ifaces {
		users[i] = iface.(*User)
	}
	return users, nil
}
