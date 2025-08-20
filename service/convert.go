package service

import (
	"fmt"
	"log"
	"time"
	"zhaoxin2025/model"

	"github.com/xuri/excelize/v2"
)

// TimeRange 定义一个时间范围结构体
// 包含开始时间和结束时间
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// DayRange 接收一个时间参数，返回同一天的开始和结束 TimeRange 结构体
func DayRange(t time.Time) TimeRange {
	return TimeRange{
		Start: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
		End:   time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()),
	}
}

// Date 接收一个时间参数，返回该天的日期
// 例如：2023-10-01 00:00:00 返回 "2023-10-01"
func Date(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func Pass(pass int) string {
	if pass == 1 {
		return "已通过"
	}
	return "未通过"
}

func Time(time time.Time) string {
	return time.Format("2006年01月02日 15:04:05")
}

// Excelize 将 Product 切片数据导出到 XLSX 文件
// data: 要导出的 Stu 切片
// filename: 输出的 XLSX 文件名
func Excelize(data []model.Stu, filename string) error {
	// 1. 创建一个新的 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("关闭 Excel 文件失败: %v", err)
		}
	}()

	// 2. 创建一个工作表 (或者使用默认的 Sheet1)
	sheetName := "学生信息表" // 自定义工作表名称
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("创建 Excel 工作表失败: %w", err)
	}
	f.SetActiveSheet(index) // 设置为当前活动工作表

	// 3. 定义表头 (手动指定，对应结构体字段)
	// 这个顺序决定了 Excel 列的顺序
	headers := []string{"OpenID", "NetID", "名字", "电话号码", "书院", "在哪里知道的", "已经会的技能", "想要掌握的技能", "主选部门", "标签", "面试时间", "面试官", "面试评价", "评级", "问题"}
	// 对应结构体字段: OpenID, NetID, Name, Phone, School, WhereKnow, Mastered, ToMaster, Depart, Tag, Date, Interviewer, Evaluation, Star, Question

	// 4. 写入表头
	for i, header := range headers {
		// 将列索引转换为 Excel 列字母 (A, B, C...)
		cell, err := excelize.CoordinatesToCellName(i+1, 1) // excelize 坐标从 1 开始
		if err != nil {
			return fmt.Errorf("生成单元格名称失败 (表头): %w", err)
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// 可选：设置表头样式 (例如：加粗)
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err == nil { // 忽略样式创建失败的错误
		startCell, _ := excelize.CoordinatesToCellName(1, 1)
		endCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
		f.SetCellStyle(sheetName, startCell, endCell, style)
	}

	// 5. 写入数据行
	for rowIndex, stu := range data {
		// 数据从第 2 行开始 (第 1 行是表头)
		excelRowIndex := rowIndex + 2

		// 准备当前行的数据 (手动从结构体字段中提取，顺序与表头一致)
		rowData := []any{
			stu.OpenID,
			stu.NetID,
			stu.Name,
			stu.Phone,
			stu.School,
			stu.WhereKnow,
			stu.Mastered,
			stu.ToMaster,
			stu.Depart,
			stu.Tag,
		}
		if stu.Interv != nil {
			rowData = append(rowData, (*stu.Interv).Time, (*stu.Interv).Interviewer, (*stu.Interv).Evaluation, (*stu.Interv).Star)
		}
		if stu.QueID != 0 {
			var question string
			if err := model.DB.Model(&model.Que{}).Where("id = ?", stu.QueID).Select("question").Scan(&question).Error; err != nil {
				rowData = append(rowData, question)
			}
		}

		// 写入一行数据到 Excel
		for colIndex, dataValue := range rowData {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, excelRowIndex)
			if err != nil {
				return fmt.Errorf("生成单元格名称失败 (数据行 %d): %w", excelRowIndex, err)
			}
			// SetCellValue 会根据值的类型自动进行转换
			f.SetCellValue(sheetName, cell, dataValue)
		}
	}

	// 可选：自动调整列宽
	for colIndex := 1; colIndex <= len(headers); colIndex++ {
		colName, _ := excelize.ColumnNumberToName(colIndex)
		f.SetColWidth(sheetName, colName, colName, 15) // 设置一个默认宽度
	}

	// 6. 保存 Excel 文件
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("保存 Excel 文件失败: %w", err)
	}

	log.Printf("数据成功导出到 %s", filename)
	return nil
}
