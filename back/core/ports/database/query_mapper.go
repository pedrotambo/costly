package database

import (
	"context"
)

type RowMapper[T any] func(rowScanner RowScanner) (T, error)

func QueryRowAndMap[T any](ctx context.Context, db RowQuerier, rowMapper RowMapper[T], query string, args ...any) (T, error) {
	row := db.QueryRowContext(ctx, query, args...)
	return rowMapper(row)
}

func QueryAndMap[T any](ctx context.Context, db RowsQuerier, rowMapper RowMapper[T], query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	ts := []T{}
	for rows.Next() {
		t, err := rowMapper(rows)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}
