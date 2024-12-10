package main

import (
	"context"
)

type Example struct {
	id   int    `db:"id"`
	item string `db:"item"`
}

const (
	queryGetExample = `select id, item
	from example
	where id=$1
`
)

func GetExample(ctx context.Context, db Querier, id int) (Example, error) {
	var item Example
	if err := db.QueryRowContext(ctx, queryGetExample, id).Scan(
		&item.id,
		&item.item,
	); err != nil {
		return Example{}, err
	}
	return item, nil
}
