package main

import (
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"os"
	"strconv"
)

func main() {

	for {
		fmt.Println("===== key-value 存储器 =====")
		fmt.Println("1. 命令行模块")
		fmt.Println("2. 程序的使用说明书")
		fmt.Println("3. 退出系统")
		fmt.Print("请输入选项 (1-3): ")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("无效输入，请输入数字 1-3")
		}

		switch choice {
		case 1:
			commandLine()
		case 2:
			usage()
		case 3:
			os.Exit(0)
		default:
			fmt.Println("无效选项，请输入数字 1-3")
		}
	}

}

func commandLine() {
	shell := ishell.New()
	loggedIn := true
	shell.Println("命令行...输入help查看帮助")
	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "登录(用户名和密码为admin)",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			c.Print("用户名: ")
			username := c.ReadLine()
			c.Print("密码: ")
			password := c.ReadPassword()
			if username == "admin" && password == "admin" {
				loggedIn = true
				c.Println("登录成功")
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "set (key) (value) (过期时间)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 3 {
				c.Println("参数错误，请输入 set (key) (value) (过期时间)")
			} else {
				i64, err := strconv.ParseInt(c.Args[2], 10, 64)
				if err != nil {
					c.Println("参数错误，请输入数字作为过期时间")
				}
				manager.setKey(c.Args[0], c.Args[1], i64)
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "get",
		Help: "get (key)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 1 {
				c.Println("参数错误，请输入 get (key)")
			} else {
				value, err := manager.getKey(c.Args[0])
				if err != nil {
					c.Println(err)
				}
				c.Println(value)
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "setnx",
		Help: "setnx (key) (value) (过期时间)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 3 {
				c.Println("参数错误，请输入 setnx (key) (value) (过期时间)")
			} else {
				i64, err := strconv.ParseInt(c.Args[2], 10, 64)
				if err != nil {
					c.Println("参数错误，请输入数字作为过期时间")
				}
				result := manager.setNXKey(c.Args[0], c.Args[1], i64)
				if result == 1 {
					c.Println("设置成功")
				} else {
					c.Println("设置失败，Key 已存在")
				}
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "del",
		Help: "del (key)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 1 {
				c.Println("参数错误，请输入 del (key)")
			} else {
				manager.delKey(c.Args[0])
			}

		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "sadd",
		Help: "sadd (key) (value)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 2 {
				c.Println("参数错误，请输入 sadd (key) (value)")
			} else {
				manager.saddKey(c.Args[0], c.Args[1])
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "smember",
		Help: "smember (key)",
		Func: func(c *ishell.Context) {
			if !loggedIn {
				c.Println("请先登录")
				return
			}
			manager := NewJSONManager(c)
			if len(c.Args) != 1 {
				c.Println("参数错误，请输入 smember (key)")
			} else {
				c.Println(manager.smember(c.Args[0]))
			}
		},
	})
	shell.Run()
}

func usage() {
	fmt.Println("===== 使用说明书 =====")
	fmt.Println("本程序是一个简单的键值存储器，支持以下操作：")
	fmt.Println()
	fmt.Println("1. 选择命令行模块进行操作")
	fmt.Println("   - 进入命令行模块后，可以使用以下命令:")
	fmt.Println("     - login: 登录系统，用户名和密码均为 'admin'")
	fmt.Println("     - set (key) (value) (expire_time): 设置一个键值对，expire_time 为过期时间的时间戳（单位：秒）")
	fmt.Println("     - get (key): 获取指定键的值，如果键不存在，将返回错误信息")
	fmt.Println("     - setnx (key) (value) (expire_time): 仅在键不存在时设置键值对，防止覆盖已有的键")
	fmt.Println("     - del (key): 删除指定的键，如果键不存在，将返回错误信息")
	fmt.Println("     - sadd (key) (value): 向指定的集合添加元素，如果集合不存在，将自动创建")
	fmt.Println("     - smember (key): 获取指定集合的所有元素")
	fmt.Println()
	fmt.Println("2. 选择使用说明书查看程序的使用方法")
	fmt.Println()
	fmt.Println("3. 选择退出系统结束程序")
}
