package streaming

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/viper"

	"github.com/ahab94/streaming-service/config"
	"github.com/ahab94/streaming-service/service"
)

// StreamServer - packet stream, processing and distribution
type StreamServer struct {
	ctx      context.Context
	svc      service.Svc
	listener net.Listener
	rmChan   chan<- *User
	addChan  chan<- *User
}

func NewStreamServer(ctx context.Context, svc service.Svc, addChan, rmChan chan<- *User) (*StreamServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString(config.StreamSrvHost), viper.GetString(config.StreamSrvPort)))
	if err != nil {
		logger.Errorf("tcp server listener error: %+v", err)
		return nil, err
	}

	s := &StreamServer{
		ctx:      ctx,
		svc:      svc,
		listener: listener,
		addChan:  addChan,
		rmChan:   rmChan,
	}

	return s, nil
}

func (s *StreamServer) Addr() string {
	return s.listener.Addr().String()
}

func (s *StreamServer) ListenAndServe() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.Errorf("error while accepting connections err: %+v", err)
			return
		}

		logger.Debugf("setting up connection %s", conn.RemoteAddr().String())

		go s.handle(conn)
	}
}

func (s *StreamServer) Stop() {
	if err := s.listener.Close(); err != nil {
		logger.Warnf("error while closing server err:%+v", err)
	}
}

func (s *StreamServer) handle(conn net.Conn) {
	defer recoverPanic()
	line, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		logger.Errorf("failed to setup user err: %+v", err)
		return
	}

	logger.Debugf("payload %s from %s", string(line), conn.RemoteAddr().String())

	words := strings.Split(string(line), " ")
	if len(words) < 2 {
		logger.Errorf("failed to setup user err: %+v", errors.New("invalid payload"))
		return
	}

	if _, err := s.svc.GetUser(s.ctx, words[0]); err != nil {
		logger.Errorf("failed to setup user err: %+v", errors.New("user validation failed"))
		return
	}

	u := NewUser(s.svc, words[0], words[1], conn, s.rmChan)
	s.addChan <- u

	u.stream()
}
