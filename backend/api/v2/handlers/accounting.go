package handlers

import (
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"

	"server/database"
)

type GetYearsResponse struct {
	Years []string `json:"years"`
}

type GetExpenseTypeResponse struct {
	Types []database.ExpenseType `json:"types"`
}

type GetAccountingStatisticsRequest struct {
	// Year 0 stand for 'any' year.
	Year uint64 `json:"year"`
}

var GetAccountingStatisticsValidationErrors = map[string]string{
	"Year": "Please provide a valid year or 0.",
}

type GetAccountingStatisticsResponse struct {
	RideMinutesCashCredit             *decimal.Decimal `json:"total_open"`
	RideMinutesCashUsed               *decimal.Decimal `json:"total_used"`
	RideMinutesCashUsedSelectedYear   *decimal.Decimal `json:"total_used_selected_year"`
	Balance                           *decimal.Decimal `json:"current_balance"`
	SessionProfit                     *decimal.Decimal `json:"current_session_profit"`
	SessionProfitSelectedYear         *decimal.Decimal `json:"current_session_profit_selected_year"`
	ExpenseTotal                      *decimal.Decimal `json:"total_expenditure"`
	ExpenseTotalSelectedYear          *decimal.Decimal `json:"total_expenditure_selected_year"`
	ExpenseNoRefundsTotal             *decimal.Decimal `json:"total_expenditure_no_refunds"`
	ExpenseNoRefundsTotalSelectedYear *decimal.Decimal `json:"total_expenditure_selected_year_no_refunds"`
	IncomeTotal                       *decimal.Decimal `json:"total_payment"`
	IncomeTotalSelectedYear           *decimal.Decimal `json:"total_payment_selected_year"`
	IncomeSessionPayment              *decimal.Decimal `json:"total_session_payment"`
	IncomeSessionPaymentSelectedYear  *decimal.Decimal `json:"total_session_payment_selected_year"`
	RefundsSessionTotal               *decimal.Decimal `json:"total_session_refunds"`
}

func (h *Handler) GetAccountingYears(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	years, err := h.GetDB().GetYears()
	if err != nil {
		slog.Warn("Cannot get accounting years.", nil)
		WriteFailureResponse(err.Error(), w)
		return
	}
	resp := GetYearsResponse{}
	for _, y := range years {
		resp.Years = append(resp.Years, strconv.Itoa(int(y)))
	}
	resp.Years = append(resp.Years, "any")
	WriteSuccessResponse("accounting years", resp, w)
}

