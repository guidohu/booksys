package database

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRow struct {
	ID        uint64           `json:"id"`
	Amount    *decimal.Decimal `json:"amount"`
	Comment   string           `json:"comment"`
	UserID    uint64           `json:"user_id"`
	FirstName string           `json:"fn"`
	LastName  string           `json:"ln"`
	TableID   uint64           `json:"tbl"`
	Timestamp time.Time        `json:"timestamp"`
	TypeID    uint64           `json:"type_id"`
	TypeName  string           `json:"session"`
}

const (
	TableIdExpenditure = iota
	TableIdPayment
	TableIdBoatFuel
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

func (d *DBMysql) GetTransactions(year uint64) ([]TransactionRow, error) {
	r := []TransactionRow{}
	err := d.orm.Debug().Raw(`
		SELECT acc.tbl as table_id, acc.id as id, acc.user_id as user_id, u.first_name as first_name, u.last_name as last_name, 
			acc.type_id as type_id, et.name as type_name, acc.timestamp as timestamp, 
			acc.amount as amount, acc.comment as comment
		FROM 
		(
			SELECT 0 as tbl, e.id, e.user_id, e.type_id, 
				e.timestamp, e.amount_chf * -1 as amount, e.comment 
			FROM expenditure e
			UNION ALL
			SELECT 1 as tbl, p.id, p.user_id, p.type_id as type_id, 
				p.timestamp, p.amount_chf as amount, p.comment as comment 
			FROM payment p
			UNION ALL
			SELECT 2 as tbl, b.id, b.user_id, 0 as type_id,
				b.timestamp, b.cost_chf * -1 as amount, 
				CONCAT(b.liters, 'L fuel') as comment
			FROM boat_fuel b WHERE contributes_to_balance = 1
		) as acc 
		LEFT JOIN user u ON u.id = acc.user_id
		LEFT JOIN expenditure_type et ON et.id = acc.type_id 
		WHERE year(acc.timestamp) = ? OR 0 = ?
		ORDER BY acc.timestamp DESC
	`, year, year).Scan(&r).Error
	return r, err
}

func (d *DBMysql) DeleteTransaction(tableID uint64, rowID uint64) error {
	switch tableID {
	case TableIdExpenditure:
		e := Expense{
			ID: uint(rowID),
		}
		return d.orm.Delete(&e).Error
	case TableIdPayment:
		p := Income{
			ID: uint(rowID),
		}
		return d.orm.Delete(&p).Error
	case TableIdBoatFuel:
		b := BoatFuel{
			ID: uint(rowID),
		}
		return d.orm.Delete(&b).Error
	default:
		return fmt.Errorf("table ID %d unknown, cannot delete row %d", tableID, rowID)
	}
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
