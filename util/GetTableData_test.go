package util

import "testing"

func TestTablePayloadBounds(t *testing.T) {
	for _, data := range [][]byte{nil, {1, 2, 3}, {255, 255, 255, 255}, {1, 0, 0, 0}} {
		if _, err := tablePayload(data); err == nil {
			t.Fatalf("accepted invalid header %x", data)
		}
	}
	data, err := tablePayload([]byte{1, 0, 0, 0, 99, 42})
	if err != nil || len(data) != 1 || data[0] != 42 {
		t.Fatalf("%v %v", data, err)
	}
}