func (h *Handler) GetExpenseTypes(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	expenseTypes, err := h.GetDB().GetExpenseTypes()
	if err != nil {
		slog.Error("Cannot get expense types", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get expense types.", w)
		return
	}
	resp := GetExpenseTypeResponse{
		Types: expenseTypes,
	}
	WriteSuccessResponse("expense types", resp, w)
}

func (h *Handler) GetIncomeTypes(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	expenseTypes, err := h.GetDB().GetExpenseTypes()
	if err != nil {
		slog.Error("Cannot get expense types", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot get expense types.", w)
		return
	}

	incomeTypes := []database.ExpenseType{}
	for _, e := range expenseTypes {
		if e.ID == database.ExpenseTypeFuelDirect {
			continue
		}
		if e.ID == database.ExpenseTypeMaintenance {
			continue
		}
		if e.ID == database.ExpenseTypeMaterial {
			continue
		}
		incomeTypes = append(incomeTypes, e)
	}

	resp := GetExpenseTypeResponse{
		Types: incomeTypes,
	}
	WriteSuccessResponse("expense types", resp, w)
}

func (h *Handler) GetAccountingStatistics(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &GetAccountingStatisticsRequest{}
	err := ReadBodyAndValidate(r, req, GetAccountingStatisticsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}
	slog.Info("GetAccountingStatistics for", slog.Uint64("year", req.Year))

	resp := &GetAccountingStatisticsResponse{
		RideMinutesCashCredit:             &decimal.Decimal{},
		RideMinutesCashUsed:               &decimal.Decimal{},
		RideMinutesCashUsedSelectedYear:   &decimal.Decimal{},
		Balance:                           &decimal.Decimal{},
		SessionProfit:                     &decimal.Decimal{},
		SessionProfitSelectedYear:         &decimal.Decimal{},
		ExpenseTotal:                      &decimal.Decimal{},
		ExpenseTotalSelectedYear:          &decimal.Decimal{},
		ExpenseNoRefundsTotal:             &decimal.Decimal{},
		ExpenseNoRefundsTotalSelectedYear: &decimal.Decimal{},
		IncomeTotal:                       &decimal.Decimal{},
		IncomeTotalSelectedYear:           &decimal.Decimal{},
		IncomeSessionPayment:              &decimal.Decimal{},
		IncomeSessionPaymentSelectedYear:  &decimal.Decimal{},
		RefundsSessionTotal:               &decimal.Decimal{},
	}

	// Get total payments.
	payments, err := h.GetDB().GetPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get payments total", slog.String("error", err.Error()))
	}
	resp.IncomeTotal = &payments

	// Get total payments for given year.
	yearPayments, err := h.GetDB().GetPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearPayments total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.IncomeTotalSelectedYear = &yearPayments

	// Get total expenses.
	expenses, err := h.GetDB().GetExpenseTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total", slog.String("error", err.Error()))
	}
	resp.ExpenseTotal = &expenses

	// Get total payments for given year.
	yearExpenses, err := h.GetDB().GetExpenseTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.ExpenseTotalSelectedYear = &yearExpenses

	// Get total expenses without refunds.
	expensesNoRefund, err := h.GetDB().GetExpenseNoRefundsTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total without refunds", slog.String("error", err.Error()))
	}
	resp.ExpenseNoRefundsTotal = &expensesNoRefund

	// Get total payments for given year without refunds.
	yearExpensesNoRefund, err := h.GetDB().GetExpenseNoRefundsTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total without refunds for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.ExpenseNoRefundsTotalSelectedYear = &yearExpensesNoRefund

	// Get total heat costs.
	heatCosts, err := h.GetDB().GetHeatCostTotal(0)
	if err != nil {
		slog.Warn("Cannot get heat cost total", slog.String("error", err.Error()))
	}
	resp.RideMinutesCashUsed = &heatCosts

	// Get total heat costs for given year.
	yearHeatCosts, err := h.GetDB().GetHeatCostTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearHeatCosts total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.RideMinutesCashUsedSelectedYear = &yearHeatCosts

	// Get session payment total.
	sessionPayments, err := h.GetDB().GetSessionPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get session payment total", slog.String("error", err.Error()))
	}
	resp.IncomeSessionPayment = &sessionPayments

	// Get session payment total for given year.
	yearSessionPayments, err := h.GetDB().GetSessionPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get session payment total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.IncomeSessionPaymentSelectedYear = &yearSessionPayments

	// Get session payment total.
	sessionRefunds, err := h.GetDB().GetSessionRefundsTotal(0)
	if err != nil {
		slog.Warn("Cannot get session refunds total", slog.String("error", err.Error()))
	}
	resp.RefundsSessionTotal = &sessionRefunds

	// Get current total open cost (payments for sessions - heat costs - refunds).
	currentRidesCredit := resp.IncomeSessionPayment.Sub(*resp.RideMinutesCashUsed).Sub(*resp.RefundsSessionTotal)
	resp.RideMinutesCashCredit = &currentRidesCredit

	// Get current balance.
	currentBalance := resp.IncomeTotal.Sub(*resp.ExpenseTotal)
	resp.Balance = &currentBalance

	// Get current profit and profit for selected year.
	currentProfit := resp.RideMinutesCashUsed.Sub(*resp.ExpenseNoRefundsTotal)
	resp.SessionProfit = &currentProfit
	yearCurrentProfit := resp.RideMinutesCashUsedSelectedYear.Sub(*resp.ExpenseNoRefundsTotalSelectedYear)
	resp.SessionProfitSelectedYear = &yearCurrentProfit

	WriteSuccessResponse("accounting statistics", resp, w)
}
