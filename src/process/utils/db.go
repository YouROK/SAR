package utils

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DBUtils struct {
	db *bolt.DB
}

func OpenDB(path string) (utils *DBUtils, err error) {
	db, err := bolt.Open(path, 0666, nil)
	dbu := new(DBUtils)
	dbu.db = db
	return dbu, nil
}

func (dbu *DBUtils) Close() error {
	if dbu.db != nil {
		return dbu.db.Close()
	}
	return nil
}

func (dbu *DBUtils) put(path, key string, data []byte) error {
	if len(path) == 0 || len(key) == 0 || len(data) == 0 {
		return os.ErrNotExist
	}
	path = filepath.Clean(path)
	tmpList := strings.Split(path, "/")
	var treelist []string
	for _, t := range tmpList {
		if t != "" {
			treelist = append(treelist, t)
		}
	}
	db := dbu.db
	var b *bolt.Bucket
	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucketIfNotExists([]byte("/"))
		if err != nil {
			return err
		}
		for _, d := range treelist {
			b, err = b.CreateBucketIfNotExists([]byte(d))
			if err != nil {
				return err
			}
		}
		return b.Put([]byte(key), data)
	})
	return err
}

//after get, if bucket not nil mast be commited
func (dbu *DBUtils) get(path string) (*bolt.Bucket, error) {
	if len(path) == 0 {
		return nil, os.ErrNotExist
	}
	path = filepath.Clean(path)
	tmpList := strings.Split(path, "/")
	var treelist []string
	for _, t := range tmpList {
		if t != "" {
			treelist = append(treelist, t)
		}
	}
	db := dbu.db
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}

	b := tx.Bucket([]byte("/"))
	if b == nil {
		return b, os.ErrNotExist
	}
	for _, d := range treelist {
		b = b.Bucket([]byte(d))
		if b == nil {
			return b, os.ErrNotExist
		}
	}
	return b, err
}

func (dbu *DBUtils) AddDir(fullpath string) error {
	return filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			buf, err := ioutil.ReadFile(path)
			if err == nil {
				err = dbu.put(path, "data", buf)
				if err == nil {
					err = dbu.put(path, "type", []byte{1})
				}
			}
			return err
		} else if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err == nil {
				err = dbu.put(path, "data", []byte(link))
				if err == nil {
					err = dbu.put(path, "type", []byte{2})
				}
			}
			return err
		}
		return nil
	})
}

func (dbu *DBUtils) GetNames(path string) ([]string, error) {
	b, err := dbu.get(path)
	if b != nil {
		defer b.Tx().Commit()
	}
	if err != nil {
		return nil, err
	}

	c := b.Cursor()
	var names []string
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		names = append(names, filepath.Join(path, string(k)))
	}
	return names, nil
}

func (dbu *DBUtils) PrintAll(b *bolt.Bucket, pos int) {
	var tx *bolt.Tx
	if b == nil {
		tx, _ = dbu.db.Begin(false)
		if tx != nil {
			b = tx.Bucket([]byte("/"))
		}
		if b == nil {
			return
		}
	}
	pos++
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		ind := ""
		for i := 0; i < pos; i++ {
			ind += "\t"
		}
		fmt.Printf(ind+"/%s\n", k)
		if string(k) == "type" {
			fmt.Printf(ind+" value=%d\n", v[0])
		} else if len(v) > 0 {
			ss := ""
			if len(v) > 10 {
				ss = string(v[:5]) + "..."
				ss = strings.Replace(ss, "\n", " ", -1)
			} else {
				ss = string(v)
			}
			fmt.Printf(ind+" value=%s\n", ss)
		}
		b := c.Bucket().Bucket(k)
		if b != nil {
			dbu.PrintAll(b, pos)
		}
	}

	if tx != nil {
		tx.Commit()
	}
}
