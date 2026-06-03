package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Sale struct {
	Product int
	Volume  int
	Date    string
}

// String реализует метод интерфейса fmt.Stringer для Sale, возвращает строковое представление объекта Sale.
// Теперь, если передать объект Sale в fmt.Println(), то выведется строка, которую вернёт эта функция.
func (s Sale) String() string {
	return fmt.Sprintf("Product: %d Volume: %d Date:%s", s.Product, s.Volume, s.Date)
}

// контекст пробрасывается из вызывающей стороны, чтобы запрос можно было отменить, если клиент ушёл
func selectSales(ctx context.Context, client int) ([]Sale, error) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	var sales []Sale

	rows, err := db.QueryContext(ctx, "SELECT product, volume, date FROM sales WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("select sales: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		sale := Sale{}

		err := rows.Scan(&sale.Product, &sale.Volume, &sale.Date)
		if err != nil {
			return nil, fmt.Errorf("scan sale: %w", err)
		}

		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return sales, nil
}

func main() {
	client := 208
	ctx := context.Background()

	sales, err := selectSales(ctx, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, sale := range sales {
		fmt.Println(sale)
	}
}
