package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Cooooing/cutil/common/logger"
	"time"
)

type PageReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type PageResp[T any] struct {
	PageReq
	Total int `json:"total"`
	List  []T `json:"list"`
}

func (r *PageReq) Validate() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
}

// QueryCount 查询总数
//
// 参数:
//   - db: 数据库连接
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - int: 总数
//   - error: 校验失败的错误信息
func QueryCount(db *sql.DB, query string, args ...any) (int, error) {
	var total int
	totalSql := fmt.Sprintf("select count(*) as total from (%s) as t", query)
	logger.Info("\nsql:\n%s\nargs:%+v", totalSql, args)
	row := db.QueryRow(totalSql, args...)
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// PageQueryForStruct 通用分页查询，返回封装的结构体列表。使用json进行结构体映射（深度分页效率较低）
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStruct[T any](db *sql.DB, page PageReq, query string, args ...any) (*PageResp[T], error) {
	var err error
	page.Validate()
	pageResp := &PageResp[T]{PageReq: page}
	pageResp.Total, err = QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.List, err = raw2StructByPage[T](db, page, query, args...)
	if err != nil {
		return nil, err
	}
	return pageResp, nil
}

// PageQueryForMap 通用分页查询，返回封装的map集合列表（深度分页效率较低）
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMap(db *sql.DB, page PageReq, query string, args ...any) (*PageResp[map[string]any], error) {
	var err error
	page.Validate()
	pageResp := &PageResp[map[string]any]{PageReq: page}
	pageResp.Total, err = QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.List, err = raw2MapByPage(db, page, query, args...)
	if err != nil {
		return nil, err
	}
	return pageResp, nil
}

func pageQueryForStruct[T any](db *sql.DB, page PageReq, countQuery string, query string, args ...any) (*PageResp[T], error) {
	var err error
	page.Validate()
	pageResp := &PageResp[T]{PageReq: page}
	pageResp.Total, err = QueryCount(db, countQuery, args...)
	if err != nil {
		return nil, err
	}
	pageResp.List, err = raw2Struct[T](db, query, args...)
	if err != nil {
		return nil, err
	}
	return pageResp, nil
}

func pageQueryForMap(db *sql.DB, page PageReq, countQuery string, query string, args ...any) (*PageResp[map[string]any], error) {
	var err error
	page.Validate()
	pageResp := &PageResp[map[string]any]{PageReq: page}
	pageResp.Total, err = QueryCount(db, countQuery, args...)
	if err != nil {
		return nil, err
	}
	pageResp.List, err = raw2Map(db, query, args...)
	if err != nil {
		return nil, err
	}
	return pageResp, nil
}

// PageQueryForStructWithLimitOffset 使用 Limit/Offset 分页查询，返回封装的结构体列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithLimitOffset[T any](db *sql.DB, page PageReq, query string, args ...any) (*PageResp[T], error) {
	return pageQueryForStruct[T](db, page, query, getLimitOffsetQuery(page, query), args...)
}

// PageQueryForMapWithLimitOffset 使用 Limit/Offset 分页查询，返回封装的map集合列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithLimitOffset(db *sql.DB, page PageReq, query string, args ...any) (*PageResp[map[string]any], error) {
	return pageQueryForMap(db, page, query, getLimitOffsetQuery(page, query), args...)
}

func getLimitOffsetQuery(page PageReq, query string) string {
	return fmt.Sprintf(`SELECT t.* FROM (%s) AS t LIMIT %d OFFSET %d`, query, page.Size, (page.Page-1)*page.Size)
}

// PageQueryForStructWithRowNumber 使用 ROW_NUMBER() 窗口函数 分页查询，返回封装的结构体列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithRowNumber[T any](db *sql.DB, page PageReq, query string, args ...any) (*PageResp[T], error) {
	return pageQueryForStruct[T](db, page, query, getRowNumberQuery(page, query), args...)
}

// PageQueryForMapWithRowNumber 使用 ROW_NUMBER() 窗口函数 分页查询，返回封装的map集合列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithRowNumber(db *sql.DB, page PageReq, query string, args ...any) (*PageResp[map[string]any], error) {
	return pageQueryForMap(db, page, query, getRowNumberQuery(page, query), args...)
}

func getRowNumberQuery(page PageReq, query string) string {
	return fmt.Sprintf(`SELECT * FROM (SELECT t.*, ROW_NUMBER() OVER () AS rn FROM (%s) AS t ) AS sub WHERE rn BETWEEN %d AND %d`, query, (page.Page-1)*page.Size+1, page.Page*page.Size)
}

