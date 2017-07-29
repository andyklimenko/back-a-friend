package db

import (
	"testing"
	"io/ioutil"
	"os"
	"path"
)

func TestDb_Create(t *testing.T) {
	var mydb Db

	tmpDir, err := ioutil.TempDir("", "dbTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpDir)

	dbDir := path.Join(tmpDir, "db")
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		t.Fatal(err)
	}

	dbPath := path.Join(dbDir, "testdb.db")
	if err := mydb.Create(dbPath); err != nil {
		t.Error(err)
	}
}