package main

import "context"

// example

const queryGetExchange = `select id, item from items where id = $1`

type ExampleItem struct {
	ID   int    `db:"id"`
	Item string `db:"item"`
}

func GetExampleItem(ctx context.Context, db Querier, id int) (ExampleItem, error) {
	var item ExampleItem
	if err := db.QueryRowContext(ctx, queryGetExchange, id).Scan(&item.ID, &item.Item); err != nil {
		return ExampleItem{}, err
	}
	return item, nil
}
