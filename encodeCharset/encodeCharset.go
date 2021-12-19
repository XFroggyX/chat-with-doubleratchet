package encodeCharset

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// FromBytes - convert binary format of a 4 byte integer to int32
func FromBytes(b []byte) (int32, error) {
	buf := bytes.NewReader(b)
	var result int32
	err := binary.Read(buf, binary.BigEndian, &result)
	return result, err
}

// ToBytes - convert an int32 to a 4 byte to binary format
func ToBytes(i int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	return buf.Bytes(), err
}

// WriteMsg - функция для отправки сообщения
func WriteMsg(conn net.Conn, msg string) error {
	// Отправляем размер сообщения
	bytes, err := ToBytes(int32(len([]byte(msg))))
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	if err != nil {
		return err
	}
	// Само сообщение
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// ReadMsg - функция для принятия сообщения
func ReadMsg(conn net.Conn) (string, error) {
	// храним длинну входных данных
	lenBuf := make([]byte, 4)
	_, err := conn.Read(lenBuf)
	if err != nil {
		return "", err
	}

	// длинна сообщения
	lenData, err := FromBytes(lenBuf)
	if err != nil {
		return "", err
	}

	buf := make([]byte, lenData)
	reqLen := 0

	for reqLen < int(lenData) {
		readMsg, err := conn.Read(buf[reqLen:])
		reqLen += readMsg
		if err == io.EOF {
			return "", fmt.Errorf("received EOF before receiving all promised data")
		}
		if err != nil {
			return "", fmt.Errorf("error reading: %s", err.Error())
		}
	}
	return string(buf), nil
}
