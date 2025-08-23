package excel

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Cooooing/cutil/common/logger"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

// CheckFunc 校验函数
//
// 参数:
//   - value: 写入单元格的值
//   - key: 写入单元格对应的键
//
// 返回:
//   - error: 校验失败的错误信息
type CheckFunc func(key string, value any) error

// File excel文件对象
type File struct {
	fileName     string
	file         *excelize.File
	checkFunc    CheckFunc
	errorStyleId *int
	titleStyleId *int
}

func NewFile(fileName string) *File {
	f := &File{fileName: fileName, file: excelize.NewFile()}
	// 添加默认样式
	if f.errorStyleId == nil {
		_ = f.SetErrorStyle(nil)
	}
	return f
}

func (f *File) SetCheckFunc(checkFunc CheckFunc) {
	f.checkFunc = checkFunc
}

func (f *File) SetErrorStyle(style *excelize.Style) error {
	if style == nil {
		styleId, err := f.file.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#FFFF00"}, // 黄色背景
				Pattern: 1,                   // 实心填充
			},
		})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create style: %v", err))
		}
		f.errorStyleId = &styleId
	} else {
		styleId, err := f.file.NewStyle(style)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to create style: %v", err))
		}
		f.errorStyleId = &styleId
	}
	return nil
}

func (f *File) SetTitleStyle(style *excelize.Style) error {
	if style == nil {
		return errors.New("style cannot be nil")
	}
	styleId, err := f.file.NewStyle(style)
	if err != nil {
		return err
	}
	f.titleStyleId = &styleId
	return nil
}

func (f *File) GetExcelizeFile() *excelize.File {
	return f.file
}

func (f *File) checkParams(titles []string, keys []string) error {
	if len(titles) != len(keys) {
		return errors.New("titles and keys length must be equal")
	}
	return nil
}

// WriteToFile 将excel数据写入指定位置文件
//
// 参数:
//   - path: 文件路径
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) WriteToFile(path string) error {
	if f.fileName == "" {
		return errors.New("file name is empty")
	}
	file, err := os.Create(filepath.Join(path, f.fileName))
	if err != nil {
		return err
	}
	return f.file.Write(file)
}

// WriteToCellByCheck 将值写入 Excel 单元格，并对写入值进行校验
//
// 参数:
//   - sheetName: 工作表名称
//   - colIndex: 列索引（从 1 开始）
//   - rowIndex: 行索引（从 1 开始）
//   - value: 要写入的值
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) WriteToCellByCheck(sheetName string, rowIndex int, colIndex int, key string, value any) error {
	if f.checkFunc != nil {
		if checkErr := f.checkFunc(key, value); checkErr != nil {
			colLetter, err := excelize.ColumnNumberToName(colIndex)
			if err != nil {
				return errors.New(fmt.Sprintf("Conversion column index %d failed: %v", colIndex, err))
			}
			cell := fmt.Sprintf("%s%d", colLetter, rowIndex)
			if err := f.file.SetCellStyle(sheetName, cell, cell, *f.errorStyleId); err != nil {
				return errors.New(fmt.Sprintf("Failed to set cell style: %v", err))
			}
			// 添加批注
			if err := f.file.AddComment(sheetName, excelize.Comment{
				Cell:   cell,
				Author: "System",
				Text:   fmt.Sprintf("%v", checkErr.Error()),
			}); err != nil {
				return errors.New(fmt.Sprintf("Failed to add comment: %v", err))
			}

			if err := f.file.SetCellValue(sheetName, cell, value); err != nil {
				return errors.New(fmt.Sprintf("Failed to set cell %s: %v", cell, err))
			}
		}
	}
	return f.WriteToCell(sheetName, rowIndex, colIndex, value)
}

// WriteToCell 将值写入 Excel 单元格
//
// 参数:
//   - sheetName: 工作表名称
//   - colIndex: 列索引（从 1 开始）
//   - rowIndex: 行索引（从 1 开始）
//   - value: 要写入的值
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) WriteToCell(sheetName string, rowIndex int, colIndex int, value any) error {
	colLetter, err := excelize.ColumnNumberToName(colIndex)
	if err != nil {
		return errors.New(fmt.Sprintf("Conversion column index %d failed: %v", colIndex, err))
	}
	cell := fmt.Sprintf("%s%d", colLetter, rowIndex)
	if err := f.file.SetCellValue(sheetName, cell, value); err != nil {
		return errors.New(fmt.Sprintf("Failed to set cell %s: %v", cell, err))
	}
	return nil
}

