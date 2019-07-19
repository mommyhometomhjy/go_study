package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
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
		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("货号:%s,价格:%f,图片地址:%s,标题:%s,关键字:%s,材质:%s", row[0], price, row[2], row[4], row[5], row[6])
	}
}
