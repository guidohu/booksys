package handlers

import (
	"net/http"
	"strconv"
	"time"

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

type GetAccountingTransactionsRequest struct {
	// Year 0 stand for 'any' year.
	Year uint64 `json:"year"`
}

var GetAccountingTransactionsValidationErrors = map[string]string{
	"Year": "Please provide a valid year or 0.",
}

type GetAccountingTransactionsResponse []database.TransactionRow

type DeleteTransactionRequest struct {
	TableID uint64 `json:"table_id" validate:"tableid"`
	RowID   uint64 `json:"row_id"`
}

var DeleteTransactionValidationErrors = map[string]string{
	"TableID": "Please provide a valid table_id.",
	"RowID":   "Please provide a valid row id.",
}

type AddTransactionRequest struct {
	Amount  *decimal.Decimal `json:"amount" validate:"numeric"`
	Comment string           `json:"comment"`
	Date    string           `json:"date" validate:"required"`
	TypeID  uint64           `json:"type_id" validate:"expensetype"`
	UserID  uint64           `json:"user_id" validate:"required,numeric"`
}

var AddTransactionValidationErrors = map[string]string{
	"Amount":  "Please provide an amount.",
	"Comment": "Please provide a valid comment.",
	"Date":    "Please provide a valid date.",
	"TypeID":  "Please provide a valid income type.",
	"UserID":  "Please provide a valid user id.",
}

func (h *Handler) GetAccountingYears(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	years, err := dbh.GetYears()
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
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	expenseTypes, err := dbh.GetExpenseTypes()
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
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	expenseTypes, err := dbh.GetExpenseTypes()
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
	WriteSuccessResponse("income types", resp, w)
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

	resp := &GetAccountingStatisticsResponse{}

	// Get total payments.
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	payments, err := dbh.GetPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get payments total", slog.String("error", err.Error()))
	}
	resp.IncomeTotal = &payments

	// Get total payments for given year.
	yearPayments, err := dbh.GetPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearPayments total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.IncomeTotalSelectedYear = &yearPayments

	// Get total expenses.
	expenses, err := dbh.GetExpenseTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total", slog.String("error", err.Error()))
	}
	resp.ExpenseTotal = &expenses

	// Get total payments for given year.
	yearExpenses, err := dbh.GetExpenseTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.ExpenseTotalSelectedYear = &yearExpenses

	// Get total expenses without refunds.
	expensesNoRefund, err := dbh.GetExpenseNoRefundsTotal(0)
	if err != nil {
		slog.Warn("Cannot get expense total without refunds", slog.String("error", err.Error()))
	}
	resp.ExpenseNoRefundsTotal = &expensesNoRefund

	// Get total payments for given year without refunds.
	yearExpensesNoRefund, err := dbh.GetExpenseNoRefundsTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearExpenses total without refunds for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.ExpenseNoRefundsTotalSelectedYear = &yearExpensesNoRefund

	// Get total heat costs.
	heatCosts, err := dbh.GetHeatCostTotal(0)
	if err != nil {
		slog.Warn("Cannot get heat cost total", slog.String("error", err.Error()))
	}
	resp.RideMinutesCashUsed = &heatCosts

	// Get total heat costs for given year.
	yearHeatCosts, err := dbh.GetHeatCostTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get yearHeatCosts total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.RideMinutesCashUsedSelectedYear = &yearHeatCosts

	// Get session payment total.
	sessionPayments, err := dbh.GetSessionPaymentTotal(0)
	if err != nil {
		slog.Warn("Cannot get session payment total", slog.String("error", err.Error()))
	}
	resp.IncomeSessionPayment = &sessionPayments

	// Get session payment total for given year.
	yearSessionPayments, err := dbh.GetSessionPaymentTotal(req.Year)
	if err != nil {
		slog.Warn("Cannot get session payment total for", slog.Uint64("year", req.Year), slog.String("error", err.Error()))
	}
	resp.IncomeSessionPaymentSelectedYear = &yearSessionPayments

	// Get session payment total.
	sessionRefunds, err := dbh.GetSessionRefundsTotal(0)
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

func (h *Handler) GetAccountingTransactions(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &GetAccountingTransactionsRequest{}
	err := ReadBodyAndValidate(r, req, GetAccountingTransactionsValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}
	slog.Info("GetAccountingTransactions for", slog.Uint64("year", req.Year))

	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	var resp GetAccountingTransactionsResponse
	resp, err = dbh.GetTransactions(req.Year)
	if err != nil {
		slog.Warn("Cannot get transactions", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	WriteSuccessResponse("accounting transactions", resp, w)
}

func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &DeleteTransactionRequest{}
	err := ReadBodyAndValidate(r, req, DeleteTransactionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	err = dbh.DeleteTransaction(req.TableID, req.RowID)
	if err != nil {
		slog.Warn("Cannot get transactions", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot delete transaction because of an error.", w)
		return
	}

	WriteSuccessResponse("transaction deleted", nil, w)
}

func (h *Handler) AddIncome(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &AddTransactionRequest{}
	err := ReadBodyAndValidate(r, req, AddTransactionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	incomeEntry := database.Income{
		UserID:        uint(req.UserID),
		Amount:        *req.Amount,
		ExpenseTypeID: uint(req.TypeID),
		Comment:       req.Comment,
	}

	// Check if user exists.
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	_, err = dbh.GetUserById(uint(req.UserID))
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
		slog.Warn("Cannot add income transaction", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot write income transaction to database because of an error.", w)
		return
	}

	// TODO write an email to the user for Session Payments.

	WriteSuccessResponse("income added", nil, w)
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	session := GetSessionFromContext(r)
	if AuthenticatedAsAdminOrFailure(session, w) != nil {
		return
	}
	req := &AddTransactionRequest{}
	err := ReadBodyAndValidate(r, req, AddTransactionValidationErrors)
	if err != nil {
		slog.Warn("Request payload is not valid", slog.String("error", err.Error()))
		WriteFailureResponse(err.Error(), w)
		return
	}

	expenseEntry := database.Income{
		UserID:        uint(req.UserID),
		Amount:        *req.Amount,
		ExpenseTypeID: uint(req.TypeID),
		Comment:       req.Comment,
	}

	// Check if user exists.
	dbh, done := h.Database.GetHandler()
	defer done()
	if dbh == nil {
		slog.Warn("No database connection is available.")
		WriteFailureResponse("Operation cannot be performed. Database connection is not established properly.", w)
		return
	}
	_, err = dbh.GetUserById(uint(req.UserID))
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

	err = dbh.AddIncome(expenseEntry)
	if err != nil {
		slog.Warn("Cannot add expense transaction", slog.String("error", err.Error()))
		WriteFailureResponse("Cannot write expense transaction to database because of an error.", w)
		return
	}

	WriteSuccessResponse("expense added", nil, w)
}
