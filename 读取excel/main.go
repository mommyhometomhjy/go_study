package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Imagestrut struct {
	ImageUrl string `json:imageUrl`
	Name     string `json:name`
}

func main() {
	//建立通道,保证进程全部执行完
	ch := make(chan string)
	//打开表格
	var count int

	f, err := excelize.OpenFile("./资料.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取第一个工作表的所有行
	rows, err := f.GetRows("Sheet1")

	for index, row := range rows {

		if index == 0 {
			continue
		}
		//解析图片地址为json
		var s []Imagestrut
		// var str string
		json.Unmarshal([]byte(row[2]), &s)

		for i, img := range s {
			//并发下载图片并命名
			if img.ImageUrl == "" {
				continue
			}
			count++
			go getImg(row[0], img.ImageUrl, strconv.Itoa(i)+img.Name+".jpg", ch)

		}

	}

	for a := 0; a < count; a++ {
		<-ch
	}
	fmt.Print("执行完毕")
}

func getImg(dirName string, url string, colorName string, ch chan (string)) (n int64, err error) {

	os.Mkdir(dirName, os.ModePerm)
	out, err := os.Create(dirName + "\\" + colorName)

	defer out.Close()
	// fmt.Println(url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	ch <- "1"
	return
}
