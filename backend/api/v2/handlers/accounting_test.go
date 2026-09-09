package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"server/database"
	"server/database/dbtest"
)

func TestGetAccountingYears(t *testing.T) {
	t.Run("returns the years plus the 'any' entry", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetYearsFn: func() ([]uint64, error) { return []uint64{2023, 2024}, nil },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetAccountingYears(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetYearsResponse
		decodeData(t, resp, &data)
		want := []string{"2023", "2024", "any"}
		if len(data.Years) != len(want) {
			t.Fatalf("years = %v, want %v", data.Years, want)
		}
		for i := range want {
			if data.Years[i] != want[i] {
				t.Errorf("years[%d] = %q, want %q", i, data.Years[i], want[i])
			}
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{GetYearsFn: func() ([]uint64, error) { return nil, errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetAccountingYears(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetExpenseTypes(t *testing.T) {
	t.Run("returns all expense types", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetExpenseTypesFn: func() ([]database.ExpenseType, error) { return database.DefaultExpenseTypes, nil },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetExpenseTypes(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetExpenseTypeResponse
		decodeData(t, resp, &data)
		if len(data.Types) != len(database.DefaultExpenseTypes) {
			t.Errorf("types = %d entries, want %d", len(data.Types), len(database.DefaultExpenseTypes))
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{GetExpenseTypesFn: func() ([]database.ExpenseType, error) { return nil, errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetExpenseTypes(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetIncomeTypes(t *testing.T) {
	t.Run("filters the expense only types", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetExpenseTypesFn: func() ([]database.ExpenseType, error) { return database.DefaultExpenseTypes, nil },
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		h.GetIncomeTypes(rec, newRequest("", &HandlerCtx{Database: db}))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		var data GetExpenseTypeResponse
		decodeData(t, resp, &data)

		excluded := map[uint]bool{
			database.ExpenseTypeFuelDirect:  true,
			database.ExpenseTypeMaintenance: true,
			database.ExpenseTypeMaterial:    true,
		}
		for _, tp := range data.Types {
			if excluded[tp.ID] {
				t.Errorf("type %d (%s) must not be offered as an income type", tp.ID, tp.Name)
			}
		}
		if len(data.Types) != len(database.DefaultExpenseTypes)-len(excluded) {
			t.Errorf("types = %d entries, want %d", len(data.Types), len(database.DefaultExpenseTypes)-len(excluded))
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{GetExpenseTypesFn: func() ([]database.ExpenseType, error) { return nil, errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		h.GetIncomeTypes(rec, newRequest("", &HandlerCtx{Database: db}))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestGetAccountingStatistics(t *testing.T) {
	// year 0 is the "all time" query, 2024 the selected year.
	yearly := func(all, selected string) func(uint64) (decimal.Decimal, error) {
		return func(year uint64) (decimal.Decimal, error) {
			if year == 0 {
				return mustDecimal(t, all), nil
			}
			return mustDecimal(t, selected), nil
		}
	}

	db := &dbtest.FakeDB{
		GetPaymentTotalFn:          yearly("1000.00", "400.00"),
		GetExpenseTotalFn:          yearly("600.00", "200.00"),
		GetExpenseNoRefundsTotalFn: yearly("500.00", "150.00"),
		GetHeatCostTotalFn:         yearly("800.00", "300.00"),
		GetSessionPaymentTotalFn:   yearly("900.00", "350.00"),
		GetSessionRefundsTotalFn:   yearly("50.00", "20.00"),
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	r := newRequest("", &HandlerCtx{Database: db})
	h.GetAccountingStatistics(rec, r, GetAccountingStatisticsRequest{Year: 2024}, GetHandlerContext(r))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetAccountingStatisticsResponse
	decodeData(t, resp, &data)

	checks := []struct {
		name string
		got  *decimal.Decimal
		want string
	}{
		{"total_payment", data.IncomeTotal, "1000.00"},
		{"total_payment_selected_year", data.IncomeTotalSelectedYear, "400.00"},
		{"total_expenditure", data.ExpenseTotal, "600.00"},
		{"total_expenditure_selected_year", data.ExpenseTotalSelectedYear, "200.00"},
		{"total_expenditure_no_refunds", data.ExpenseNoRefundsTotal, "500.00"},
		{"total_expenditure_selected_year_no_refunds", data.ExpenseNoRefundsTotalSelectedYear, "150.00"},
		{"total_used", data.RideMinutesCashUsed, "800.00"},
		{"total_used_selected_year", data.RideMinutesCashUsedSelectedYear, "300.00"},
		{"total_session_payment", data.IncomeSessionPayment, "900.00"},
		{"total_session_payment_selected_year", data.IncomeSessionPaymentSelectedYear, "350.00"},
		{"total_session_refunds", data.RefundsSessionTotal, "50.00"},
		// session payments (900) - heat costs (800) - refunds (50)
		{"total_open", data.RideMinutesCashCredit, "50.00"},
		// income (1000) - expenses (600)
		{"current_balance", data.Balance, "400.00"},
		// heat costs (800) - expenses without refunds (500)
		{"current_session_profit", data.SessionProfit, "300.00"},
		// heat costs 2024 (300) - expenses without refunds 2024 (150)
		{"current_session_profit_selected_year", data.SessionProfitSelectedYear, "150.00"},
	}
	for _, c := range checks {
		if c.got == nil {
			t.Errorf("%s is nil", c.name)
			continue
		}
		if !c.got.Equal(mustDecimal(t, c.want)) {
			t.Errorf("%s = %s, want %s", c.name, c.got, c.want)
		}
	}
}

func TestGetAccountingStatisticsToleratesDatabaseErrors(t *testing.T) {
	// Every single total is reported as an error. The handler logs and keeps
	// going, so the response must still be a well formed success response with
	// zero values instead of a panic.
	fail := func(uint64) (decimal.Decimal, error) { return decimal.Zero, errNotFound }
	db := &dbtest.FakeDB{
		GetPaymentTotalFn:          fail,
		GetExpenseTotalFn:          fail,
		GetExpenseNoRefundsTotalFn: fail,
		GetHeatCostTotalFn:         fail,
		GetSessionPaymentTotalFn:   fail,
		GetSessionRefundsTotalFn:   fail,
	}
	h := newTestHandler(t, db, nil)

	rec := httptest.NewRecorder()
	r := newRequest("", &HandlerCtx{Database: db})
	h.GetAccountingStatistics(rec, r, GetAccountingStatisticsRequest{Year: 0}, GetHandlerContext(r))

	resp := decodeResponse(t, rec)
	if !resp.OK {
		t.Fatalf("unexpected failure: %s", resp.Msg)
	}
	var data GetAccountingStatisticsResponse
	decodeData(t, resp, &data)
	if data.Balance == nil || !data.Balance.IsZero() {
		t.Errorf("current_balance = %v, want 0", data.Balance)
	}
}

func TestGetAccountingTransactions(t *testing.T) {
	t.Run("returns the transactions of the requested year", func(t *testing.T) {
		amount := mustDecimal(t, "-42.50")
		var gotYear uint64
		db := &dbtest.FakeDB{
			GetTransactionsFn: func(year uint64) ([]database.TransactionRow, error) {
				gotYear = year
				return []database.TransactionRow{
					{ID: 1, Amount: &amount, TableID: database.TableIDExpenditure, Comment: "fuel"},
				}, nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetAccountingTransactions(rec, r, GetAccountingStatisticsRequest{Year: 2024}, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotYear != 2024 {
			t.Errorf("GetTransactions(%d), want 2024", gotYear)
		}
		var data GetAccountingTransactionsResponse
		decodeData(t, resp, &data)
		if len(data) != 1 || data[0].Comment != "fuel" {
			t.Errorf("transactions = %+v, want the single fuel transaction", data)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetTransactionsFn: func(uint64) ([]database.TransactionRow, error) { return nil, errNotFound },
		}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.GetAccountingTransactions(rec, r, GetAccountingStatisticsRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestDeleteTransaction(t *testing.T) {
	t.Run("deletes the requested row", func(t *testing.T) {
		var gotTable, gotRow uint64
		db := &dbtest.FakeDB{
			DeleteTransactionFn: func(tableID, rowID uint64) error {
				gotTable, gotRow = tableID, rowID
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteTransaction(rec, r, DeleteTransactionRequest{TableID: database.TableIDPayment, RowID: 12}, GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if gotTable != database.TableIDPayment || gotRow != 12 {
			t.Errorf("DeleteTransaction(%d, %d), want (%d, 12)", gotTable, gotRow, database.TableIDPayment)
		}
	})

	t.Run("reports a database failure", func(t *testing.T) {
		db := &dbtest.FakeDB{DeleteTransactionFn: func(uint64, uint64) error { return errNotFound }}
		h := newTestHandler(t, db, nil)
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.DeleteTransaction(rec, r, DeleteTransactionRequest{}, GetHandlerContext(r))
		if resp := decodeResponse(t, rec); resp.OK {
			t.Error("expected a failure response")
		}
	})
}

func TestDeleteTransactionValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "known table id", payload: `{"table_id":1,"row_id":3}`},
		{name: "unknown table id", payload: `{"table_id":99,"row_id":3}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &DeleteTransactionRequest{}, DeleteTransactionValidationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddIncome(t *testing.T) {
	baseReq := func() AddTransactionRequest {
		return AddTransactionRequest{
			Amount:  decimalPtr(t, "120.00"),
			Comment: "a payment",
			Date:    "2024-05-04T15:04",
			TypeID:  database.ExpenseTypeMembershipFee,
			UserID:  7,
		}
	}

	t.Run("stores the income", func(t *testing.T) {
		var stored database.Income
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
			AddIncomeFn: func(i database.Income) error {
				stored = i
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddIncome(rec, r, baseReq(), GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.UserID != 7 {
			t.Errorf("user = %d, want 7", stored.UserID)
		}
		if !stored.Amount.Equal(mustDecimal(t, "120.00")) {
			t.Errorf("amount = %s, want 120.00", stored.Amount)
		}
		if stored.ExpenseTypeID != database.ExpenseTypeMembershipFee {
			t.Errorf("type = %d, want %d", stored.ExpenseTypeID, database.ExpenseTypeMembershipFee)
		}
		want := time.Date(2024, time.May, 4, 15, 4, 0, 0, time.UTC)
		if !stored.Timestamp.Equal(want) {
			t.Errorf("timestamp = %v, want %v", stored.Timestamp, want)
		}
	})

	t.Run("fills in a default comment for session and membership payments", func(t *testing.T) {
		tests := []struct {
			typeID  uint64
			comment string
		}{
			{database.ExpenseTypeSession, "Session Payment"},
			{database.ExpenseTypeMembershipFee, "Membership Fee"},
		}
		for _, tt := range tests {
			var stored database.Income
			db := &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
				AddIncomeFn: func(i database.Income) error {
					stored = i
					return nil
				},
			}
			h := newTestHandler(t, db, nil)

			req := baseReq()
			req.Comment = ""
			req.TypeID = tt.typeID
			rec := httptest.NewRecorder()
			r := newRequest("", &HandlerCtx{Database: db})
			h.AddIncome(rec, r, req, GetHandlerContext(r))

			if resp := decodeResponse(t, rec); !resp.OK {
				t.Fatalf("unexpected failure: %s", resp.Msg)
			}
			if stored.Comment != tt.comment {
				t.Errorf("comment = %q, want %q", stored.Comment, tt.comment)
			}
		}
	})

	t.Run("failures", func(t *testing.T) {
		withUser := func() *dbtest.FakeDB {
			return &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
			}
		}

		noCommentReq := baseReq()
		noCommentReq.Comment = ""
		noCommentReq.TypeID = database.ExpenseTypeOther

		badDateReq := baseReq()
		badDateReq.Date = "04.05.2024"

		storeFailDB := withUser()
		storeFailDB.AddIncomeFn = func(database.Income) error { return errNotFound }

		tests := []struct {
			name string
			req  AddTransactionRequest
			db   *dbtest.FakeDB
		}{
			{name: "unknown user", req: baseReq(), db: &dbtest.FakeDB{}},
			{name: "unparsable date", req: badDateReq, db: withUser()},
			{name: "missing comment", req: noCommentReq, db: withUser()},
			{name: "store fails", req: baseReq(), db: storeFailDB},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddIncome(rec, r, tt.req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestAddExpense(t *testing.T) {
	baseReq := func() AddTransactionRequest {
		return AddTransactionRequest{
			Amount:  decimalPtr(t, "80.00"),
			Comment: "new rope",
			Date:    "2024-05-04T15:04",
			TypeID:  database.ExpenseTypeMaterial,
			UserID:  7,
		}
	}

	t.Run("stores the entry in the expenditure table", func(t *testing.T) {
		var stored database.Expense
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
			AddExpenseFn: func(e database.Expense) error {
				stored = e
				return nil
			},
			AddIncomeFn: func(database.Income) error {
				t.Error("an expense must not be written to the payment (income) table")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddExpense(rec, r, baseReq(), GetHandlerContext(r))

		if resp := decodeResponse(t, rec); !resp.OK {
			t.Fatalf("unexpected failure: %s", resp.Msg)
		}
		if stored.UserID != 7 {
			t.Errorf("user = %d, want 7", stored.UserID)
		}
		if !stored.Amount.Equal(mustDecimal(t, "80.00")) {
			t.Errorf("amount = %s, want 80.00", stored.Amount)
		}
		if stored.ExpenseTypeID != database.ExpenseTypeMaterial {
			t.Errorf("type = %d, want %d", stored.ExpenseTypeID, database.ExpenseTypeMaterial)
		}
		if stored.Comment != "new rope" {
			t.Errorf("comment = %q, want %q", stored.Comment, "new rope")
		}
		want := time.Date(2024, time.May, 4, 15, 4, 0, 0, time.UTC)
		if !stored.Timestamp.Equal(want) {
			t.Errorf("timestamp = %v, want %v", stored.Timestamp, want)
		}
	})

	t.Run("refuses direct fuel expenses", func(t *testing.T) {
		db := &dbtest.FakeDB{
			GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
			AddExpenseFn: func(database.Expense) error {
				t.Error("fuel expenses have to go through the boat API")
				return nil
			},
		}
		h := newTestHandler(t, db, nil)

		req := baseReq()
		req.TypeID = database.ExpenseTypeFuelDirect
		rec := httptest.NewRecorder()
		r := newRequest("", &HandlerCtx{Database: db})
		h.AddExpense(rec, r, req, GetHandlerContext(r))

		resp := decodeResponse(t, rec)
		if resp.OK || resp.Msg != "Please use boat API for fuel expenses." {
			t.Errorf("response = %+v, want the fuel expense rejection", resp)
		}
	})

	t.Run("failures", func(t *testing.T) {
		withUser := func() *dbtest.FakeDB {
			return &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 7}, nil },
			}
		}

		noCommentReq := baseReq()
		noCommentReq.Comment = ""

		badDateReq := baseReq()
		badDateReq.Date = "not-a-date"

		storeFailDB := withUser()
		storeFailDB.AddExpenseFn = func(database.Expense) error { return errNotFound }

		tests := []struct {
			name string
			req  AddTransactionRequest
			db   *dbtest.FakeDB
		}{
			{name: "unknown user", req: baseReq(), db: &dbtest.FakeDB{}},
			{name: "unparsable date", req: badDateReq, db: withUser()},
			{name: "missing comment", req: noCommentReq, db: withUser()},
			{name: "store fails", req: baseReq(), db: storeFailDB},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := newTestHandler(t, tt.db, nil)
				rec := httptest.NewRecorder()
				r := newRequest("", &HandlerCtx{Database: tt.db})
				h.AddExpense(rec, r, tt.req, GetHandlerContext(r))
				if resp := decodeResponse(t, rec); resp.OK {
					t.Error("expected a failure response")
				}
			})
		}
	})
}

func TestAddTransactionValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{name: "valid payload", payload: `{"amount":"12.50","date":"2024-05-04T15:04","type_id":5,"user_id":1,"comment":"c"}`},
		{name: "missing amount", payload: `{"date":"2024-05-04T15:04","type_id":5,"user_id":1,"comment":"c"}`, wantErr: true},
		{name: "missing date", payload: `{"amount":"12.50","type_id":5,"user_id":1,"comment":"c"}`, wantErr: true},
		{name: "unknown expense type", payload: `{"amount":"12.50","date":"2024-05-04T15:04","type_id":99,"user_id":1,"comment":"c"}`, wantErr: true},
		{name: "missing user", payload: `{"amount":"12.50","date":"2024-05-04T15:04","type_id":5,"comment":"c"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReadBodyAndValidate(newRequest(tt.payload, nil), &AddTransactionRequest{}, AddTransactionValidationErrors)
			if (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// A transaction request without an amount must never reach the handler with a
// nil Amount, because the handler dereferences it unconditionally.
func TestAddTransactionWithoutAmountIsRejectedBeforeTheHandler(t *testing.T) {
	payload := `{"date":"2024-05-04T15:04","type_id":5,"user_id":1,"comment":"c"}`

	for _, name := range []string{"income", "expense"} {
		t.Run(name, func(t *testing.T) {
			db := &dbtest.FakeDB{
				GetUserByIdFn: func(uint) (database.User, error) { return database.User{ID: 1}, nil },
			}
			h := newTestHandler(t, db, nil)

			next := h.AddIncome
			if name == "expense" {
				next = h.AddExpense
			}
			handler := WithRequestBody(next, AddTransactionValidationErrors)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, newRequest(payload, &HandlerCtx{Database: db}))

			if resp := decodeResponse(t, rec); resp.OK {
				t.Error("a transaction without an amount has to be rejected")
			}
		})
	}
}
