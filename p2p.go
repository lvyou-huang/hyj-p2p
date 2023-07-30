package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
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
type Hub struct {
	Clients    map[string]Client
	Register   chan Client
	Unregister chan Client
}

var hub Hub

func init() {
	hub.Clients = make(map[string]Client)
	hub.Register = make(chan Client)
	hub.Unregister = make(chan Client)
}
func Start() {
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
	r.Run(":8080")
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
