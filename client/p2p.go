package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	mymd5 "p2p/hash"
	"p2p/processfile"
	"strconv"
	"strings"
)

type Client struct {
	Name    string
	Pass    string
	Files   []string
	UDPaddr *net.UDPAddr
}

func main() {
	if len(os.Args) < 7 {
		fmt.Println("./clients tag remoteIP remotePort")
		return
	}
	form := url.Values{}
	//登录或注册用户名
	name := os.Args[1]
	form.Set("name", name)
	//密码
	pass := os.Args[2]

	form.Set("pwd", pass)
	//服务器IP
	remoteIP := os.Args[3]
	//服务器端口
	remotePort := os.Args[4]
	//本地绑定端口
	port, _ := strconv.Atoi(os.Args[5])
	//想下载的文件
	file := os.Args[6]
	localAddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}
	form.Set("udpPort", os.Args[3])
	form.Set("udpIp", "127.0.0.1")
	resp, err := http.PostForm("http://"+remoteIP+":"+remotePort, form)
	if err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	var m map[string]Client
	err = json.Unmarshal(buf[:n], &m)
	if err != nil {
		log.Println(err)
		return
	}
	if len(m) <= 1 {
		log.Println("暂时没有其它节点")
		return
	}
	//获得最佳节点
	var clients []Client
	delete(m, name)

	for _, v := range m {
		clients = append(clients, v)
	}

	log.Println(clients)
	for _, c := range clients {
		has := c.has(file)
		toaddr := c.UDPaddr
		log.Println(has)
		//进行p2p链接
		// 验证是否失效

		//进行p2p下载
		if len(has) == 0 {
			continue
		} else {
			//会导致重复下载相同的片

			for _, fileName := range has {
				f, err := os.Open(fileName)
				f.Close()
				if len(fileName) == 0 || err == nil {
					continue
				}
				p2p(&localAddr, toaddr, fileName)
			}
		}

	}
	err = processfile.MakeFile(file)
	if err != nil {
		log.Println(err)
	}
}
func FindBestClient(clients []Client, file string) map[string]string {
	var count = [10]int{0}
	//键：客户端的名字，值：需要的file的名字
	needClientFiles := make(map[string]string)
	//优先下载最稀有的切片
	for _, c := range clients {
		files := c.has(file)
		for _, f := range files {
			split := strings.Split(f, ".")
			name := split[0]
			nameslice := strings.Split(name, "-")
			if nameslice[0] == file {
				switch nameslice[1] {
				case "0":
					count[0] += 1
					if count[0] == 1 {
						needClientFiles[c.Name] = f
					}
				case "1":
					count[1] += 1
					if count[1] == 1 {
						needClientFiles[c.Name] = f
					}
				case "2":
					count[2] += 1
					if count[2] == 1 {
						needClientFiles[c.Name] = f
					}
				case "3":
					count[3] += 1
					if count[3] == 1 {
						needClientFiles[c.Name] = f
					}
				case "4":
					count[4] += 1
					if count[4] == 1 {
						needClientFiles[c.Name] = f
					}
				case "5":
					count[5] += 1
					if count[5] == 1 {
						needClientFiles[c.Name] = f
					}
				case "6":
					count[6] += 1
					if count[6] == 1 {
						needClientFiles[c.Name] = f
					}
				case "7":
					count[7] += 1
					if count[7] == 1 {
						needClientFiles[c.Name] = f
					}
				case "8":
					count[8] += 1
					if count[8] == 1 {
						needClientFiles[c.Name] = f
					}
				case "9":
					count[9] += 1
					if count[9] == 1 {
						needClientFiles[c.Name] = f
					}
				default:
					continue
				}

			}
		}
	}

	return needClientFiles
}
func isNeed(f string, file string) bool {
	fi := strings.Split(f, ".")
	fiwant := strings.Split(file, ".")
	//存在漏洞
	if len(fi[0]) <= 2 {
		return false
	}
	if fi[0][:len(fi[0])-2] == fiwant[0][:len(fiwant[0])] {
		return true
	}
	return false
}
func (client Client) has(file string) []string {
	files := make([]string, 1)
	for _, f := range client.Files {
		if isNeed(f, file) {
			files = append(files, f)
		}
	}
	return files
}
func p2p(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr, fileName string) {
	//请求建立联系

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
flag:
	log.Println("开始下载", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		file.Close()
		return
	}

	defer func() {
		file.Close()
		r := recover()
		if r != nil {
			log.Println(r)
			return
		}
	}()
	//设置超时
	if err != nil {
		log.Panic(err)
	}

	//告知想要哪个文件
	conn.Write([]byte(fileName))
	//启动goroutine监控标准输入
	buf := make([]byte, 1024)
	oob := make([]byte, 1024)

	var (
		i     = 0
		bar   *progressbar.ProgressBar
		calmd []byte
	)

	for {
		//接受UDP消息打印
		n, _, _, _, err := conn.ReadMsgUDP(buf, oob)
		if err != nil {
			log.Println(err)
			return
		}
		if i == 0 {
			b1 := buf[0]
			b2 := buf[1]
			b3 := buf[2]
			b4 := buf[3]

			calmd = buf[4:n]
			// 开始进度条
			size := int32(b1)<<24 | int32(b2)<<16 | int32(b3)<<8 | int32(b4)
			//log.Println(size)
			bar = progressbar.DefaultBytes(
				int64(size),
				"downloading",
			)

			_, err = conn.Write([]byte("ACK"))
			if err != nil {
				log.Println(err)
				return
			}
			i++
			continue
		}
		if n > 0 {
			//file.Write(buf[:n])
			io.Copy(io.MultiWriter(file, bar), bytes.NewReader(buf[:n]))

		}
		_, err = conn.Write([]byte("ACK"))
		if err != nil {
			log.Println(err)
			return
		}
		if n < 1024 {
			break
		}
		i++
	}

	hash, _ := mymd5.HashSHA256File(fileName)
	if hash == string(calmd) {
		conn.Write([]byte("WRONG"))
		log.Println("md5校验失败，重新下载")
		goto flag
	} else {
		conn.Write([]byte("OK"))
	}
	file.Close()
	conn.Close()
}
