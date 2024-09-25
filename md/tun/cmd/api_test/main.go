package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const Cookie = "oauth=%7B%22access_token%22%3A%22cfb24c1b-720b-4b85-a13e-dc75e31c6422%22%2C%22token_type%22%3A%22bearer%22%2C%22refresh_token%22%3A%22bcf98057-0f11-4835-99ee-fc3caf305e1d%22%2C%22scope%22%3A%22read%20write%20trust%22%2C%22expires_date%22%3A%222024-08-02T11%3A41%3A33.169Z%22%7D; username=%22admin%22"

const Host = "http://10.82.44.106"

const access_token = "cfb24c1b-720b-4b85-a13e-dc75e31c6422"

type SSIDList struct {
	Description string `json:"description"`
	Id          int    `json:"id"`
	SSID        string `json:"ssid"`
	UpdateTime  string `json:"update_time"`
}

type AuthObjects struct {
	AuthType     string     `json:"authType"`
	DB1SyncValue int        `json:"db1SyncValue"`
	DB2SyncValue int        `json:"db2SyncValue"`
	GroupName    string     `json:"groupName"`
	Id           int        `json:"id"`
	Location     string     `json:"location"`
	Mac          string     `json:"mac"`
	SSIDList     []SSIDList `json:"ssidList"`
	Status       bool       `json:"status"`
	UserName     string     `json:"userName"`
	UserPass     string     `json:"userPassword"`
}

type UpdataFrom struct {
	AuthType string `json:"authType"`
	Location string `json:"location"`
	Mac      string `json:"mac"`
	SSID     []int  `json:"ssid"`
	Status   bool   `json:"status"`
}

type ObjListItem struct {
	AuthType     string     `json:"authType"`
	Id           int        `json:"id"`
	Location     string     `json:"location"`
	Mac          string     `json:"mac"`
	SSIDList     []SSIDList `json:"ssidList"`
	Status       bool       `json:"status"`
	UserName     string     `json:"userName"`
	UserPassword string     `json:"userPassword"`
}

type ObjList struct {
	Data  []ObjListItem `json:"data"`
	Total int           `json:"total"`
}

// func main() {
// 	access_token := "cfb24c1b-720b-4b85-a13e-dc75e31c6422"
// 	page := 1
// 	for {
// 		fmt.Fprintf(os.Stderr, "page:%d\n", page)
// 		GetObjList, err := GetObjList(page, 1000, access_token)
// 		if err != nil {
// 			println(err.Error())
// 			return
// 		}
// 		if len(GetObjList.Data) == 0 {
// 			break
// 		}
// 		page++
// 		for _, item := range GetObjList.Data {
// 			if len(item.SSIDList) == 1 && item.SSIDList[0].Id == 2 {
// 				fmt.Printf("id:%d\t\t\tMac:%s \n", item.Id, item.Mac)
// 			}
// 		}
// 	}

// }

type updateItem struct {
	Mac string
	ID  int
}

func main() {
	path := "/home/code/tunnel_go/cmd/api_test/mac"
	f, err := os.Open(path)
	if err != nil {
		println(err.Error())
		return
	}
	updateArr := make([]updateItem, 0, 10000)
	regID := regexp.MustCompile(`id:([0-9]+)`)
	regMac := regexp.MustCompile(`Mac:([0-9A-Za-z]+)`)
	// var lines []string
	r := bufio.NewReader(f)
	for {
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		// str := string(bytes)
		regIDStr := regID.FindStringSubmatch(string(bytes))
		regMacStr := regMac.FindStringSubmatch(string(bytes))
		if len(regIDStr) != 2 && len(regMacStr) != 2 {
			println("not match str:", string(bytes))
			continue
		}

		id, err := strconv.Atoi(regIDStr[1])
		if err != nil {
			println(err.Error())
			continue
		}
		if id <= 0 {
			println("id less 0 str:", string(bytes))
			continue
		}

		updateArr = append(updateArr, updateItem{
			Mac: regMacStr[1],
			ID:  id,
		})
	}
	total := len(updateArr)
	for index, item := range updateArr {
		fmt.Printf("%d/%d start mac:%s id:%d \n", index, total, item.Mac, item.ID)
		obj, err := GetAuthObjects(item.ID)
		if err != nil {
			fmt.Printf("error getObjcet mac:%s  id:%d  err:%s  \n", item.Mac, item.ID, err.Error())
			return
		}
		if obj.Mac != item.Mac {
			fmt.Printf("error getObjcet mac not same apiMac:%s  fileMac:%s id:%d    \n", obj.Mac, item.Mac, item.ID)
			return
		}
		if len(obj.SSIDList) != 1 {
			fmt.Printf("warn getObjcet len is not 1 mac:%s  id:%d  \n", item.Mac, item.ID)
			continue
		}
		if obj.SSIDList[0].Id != 2 {
			fmt.Printf("warn getObjcet id not 2 mac:%s  id:%d   \n", item.Mac, item.ID)
			continue
		}
		if obj.Id != item.ID {
			fmt.Printf("error getObjcet id not same fileID:%d  apiId:%d   mac:%s    \n", item.ID, obj.Id, obj.Mac)
			return
		}

		// 获取提交表单的数据
		UpdateFrom := obj.GetUpdateFrom()
		if err = UpdateFrom.Update(obj.Id); err != nil {
			fmt.Printf("error Update err mac:%s  id:%d  err:%s  \n", item.Mac, item.ID, err.Error())
			return
		}
		// println(UpdateFrom.Mac)

	}

}

func GetAuthObjects(id int) (r AuthObjects, err error) {
	client := &http.Client{}
	Url := fmt.Sprintf("%s/device_authority/auth_objects/%d?access_token=%s&_=1722564529509", Host, id, access_token)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return
	}
	// req.Header.Add("Accept", "text/html, */*; q=0.01")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", Cookie)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &r)
	return
}

func GetObjList(page int, pageSize int) (r ObjList, err error) {
	client := &http.Client{}
	Url := fmt.Sprintf("%s/device_authority/auth_objects?page=%d&pageSize=%d&access_token=%s&ssid=2&_=1722564529509", Host, page, pageSize, access_token)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return
	}
	// req.Header.Add("Accept", "text/html, */*; q=0.01")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", Cookie)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &r)
	return
}

func (r AuthObjects) GetUpdateFrom() (u UpdataFrom) {
	u.SSID = make([]int, 0)
	u.AuthType = r.AuthType
	u.Location = r.Location
	u.Mac = r.Mac
	u.Status = r.Status
	u.SSID = append(u.SSID, r.SSIDList[0].Id, 17)
	return
}

func (f UpdataFrom) Update(id int) (err error) {
	// return
	postStrByte, err := json.Marshal(f)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(postStrByte))
	client := &http.Client{}
	Url := fmt.Sprintf("%s/device_authority/auth_objects/%d?access_token=%s", Host, id, access_token)
	req, err := http.NewRequest("PUT", Url, payload)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", Cookie)

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("http code is not 200 is %d ", res.StatusCode)
	}

	return
}
