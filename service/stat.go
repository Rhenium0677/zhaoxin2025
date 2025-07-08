package service

import (
	"strconv"
	"zhaoxin2025/model"
)

type Stat struct {
	Province []Province `json:"province"`
	School   []School   `json:"school"`
	Gender   Gender     `json:"gender"`
	Depart   Depart     `json:"depart"`
	Total    int        `json:"total"`
}

type Province struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type School struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type Gender struct {
	Male    int `json:"male"`
	Femaile int `json:"female"`
}

type Depart struct {
	Tech      int `json:"tech"`
	Art       int `json:"art"`
	Video     int `json:"video"`
	TechPass  int `json:"tech_pass"`
	ArtPass   int `json:"art_pass"`
	VideoPass int `json:"video_pass"`
}

var provinceDict = map[int]string{
	11: "北京市",
	12: "天津市",
	13: "河北省",
	14: "山西省",
	15: "内蒙古自治区",
	21: "辽宁省",
	22: "吉林省",
	23: "黑龙江省",
	31: "上海市",
	32: "江苏省",
	33: "浙江省",
	34: "安徽省",
	35: "福建省",
	36: "江西省",
	37: "山东省",
	41: "河南省",
	42: "湖北省",
	43: "湖南省",
	44: "广东省",
	45: "广西壮族自治区",
	46: "海南省",
	51: "四川省",
	52: "贵州省",
	53: "云南省",
	54: "西藏省",
	50: "重庆市",
	61: "陕西省",
	62: "甘肃省",
	63: "青海省",
	64: "宁夏回族自治区",
	65: "新疆维吾尔自治区",
	83: "台湾省",
	81: "香港特别行政区",
	82: "澳门特别行政区",
	99: "外国",
}

func GetStat(stus []model.Stu, intervs []model.Interv) Stat {
	res := Stat{}
	res.Total = len(stus)

	province := make(map[string]int)
	school := make(map[string]int)
	gender := make(map[string]int)
	depart := make(map[string]int)
	for _, stu := range stus {
		province[provinceDict[toProvince(stu.NetID)]]++
		school[stu.School]++
		depart[string(stu.Depart)]++
		if toGender(stu.NetID) == 1 {
			gender["male"]++
		} else {
			gender["female"]++
		}
	}
	for _, interv := range intervs {
		if interv.Department == model.Tech && interv.Pass == 1 {
			depart["tech_pass"]++
		}
		if interv.Department == model.Art && interv.Pass == 1 {
			depart["art_pass"]++
		}
		if interv.Department == model.Video && interv.Pass == 1 {
			depart["video_pass"]++
		}
	}
	res.Province = make([]Province, 0, res.Total)
	for k, v := range province {
		res.Province = append(res.Province, Province{
			Name:   k,
			Number: v,
		})
	}
	res.School = make([]School, 0, res.Total)
	for k, v := range school {
		res.School = append(res.School, School{
			Name:   k,
			Number: v,
		})
	}
	res.Gender = Gender{
		Male:    gender["male"],
		Femaile: gender["female"],
	}
	res.Depart = Depart{
		Tech:      depart["tech"],
		Art:       depart["art"],
		Video:     depart["video"],
		TechPass:  depart["tech_pass"],
		ArtPass:   depart["art_pass"],
		VideoPass: depart["video_pass"],
	}
	return res
}

// 输入学号,输出省份
func toProvince(netid string) int {
	num := netid[3:5]
	numInt, _ := strconv.Atoi(num)
	return numInt
}

func toGender(netid string) int {
	num := netid[5:6]
	numInt, _ := strconv.Atoi(num)
	if numInt%2 == 1 {
		return 1
	} else {
		return 0
	}
}
