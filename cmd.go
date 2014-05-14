package main

import (
	"fmt"
	"github.com/jmhodges/levigo"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	OK = iota
	UnknowCmd
	InvalidCmd
	OpenDBFailed
	UnknowError
)

type Conn struct {
	dbs    map[string]*levigo.DB
	client net.Conn
}

var ro = levigo.NewReadOptions()
var wo = levigo.NewWriteOptions()

func dbget(db *levigo.DB, key []byte) (val []byte, err error) { return db.Get(ro, key) }

func dbset(db *levigo.DB, key, val []byte) (err error) { return db.Put(wo, key, val) }

func dbdel(db *levigo.DB, key []byte) (err error) { return db.Delete(wo, key) }

func mkdb(name string) (db *levigo.DB, err error) {

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3 << 30))
	opts.SetCreateIfMissing(true)
	db, err = levigo.Open(name, opts)
	if err != nil {
		err = ErrorOpenDB
	}
	return
}

type CmdError struct {
	code int
}

func (c *CmdError) Error() string {
	return strconv.Itoa(c.code)
}

var ErrorInvalid = &CmdError{code: InvalidCmd}
var ErrorNoSuchCmd = &CmdError{code: UnknowCmd}
var ErrorOpenDB = &CmdError{code: OpenDBFailed}

func parseCmd(cmd string) (verb, dbname string, err error) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) != 2 {
		err = ErrorInvalid
		return
	}
	switch cmds[0] {
	case "g", "s", "d", "l":
		verb = cmds[0]
	default:
		err = ErrorNoSuchCmd
		return
	}
	dbname = cmds[1]
	return
}

func ReplyData(w io.Writer, data string) {
	fmt.Fprintf(w, "ok 0\r\n%s\r\n\r\n", data)
}

func ReplyOK(w io.Writer) {
	fmt.Fprint(w, "ok 0\r\n\r\n")
}

func ReplyError(w io.Writer, err error) {
	code := UnknowError
	if err_, ok := err.(*CmdError); ok {
		code = err_.code
	}
	fmt.Fprintf(w, "e %d\r\n\r\n", code)
}
