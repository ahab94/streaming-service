package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"

	domain "github.com/ahab94/streaming-service/models"
)

// GetUsage gets hourly usages
func (m *client) GetUsage(ctx context.Context, userID string, start, end time.Time) ([]domain.HourlyUsage, error) {
	hourlyUsages := make([]domain.HourlyUsage, 0)

	stmt, names := qb.Select(usageCF).
		Where(qb.EqLit(userIDKey, stringer(userID))).
		Where(qb.GtOrEqLit(hourKey, stringer(start.Format(timeFormat)))).
		Where(qb.LtOrEqLit(hourKey, stringer(end.Format(timeFormat)))).
		ToCql()

	logger.Debugf("processed get query: %+v", stmt)

	q := gocqlx.Query(m.conn.Query(stmt), names)
	if err := q.SelectRelease(&hourlyUsages); err != nil {
		logger.Errorf("unable to get hourlyUsages by statement: %+v, err: %+v ", stmt, err)
		return nil, err
	}

	logger.Infof("filtered: %d hourly usages(s)", len(hourlyUsages))

	return hourlyUsages, nil
}

// CountUsage saves Packet in store
func (m *client) IncrementUsage(ctx context.Context, packet *domain.Packet) error {
	stmt, names := qb.Update(usageCF).SetLit("pkt_count", "pkt_count + 1").
		Where(qb.EqLit(userIDKey, fmt.Sprintf("'%s'", packet.UserID))).
		Where(qb.EqLit(hourKey, stringer(packet.Timestamp.Round(time.Hour).Format(timeFormat)))).
		Where(qb.EqLit(channelIDKey, fmt.Sprintf("'%s'", packet.ChannelID))).
		ToCql()

	q := gocqlx.Query(m.conn.Query(stmt), names)
	if err := q.ExecRelease(); err != nil {
		logger.Errorf("unable to save usage: %v ", err)
		return err
	}

	return nil
}
