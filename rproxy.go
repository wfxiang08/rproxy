package main

import (
	"flag"
	"fmt"
	"github.com/gosexy/yaml"
	"net"
	"os"
)

var cluster Cluster

func handleConnection(c net.Conn) {
	defer c.Close()
	conn := new(Connection)
	conn.SetWay(c)

	// 读取一个Connection的数据，将他转发给cluster
	req := conn.ReadRequest()
	cluster.Write(&req)
}

func main() {
	configfile := flag.String("config", "devel.yml", "config file")
	flag.Parse()
	config, err := yaml.Open(*configfile)
	if err != nil {
		fmt.Printf("readfile(%q): %s", *config, err)
		os.Exit(1)
	} else {
		// 得到一个Generic Dict
		cfg := config.Get("cluster").(map[interface{}]interface{})
		cluster.Config(cfg)
	}

	// 绑定本地的端口，然后转发请求? 似乎没有处理返回
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
