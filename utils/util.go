package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

//该包用于存储一些 类型转化 错误判断等和账本无直接关系但是要用到的函数

// 用于处理错误的通用函数 Handle(err error)
// 作用是在出现错误时，记录错误信息并终止程序的执行
func Handle(err error) {
	if err != nil {
		//在Go语言中，nil是一个预定义的常量，用于表示指针、接口、切片、map、函数和通道的零值或空值。
		//它是一种特殊的零值，表示该类型的变量不指向任何有效的内存地址或不包含任何有效的数值。
		log.Panic(err)
	}
	//如果 err 不为 nil，意味着有错误发生
}

// 函数 将整数转化为16进制 并且最后返回字节切片
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	//如果在写入过程中出现任何问题，会将错误存储在变量 err 中
	Handle(err)         //相当于这两行实现了判断写入缓存区是否报错
	return buff.Bytes() //按字节返回
}
