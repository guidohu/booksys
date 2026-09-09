package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"server/database"
)

// GetYearsResponse is the payload returned by GetAccountingYears.
type GetYearsResponse struct {
	Years []string `json:"years"`
}

// GetExpenseTypeResponse is the payload returned by GetExpenseTypes and
// GetIncomeTypes.
type GetExpenseTypeResponse struct {
	Types []database.ExpenseType `json:"types"`
}

// GetAccountingStatisticsRequest is the request body of GetAccountingStatistics, which serves
// /api/v2/accounting/statistics/get.
type GetAccountingStatisticsRequest struct {
	// Year 0 stand for 'any' year.
	Year uint64 `json:"year"`
}

// GetAccountingStatisticsValidationErrors maps the field names of GetAccountingStatisticsRequest to the message the API
// returns when that field fails validation.
var GetAccountingStatisticsValidationErrors = map[string]string{
	"Year": "Please provide a valid year or 0.",
}

// GetAccountingStatisticsResponse is the payload returned by GetAccountingStatistics, which serves
// /api/v2/accounting/statistics/get.
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

// GetAccountingTransactionsRequest is the request body of GetAccountingTransactions, which serves
// /api/v2/accounting/transactions/get.
type GetAccountingTransactionsRequest struct {
	// Year 0 stand for 'any' year.
	Year uint64 `json:"year"`
}

// GetAccountingTransactionsValidationErrors maps the field names of GetAccountingTransactionsRequest to the message the API
// returns when that field fails validation.
var GetAccountingTransactionsValidationErrors = map[string]string{
	"Year": "Please provide a valid year or 0.",
}

// GetAccountingTransactionsResponse is the payload returned by GetAccountingTransactions, which serves
// /api/v2/accounting/transactions/get.
type GetAccountingTransactionsResponse []database.TransactionRow

// DeleteTransactionRequest is the request body of DeleteTransaction, which serves
// /api/v2/accounting/transactions/delete.
type DeleteTransactionRequest struct {
	TableID uint64 `json:"table_id" validate:"tableid"`
	RowID   uint64 `json:"row_id"`
}

// DeleteTransactionValidationErrors maps the field names of DeleteTransactionRequest to the message the API
// returns when that field fails validation.
var DeleteTransactionValidationErrors = map[string]string{
	"TableID": "Please provide a valid table_id.",
	"RowID":   "Please provide a valid row id.",
}

// AddTransactionRequest is the request body of AddIncome and AddExpense.
type AddTransactionRequest struct {
	Amount  *decimal.Decimal `json:"amount" validate:"required"`
	Comment string           `json:"comment"`
	Date    string           `json:"date" validate:"required"`
	TypeID  uint64           `json:"type_id" validate:"expensetype"`
	UserID  uint64           `json:"user_id" validate:"required,numeric"`
}

// AddTransactionValidationErrors maps the field names of AddTransactionRequest to the message the API
// returns when that field fails validation.
var AddTransactionValidationErrors = map[string]string{
	"Amount":  "Please provide an amount.",
	"Comment": "Please provide a valid comment.",
	"Date":    "Please provide a valid date.",
	"TypeID":  "Please provide a valid income type.",
	"UserID":  "Please provide a valid user id.",
}

