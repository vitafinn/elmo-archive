package util

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"reflect"
)

func ReadTableData(gameDataDir, gameServer, tableName string) ([]byte, error) {
	tablePath := filepath.Join(gameDataDir, "Table", tableName+".bytes")
	data, err := os.ReadFile(tablePath)
	if os.IsNotExist(err) && gameServer != "" {
		tablePath = filepath.Join(gameDataDir, "Table", gameServer, tableName+".bytes")
		data, err = os.ReadFile(tablePath)
	}
	if err != nil {
		return nil, err
	}
	return tablePayload(data)
}

func tablePayload(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("table header is truncated")
	}
	offset := uint64(binary.LittleEndian.Uint32(data[:4])) + 4
	if offset > uint64(len(data)) {
		return nil, fmt.Errorf("table header exceeds file length")
	}
	return data[offset:], nil
}

func GetTableData(gameDataDir, gameServer string, payload proto.Message) error {
	data, err := ReadTableData(gameDataDir, gameServer, reflect.TypeOf(payload).Elem().Name())
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, payload)
}
