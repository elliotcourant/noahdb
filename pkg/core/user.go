package core

import (
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
)

type userContext struct {
	*base
}

// UserContext is just a wrapper interface for user metadata.
type UserContext interface {
	GetUsers() ([]User, error)
}

func (ctx *base) Users() UserContext {
	return &userContext{
		ctx,
	}
}

func (ctx *userContext) GetUsers() ([]User, error) {
	response, err := ctx.db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	return ctx.userFromRows(response)
}

func (ctx *userContext) userFromRows(response *frunk.QueryResponse) ([]User, error) {
	rows := rqliter.NewRqlRows(response)
	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(
			&user.UserID,
			&user.UserName,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
