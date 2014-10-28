package main

import (
	"crmhelper/crm2"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	erp      = flag.String("erp", "noErp", "crm erp")
	pwd      = flag.String("pwd", "noPwd", "crm password")
	file     = flag.String("file", "data.csv", "file")
	n        = flag.Uint("n", 100, "how many goroutine")
	proxy    = flag.Bool("proxy", false, "proxy?") // Todo: --> False
	proxyurl = flag.String("proxyurl", "", "proxy format: http://127.0.0.1:1022")
)

func main() {
	flag.Parse()

	//==================================================
	// login
	// n, bufferSize
	crm := crm2.New(*erp, *pwd, int(*n), 23000)
	if *proxy {
		proxyUrl, err := url.Parse(*proxyurl)
		if err != nil {
			fmt.Println("PROXYURL PARSE ERROR")
			flag.Usage()
			os.Exit(-1)
		}
		crm.SetProxy(proxyUrl)
	}
	err := crm.Login()
	fatalErrCheck(err)
	fmt.Println("LOGIN SUCCESS")
	//==================================================

	chr := make(chan string, int(*n))
	chw := make(chan string, int(*n))

	// io & statis
	chread, chsucc := FileDuplex(*file, chr, chw)
	var hasread, hassucc uint64
	go func() {
		for hasread = range chread {
			//fmt.Printf("READ-->%d\n", hasread)
		}
	}()
	go func() {
		for hassucc = range chsucc {
			//fmt.Printf("SUCC-->%d\n", hassucc)
		}
	}()
	fmt.Println("INITIALIZING CSV READER AND WRITER")

	//handle
	var wg sync.WaitGroup
	var i uint
	var empty string
	for ; i < *n; i++ {
		wg.Add(1)
		go func(goid uint) {
			buf := crm.GetBuf()
			defer crm.PutBuf(buf)
			for caseid := range chr {
				html, err := crm.CaseInfo(caseid)
				if err != nil {
					if err != io.EOF {
						fmt.Printf("CASEID[%s] --> %s\n", caseid, err.Error())
					}
					chr <- caseid // 失败后重新尝试...直到成功
				} else {
					tel := crm.FilterTel(*html)
					if tel == nil {
						tel = &empty
					}

					buf.Reset()
					buf.WriteString(caseid)
					buf.WriteByte(',')
					buf.WriteString(*tel)
					buf.WriteByte('\n')
					chw <- buf.String()
				}
			}
			wg.Done()

		}(i)
	}
	fmt.Println("INITIALIZING [", *n, "] GOROUTINES")

	//==================================================
	// 速率统计
	starttime := time.Now()
	lasttime := time.Now()
	var lasthassucc uint64
	// qps统计,不准确,无锁
	go func() {
		chT := time.Tick(time.Second * 3)
		for {
			<-chT
			since := time.Since(lasttime).Seconds() // float64
			if since != 0 {
				fmt.Printf("%d SUCC    %d Q/S\n", hassucc, uint64(hassucc-lasthassucc)/uint64(since))
			}
			lasthassucc = hassucc
			lasttime = time.Now()
		}
	}()

	fmt.Println("DONT CLOSE!!!")
	wg.Wait()
	fmt.Println("COST --> ", time.Since(starttime).Seconds(), "Sec")
	fmt.Println("COMPLETED!!!")
}

func fatalErrCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
