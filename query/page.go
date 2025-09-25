package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Cooooing/cutil/common/logger"
	"github.com/Cooooing/cutil/query/base"
)

// ---------------- PageReq ----------------

type PageReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p *PageReq) Validate() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	return nil
}
func (p *PageReq) GetPage() int {
	return p.Page
}

func (p *PageReq) GetSize() int {
	return p.Size
}

// ---------------- PageResp ----------------

type PageResp[T any] struct {
	Page  int  `json:"page"`
	Size  int  `json:"size"`
	Total int  `json:"total"`
	List  []*T `json:"list"`
}

func (p *PageResp[T]) SetList(data []*T) {
	p.List = data
}

func (p *PageResp[T]) SetTotal(total int) {
	p.Total = total
}

func (p *PageResp[T]) SetPageReq(pageReq base.PageReqInterface) {
	p.Page = pageReq.GetPage()
	p.Size = pageReq.GetSize()
}

func (p *PageResp[T]) GetList() []*T {
	return p.List
}

func (p *PageResp[T]) GetTotal() int {
	return p.Total
}
func (p *PageResp[T]) GetPage() int {
	return p.Page
}

func (p *PageResp[T]) GetSize() int {
	return p.Size
}

// ---------------- PageRespFactory ----------------

type PageReqFactory func() base.PageReqInterface
type PageRespFactory[T any] func() base.PageRespInterface[T]

var (
	mu                     sync.RWMutex
	defaultPageReqFactory  PageReqFactory
	defaultPageRespFactory any // 泛型无法直接存储，存为 any
)

// SetDefaultPageReqFactory 设置全局 PageReqInterface 工厂
func SetDefaultPageReqFactory(factory PageReqFactory) {
	mu.Lock()
	defer mu.Unlock()
	defaultPageReqFactory = factory
}

// SetDefaultPageRespFactory 设置全局 PageRespInterface 工厂
func SetDefaultPageRespFactory[T any](factory PageRespFactory[T]) {
	mu.Lock()
	defer mu.Unlock()
	defaultPageRespFactory = factory
}

func getDefaultPageReq() base.PageReqInterface {
	mu.RLock()
	defer mu.RUnlock()
	if defaultPageReqFactory != nil {
		return defaultPageReqFactory()
	}
	return &PageReq{}
}

func getDefaultPageResp[T any]() base.PageRespInterface[T] {
	mu.RLock()
	defer mu.RUnlock()
	if defaultPageRespFactory != nil {
		return defaultPageRespFactory.(PageRespFactory[T])()
	}
	return &PageResp[T]{}
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
//   - PageRespInterface[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStruct[T any](db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[T], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[T]()
	total, err := QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(total)
	pageResp.SetPageReq(page)
	list, err := base.Raw2StructByPage[T](db, page, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetList(list)
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
//   - PageRespInterface[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMap(db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[map[string]any]()
	total, err := QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(total)
	pageResp.SetPageReq(page)
	list, err := base.Raw2MapByPage(db, page, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetList(list)
	return pageResp, nil
}

func pageQueryForStruct[T any](db *sql.DB, page base.PageReqInterface, countQuery string, query string, args ...any) (base.PageRespInterface[T], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[T]()
	total, err := QueryCount(db, countQuery, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(total)
	pageResp.SetPageReq(page)
	list, err := base.Raws2Struct[T](db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetList(list)
	return pageResp, nil
}

func pageQueryForMap(db *sql.DB, page base.PageReqInterface, countQuery string, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[map[string]any]()
	total, err := QueryCount(db, countQuery, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(total)
	pageResp.SetPageReq(page)
	list, err := base.Raw2Map(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetList(list)
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
//   - PageRespInterface[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithLimitOffset[T any](db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[T], error) {
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
//   - PageRespInterface[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithLimitOffset(db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	return pageQueryForMap(db, page, query, getLimitOffsetQuery(page, query), args...)
}

func getLimitOffsetQuery(page base.PageReqInterface, query string) string {
	return fmt.Sprintf(`SELECT t.* FROM (%s) AS t LIMIT %d OFFSET %d`, query, page.GetSize(), (page.GetPage()-1)*page.GetSize())
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
//   - PageRespInterface[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithRowNumber[T any](db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[T], error) {
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
//   - PageRespInterface[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithRowNumber(db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	return pageQueryForMap(db, page, query, getRowNumberQuery(page, query), args...)
}

func getRowNumberQuery(page base.PageReqInterface, query string) string {
	return fmt.Sprintf(`SELECT * FROM (SELECT t.*, ROW_NUMBER() OVER () AS rn FROM (%s) AS t ) AS sub WHERE rn BETWEEN %d AND %d`, query, (page.GetPage()-1)*page.GetSize()+1, page.GetPage()*page.GetSize())
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
//   - PageRespInterface[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithFetchOffset[T any](db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[T], error) {
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
//   - PageRespInterface[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithFetchOffset(db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	return pageQueryForMap(db, page, query, getFetchOffsetQuery(page, query), args...)
}

func getFetchOffsetQuery(page base.PageReqInterface, query string) string {
	return fmt.Sprintf(`SELECT t.* FROM (%s) AS t OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, query, (page.GetPage()-1)*page.GetSize(), page.GetSize())
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
//   - PageRespInterface[T]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForStructWithDeclareCursor[T any](db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[T], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[T]()
	PageMap, err := PageQueryForMapWithDeclareCursor(db, page, query, args...)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(PageMap.GetList())
	if err != nil {
		return nil, err
	}
	var list []*T
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(PageMap.GetTotal())
	pageResp.SetPageReq(page)
	pageResp.SetList(list)
	return pageResp, err
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
//   - PageRespInterface[map[string]any]: 分页结果
//   - error: 校验失败的错误信息
func PageQueryForMapWithDeclareCursor(db *sql.DB, page base.PageReqInterface, query string, args ...any) (base.PageRespInterface[map[string]any], error) {
	var err error
	if page == nil {
		page = getDefaultPageReq()
	}
	if err = page.Validate(); err != nil {
		return nil, err
	}
	pageResp := getDefaultPageResp[map[string]any]()
	total, err := QueryCount(db, query, args...)
	if err != nil {
		return nil, err
	}
	pageResp.SetTotal(total)
	pageResp.SetPageReq(page)

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
	moveSQL := fmt.Sprintf("MOVE FORWARD %d IN %s", (page.GetPage()-1)*page.GetSize(), cursorName)
	if _, err := tx.Exec(moveSQL); err != nil {
		return nil, fmt.Errorf("move cursor failed: %w", err)
	}

	// 获取数据
	fetchSQL := fmt.Sprintf("FETCH %d FROM %s", page.GetSize(), cursorName)
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
	list := make([]*map[string]any, 0, page.GetSize())
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
	pageResp.SetList(list)

	// 关闭游标
	closeSQL := fmt.Sprintf("CLOSE %s", cursorName)
	if _, err := tx.Exec(closeSQL); err != nil {
		return nil, fmt.Errorf("close cursor failed: %w", err)
	}

	return pageResp, nil
}
