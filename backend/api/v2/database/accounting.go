package database

import (
	"github.com/shopspring/decimal"
)

func (d *DBMysql) GetYears() ([]uint64, error) {
	var years []uint64
	err := d.orm.Raw(`
		SELECT year(timestamp) AS year FROM payment
		UNION
		SELECT year(timestamp) AS year FROM expenditure
		UNION
		SELECT year(timestamp) AS year FROM boat_fuel
		ORDER BY year ASC`).Scan(&years).Error
	return years, err
}

func (d *DBMysql) GetPaymentTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf), 0) as result
		FROM payment
		WHERE year(timestamp) = ? OR 0 = ?
	`, year, year)
}

func (d *DBMysql) GetExpenseTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT SUM(t.tot) as result
		FROM
		(
			SELECT coalesce(sum(e.amount_chf),0) as tot
			FROM expenditure e
			WHERE year(timestamp) = ? OR 0 = ?
			UNION ALL
			SELECT coalesce(sum(b.cost_chf), 0) as tot
			FROM boat_fuel b
			WHERE (year(b.timestamp) = ? OR 0 = ?)
		      AND b.contributes_to_balance = 1
		) as t
	`, year, year, year, year)
}

func (d *DBMysql) GetExpenseNoRefundsTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT SUM(t.tot) as result
		FROM
		(
			SELECT coalesce(sum(e.amount_chf),0) as tot
			FROM expenditure e
			WHERE (year(timestamp) = ? OR 0 = ?)
			  AND e.type_id <> ?
			UNION ALL
			SELECT coalesce(sum(b.cost_chf), 0) as tot
			FROM boat_fuel b
			WHERE (year(b.timestamp) = ? OR 0 = ?)
		      AND b.contributes_to_balance = 1
		) as t
	`, year, year, ExpenseTypeOwnersRefund, year, year)
}

func (d *DBMysql) GetHeatCostTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(cost_chf),0) as result
		FROM heat WHERE year(timestamp) = ? OR 0 = ?
	`, year, year)
}

func (d *DBMysql) GetSessionPaymentTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf),0) as result
		FROM payment WHERE type_id = ? 
		AND (year(timestamp) = ? OR 0 = ?)
	`, ExpenseTypeSession, year, year)
}

func (d *DBMysql) GetSessionRefundsTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf),0) as result
		FROM expenditure WHERE type_id = ? 
		AND (year(timestamp) = ? OR 0 = ?)
	`, ExpenseTypeSession, year, year)
}

func (d *DBMysql) getSingleDecimalResult(rawQuery string, values ...interface{}) (decimal.Decimal, error) {
	p := struct {
		Result *decimal.Decimal
	}{
		Result: &decimal.Zero,
	}
	err := d.orm.Debug().Raw(rawQuery, values...).First(&p).Error
	return *p.Result, err
}
