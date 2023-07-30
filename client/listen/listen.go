package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	mymd5 "p2p/hash"
	"strconv"
	"strings"
)

type Client struct {
	Name string
	//Pass    string
	Files   []string
	UDPaddr *net.UDPAddr
}

func main() {
	if len(os.Args) < 5 {
		log.Println("参数不够")
		return
	}
	//注册名字
	name := os.Args[1]
	//密码
	pwd := os.Args[2]
	//监听端口
	udpport := os.Args[3]
	//服务器ip
	ip := os.Args[4]
	//服务器端口
	serverport := os.Args[5]
	err := Register("http://"+ip+":"+serverport, udpport, name, pwd)
	if err != nil {
		return
	}
	ListenUdp("127.0.0.1" + ":" + udpport)
}

// 注册一个客户端0
func Register(toaddress string, udpport string, name string, pwd string) error {
	form := url.Values{}
	form.Set("udpPort", udpport)
	form.Set("name", name)

	form.Set("pwd", pwd)
	form.Set("udpIp", "127.0.0.1")
	files, err := os.ReadDir(".")
	if err != nil {
		log.Println(err)
		return err
	}
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	for i, v := range fileNames {
		if v == "listen.go" {
			fileNames = append(fileNames[:i], fileNames[i+1:]...)
		}
	}
	fileNameSep := strings.Join(fileNames, ",")
	form.Set("files", fileNameSep)
	resp, err := http.PostForm(toaddress, form)
	if err != nil {
		log.Println(err)
		return err
	}

	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	var m map[string]Client
	err = json.Unmarshal(buf[:n], &m)
	if err != nil {
		log.Println("账号或者密码错误")
		return err
	}
	if len(m) <= 1 {
		log.Println("暂时没有其它节点")
		return err
	}
	return nil
}
func ListenUdp(address string) {
	listenUdp(address)
}

func listenUdp(address string) {
	ip, portstr, err := net.SplitHostPort(address)
	if err != nil {
		log.Println(err)
		return
	}
	port, _ := strconv.Atoi(portstr)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
		Zone: "",
	})
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Listening on", conn.LocalAddr().String())
	for {

		buffer := make([]byte, 1024)
	flag:
		n, addr, err := conn.ReadFrom(buffer)

		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, addr.String())
		ip, portstr, err = net.SplitHostPort(addr.String())
		port, _ = strconv.Atoi(portstr)
		//判断buffer是否是一个文件名
		if string(buffer[:n]) == "ping" {
			_, _, err = conn.WriteMsgUDP([]byte("pings"), []byte("pings"), &net.UDPAddr{
				IP:   net.ParseIP(ip),
				Port: port,
				Zone: "",
			})
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			continue
		}

		file, err := os.Open(string(buffer[:n]))

		if err != nil {
			log.Println(err)

			continue
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			log.Println(err)
			conn.Write([]byte("not a filename"))
			continue
		}
		calcMd5, err := mymd5.HashSHA256File(info.Name())
		if err != nil {
			log.Println(err)
			continue
		}
		//逐帧发送，使用ACK协议验证对方是否收到
		i := 0

		for {
			if i == 0 {
				x := int32(info.Size())
				b1 := byte(x >> 24)
				b2 := byte(x >> 16)
				b3 := byte(x >> 8)
				b4 := byte(x & 0xff)
				join := BytesCombine3([]byte{b1, b2, b3, b4}, []byte(calcMd5))
				_, _, err = conn.WriteMsgUDP(join, []byte{}, &net.UDPAddr{
					IP:   net.ParseIP(ip),
					Port: port,
					Zone: "",
				})
				for {

					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						log.Println(err)
						continue
					}
					if string(buf[:n]) == "ACK" {
						//log.Println("Recive ACK", i)
						break
					}

				}
				i++
				continue
			}
			n, err := file.Read(buffer)
			//设置超时，超时则重传
			_, _, err = conn.WriteMsgUDP(buffer[:n], []byte{}, &net.UDPAddr{
				IP:   net.ParseIP(ip),
				Port: port,
				Zone: "",
			})
			if err != nil {
				fmt.Println("Error reading:", err)
				continue
			}
			for {

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					log.Println(err)
					continue
				}
				if string(buf[:n]) == "ACK" {
					//log.Println("Recive ACK", i)
					break
				}

			}
			if n < 1024 {
				log.Println("传输完成")
				break
			}
			i += 1
			//fmt.Println(string(buffer[:n]))
		}
		buf := make([]byte, 1024)
		n, _ = conn.Read(buf)
		if string(buf[:n]) == "WRONG" {
			log.Println("md5校验失败,重新传输")
			file.Close()
			goto flag
		}
	}

}
func BytesCombine3(pBytes ...[]byte) []byte {
	var buffer bytes.Buffer
	for index := 0; index < len(pBytes); index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}
