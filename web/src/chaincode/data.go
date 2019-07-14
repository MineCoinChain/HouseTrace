package main

// 房源信息
type RentingHouseInfo struct {
	RentingID string    `json:RentinfID`  //　统一编码
	HouseInfo `json:RentHouseInfo`    			//　房屋信息， fgj
	AreaInfo  `json:RentAreaInfo`     		//　区域信息, tgj
	OrderInfo `json:RentingOrderInfo` 		//　订单信息, zfpt
}

// 社区信息 otgj
type AreaInfo struct {
	AreaID       string `json:"area_id"`        //社区编号
	AreaAddress  string `json:"area_address"`   // 房源所在区域
	BasicNetWork string `json:"basic_net_work"` // 区域基础网络编号
	CPoliceName  string `json:"c_police_name"`  // 社区民警姓名
	CPoliceNum   string `json:"c_police_num"`   //社区民警工号
}

// 房屋信息  ofgj
type HouseInfo struct {
	HouseID    string `json:"house_id"`    // 房产证编号
	HouseOwner string `json:"house_owner"` // 房主
	RegDate    string `json:"reg_date"`    // 登记日期
	HouseArea  string `json:"house_area"`  // 住房面积
	HouseUsed  string `json:"house_used"`  // 房屋设计用途
	IsMortgage string `json:"is_mortgage"` // 是否抵押
}

// 订单信息	oagency
type OrderInfo struct {
	DocHash   string `json:"doc_hash"`   // 电子合同Hash
	OrderID   string `json:"order_id"`   // 订单编号
	RenterID  string `json:"renter_id"`  // 承租人信息
	RentMoney string `json:"rent_money"` // 租金
	BeginDate string `json:"begin_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	Note      string `json:"note"`       // 备注
}

// 承租人信息  ozxzx
type RenterInfo struct {
	RenterID   string `json:"renter_id"`   // 承租人ID
	RenterDesc string `json:"renter_desc"` //承租人描述 【可以由征信部门提供数据】
}