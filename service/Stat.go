package service

import (
	"strconv"
	"zhaoxin2025/model"
)

type Stat struct {
	Province map[string]int `json:"province"`
	School   map[string]int `json:"school"`
	Gender   map[Gender]int `json:"gender"`
	Total    int            `json:"total"`
}
type Gender string

const (
	Male   Gender = "男"
	Female Gender = "女"
)

var provinceDict = map[int]string{
	11: "北京",
	12: "天津",
	13: "河北",
	14: "山西",
	15: "内蒙古",
	21: "辽宁",
	22: "吉林",
	23: "黑龙江",
	31: "上海",
	32: "江苏",
	33: "浙江",
	34: "安徽",
	35: "福建",
	36: "江西",
	37: "山东",
	41: "河南",
	42: "湖北",
	43: "湖南",
	44: "广东",
	45: "广西",
	46: "海南",
	51: "四川",
	52: "贵州",
	53: "云南",
	54: "西藏",
	50: "重庆",
	61: "陕西",
	62: "甘肃",
	63: "青海",
	64: "宁夏",
	65: "新疆",
	83: "台湾",
	81: "香港",
	82: "澳门",
	99: "外国",
}

func GetStat(data []model.Stu) Stat {
	var res Stat
	res.Total = len(data)
	province := make(map[string]int)
	school := make(map[string]int)
	gender := make(map[Gender]int)
	for _, stu := range data {
		province[toProvince(stu.NetID)]++
		school[stu.School]++
		gender[toGender(stu.NetID)]++
	}
	res.Province = province
	res.School = school
	res.Gender = gender
	return res
}

// 输入学号,输出省份
func toProvince(netid string) string {
	num := netid[3:5]
	numInt, _ := strconv.Atoi(num)
	return provinceDict[numInt]
}

func toGender(netid string) Gender {
	num := netid[5:6]
	numInt, _ := strconv.Atoi(num)
	if numInt%2 == 1 {
		return Male
	} else {
		return Female
	}
}
