package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"

	domain "github.com/ahab94/streaming-service/models"
)

// ListPackets gets Packet by id
func (m *client) ListPackets(ctx context.Context, filter map[string]interface{}) ([]domain.Packet, error) {
	packets := make([]domain.Packet, 0)

	stmt, names := mkListQuery(filter)
	logger.Debugf("processed list query: %+v", stmt)

	q := gocqlx.Query(m.conn.Query(stmt), names)
	if err := q.SelectRelease(&packets); err != nil {
		logger.Errorf("unable to get packets by filter: %+v, err: %+v ", filter, err)
		return nil, err
	}

	logger.Infof("filtered: %d packet(s)", len(packets))

	return packets, nil
}

// SavePacket saves Packet in store
func (m *client) SavePacket(ctx context.Context, pkt *domain.Packet) error {
	stmt, names := qb.Insert(packetCF).Columns(pkt.Names()...).ToCql()
	q := gocqlx.Query(m.conn.Query(stmt), names).BindStruct(pkt)
	if err := q.ExecRelease(); err != nil {
		logger.Errorf("unable to save packet: %v ", err)
		return err
	}

	return nil
}

func mkListQuery(filter map[string]interface{}) (string, []string) {
	q := qb.Select(packetCF)

	if v, ok := filter[userIDKey]; ok {
		q = q.Where(qb.EqLit(userIDKey, stringer(v)))
	}

	if v, ok := filter[startDateKey]; ok {
		if ts, ok := v.(time.Time); ok {
			q = q.Where(qb.GtOrEqLit(timestampKey, stringer(ts.Format(timeFormat))))
		}
	}

	if v, ok := filter[endDateKey]; ok {
		if ts, ok := v.(time.Time); ok {
			q = q.Where(qb.LtOrEqLit(timestampKey, stringer(ts.Format(timeFormat))))
		}
	}

	return q.ToCql()
}

func stringer(v interface{}) string {
	return fmt.Sprintf("'%v'", v)
}
