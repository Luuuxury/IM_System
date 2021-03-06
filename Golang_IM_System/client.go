package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// Client 客户端类
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// NewClient 创建新的客户端对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net Dial error", err)
		return nil
	}
	client.conn = conn
	return client
}

//开始菜单
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.共聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围的数字<<<<<")
		return false
	}
}

// PublicChat 公共聊天
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>请输入聊天内容， exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// SelectUsers 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err", err)
		return
	}
}

func (client *Client) PrivatChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名], exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入发送消息内容, exit退出:")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>> 请输入消息内容,exit退出：")
		fmt.Scanln(&remoteName)
	}

}

// UpdateName 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// DealResponse 处理server回应的消息，直接现实到标准输出即可
func (client *Client) DealResponse() {
	// 一旦client.conn 有数据，就直接copy 到 stdout 标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			// 共聊模式
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			client.PrivatChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置默认服务器IP地址(默认 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口 默认为（8888）")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>> 连接服务器失败....")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>> 连接服务器成功....")

	//
	client.Run()
}