// ExportFromDataMap 从数据集合导出数据
//
// 参数:
//   - sheetName: 工作表名称
//   - titles: 标题行
//   - keys: 数据行对应的键
//   - values: 数据键值map集合数组
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) ExportFromDataMap(sheetName string, titles []string, keys []string, values []map[string]any) error {
	if err := f.checkParams(titles, keys); err != nil {
		return err
	}

	// 写入标题行
	for i, title := range titles {
		colIndex, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}
		_ = f.file.SetCellStr(sheetName, fmt.Sprintf("%s%d", colIndex, 1), title)
	}
	if f.titleStyleId != nil {
		titleEndCell, err := excelize.ColumnNumberToName(len(titles))
		if err != nil {
			return err
		}
		_ = f.file.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", titleEndCell, 1), *f.errorStyleId)
	}

	// 写入数据
	for i, value := range values {
		for j, key := range keys {
			err := f.WriteToCell(sheetName, i+2, j+1, value[key])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ExportFromQuery 从数据库查询数据并导出
//
// 参数:
//   - sheetName: 工作表名称
//   - titles: 标题行
//   - keys: 数据行对应的键（列名）
//   - db: 数据源
//   - query: 查询sql
//   - args: 查询参数
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) ExportFromQuery(sheetName string, titles []string, keys []string, db *sql.DB, query string, args ...any) error {
	if err := f.checkParams(titles, keys); err != nil {
		return err
	}

	// 写入标题行
	for i, title := range titles {
		colIndex, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.file.SetCellStr(sheetName, fmt.Sprintf("%s%d", colIndex, 1), title)
	}
	if f.titleStyleId != nil {
		titleEndCell, err := excelize.ColumnNumberToName(len(titles))
		if err != nil {
			return err
		}
		_ = f.file.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s%d", titleEndCell, 1), *f.errorStyleId)
	}

	// 写入数据
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}
	for _, key := range keys {
		if _, exists := columnMap[key]; !exists {
			logger.Warn("column %s not exists", key)
		}
	}
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	rowIndex := 1
	for rows.Next() {
		rowIndex++
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}
		for i, key := range keys {
			value := values[columnMap[key]]
			if err := f.WriteToCellByCheck(sheetName, rowIndex, i+1, key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExportStreamFromDataMap 从数据源导出数据（流式，适合大规模数据）
//
// 参数:
//   - sheetName: 工作表名称
//   - titles: 标题行
//   - keys: 数据行对应的键（列名）
//   - values: 数据键值map集合数组
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) ExportStreamFromDataMap(sheetName string, titles []string, keys []string, values []map[string]any) error {
	if err := f.checkParams(titles, keys); err != nil {
		return err
	}
	sw, err := f.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	// 写入标题行
	titleData := make([]any, len(titles))
	for i, title := range titles {
		cell := excelize.Cell{Value: title}
		if f.titleStyleId != nil {
			cell.StyleID = *f.titleStyleId
		}
		titleData[i] = cell
	}
	err = sw.SetRow("A1", titleData)
	if err != nil {
		return err
	}

	// 按行写入数据
	for i, row := range values {
		rowIndex := 2 + i
		cellIndex, err := excelize.CoordinatesToCellName(1, rowIndex)
		if err != nil {
			return err
		}
		data := make([]any, len(keys))
		for j, key := range keys {
			cell := excelize.Cell{Value: row[key]}
			if f.checkFunc != nil && f.errorStyleId != nil {
				if err := f.checkFunc(key, row[key]); err != nil {
					cell.StyleID = *f.errorStyleId
				}
			}
			data[j] = cell
		}
		// 写入一行
		if err := sw.SetRow(cellIndex, data); err != nil {
			return err
		}
	}

	// 结束流式写入
	if err := sw.Flush(); err != nil {
		return err
	}
	return nil
}

// ExportStreamFromQuery 从数据库查询数据并导出（流式，适合大规模数据）
//
// 参数:
//   - sheetName: 工作表名称
//   - titles: 标题行
//   - keys: 数据行对应的键（列名）
//   - db: 数据源
//   - query: 查询sql
//   - args: 查询参数
//
// 返回:
//   - error: 写入失败的错误信息
func (f *File) ExportStreamFromQuery(sheetName string, titles []string, keys []string, db *sql.DB, query string, args ...any) error {
	if err := f.checkParams(titles, keys); err != nil {
		return err
	}
	sw, err := f.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	// 写入标题行
	titleData := make([]any, len(titles))
	for i, title := range titles {
		cell := excelize.Cell{Value: title}
		if f.titleStyleId != nil {
			cell.StyleID = *f.titleStyleId
		}
		titleData[i] = cell
	}
	err = sw.SetRow("A1", titleData)
	if err != nil {
		return err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}
	for _, key := range keys {
		if _, exists := columnMap[key]; !exists {
			logger.Warn("column %s not exists", key)
		}
	}
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 按行写入数据
	rowIndex := 1
	for rows.Next() {
		rowIndex++
		cellIndex, err := excelize.CoordinatesToCellName(1, rowIndex)
		if err != nil {
			return err
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}
		data := make([]any, len(keys))
		for j, key := range keys {
			cell := excelize.Cell{Value: values[columnMap[key]]}
			if f.checkFunc != nil && f.errorStyleId != nil {
				if err := f.checkFunc(key, values[columnMap[key]]); err != nil {
					cell.StyleID = *f.errorStyleId
				}
			}
			data[j] = cell
		}
		// 写入一行
		if err := sw.SetRow(cellIndex, data); err != nil {
			return err
		}
	}

	// 结束流式写入
	if err := sw.Flush(); err != nil {
		return err
	}
	return nil
}
