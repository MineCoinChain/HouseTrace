package controllers

import (
	"web/models/application"
	"encoding/json"
	"github.com/astaxie/beego"
	"fmt"
	"web/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}
func (c *MainController) SetHouseInfo() {
	c.TplName = "setHouseInfo.html"
	fmt.Println("SetHouseInfo")

	// 获取用户输入
	rentingID := c.GetString("rentingID")
	fczbh := c.GetString("fczbh")
	fzxm := c.GetString("fzxm")
	djrq := c.GetString("djrq")
	zfmj := c.GetString("zfmj")
	fwsjyt := c.GetString("fwsjyt")
	sfdy := c.GetString("sfdy")

	if rentingID == "" {
		fmt.Println("rentingID不应为空!")
		return
	}

	// 组织参数
	var args []string
	args = append(args, "addHouseInfo")
	args = append(args, rentingID)
	args = append(args, fczbh)
	args = append(args, fzxm)
	args = append(args, djrq)
	args = append(args, zfmj)
	args = append(args, fwsjyt)
	args = append(args, sfdy)

	fmt.Printf("添加房屋信息输入的数据args :%s\n", args)

	//TODO 添加到fabric账本中
	ret, err := application.App.AddHouseItem(args)
	if err != nil {
		fmt.Println("AddHouseItem err...")
	}

	fmt.Println("<--- 添加房源信息结果　--->：", ret)
}
func (hcc *MainController) GetHouseInfo() {
	hcc.TplName = "getHouseInfo.html"

	//　１、获取用户输入的  rentingID
	key := hcc.GetString("rentingId")

	//　２、组织　chiancode 所需要的参数
	var args []string
	//连代码调用参数
	args = append(args, "getHouseInfo")
	//rentingID
	args = append(args, key)

	// 3、调用　model 层函数，查询数据
	response, err := application.App.GetHouseInfo(args)
	if err != nil {
		fmt.Println("models.App.GetHouseInfo err....")
	}

	//解析结果集
	var jsonData models.HouseInfo
	err = json.Unmarshal([]byte(response), &jsonData)
	if err != nil {
		fmt.Println("json.Unmarshal err....")
	}

	fmt.Println("----------- jsonData", jsonData)

	// 5、将数据展示在前端界面
	hcc.Data["houseId"] = jsonData.HouseID
	hcc.Data["houseOwner"] = jsonData.HouseOwner
	hcc.Data["regDate"] = jsonData.RegDate
	hcc.Data["houseArea"] = jsonData.HouseArea
	hcc.Data["houseUsed"] = jsonData.HouseUsed
	hcc.Data["isMortgage"] = jsonData.IsMortgage
}
