package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"zhaoxin2025/common"
	"zhaoxin2025/model"
	"zhaoxin2025/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Interv struct{}

// 查询面试记录
func (*Interv) Get(c *gin.Context) {
	// 绑定查询参数
	var info struct {
		ID          int              `form:"id" binding:"omitempty"`
		Department  model.Department `form:"department,omitempty"`
		Interviewer string           `form:"interviewer,omitempty"`
		Pass        int              `form:"pass" binding:"omitempty,oneof=0 1"`
		Date        time.Time        `form:"date" binding:"omitempty"`
		common.PagerForm
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var interv model.Interv
	if err := copier.Copy(&interv, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	// 调用服务层获取数据
	total, data, err := srv.Interv.Get(interv, info.Page, info.Limit)
	if err != nil {
		c.Error(err)
		return
	}
	// 返回查询结果
	c.JSON(http.StatusOK, ResponseNew(c, struct {
		Total int64          `json:"total"`
		Data  []model.Interv `json:"data"`
	}{
		Total: total,
		Data:  data,
	}))
}

// New 新建面试时间
func (*Interv) New(c *gin.Context) {
	// 绑定请求体
	var info struct {
		TimeRange service.TimeRange `json:"timerange" binding:"required"`
		Interval  int               `json:"interval" binding:"required"`
		// 以分钟为单位计算
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if info.TimeRange.Start.Add(time.Duration(info.Interval) * time.Minute).After(info.TimeRange.End) {
		c.Error(common.ErrNew(errors.New("时间范围不够一个间隔"), common.ParamErr))
	}
	// 创建面试时间
	data, err := srv.Interv.New(info.TimeRange, info.Interval)
	if len(data) > 0 {
		c.JSON(http.StatusConflict, ResponseNew(c, struct {
			Conflict []time.Time `json:"conflict"`
			Error    error       `json:"error"`
		}{
			Conflict: data,
			Error:    errors.New("有时间冲突的面试记录"),
		}))
		return
	} else if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 更新面试记录
// 更新评价、是否通过、星级等信息
func (*Interv) Update(c *gin.Context) {
	// 绑定请求体
	var info struct {
		ID          int              `json:"id" binding:"required"`
		NetID       string           `json:"netid" binding:"required,numeric,len=10"`
		Interviewer string           `json:"interviewer" binding:"omitempty"`
		Time        time.Time        `json:"time" binding:"omitempty"`
		Department  model.Department `json:"department" binding:"omitempty,oneof=tech video art"`
		Star        int              `json:"star" binding:"omitempty"`
		Evaluation  string           `json:"evaluation" binding:"omitempty"`
		Pass        int              `json:"pass" binding:"omitempty,oneof=0 1"`
		QueID       int              `json:"queid" binding:"omitempty"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	var interv model.Interv
	if err := copier.Copy(&interv, &info); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	// 更新面试记录
	if err := srv.Interv.Update(interv); err != nil {
		c.Error(err)
		return
	}
	// 响应
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// 删除面试记录
func (*Interv) Delete(c *gin.Context) {
	// 绑定请求体
	var info struct {
		ID []int `form:"id" binding:"required,min=1,dive,gte=1"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 删除面试记录
	err := srv.Interv.Delete(info.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

// Swap 交换面试记录
func (*Interv) Swap(c *gin.Context) {
	var info struct {
		ID1 int `json:"id1" binding:"required"`
		ID2 int `json:"id2" binding:"required"`
	}
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	if info.ID1 == info.ID2 {
		c.Error(common.ErrNew(errors.New("不能交换同一条记录"), common.ParamErr))
		return
	}
	if err := srv.Interv.Swap(info.ID1, info.ID2); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

func (*Interv) GetQue(c *gin.Context) {
	var info struct {
		NetID      string           `form:"netid" binding:"required,numeric,len=10"`
		Department model.Department `form:"department" binding:"required,oneof=tech video art"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	data, err := srv.Interv.GetQue(info.NetID, info.Department)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, data))
}

func (*Interv) BlockAndRecover(c *gin.Context) {
	var info struct {
		TimeRange service.TimeRange `json:"timerange" binding:"required"`
		Block     int               `json:"block" binding:"oneof=0 1"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	// 调用服务层获取数据
	err := srv.Interv.BlockAndRecover(info.TimeRange, info.Block == 1)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}

func (*Interv) QQGroup(c *gin.Context) {
	var info struct {
		URL        string           `json:"url" binding:"required,url"`
		QQGroup    string           `json:"qqgroup" binding:"required,numeric"`
		Department model.Department `json:"department" binding:"omitempty,oneof=tech video art none"`
	}
	if err := c.ShouldBind(&info); err != nil {
		c.Error(common.ErrNew(err, common.ParamErr))
		return
	}
	targetDir := "QQGroup"
	fileName := fmt.Sprintf("%s.json", info.Department)
	path := filepath.Join(targetDir, fileName)
	fmt.Printf("尝试在路径 %s 创建文件\n", path)

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		c.Error(common.ErrNew(err, common.SysErr))
		return
	}
	fmt.Printf("目录 %s 已存在或创建成功\n", targetDir)

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	") // 设置缩进
	err = encoder.Encode(info)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return
	}

	fmt.Printf("Successfully encoded JSON to %s\n", fileName)
	c.JSON(http.StatusOK, ResponseNew(c, nil))
}
