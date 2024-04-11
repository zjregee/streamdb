package storage

import (
	"os"
	"testing"
)

func TestBolt(t *testing.T) {
	db := NewDBWithPath("test.db")
	err := db.Open()
	if err != nil {
		t.Errorf("failed to open database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove("test.db")
	}()
	err = db.WriteData("bucket", "key", "value")
	if err != nil {
		t.Errorf("failed to write data: %v", err)
	}
	data, err := db.ReadData("bucket", "key")
	if err != nil {
		t.Errorf("failed to read data: %v", err)
	}
	if data != "value" {
		t.Errorf("expected value 'value', got '%s'", data)
	}
}