// GetAccountingYears returns every year that has transactions.
//
// It serves /api/v2/accounting/years/list and is open to administrators.
func (h *Handler) GetAccountingYears(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	years, err := dbh.GetYears()
	if err != nil {
		slog.Warn("Cannot get accounting years", slog.Any("error", err))
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

// GetExpenseTypes returns the expense categories.
//
// It serves /api/v2/accounting/expense_types/list and is open to administrators.
func (h *Handler) GetExpenseTypes(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	expenseTypes, err := dbh.GetExpenseTypes()
	if err != nil {
		slog.Error("Cannot get expense types", slog.Any("error", err))
		WriteFailureResponse("Cannot get expense types.", w)
		return
	}
	resp := GetExpenseTypeResponse{
		Types: expenseTypes,
	}
	WriteSuccessResponse("expense types", resp, w)
}

// GetIncomeTypes returns the income categories.
//
// It serves /api/v2/accounting/income_types/list and is open to administrators.
func (h *Handler) GetIncomeTypes(w http.ResponseWriter, r *http.Request) {
	hCtx := GetHandlerContext(r)
	dbh := hCtx.Database
	expenseTypes, err := dbh.GetExpenseTypes()
	if err != nil {
		slog.Error("Cannot get expense types", slog.Any("error", err))
		WriteFailureResponse("Cannot get expense types.", w)
		return
	}

	incomeTypes := make([]database.ExpenseType, 0, len(expenseTypes))
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
	WriteSuccessResponse("income types", resp, w)
}

// GetAccountingStatistics returns the income, expense and balance totals for a year.
//
// It serves /api/v2/accounting/statistics/get and is open to administrators.
func (h *Handler) GetAccountingStatistics(w http.ResponseWriter, r *http.Request, req GetAccountingStatisticsRequest, hCtx *HandlerCtx) {
	resp := &GetAccountingStatisticsResponse{}

	// Get total payments.
	dbh := hCtx.Database
	payments, err := dbh.GetPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get payments total", slog.Any("error", err))
	}
	resp.IncomeTotal = &payments

	// Get total payments for given year.
	yearPayments, err := dbh.GetPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearPayments total for", slog.Uint64("year", req.Year), slog.Any("error", err))
	}
	resp.IncomeTotalSelectedYear = &yearPayments

	// Get total expenses.
	expenses, err := dbh.GetExpenseTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total", slog.Any("error", err))
	}
	resp.ExpenseTotal = &expenses

	// Get total payments for given year.
	yearExpenses, err := dbh.GetExpenseTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total for", slog.Uint64("year", req.Year), slog.Any("error", err))
	}
	resp.ExpenseTotalSelectedYear = &yearExpenses

	// Get total expenses without refunds.
	expensesNoRefund, err := dbh.GetExpenseNoRefundsTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total without refunds", slog.Any("error", err))
	}
	resp.ExpenseNoRefundsTotal = &expensesNoRefund

	// Get total payments for given year without refunds.
	yearExpensesNoRefund, err := dbh.GetExpenseNoRefundsTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total without refunds for", slog.Uint64("year", req.Year), slog.Any("error", err))
	}
	resp.ExpenseNoRefundsTotalSelectedYear = &yearExpensesNoRefund

	// Get total heat costs.
	heatCosts, err := dbh.GetHeatCostTotal(0)
	if err != nil {
		slog.Warn("Cannot get heat cost total", slog.Any("error", err))
	}
	resp.RideMinutesCashUsed = &heatCosts

	// Get total heat costs for given year.
	yearHeatCosts, err := dbh.GetHeatCostTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearHeatCosts total for", slog.Uint64("year", req.Year), slog.Any("error", err))
	}
	resp.RideMinutesCashUsedSelectedYear = &yearHeatCosts

	// Get session payment total.
	sessionPayments, err := dbh.GetSessionPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get session payment total", slog.Any("error", err))
	}
	resp.IncomeSessionPayment = &sessionPayments

	// Get session payment total for given year.
	yearSessionPayments, err := dbh.GetSessionPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get session payment total for", slog.Uint64("year", req.Year), slog.Any("error", err))
	}
	resp.IncomeSessionPaymentSelectedYear = &yearSessionPayments

	// Get session payment total.
	sessionRefunds, err := dbh.GetSessionRefundsTotal(0)
	if err != nil {
		slog.Warn("Cannot get session refunds total", slog.Any("error", err))
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

// GetAccountingTransactions returns the transactions of a year.
//
// It serves /api/v2/accounting/transactions/get and is open to administrators.
func (h *Handler) GetAccountingTransactions(w http.ResponseWriter, r *http.Request, req GetAccountingStatisticsRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	var resp GetAccountingTransactionsResponse
	resp, err := dbh.GetTransactions(req.Year)
	if err != nil {
		slog.Warn("Cannot get transactions", slog.Any("error", err))
		WriteFailureResponse(err.Error(), w)
		return
	}

	WriteSuccessResponse("accounting transactions", resp, w)
}

// DeleteTransaction removes a single transaction.
//
// It serves /api/v2/accounting/transactions/delete and is open to administrators.
func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request, req DeleteTransactionRequest, hCtx *HandlerCtx) {
	dbh := hCtx.Database
	err := dbh.DeleteTransaction(req.TableID, req.RowID)
	if err != nil {
		slog.Warn("Cannot get transactions", slog.Any("error", err))
		WriteFailureResponse("Cannot delete transaction because of an error.", w)
		return
	}

	WriteSuccessResponse("transaction deleted", nil, w)
}

