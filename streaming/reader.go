package streaming

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	pq "github.com/kyroy/priority-queue"
	"github.com/spf13/viper"

	"github.com/ahab94/streaming-service/config"
	"github.com/ahab94/streaming-service/models"
)

type Reader struct {
	conn          net.Conn
	data          chan<- map[string][]models.Packet
	buffer        map[string]*pq.PriorityQueue
	CurrentSeqNum map[string]int64
}

func NewReader(data chan<- map[string][]models.Packet) *Reader {
	return &Reader{
		data:          data,
		buffer:        make(map[string]*pq.PriorityQueue),
		CurrentSeqNum: make(map[string]int64),
	}
}

func (r *Reader) Start() {
	for {
		if err := r.connect(); err != nil {
			logger.Errorf("failed to connect to channel source err: %+v backing off for %d seconds", err, 5)
			time.Sleep(5 * time.Second)
			continue
		}
		r.readBatch()
	}
}

func (r *Reader) connect() error {
	defer recoverPanic()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", viper.GetString(config.ChannelSourceHost),
		viper.GetString(config.ChannelSourcePort)))
	if err != nil {
		return err
	}

	r.conn = conn

	return nil
}

// readBatch reads n number of packet and push to publisher
func (r *Reader) readBatch() {
	for {
		packetsByChannel := make(map[string][]models.Packet)
		for i := 0; i < viper.GetInt(config.PacketBatchSize); i++ {
			packet, err := r.ReadPacket()
			if err != nil {
				if err == io.EOF {
					r.data <- packetsByChannel
					return
				}
				logger.Errorf("error encountered while reading packet err: %+v", err)
			}

			if _, ok := r.buffer[packet.ChannelID]; !ok {
				r.buffer[packet.ChannelID] = pq.NewPriorityQueue()
			}

			if packet.SequenceNum != r.CurrentSeqNum[packet.ChannelID] {
				r.buffer[packet.ChannelID].Insert(packet, float64(packet.SequenceNum))

				pkt := r.buffer[packet.ChannelID].PopLowest().(models.Packet)
				if pkt.SequenceNum != r.CurrentSeqNum[packet.ChannelID] {
					r.buffer[packet.ChannelID].Insert(pkt, float64(pkt.SequenceNum))
					continue
				}
				packet = pkt
			}

			packetsByChannel[packet.ChannelID] = append(packetsByChannel[packet.ChannelID], packet)
			r.CurrentSeqNum[packet.ChannelID] = packet.SequenceNum + 1
		}

		r.data <- packetsByChannel
	}
}

// ReadPacket reads packet from channel source
func (r *Reader) ReadPacket() (models.Packet, error) {
	defer recoverPanic()

	pkt := models.Packet{}
	metadata := make([]byte, 23)
	if _, err := io.ReadFull(r.conn, metadata); err != nil {
		return pkt, err
	}

	sequence, err := strconv.ParseInt(fmt.Sprint(binary.BigEndian.Uint64(metadata[11:19])), 10, 64)
	if err != nil {
		return pkt, err
	}

	pSize, err := strconv.Atoi(fmt.Sprint(binary.BigEndian.Uint32(metadata[19:23])))
	if err != nil {
		return pkt, err
	}

	data := make([]byte, pSize)
	if _, err := io.ReadFull(r.conn, data); err != nil {
		return pkt, err
	}

	pkt.SequenceNum = sequence
	pkt.ChannelID = string(metadata[5:11])
	pkt.Raw = data

	return pkt, nil
}
