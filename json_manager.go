package main

import (
	"encoding/json"
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"os"
	"time"
)

// JSONItem 定义一个结构体来表示数据模型
type JSONItem struct {
	Value      string   `json:"value,omitempty"`
	SetValues  []string `json:"set_values,omitempty"`
	ExpireTime int64    `json:"expire_time"`
}

// JSONManager 结构体
type JSONManager struct {
	items map[string]JSONItem // 使用 map 存储 JSONItem
	c     *ishell.Context
}

// NewJSONManager 创建一个新的 JSONManager
func NewJSONManager(c *ishell.Context) *JSONManager {
	manager := &JSONManager{
		items: make(map[string]JSONItem), // 初始化 map
		c:     c,
	}
	manager.readItemsFromFile()
	return manager
}

// 从文件读取数据
func (m *JSONManager) readItemsFromFile() {
	file, _ := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			m.c.Println("关闭文件失败", err)
		}
	}(file)
	err := json.NewDecoder(file).Decode(&m.items)
	if err != nil {
		return
	}
}

// 保存数据到文件
func (m *JSONManager) saveItemsToFile() {
	file, _ := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			m.c.Println("关闭文件失败", err)
		}
	}(file)
	m.c.Println(m.items)
	data, err := json.MarshalIndent(m.items, "", "  ")
	if err != nil {
		m.c.Println("JSON 序列化失败", err)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		m.c.Println("写入文件失败", err)
	}
}

// 取 Key
func (m *JSONManager) getKey(key string) (string, error) {
	m.readItemsFromFile()
	item, exists := m.items[key]
	if !exists {
		return "", fmt.Errorf("未找到Key: %s", key)
	}
	if time.Now().Unix() > item.ExpireTime {
		m.delKey(key)
		return "", fmt.Errorf("已过期，已删除")
	}
	return item.Value, nil
}

// 删除 Key
func (m *JSONManager) delKey(key string) {
	m.readItemsFromFile()
	delete(m.items, key)
	m.c.Println("删除成功")
	m.saveItemsToFile()
}

// 设置 Key
func (m *JSONManager) setKey(key string, value string, expireTime int64) {
	m.readItemsFromFile()
	m.items[key] = JSONItem{
		Value:      value,
		ExpireTime: expireTime*60 + time.Now().Unix(),
	}
	m.saveItemsToFile()
	m.c.Println("设置成功")
}

// 设置 NX Key
func (m *JSONManager) setNXKey(key string, value string, expireTime int64) int {
	m.readItemsFromFile()
	if _, exists := m.items[key]; !exists {
		m.items[key] = JSONItem{
			Value:      value,
			ExpireTime: expireTime*60 + time.Now().Unix(),
		}
		m.saveItemsToFile()
		return 1
	}
	return 0
}

// 添加元素到 Set
func (m *JSONManager) saddKey(setName string, value string) {
	m.readItemsFromFile()
	item, exists := m.items[setName]
	if exists {
		for _, v := range item.SetValues {
			if v == value {
				m.c.Println("元素已经存在于 Set 中")
				return
			}
		}
		item.SetValues = append(item.SetValues, value)
		m.items[setName] = item
		m.saveItemsToFile()
		m.c.Println("添加成功")
		return
	}

	// Set 不存在，创建 Set
	newSetItem := JSONItem{
		SetValues: []string{value},
	}
	m.items[setName] = newSetItem
	m.saveItemsToFile()
	m.c.Println("新 Set 创建成功并添加元素")
}

// 获取 Set 中的所有元素
func (m *JSONManager) smember(setName string) []string {
	m.readItemsFromFile()
	item, exists := m.items[setName]
	if exists {
		return item.SetValues
	}
	m.c.Println("Set 不存在")
	return []string{}
}
