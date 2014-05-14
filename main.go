package main

import (
	"flag"
	"github.com/jmhodges/levigo"
	"io"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"sync/atomic"
)

var ConnCount int64

func Count(args []string) (ret string, err error) {
	ret = strconv.FormatInt(ConnCount, 10)
	return
}

var port = flag.String("port", "23200", "leveld listen port")
var host = flag.String("host", "localhost", "leveld bind host")

func main() {
	flag.Parse()

	dbs := map[string]*levigo.DB{}

	l, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalln("net.Listen", err)
	}
	log.Println("leveldb listen on ", *port)

	for {
		cli, err := l.Accept()
		conn := Conn{
			dbs:    dbs,
			client: cli,
		}
		if err != nil {
			log.Println("net.Listener.Accept", err)
			continue
		}
		log.Println("conn.Accepted:", cli.RemoteAddr().String())
		go HandleConnect(conn)
	}
}

func HandleConnect(conn Conn) {

	atomic.AddInt64(&ConnCount, 1)
	buf := textproto.NewConn(conn.client)

	defer (func() {
		conn.client.Close()
		atomic.AddInt64(&ConnCount, -1)
	})()
	for {
		msg := []string{}
		for {
			if str, err := buf.ReadLine(); err == nil {
				if len(str) == 0 { //END
					break
				}
				msg = append(msg, str)
			} else if err == io.EOF {
				return
			} else {
				log.Println(err)
			}
		}
		if len(msg) < 2 { //BAD
			ReplyError(conn.client, ErrorNoSuchCmd)
			continue
		}
		cmd, dbName, err := parseCmd(msg[0])
		if err != nil {
			ReplyError(conn.client, err)
			continue
		}
		key := []byte(msg[1])

		idb, ok := conn.dbs[dbName]
		if !ok {
			if idb, err = mkdb(dbName); err != nil {
				ReplyError(conn.client, err)
				continue
			}
			conn.dbs[dbName] = idb
		}
		switch cmd {
		case "g":
			if val, err := dbget(idb, key); err != nil {
				ReplyError(conn.client, err)
			} else {
				ReplyData(conn.client, string(val))
			}
		case "s":
			if len(msg) != 3 {
				ReplyError(conn.client, ErrorInvalid)
				continue
			}
			val := msg[2]
			if err = dbset(idb, key, []byte(val)); err == nil {
				ReplyOK(conn.client)
			} else {
				ReplyError(conn.client, err)
			}
		case "d":
			if err = dbdel(idb, key); err == nil {
				ReplyOK(conn.client)
			} else {
				ReplyError(conn.client, err)
			}
		}
	}
	return
}