// AddIncome records an incoming payment.
//
// It serves /api/v2/accounting/income/add and is open to administrators.
func (h *Handler) AddIncome(w http.ResponseWriter, r *http.Request, req AddTransactionRequest, hCtx *HandlerCtx) {
	incomeEntry := database.Income{
		UserID:        uint(req.UserID),
		Amount:        *req.Amount,
		ExpenseTypeID: uint(req.TypeID),
		Comment:       req.Comment,
	}

	// Check if user exists.
	dbh := hCtx.Database
	_, err := dbh.GetUserByID(uint(req.UserID))
	if err != nil {
		slog.Warn("Cannot find user with", slog.Uint64("user_id", req.UserID))
		WriteFailureResponse("Cannot find the selected user", w)
		return
	}
	incomeEntry.UserID = uint(req.UserID)

	// Parse time.
	date, err := time.Parse("2006-01-02T15:04", req.Date)
	if err != nil {
		slog.Warn("Cannot parse date information from", slog.String("date", req.Date))
		WriteFailureResponse("Cannot parse date, has to be of the form 2006-01-02T15:04", w)
		return
	}
	incomeEntry.Timestamp = date

	// Some require a comment and for some types we set
	// a default one.
	switch req.TypeID {
	case database.ExpenseTypeSession:
		if req.Comment == "" {
			incomeEntry.Comment = "Session Payment"
		}
	case database.ExpenseTypeMembershipFee:
		if req.Comment == "" {
			incomeEntry.Comment = "Membership Fee"
		}
	default:
		if incomeEntry.Comment == "" {
			slog.Warn("No comment specified for new income entry.")
			WriteFailureResponse("No comment specified for new income entry.", w)
			return
		}
	}

	err = dbh.AddIncome(incomeEntry)
	if err != nil {
		slog.Warn("Cannot add income transaction", slog.Any("error", err))
		WriteFailureResponse("Cannot write income transaction to database because of an error.", w)
		return
	}

	// TODO write an email to the user for Session Payments.

	WriteSuccessResponse("income added", nil, w)
}

// AddExpense records an expense.
//
// It serves /api/v2/accounting/expense/add and is open to administrators.
func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request, req AddTransactionRequest, hCtx *HandlerCtx) {
	expenseEntry := database.Expense{
		UserID:        uint(req.UserID),
		Amount:        *req.Amount,
		ExpenseTypeID: uint(req.TypeID),
		Comment:       req.Comment,
	}

	// Check if user exists.
	dbh := hCtx.Database
	_, err := dbh.GetUserByID(uint(req.UserID))
	if err != nil {
		slog.Warn("Cannot find user with", slog.Uint64("user_id", req.UserID))
		WriteFailureResponse("Cannot find the selected user", w)
		return
	}
	expenseEntry.UserID = uint(req.UserID)

	// Parse time.
	date, err := time.Parse("2006-01-02T15:04", req.Date)
	if err != nil {
		slog.Warn("Cannot parse date information from", slog.String("date", req.Date))
		WriteFailureResponse("Cannot parse date, has to be of the form 2006-01-02T15:04", w)
		return
	}
	expenseEntry.Timestamp = date

	// Some require a comment and for some types we set
	// a default one.
	switch req.TypeID {
	case database.ExpenseTypeFuelDirect:
		slog.Warn("Attempt to use accounting API to add fuel expenses.")
		WriteFailureResponse("Please use boat API for fuel expenses.", w)
		return
	default:
		if expenseEntry.Comment == "" {
			slog.Warn("No comment specified for new expense entry.")
			WriteFailureResponse("No comment specified for new expense entry.", w)
			return
		}
	}

	err = dbh.AddExpense(expenseEntry)
	if err != nil {
		slog.Warn("Cannot add expense transaction", slog.Any("error", err))
		WriteFailureResponse("Cannot write expense transaction to database because of an error.", w)
		return
	}

	WriteSuccessResponse("expense added", nil, w)
}
