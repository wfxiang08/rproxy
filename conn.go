package main



import (
"io"
"io/ioutil"
"bufio"
"fmt"
"strconv"
"strings"
"net"
)



type Connection struct{
  way net.Conn
  reader *bufio.Reader
  connected bool
}

type Request struct{
  data []string
} 

type Response struct{
  data string
}

func (conn *Connection) Open(connstr string) {
  arr := strings.Split(connstr,":")
  addr := fmt.Sprintf("%s:%s",arr[0],arr[1])
  c,err := net.Dial("tcp",addr)
  if err!=nil{
    return 
  }
  conn.way = c
  conn.reader = bufio.NewReader(conn.way)
  conn.connected = true
  _ = fmt.Sprintf("SELECT %s\r\n",arr[2])
}

func (conn *Connection) SetWay(way net.Conn){
  conn.way = way
  conn.reader = bufio.NewReader(conn.way)
  conn.connected = true
}



func (conn *Connection) ReadRequest() Request{
  var req Request
  count := conn.readBulkNum()
  for i:= 0; i<count; i++{
    bulk := conn.readBulk()
    if bulk != nil {
      req.data = append(req.data,string(bulk))
      fmt.Printf("%s\n",string(bulk))
    }
  }
  return req
}


func (conn *Connection) readBulkNum() int{
  var count int
  line, _ := conn.reader.ReadString('\n')
  fmt.Sscanf(line,"*%d\r\n",&count)
  return count
}

func(conn *Connection) readBulk() []byte{
  var data []byte
  line, err:= conn.reader.ReadString('\n')
  if err != nil {
      fmt.Printf(err.Error())
  }
  switch line[0]{
  case '$':
    size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
    if err != nil {
      return nil
    }
    if size == -1 {
      return nil
    }
    lr := io.LimitReader(conn.reader,int64(size))
    data,err = ioutil.ReadAll(lr)
    if err == nil{
      _,err = conn.reader.ReadString('\n')
    }
    return data
  default:
    return nil
  }
  return nil
}


