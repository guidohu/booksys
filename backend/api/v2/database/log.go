package database

import (
	"fmt"
	"regexp"

	"golang.org/x/exp/slog"
)

type Log struct {
	ID         uint   `json:"id"`
	LogMessage string `json:"log"`
	Time       string `json:"time"`
	Type       string `json:"type"`
}

func (d *DBMysql) GetLogs() ([]Log, error) {
	currency, err := d.GetPropertyValue("currency")
	if err != nil {
		slog.Warn("Cannot get currency from configuration table", slog.String("error", err.Error()))
		return []Log{}, err
	}
	// make sure the currency is a valid value
	// TODO make this a validator and only allow the same values
	// to be written to configuration.
	r, _ := regexp.Compile(`[A-Za-z]{2,5}|\$`)
	if !r.MatchString(currency.Value) {
		slog.Warn("Currency is not a safe and valid string to build the SQL statement for retrieving logs")
		return []Log{}, fmt.Errorf("currency is no valid value log text cannot be built")
	}
	query := fmt.Sprintf(`SELECT c_log.id as id, c_log.type as type, c_log.time as time, c_log.log_message as log_message FROM (
		SELECT h.id as id, "heat" as type, h.timestamp as time, 
			   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
					  ' was riding for ', FORMAT(h.duration_s/60,1), 
					  ' min and ', FORMAT(h.cost_chf, 2), ' %s') as log_message 
			FROM user u, heat h 
			WHERE u.id = h.user_id AND h.comment IS NULL
		UNION ALL
		SELECT h.id as id, "heat" as type, h.timestamp as time, 
			   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
					  ' was riding for ', FORMAT(h.duration_s/60,1), 
					  ' min and ', FORMAT(h.cost_chf, 2), ' %s (Comment: ', h.comment, ')') as log_message 
			FROM user u, heat h 
			WHERE u.id = h.user_id AND NOT h.comment IS NULL
		UNION ALL
		SELECT p.id as id, "payment" as type, p.timestamp as time, 
			   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, 
					  ' paid for sessions ', FORMAT(p.amount_chf, 2), ' %s') as log_message 
			FROM user u, payment p 
			WHERE u.id = p.user_id AND p.type_id = 4
		UNION ALL
		SELECT bf.id as id, "boat_fuel" as type, bf.timestamp as time, 
			   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, ' added ', 
					  FORMAT(bf.liters, 2), 'L fuel for ', FORMAT(cost_chf, 2), ' %s') as log_message 
			FROM user u, boat_fuel bf 
			WHERE bf.user_id = u.id AND bf.contributes_to_balance = 1
		UNION ALL
		SELECT bf.id as id, "boat_fuel" as type, bf.timestamp as time, 
			   CONCAT(u.first_name COLLATE utf8_general_ci, ' ', u.last_name, ' added ', 
					  FORMAT(bf.liters, 2), 'L fuel for ', FORMAT(cost_chf, 2), ' %s (on credit)') as log_message 
			FROM user u, boat_fuel bf 
			WHERE bf.user_id = u.id AND bf.contributes_to_balance = 0
	  ) as c_log ORDER BY time DESC`, currency.Value, currency.Value, currency.Value, currency.Value, currency.Value)
	var logs []Log
	err = d.orm.Raw(query).Find(&logs).Error
	return logs, err
}
