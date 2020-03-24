package streaming

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/ahab94/streaming-service/models"
	"github.com/ahab94/streaming-service/service"
)

type User struct {
	ID          string
	SessionID   string
	channelID   string
	errCount    int
	conn        net.Conn
	removeInput chan<- *User
	done        chan struct{}
	input       chan []models.Packet
	svc         service.Svc
}

func NewUser(svc service.Svc, userID, channel string, conn net.Conn, rmChan chan<- *User) *User {
	return &User{
		svc:         svc,
		removeInput: rmChan,
		ID:          userID,
		conn:        conn,
		channelID:   channel,
		SessionID:   uuid.New().String(),
		done:        make(chan struct{}),
		input:       make(chan []models.Packet),
	}
}

func (u *User) Stop() {
	defer recoverPanic()
	if err := u.conn.Close(); err != nil {
		logger.Warnf("error while closing connection for user %s err:%+v", u.ID, err)
	}

	u.removeInput <- u
	close(u.done)
}

func (u *User) Input() chan []models.Packet {
	return u.input
}

func (u *User) stream() {
	for {
		select {
		case packets := <-u.input:
			logger.Debugf("sending %d packets to user %s", len(packets), u.ID)
			u.writePackets(packets)
		case <-u.done:
			logger.Debugf("stopping user %s session %s...", u.ID, u.SessionID)
			return
		case <-time.After(3 * time.Second):
			if u.IsClosed() {
				logger.Debugf("user %s session %s connection closed", u.ID, u.SessionID)
				u.Stop()
			}
		}
	}
}

func (u *User) writePackets(packets []models.Packet) {
	defer recoverPanic()
	for _, packet := range packets {
		_, err := u.conn.Write(packet.Raw)
		if u.isFatal(err) {
			return
		}

		u.savePacket(packet)
	}
}

func (u *User) savePacket(packet models.Packet) {
	packet.UserID = u.ID
	packet.Timestamp = time.Now()

	if err := u.svc.SavePacket(context.TODO(), &packet); err != nil {
		logger.Errorf("error while storing packet for user %s, err: %+v...", u.ID, err)
	}
}

func (u *User) isFatal(err error) bool {
	if err == nil {
		u.errCount = 0
		return false
	}

	logger.Errorf("error while writing packet on user %s, err: %+v...", u.ID, err)

	u.errCount += 1
	if u.errCount == 10 {
		u.Stop()
		return true
	}

	return false
}

func (u *User) IsClosed() bool {
	if err := u.conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		logger.Debugf("user %s error in setting read deadline %+v", u.ID, err)
	}
	if _, err := u.conn.Read(make([]byte, 1)); err != nil {
		if err == io.EOF {
			return true
		}
		if nerr, ok := err.(net.Error); ok && !nerr.Temporary() {
			return true
		}
	}
	return false
}
