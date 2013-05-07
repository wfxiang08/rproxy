package main
import (
//"fmt"
)

type Cluster struct{
  master *Connection
  slaves []*Connection
}

func (cluster *Cluster) Config(cfg map[interface{}] interface{}){
  master := cfg["master"].(string)
  conn := new(Connection)
  conn.Open(master)
  cluster.master = conn
  slaves := cfg["slaves"].([]interface{})
  for _,slave := range(slaves){
    slaveConn := new(Connection)
    slaveConn.Open(slave.(string))
    cluster.slaves = append(cluster.slaves,slaveConn)
  }
}


func (cluster *Cluster) Write(req *Request){
  
}