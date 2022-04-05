package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	princessFieldNames          = builder.RawFieldNames(&Princess{})
	princessRows                = strings.Join(princessFieldNames, ",")
	princessRowsExpectAutoSet   = strings.Join(stringx.Remove(princessFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	princessRowsWithPlaceHolder = strings.Join(stringx.Remove(princessFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	PrincessModel interface {
		Insert(data *Princess) (sql.Result, error)
		FindOne(id int64) (*Princess, error)
		FindOneByName(name string) (*Princess, error)
		Update(data *Princess) error
		Delete(id int64) error
		// 注意，这里也要加上！
		TransactInsert(ctx context.Context, session sqlx.Session, data *Princess) (sql.Result, error)
		TransactCtxSelf(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
	}

	defaultPrincessModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Princess struct {
		Id              int64   `db:"id"`
		Name            string  `db:"name"`
		Cv              string  `db:"cv"`
		Hp              float64 `db:"hp"`
		Attack          float64 `db:"attack"`
		PhysicalDefense float64 `db:"physical_defense"`
		MagicDefense    float64 `db:"magic_defense"`
		AttackType      int64   `db:"attack_type"`
		Nickname        string  `db:"nickname"`
	}
)

func NewPrincessModel(conn sqlx.SqlConn) PrincessModel {
	return &defaultPrincessModel{
		conn:  conn,
		table: "`princess`",
	}
}

func (m *defaultPrincessModel) Insert(data *Princess) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, princessRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Name, data.Cv, data.Hp, data.Attack, data.PhysicalDefense, data.MagicDefense, data.AttackType, data.Nickname)
	return ret, err
}

func (m *defaultPrincessModel) FindOne(id int64) (*Princess, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", princessRows, m.table)
	var resp Princess
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPrincessModel) FindOneByName(name string) (*Princess, error) {
	var resp Princess
	query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", princessRows, m.table)
	err := m.conn.QueryRow(&resp, query, name)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultPrincessModel) Update(data *Princess) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, princessRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Name, data.Cv, data.Hp, data.Attack, data.PhysicalDefense, data.MagicDefense, data.AttackType, data.Nickname, data.Id)
	return err
}

func (m *defaultPrincessModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultPrincessModel) TransactInsert(ctx context.Context, session sqlx.Session, data *Princess) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, princessRowsExpectAutoSet)
	// 这里要改为使用传入的 session 执行！
	ret, err := session.ExecCtx(ctx, query, data.Name, data.Cv, data.Hp, data.Attack, data.PhysicalDefense, data.MagicDefense, data.AttackType, data.Nickname)
	return ret, err
}

// TransCtx 暴露给 logic 来开启事务
func (m *defaultPrincessModel) TransactCtxSelf(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		return fn(ctx, s)
	})
}
