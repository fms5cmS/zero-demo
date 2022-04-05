package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	princessDataFieldNames          = builder.RawFieldNames(&PrincessData{})
	princessDataRows                = strings.Join(princessDataFieldNames, ",")
	princessDataRowsExpectAutoSet   = strings.Join(stringx.Remove(princessDataFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	princessDataRowsWithPlaceHolder = strings.Join(stringx.Remove(princessDataFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheLearnPrincessDataIdPrefix         = "cache:learn:princessData:id:"
	cacheLearnPrincessDataPrincessIdPrefix = "cache:learn:princessData:princessId:"
)

type (
	PrincessDataModel interface {
		Insert(data *PrincessData) (sql.Result, error)
		FindOne(id int64) (*PrincessData, error)
		FindOneByPrincessId(princessId int64) (*PrincessData, error)
		Update(data *PrincessData) error
		Delete(id int64) error
		// 注意，这里也要加上！
		TransactInsert(ctx context.Context, session sqlx.Session, data *PrincessData) (sql.Result, error)
		TransactCtxSelf(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
	}

	defaultPrincessDataModel struct {
		sqlc.CachedConn
		table string
	}

	PrincessData struct {
		Id         int64     `db:"id"`
		PrincessId int64     `db:"princess_id"`
		Data       string    `db:"data"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewPrincessDataModel(conn sqlx.SqlConn, c cache.CacheConf) PrincessDataModel {
	return &defaultPrincessDataModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`princess_data`",
	}
}

func (m *defaultPrincessDataModel) Insert(data *PrincessData) (sql.Result, error) {
	learnPrincessDataPrincessIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataPrincessIdPrefix, data.PrincessId)
	learnPrincessDataIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, data.Id)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, princessDataRowsExpectAutoSet)
		return conn.Exec(query, data.PrincessId, data.Data)
	}, learnPrincessDataIdKey, learnPrincessDataPrincessIdKey)
	return ret, err
}

func (m *defaultPrincessDataModel) FindOne(id int64) (*PrincessData, error) {
	learnPrincessDataIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, id)
	var resp PrincessData
	err := m.QueryRow(&resp, learnPrincessDataIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", princessDataRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPrincessDataModel) FindOneByPrincessId(princessId int64) (*PrincessData, error) {
	learnPrincessDataPrincessIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataPrincessIdPrefix, princessId)
	var resp PrincessData
	err := m.QueryRowIndex(&resp, learnPrincessDataPrincessIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `princess_id` = ? limit 1", princessDataRows, m.table)
		if err := conn.QueryRow(&resp, query, princessId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPrincessDataModel) Update(data *PrincessData) error {
	learnPrincessDataIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, data.Id)
	learnPrincessDataPrincessIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataPrincessIdPrefix, data.PrincessId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, princessDataRowsWithPlaceHolder)
		return conn.Exec(query, data.PrincessId, data.Data, data.Id)
	}, learnPrincessDataIdKey, learnPrincessDataPrincessIdKey)
	return err
}

func (m *defaultPrincessDataModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	learnPrincessDataIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, id)
	learnPrincessDataPrincessIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataPrincessIdPrefix, data.PrincessId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, learnPrincessDataIdKey, learnPrincessDataPrincessIdKey)
	return err
}

func (m *defaultPrincessDataModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, primary)
}

func (m *defaultPrincessDataModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", princessDataRows, m.table)
	return conn.QueryRow(v, query, primary)
}

func (m *defaultPrincessDataModel) TransactInsert(ctx context.Context, session sqlx.Session, data *PrincessData) (sql.Result, error) {
	learnPrincessDataPrincessIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataPrincessIdPrefix, data.PrincessId)
	learnPrincessDataIdKey := fmt.Sprintf("%s%v", cacheLearnPrincessDataIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, princessDataRowsExpectAutoSet)
		// 这里要改为使用传入的 session 执行！
		return session.ExecCtx(ctx, query, data.PrincessId, data.Data)
	}, learnPrincessDataIdKey, learnPrincessDataPrincessIdKey)
	return ret, err
}

// TransCtx 暴露给 logic 来开启事务
func (m *defaultPrincessDataModel) TransactCtxSelf(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}
