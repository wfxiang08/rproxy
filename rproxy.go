package main
import (
"net"
"fmt"
"flag"
"os"
"github.com/gosexy/yaml"
)

var cluster Cluster



func handleConnection(c net.Conn){
  defer c.Close()
  conn := new(Connection)
  conn.SetWay(c) 
  req  := conn.ReadRequest()
  cluster.Write(&req)
}

func main(){
  configfile := flag.String("config","devel.yml","config file")
  flag.Parse()
  config,err := yaml.Open(*configfile)
  if err != nil{
    fmt.Printf("readfile(%q): %s", *config, err)
    os.Exit(1)
  }else{
    cfg  := config.Get("cluster").(map[interface {}]interface{})
    cluster.Config(cfg)
  }
  
  ln, err := net.Listen("tcp", "0.0.0.0:7000")
  if err != nil {
  	// handle error
  }
  for {
  	conn, err := ln.Accept()
  	if err != nil {
  		// handle error
  		continue
  	}
  	go handleConnection(conn)
  }
}