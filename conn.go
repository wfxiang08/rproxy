package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	way       net.Conn
	reader    *bufio.Reader
	connected bool
}

type Request struct {
	data []string
}

type Response struct {
	data string
}

func (conn *Connection) Open(connstr string) {
	arr := strings.Split(connstr, ":")
	addr := fmt.Sprintf("%s:%s", arr[0], arr[1])

	// 建立tcp连接
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	conn.way = c
	conn.reader = bufio.NewReader(conn.way)
	conn.connected = true
	_ = fmt.Sprintf("SELECT %s\r\n", arr[2])
}

func (conn *Connection) SetWay(way net.Conn) {
	conn.way = way
	conn.reader = bufio.NewReader(conn.way)
	conn.connected = true
}

// Redis协议
// 参考: http://redis.readthedocs.org/en/latest/topic/protocol.html
//*<参数数量> CR LF
//$<参数 1 的字节数量> CR LF
//<参数 1 的数据> CR LF
//...
//$<参数 N 的字节数量> CR LF
//<参数 N 的数据> CR LF
//
func (conn *Connection) WriteRequest(req *Request) {
	length := len(req.data)
	var line string
	line = fmt.Sprintf("*%d\r\n", length) // 参数格式: *%d
	conn.way.Write([]byte(line))

	for _, arg := range req.data {
		line = fmt.Sprintf("$%d\r\n", len(arg)) // $<参数 1 的字节数量> CR LF
		conn.way.Write([]byte(line))
		line = fmt.Sprintf("%s\r\n", arg) // 参数1的数据
		conn.way.Write([]byte(line))
	}
}

func (conn *Connection) ReadRequest() Request {
	var req Request
	count := conn.readBulkNum()

	for i := 0; i < count; i++ {
		// 一口气读取一个参数
		bulk := conn.readBulk()
		if bulk != nil {
			// 所有的参数放在: req.data中
			req.data = append(req.data, string(bulk))
			fmt.Printf("%s\n", string(bulk))
		}
	}
	return req
}

func (conn *Connection) readBulkNum() int {
	var count int
	// 从 *%d\r\n中读取int
	line, _ := conn.reader.ReadString('\n')
	fmt.Sscanf(line, "*%d\r\n", &count)
	return count
}

func (conn *Connection) readBulk() []byte {
	var data []byte
	line, err := conn.reader.ReadString('\n')
	if err != nil {
		fmt.Printf(err.Error())
	}
	switch line[0] {
	case '$':
		// 读取size
		size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil
		}
		if size == -1 {
			return nil
		}

		// 继续读取data，限定长度
		lr := io.LimitReader(conn.reader, int64(size))
		data, err = ioutil.ReadAll(lr)

		// 读取完毕 \r\n
		if err == nil {
			_, err = conn.reader.ReadString('\n')
		}
		return data
	default:
		return nil
	}
	return nil
}
