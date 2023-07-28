package main

import (
	"github.com/shirou/gopsutil/disk"
	"time"	
	"net/http"
	"encoding/json"
	"strings"
	"io/ioutil"
	"log"
	"fmt"
	"strconv"
	"github.com/shirou/gopsutil/host"
	"os"
	"errors"
)


type Message struct {
	MsgType string `json:"msgtype"`
	Markdown struct {
		
		Title string `json:"content"`
		

	} `json:"markdown"`
}




func parReady(dir string)  error {
	//fmt.Println(dir+"/1.txt")
	file, err := os.Create(dir+"/1.txt")
	file.Close()
	//fmt.Println(err)
	if err != nil {
		return errors.New("read-only file system") 
		
	} else {
		//fmt.Println(errors.New("Ready"))
		return  errors.New("Ready")

	}	
}


func sendMessage (title,par,unum,hostname,sinfo string) {
	//https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxx
	//http://<proxy_ip>/cgi-bin/webhook/send?key=xxx
	hook := "http://<proxy_ip>/cgi-bin/webhook/send?key=89630204-f3ef-458b-a23b-9213842b98aa"
	var m Message
	m.MsgType = "markdown"
	m.Markdown.Title= title +"\n>分区:" + par + "\n> 当前使用率(%):" + unum + "\n主机名:" + hostname + "\n分区信息:" + sinfo
	jsons, err := json.Marshal(m)
	resp := string(jsons)
	var client = &http.Client{
		Timeout: time.Second * 5,
	}


	rqst, err := http.NewRequest("POST",hook,strings.NewReader(resp))
	if err != nil {
		
		log.Printf("%s",err)
		return
	}
	rqst.Host = "qyapi.weixin.qq.com"
	rqst.Header.Set("Content-Type", "application/json")


	r, err := client.Do(rqst)
	if err != nil {
		log.Printf("%s",err)
		
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s",err)
		
		return
	}
	log.Printf("%s",body)
}

func main() {
	infos, err := disk.Partitions(false)

	fmt.Println()

	if err != nil {
		log.Println(err)
		
	}
	
	for _, info := range infos {

		var par string = info.Mountpoint

		u , _ := disk.Usage(par)

		unum  := int(u.UsedPercent)
		//fmt.Println(u)

		//分区可写检查
		r := parReady(par).Error() //bug

		hostinfo, err := host.Info()
		if err != nil {
			log.Println(err)
		}


		h := hostinfo.Hostname

		var use string = strconv.Itoa(unum)

		if  unum > 85 {
			
			sendMessage("# <font color='warning'>存储空间不足</font>",par,use,h,r)
			

		} else if r != "Ready" {

			sendMessage("# <font color='warning'>存储状态异常</font>",par,use,h,r)

		} else {
			log.Println("磁盘状态正常")
		}

	}
}
