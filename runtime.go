package streaming_service

import (
	"github.com/ahab94/streaming-service/db"
	"github.com/ahab94/streaming-service/db/cassandra"
	"github.com/ahab94/streaming-service/service"
	"github.com/ahab94/streaming-service/streaming"
)

type Runtime struct {
	ds        db.DataStore
	svc       service.Svc
	publisher *streaming.Publisher
	reader    *streaming.Reader
}

func NewRuntime() (*Runtime, error) {
	ds, err := cassandra.NewStore()
	if err != nil {
		return nil, err
	}

	p := streaming.NewPublisher()
	rt := &Runtime{
		ds:        ds,
		publisher: p,
		svc:       service.NewService(ds),
		reader:    streaming.NewReader(p.DataInput()),
	}

	return rt, nil
}

func DefaultRuntime() *Runtime {
	return &Runtime{}
}

func (r *Runtime) WithDataStore(ds db.DataStore) *Runtime {
	r.ds = ds
	return r
}

func (r *Runtime) WithService(service service.Svc) *Runtime {
	r.svc = service
	return r
}

func (r *Runtime) DataStore() db.DataStore {
	return r.ds
}

func (r *Runtime) Service() service.Svc {
	return r.svc
}

func (r *Runtime) Publisher() *streaming.Publisher {
	return r.publisher
}

func (r *Runtime) Reader() *streaming.Reader {
	return r.reader
}
