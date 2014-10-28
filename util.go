package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// 从文件file读取数据,按行写入chr, 处理后写入chw
// 方法立即返回,返回值"只读",第一个接受已经读取数据,第二个接受成功写入缓存数据(尚未写入硬盘)
// 必须使用for-range 从chr读取数据, chw"只写"
// 必须使用for-range 从chread与chsucc接受统计数据
func FileDuplex(file string, chr, chw chan string) (chan uint64, chan uint64) {

	fr, err := os.Open(file)
	fatalErrCheck(err)
	fw, err := os.Create(file + "out.csv")
	fatalErrCheck(err)
	br := bufio.NewReader(fr)
	bw := bufio.NewWriter(fw)

	chread := make(chan uint64, 10)
	chsucc := make(chan uint64, 10)
	var readCompleted bool
	var readNum, succNum uint64

	// read
	go func() {

		defer fr.Close()
		defer close(chread)

		for {
			//line, err := br.ReadSlice('\n')
			line, _, err := br.ReadLine()
			if err != nil {
				if err == io.EOF {
					readCompleted = true
					fmt.Println("READING COMPLETED")
					break // --> 触发关闭读统计chan,关闭读文件句柄
				}
				fmt.Println("READ : " + err.Error())
				continue
			} else {
				chr <- string(line) // !!!COPY

				readNum++
				select {
				case chread <- readNum:
				default:
				}
			}
		}
		// 读完毕
	}()

	// write
	go func() {

		defer fw.Close()
		defer bw.Flush()
		defer close(chsucc)

		for line := range chw {
			// 写入错误可以简化掉...
			n, err := bw.WriteString(line)
			if n != len(line) {
				fmt.Println("WRITE : " + err.Error())
			} else {

				succNum++
				select {
				case chsucc <- succNum:
				default:
				}

				// !!!!程序在此退出 条件: csv全部读取完毕 && 读取数目 == 处理数目 , 与重试次数无关...
				if readCompleted && succNum == readNum {
					close(chr) // --> 结束所有 从chr for-range 读取数据的goroutine
					close(chw) // --> 结束自身 for -> 触发关闭写统计chan,flush -> close写文件句柄
				}
			}
		}
		// 写完毕
	}()

	return chread, chsucc
}
