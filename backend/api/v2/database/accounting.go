package database

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// TransactionRow is a single money transaction as reported by GetTransactions.
// It merges rows from the expenditure, payment and boat_fuel tables, TableID
// says which one a row came from.
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
	TypeName  string           `json:"type_name"`
}

// The tables that hold money transactions. They are exposed to clients so that
// a transaction can be addressed by table and row.
const (
	TableIDExpenditure = iota
	TableIDPayment
	TableIDBoatFuel
)

// TableIDMap contains the IDs of tables that are used for
// money transactions.
var TableIDMap = map[uint64]struct {
	TableID   uint64
	TableName string
}{
	TableIDExpenditure: {
		TableID:   TableIDExpenditure,
		TableName: "expenditure",
	},
	TableIDPayment: {
		TableID:   TableIDPayment,
		TableName: "payment",
	},
	TableIDBoatFuel: {
		TableID:   TableIDBoatFuel,
		TableName: "boat_fuel",
	},
}

// GetYears returns every year that has at least one money transaction.
func (d *Mysql) GetYears() ([]uint64, error) {
	var years []uint64
	err := d.orm.Raw(`
		SELECT year(timestamp) AS year FROM payment
		UNION
		SELECT year(timestamp) AS year FROM expenditure
		UNION
		SELECT year(timestamp) AS year FROM boat_fuel
		GROUP BY year
		ORDER BY year ASC`).Scan(&years).Error
	return years, err
}

// GetPaymentTotal returns the total payments for a given year. If year
// is 0 it will return the total of all payments.
func (d *Mysql) GetPaymentTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf), 0) as result
		FROM payment
		WHERE year(timestamp) = ? OR 0 = ?
	`, year, year)
}

// GetExpenseTotal returns the total expenses for a given year, fuel included.
// If year is 0 it returns the total over all years.
func (d *Mysql) GetExpenseTotal(year uint64) (decimal.Decimal, error) {
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

// GetExpenseNoRefundsTotal is GetExpenseTotal without the refunds paid back to
// owners. If year is 0 it returns the total over all years.
func (d *Mysql) GetExpenseNoRefundsTotal(year uint64) (decimal.Decimal, error) {
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

// GetHeatCostTotal returns the total cost of all rides in a given year. If
// year is 0 it returns the total over all years.
func (d *Mysql) GetHeatCostTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(cost_chf),0) as result
		FROM heat WHERE year(timestamp) = ? OR 0 = ?
	`, year, year)
}

// GetSessionPaymentTotal returns the total the users paid for sessions in a
// given year. If year is 0 it returns the total over all years.
func (d *Mysql) GetSessionPaymentTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf),0) as result
		FROM payment WHERE type_id = ? 
		AND (year(timestamp) = ? OR 0 = ?)
	`, ExpenseTypeSession, year, year)
}

// GetSessionRefundsTotal returns the total refunded for sessions in a given
// year. If year is 0 it returns the total over all years.
func (d *Mysql) GetSessionRefundsTotal(year uint64) (decimal.Decimal, error) {
	return d.getSingleDecimalResult(`
		SELECT coalesce(sum(amount_chf),0) as result
		FROM expenditure WHERE type_id = ? 
		AND (year(timestamp) = ? OR 0 = ?)
	`, ExpenseTypeSession, year, year)
}

// GetTransactions returns all money transactions of a given year, most recent
// first. If year is 0 it returns the transactions of all years.
func (d *Mysql) GetTransactions(year uint64) ([]TransactionRow, error) {
	var r []TransactionRow
	// Table IDs are
	// 0: expenditure
	// 1: payment
	// 2: boat_fuel
	err := d.orm.Raw(`
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

// DeleteTransaction removes a single transaction, addressed by the table it
// lives in and its row ID.
func (d *Mysql) DeleteTransaction(tableID uint64, rowID uint64) error {
	switch tableID {
	case TableIDExpenditure:
		e := Expense{
			ID: uint(rowID),
		}
		return d.orm.Delete(&e).Error
	case TableIDPayment:
		p := Income{
			ID: uint(rowID),
		}
		return d.orm.Delete(&p).Error
	case TableIDBoatFuel:
		b := BoatFuel{
			ID: uint(rowID),
		}
		return d.orm.Delete(&b).Error
	default:
		return fmt.Errorf("table ID %d unknown, cannot delete row %d", tableID, rowID)
	}
}

// AddIncome records an incoming payment.
func (d *Mysql) AddIncome(data Income) error {
	return d.orm.Create(&data).Error
}

// AddExpense records an expense.
func (d *Mysql) AddExpense(data Expense) error {
	return d.orm.Create(&data).Error
}

func (d *Mysql) getSingleDecimalResult(rawQuery string, values ...any) (decimal.Decimal, error) {
	p := struct {
		Result *decimal.Decimal
	}{
		Result: &decimal.Zero,
	}
	err := d.orm.Raw(rawQuery, values...).First(&p).Error
	return *p.Result, err
}
