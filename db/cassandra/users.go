package cassandra

import (
	"context"

	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"

	domain "github.com/ahab94/streaming-service/models"
)

// GetUser gets a user by id
func (m *client) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user := &domain.User{}
	stmt, names := qb.Select(usersCF).Where(qb.EqLit(userIDKey, stringer(userID))).ToCql()
	logger.Debugf("processed get query: %+v", stmt)

	if err := gocqlx.Query(m.conn.Query(stmt), names).GetRelease(user); err != nil {
		logger.Errorf("unable to get user by statement: %+v, err: %+v ", stmt, err)
		return nil, err
	}

	return user, nil
}

// SaveUser saves user in store
func (m *client) SaveUser(ctx context.Context, user *domain.User) error {
	stmt, names := qb.Insert(usersCF).Columns(user.Names()...).ToCql()
	q := gocqlx.Query(m.conn.Query(stmt), names).BindStruct(user)
	if err := q.ExecRelease(); err != nil {
		logger.Errorf("unable to save user: %v ", err)
		return err
	}

	return nil
}
