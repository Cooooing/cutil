package base

import (
	"database/sql"
	"encoding/json"
	"github.com/Cooooing/cutil/common/logger"
	"github.com/Cooooing/cutil/common/str"
	"reflect"
	"strings"
	"time"
)

// Todo 目前采用json序列化方式。后续通过自定义tag反射实现映射。

const FieldTag = "corm"
const FieldTagColumn = "column"
const FieldTagComment = "comment"
const FieldTagPrimaryKey = "primaryKey"

type FieldMeta struct {
	Field     reflect.StructField
	Column    string
	IsPrimary bool
	Comment   string
}

// parseCormTag 解析 corm:"..." 标签
func parseCormTag(sf reflect.StructField) FieldMeta {
	tag := sf.Tag.Get(FieldTag)
	meta := FieldMeta{
		Field:  sf,
		Column: str.ToSnakeCase(sf.Name), // 默认用蛇形命名字段名
	}
	if tag == "" {
		return meta
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		key := strings.TrimSpace(kv[0])

		if len(kv) == 1 {
			// 单独的flag（如 primaryKey）
			switch strings.ToLower(key) {
			case FieldTagPrimaryKey:
				meta.IsPrimary = true
			}
			continue
		}

		val := strings.TrimSpace(kv[1])
		switch strings.ToLower(key) {
		case FieldTagColumn:
			meta.Column = val
		case FieldTagComment:
			meta.Comment = val
		}
	}

	return meta
}

func Raw2Struct[T any](columns []string, rows *sql.Rows) (*T, error) {
	var (
		err  error
		item = new(T)
	)
	// 准备扫描容器
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 扫描一行
	if err = rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// 遍历 struct 字段并赋值
	v := reflect.ValueOf(item).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		sf := t.Field(i)
		meta := parseCormTag(sf)
		col := strings.ToLower(meta.Column)

		// 找到列索引
		var idx = -1
		for j, c := range columns {
			if strings.ToLower(c) == col {
				idx = j
				break
			}
		}
		if idx == -1 {
			continue
		}

		val := values[idx]
		if val == nil {
			continue
		}

		fv := reflect.ValueOf(val)

		// []byte 转 string
		if b, ok := val.([]byte); ok {
			if field.Kind() == reflect.String {
				field.SetString(string(b))
				continue
			}
			if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String {
				s := string(b)
				field.Set(reflect.ValueOf(&s))
				continue
			}
		}

		// 普通赋值（支持指针字段）
		if field.Kind() == reflect.Pointer {
			ptr := reflect.New(field.Type().Elem())
			if fv.Type().AssignableTo(field.Type().Elem()) {
				ptr.Elem().Set(fv)
			} else if fv.Type().ConvertibleTo(field.Type().Elem()) {
				ptr.Elem().Set(fv.Convert(field.Type().Elem()))
			}
			field.Set(ptr)
		} else {
			if fv.Type().AssignableTo(field.Type()) {
				field.Set(fv)
			} else if fv.Type().ConvertibleTo(field.Type()) {
				field.Set(fv.Convert(field.Type()))
			}
		}
	}
	return item, nil
}

func Raw2StructByPage[T any](db *sql.DB, page PageReqInterface, query string, args ...any) ([]*T, error) {
	list, err := Raw2MapByPage(db, page, query, args...)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	var result []*T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func Raw2MapByPage(db *sql.DB, page PageReqInterface, query string, args ...any) ([]*map[string]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	list := make([]*map[string]any, 0, page.GetSize())

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("close rows failed: %w", err)
		}
	}(rows)
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	start := (page.GetPage() - 1) * page.GetSize()
	end := start + page.GetSize()

	current := 0
	for rows.Next() {
		current++
		// 跳过不需要的记录
		if current < start {
			continue
		}
		// 达到分页结束位置，停止遍历
		if current >= end {
			break // 提前终止，减少后续数据传输
		}

		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		data := make(map[string]any, len(columns))
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		for _, key := range columns {
			switch v := values[columnMap[key]].(type) {
			case []byte:
				data[key] = string(v)
			case nil, string, int64, int32, int16, float64, float32, bool, time.Time:
				data[key] = v
			default:
				data[key] = v
			}
		}
		list = append(list, &data)
	}
	return list, nil
}

func Raws2Struct[T any](db *sql.DB, query string, args ...any) ([]*T, error) {
	var (
		err  error
		list []*T
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item, err := Raw2Struct[T](columns, rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil

	// list, err := Raw2Map(db, query, args...)
	// if err != nil {
	// 	return nil, err
	// }
	// bytes, err := json.Marshal(list)
	// if err != nil {
	// 	return nil, err
	// }
	// var result []T
	// err = json.Unmarshal(bytes, &result)
	// if err != nil {
	// 	return nil, err
	// }
	// return result, err
}

func Raw2Map(db *sql.DB, query string, args ...any) ([]*map[string]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	list := make([]*map[string]any, 0)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		data := make(map[string]any, len(columns))
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		for _, key := range columns {
			switch v := values[columnMap[key]].(type) {
			case []byte:
				data[key] = string(v)
			case nil, string, int64, int32, int16, float64, float32, bool, time.Time:
				data[key] = v
			default:
				data[key] = v
			}
		}
		list = append(list, &data)
	}
	return list, nil
}
