package mcquery

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

type packetType byte

var (
	packetStat      packetType = 0x00
	packetHandshake packetType = 0x09
)

// McQuery represents a in-progress status request to a server.
type McQuery struct {
	conn      net.Conn
	challenge []byte
}

// Dial connects and handshakes with the Minecraft server.
func Dial(address string, timeout time.Duration) (*McQuery, error) {
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return nil, err
	}

	challenge, err := getChallenge(conn)
	if err != nil {
		return nil, err
	}

	return &McQuery{conn, challenge}, nil
}

func getChallenge(conn net.Conn) ([]byte, error) {
	err := sendRequest(conn, packetHandshake, []byte{})
	if err != nil {
		return nil, err
	}
	data, err := receiveResponse(conn, packetHandshake)
	if err != nil {
		return nil, err
	}

	endIdx := 0
	for i, n := range data {
		if n == '\000' {
			endIdx = i
			break
		}
	}

	schallange := string(data[:endIdx])

	challenge, err := strconv.Atoi(schallange)
	if err != nil {
		return nil, err
	}

	out := []byte{
		byte(challenge >> 24),
		byte(challenge >> 16),
		byte(challenge >> 8),
		byte(challenge),
	}

	return out, nil
}

// GetStatus returns the Minecraft server status from the server. It must be
// called within 30 seconds of calling Dial. Otherwise, Dial has to be called again.
func (mcq *McQuery) GetStatus() (status map[string]string, players []string, err error) {
	var buf bytes.Buffer
	buf.Write(mcq.challenge)
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	err = sendRequest(mcq.conn, packetStat, buf.Bytes())
	if err != nil {
		status = nil
		players = nil
		return
	}

	var data []byte
	data, err = receiveResponse(mcq.conn, packetStat)
	if err != nil {
		status = nil
		players = nil
		return
	}
	data = data[11:] // meaningless: 'splitnum' + 2 bytes

	status = make(map[string]string)

	cursor := -1
	for cursor < len(data) {
		key := readString(data, &cursor)
		if len(key) == 0 {
			break
		}

		value := readString(data, &cursor)

		status[key] = value
	}
	readString(data, &cursor) // meaningless: 'player_' + '\000'

	for cursor < len(data) {
		name := readString(data, &cursor)
		if len(name) > 0 {
			players = append(players, name)
		}
	}

	return
}

func readString(data []byte, cursor *int) string {
	*cursor++
	start := *cursor
	for ; (*cursor < len(data)) && (data[*cursor] != 0x00); *cursor++ {
	}
	return string(data[start:*cursor])
}

func sendRequest(w io.Writer, ptype packetType, data []byte) error {
	var buf bytes.Buffer
	buf.Write([]byte{
		0xFE, 0xFD, byte(ptype), 0x01, 0x01, 0x02, 0x03,
	})
	buf.Write(data)

	_, err := w.Write(buf.Bytes())
	return err
}

func receiveResponse(r io.Reader, expected packetType) ([]byte, error) {
	buf := make([]byte, 10240)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n < 5 {
		return nil, errors.New("Not enough data to read header.")
	}
	if buf[0] != byte(expected) {
		return nil, errors.New("Unexpected packet type.")
	}
	return buf[5:], nil
}