// PageQueryForStructWithFetchOffset 使用 Fetch/Offset 分页查询，返回封装的结构体列表。（SQL 标准语法，与 Limit/Offset 用法一致）
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithFetchOffset[T any](db *sql.DB, page PageReq, query string, args ...any) (*PageResp[T], error) {
	return pageQueryForStruct[T](db, page, query, getFetchOffsetQuery(page, query), args...)
}

// PageQueryForMapWithFetchOffset 使用 Fetch/Offset 分页查询，返回封装的map集合列表。（SQL 标准语法，与 Limit/Offset 用法一致）
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithFetchOffset(db *sql.DB, page PageReq, query string, args ...any) (*PageResp[map[string]any], error) {
	return pageQueryForMap(db, page, query, getFetchOffsetQuery(page, query), args...)
}

func getFetchOffsetQuery(page PageReq, query string) string {
	return fmt.Sprintf(`SELECT t.* FROM (%s) AS t OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, query, (page.Page-1)*page.Size, page.Size)
}

// PageQueryForStructWithDeclareCursor 使用 Declare Cursor 分页查询，返回封装的结构体列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithDeclareCursor[T any](db *sql.DB, page PageReq, query string, args ...any) (*PageResp[T], error) {
	PageMap, err := PageQueryForMapWithDeclareCursor(db, page, query, args...)
	if err != nil {
		return nil, err
	}
	list := &PageResp[T]{
		PageReq: PageMap.PageReq,
		Total:   PageMap.Total,
	}
	bytes, err := json.Marshal(PageMap.List)
	if err != nil {
		return nil, err
	}
	var result []T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	list.List = result
	return list, err
}

// PageQueryForMapWithDeclareCursor 使用 Declare Cursor 分页查询，返回封装的map集合列表。
//
// 参数:
//   - db: 数据库连接
//   - page: 分页参数
//   - query: 查询语句
//   - args: 查询参数
//
// 返回:
//   - *PageResp[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithDeclareCursor(db *sql.DB, page PageReq, query string, args ...any) (*PageResp[map[string]any], error) {
	var err error
	page.Validate()
	pageResp := &PageResp[map[string]any]{PageReq: page}
	pageResp.Total, err = QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}
	// 延迟提交或回滚
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("rollback transaction failed: %w", rbErr)
			}
			return
		}
		if cmErr := tx.Commit(); cmErr != nil {
			logger.Error("commit transaction failed: %w", cmErr)
		}
	}()

	// 声明游标
	cursorName := "page_query_cursor"
	declareSQL := fmt.Sprintf("DECLARE %s CURSOR FOR %s", cursorName, query)
	if _, err := tx.Exec(declareSQL, args...); err != nil {
		return nil, fmt.Errorf("declare cursor failed: %w", err)
	}

	// 移动游标
	moveSQL := fmt.Sprintf("MOVE FORWARD %d IN %s", (page.Page-1)*page.Size, cursorName)
	if _, err := tx.Exec(moveSQL); err != nil {
		return nil, fmt.Errorf("move cursor failed: %w", err)
	}

	// 获取数据
	fetchSQL := fmt.Sprintf("FETCH %d FROM %s", page.Size, cursorName)
	rows, err := tx.Query(fetchSQL)
	if err != nil {
		return nil, fmt.Errorf("fetch data failed: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("close rows failed: %w", err)
		}
	}(rows)

	// 数据映射
	list := make([]map[string]any, 0, page.Size)
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
		list = append(list, data)
	}
	pageResp.List = list

	// 关闭游标
	closeSQL := fmt.Sprintf("CLOSE %s", cursorName)
	if _, err := tx.Exec(closeSQL); err != nil {
		return nil, fmt.Errorf("close cursor failed: %w", err)
	}

	return pageResp, nil
}

func raw2StructByPage[T any](db *sql.DB, page PageReq, query string, args ...any) ([]T, error) {
	list, err := raw2MapByPage(db, page, query, args...)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	var result []T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func raw2MapByPage(db *sql.DB, page PageReq, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	list := make([]map[string]any, 0, page.Size)

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

	start := (page.Page - 1) * page.Size
	end := start + page.Size

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
		list = append(list, data)
	}
	return list, nil
}

func raw2Struct[T any](db *sql.DB, query string, args ...any) ([]T, error) {
	list, err := raw2Map(db, query, args...)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	var result []T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func raw2Map(db *sql.DB, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	list := make([]map[string]any, 0)

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
		list = append(list, data)
	}
	return list, nil
}
