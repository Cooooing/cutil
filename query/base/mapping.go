package base

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Cooooing/cutil/base/logger"
	"github.com/Cooooing/cutil/base/str"
)

// Todo 目前采用json序列化方式。后续通过自定义tag反射实现映射。

const (
	FieldTag           = "corm"
	FieldTagColumn     = "column"
	FieldTagComment    = "comment"
	FieldTagPrimaryKey = "primaryKey"
)

type FieldMeta struct {
	Field     reflect.StructField
	Column    string
	IsPrimary bool
	Comment   string
	Index     int // 在 struct 中的索引
}

var fieldMetaCache sync.Map // map[reflect.Type][]FieldMeta

// getFieldMetas 获取结构体字段元信息（带缓存）
func getFieldMetas(t reflect.Type) []FieldMeta {
	if metas, ok := fieldMetaCache.Load(t); ok {
		return metas.([]FieldMeta)
	}

	var metas []FieldMeta
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" { // 非导出字段跳过
			continue
		}
		meta := parseCormTag(sf)
		meta.Index = i
		metas = append(metas, meta)
	}

	fieldMetaCache.Store(t, metas)
	return metas
}

// parseCormTag 解析 corm:"..." 标签
func parseCormTag(sf reflect.StructField) FieldMeta {
	tag := sf.Tag.Get(FieldTag)
	meta := FieldMeta{
		Field:  sf,
		Column: str.ToSnakeCase(sf.Name), // 默认用蛇形命名
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
			switch strings.ToLower(key) {
			case FieldTagPrimaryKey:
				meta.IsPrimary = true
			}
			continue
		}

		val := strings.TrimSpace(kv[1])
		switch strings.ToLower(key) {
		case FieldTagColumn:
			if val != "" {
				meta.Column = val
			}
		case FieldTagComment:
			meta.Comment = val
		}
	}

	return meta
}

// Raw2Struct 将当前行映射到一个结构体实例（必须保证 rows.Next() 已经被调用成功）
func Raw2Struct[T any](columns []string, rows *sql.Rows) (*T, error) {
	item := new(T)
	v := reflect.ValueOf(item).Elem()
	t := v.Type()

	// 获取字段元信息（带缓存）
	metas := getFieldMetas(t)

	// 准备容器接收 Scan 的值
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 读取当前行
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// 遍历字段
	for _, meta := range metas {
		field := v.Field(meta.Index)

		// 找到列索引
		var idx = -1
		for j, c := range columns {
			if strings.EqualFold(c, meta.Column) {
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

		if err := assignValue(field, val); err != nil {
			return nil, fmt.Errorf("assign field %s failed: %w", meta.Field.Name, err)
		}
	}

	return item, nil
}

// assignValue 负责把数据库返回的值赋给 struct 的字段
func assignValue(field reflect.Value, val any) error {
	// 处理 []byte -> string
	if b, ok := val.([]byte); ok {
		switch {
		case field.Kind() == reflect.String:
			field.SetString(string(b))
			return nil
		case field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String:
			s := string(b)
			field.Set(reflect.ValueOf(&s))
			return nil
		}
		// 可能是时间类型
		if field.Type() == reflect.TypeOf(time.Time{}) {
			tm, err := parseTime(string(b))
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(tm))
			return nil
		}
		if field.Kind() == reflect.Pointer && field.Type().Elem() == reflect.TypeOf(time.Time{}) {
			tm, err := parseTime(string(b))
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(&tm))
			return nil
		}
	}

	// 处理 sql.NullXXX
	switch nv := val.(type) {
	case sql.NullString:
		if nv.Valid {
			if field.Kind() == reflect.String {
				field.SetString(nv.String)
			} else if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf(&nv.String))
			}
		}
		return nil
	case sql.NullInt64:
		if nv.Valid {
			return setNumber(field, nv.Int64)
		}
		return nil
	case sql.NullFloat64:
		if nv.Valid {
			return setNumber(field, nv.Float64)
		}
		return nil
	case sql.NullBool:
		if nv.Valid {
			if field.Kind() == reflect.Bool {
				field.SetBool(nv.Bool)
			} else if field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.Bool {
				field.Set(reflect.ValueOf(&nv.Bool))
			}
		}
		return nil
	}

	// 通用赋值（支持指针字段）
	fv := reflect.ValueOf(val)
	if field.Kind() == reflect.Pointer {
		ptr := reflect.New(field.Type().Elem())
		if fv.Type().AssignableTo(field.Type().Elem()) {
			ptr.Elem().Set(fv)
		} else if fv.Type().ConvertibleTo(field.Type().Elem()) {
			ptr.Elem().Set(fv.Convert(field.Type().Elem()))
		} else {
			return fmt.Errorf("cannot assign %v to %v", fv.Type(), field.Type())
		}
		field.Set(ptr)
	} else {
		if fv.Type().AssignableTo(field.Type()) {
			field.Set(fv)
		} else if fv.Type().ConvertibleTo(field.Type()) {
			field.Set(fv.Convert(field.Type()))
		} else {
			return fmt.Errorf("cannot assign %v to %v", fv.Type(), field.Type())
		}
	}
	return nil
}

// setNumber 辅助方法，用于处理 int64/float64 -> int/float32 等情况
func setNumber(field reflect.Value, num any) error {
	fv := reflect.ValueOf(num)
	if field.Kind() == reflect.Pointer {
		ptr := reflect.New(field.Type().Elem())
		if fv.Type().ConvertibleTo(field.Type().Elem()) {
			ptr.Elem().Set(fv.Convert(field.Type().Elem()))
			field.Set(ptr)
			return nil
		}
	} else {
		if fv.Type().ConvertibleTo(field.Type()) {
			field.Set(fv.Convert(field.Type()))
			return nil
		}
	}
	return fmt.Errorf("cannot set number %v to %v", fv.Type(), field.Type())
}

// parseTime 尝试解析数据库返回的 datetime
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}
	for _, l := range layouts {
		if tm, err := time.Parse(l, s); err == nil {
			return tm, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", s)
}

// -----

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
