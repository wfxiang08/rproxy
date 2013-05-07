package main
import (
"testing"
//"fmt"
)

func TestWrite(t *testing.T){
  req := new(Request)
  req.data = []string{"SET","hello","world"}
  conn := new(Connection)
  conn.Open("127.0.0.1:6379:0")
  conn.WriteRequest(req)
  //resp := conn.ReadResponse()
}