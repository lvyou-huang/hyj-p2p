package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"p2p/repository"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Name    string
	Files   []string
	UDPaddr net.UDPAddr
}

// 本地化，只测试本地，公网需要获得公网ip
type Hub struct {
	Port       []string //port(在公网运行服务器应位ip：port)
	MYport     string
	Status     int //0:follower,-1:candidate,1:leader
	Clients    map[string]Client
	Register   chan Client
	Unregister chan Client
}

var hub Hub

func init() {
	//上线无法知道其它节点
	hub.Port = []string{"8080", "8081", "8082"}
	hub.Status = 0
	hub.MYport = ""
	hub.Clients = make(map[string]Client)
	hub.Register = make(chan Client)
	hub.Unregister = make(chan Client)
}

func main() {
	port := os.Args[1]
	go hub.Run()
	//心跳机制
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for range ticker.C {
			for _, c := range hub.Clients {
				go func() {
					conn, err := net.DialUDP("udp", nil, &c.UDPaddr)
					if err != nil {
						hub.Unregister <- c
					}
					conn.Write([]byte("ping"))
					conn.SetReadDeadline(time.Now().Add(3 * time.Second))
					buf := make([]byte, 1024)
					oob := make([]byte, 1024)
					_, _, _, _, err = conn.ReadMsgUDP(buf, oob)
					if err != nil {
						hub.Unregister <- c
					}
				}()
			}
		}
	}()
	r := gin.Default()
	r.POST("", func(c *gin.Context) {
		name := c.PostForm("name")
		pass := c.PostForm("pwd")
		udpIp := c.PostForm("udpIp")
		udpPort, _ := strconv.Atoi(c.PostForm("udpPort"))
		owner := c.PostForm("files")
		files := strings.Split(owner, ",")
		client := Client{
			Name:  name,
			Files: files,
			UDPaddr: net.UDPAddr{
				IP:   net.ParseIP(udpIp),
				Port: udpPort,
				Zone: "",
			},
		}

		if _, ok := hub.Clients[name]; !ok {
			//注册
			repository.Register(name, pass)
			hub.Register <- client

		} else {
			//登录
			if repository.Login(name, pass) {
				log.Println("登录成功")
				hub.Register <- client
			} else {
				log.Println("登录失败")
				c.JSON(401, "wrong")
				return
			}
		}
		//返回所有用户信息

		log.Println("发送 clients")
		c.JSON(200, hub.Clients)
		// Handle the connection here

		//....
	})
	r.Run(":" + port)
}
func vote(candiante []string, voted string) []string {
	candiante = append(candiante, voted)
	return candiante
}
func RaftLeader(hub *Hub) {
	//进入leader
	hub.Status = 1
	for _, s := range hub.Port {
		if s == hub.MYport {
			continue
		}
		conn, err := net.Dial("tcp", ":"+s)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		conn.Write([]byte("win"))
		conn.Close()
	}
	//每60向follower发送心跳，保证存活
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for range ticker.C {
			for _, s := range hub.Port {
				if s == hub.MYport {
					continue
				}
				go func() {
					conn, err := net.Dial("tcp", ":"+s)
					if err != nil {
						log.Println(err)
						conn.Close()
						return
					}
					conn.Write([]byte("ping"))
					defer conn.Close()
				}()
			}
		}
	}()
}
func RaftCandiate(hub *Hub, candiate []string) {
	//进入候选人状态
	hub.Status = -1
	candiate = vote(candiate, hub.MYport)

	for _, s := range hub.Port {
		if s == hub.MYport {
			continue
		}
		conn, err := net.Dial("tcp", ":"+s)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		// 创建一个缓冲区，用于存储编码后的数据
		var buf bytes.Buffer

		// 创建一个编码器
		enc := gob.NewEncoder(&buf)

		// 编码数据
		if err := enc.Encode(candiate); err != nil {
			fmt.Println(err)
			return
		}
		// 获取编码后的字节序列
		encodedData := buf.Bytes()
		conn.Write(encodedData)
		conn.Close()
	}
	go func() {
		conn, _ := net.Listen("tcp", ":"+hub.MYport)
		accept, _ := conn.Accept()
		var buf []byte
		n, _ := accept.Read(buf)
		dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
		var candiate []string
		if err := dec.Decode(&candiate); err != nil {
			fmt.Println(err)
			return
		}
		found := false
		for _, v := range candiate {
			if v == hub.MYport {
				found = true
				break
			}
		}
		if found {
			conn.Close()
			RaftLeader(hub)
		} else {
			conn.Close()
			RaftCandiate(hub, candiate)
		}
	}()
}
func Raftfollower(hub *Hub) {
	//刚开始是follower，需要leader的定时心跳
	hub.Status = 0
	conn, err := net.Listen("tcp", ":"+hub.MYport)
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	pongReceived := make(chan bool)
	//心跳
	go func() {
		for {
			select {
			case <-ticker.C:
				if !<-pongReceived {
					conn.Close()
					RaftCandiate(hub, []string{})
				}
			}
		}
	}()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		accept, err := conn.Accept()
		if err != nil {
			log.Println(err)
		}
		pongReceived <- true
		accept.Write([]byte("pong"))
	}

}
func (hub *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	go func() {
		for {
			select {
			case client := <-hub.Register:
				log.Println("添加", client.Name)
				hub.Clients[client.Name] = client

			case client := <-hub.Unregister:
				log.Println("删除", client.Name)
				delete(hub.Clients, client.Name)
			}
		}
	}()

}
