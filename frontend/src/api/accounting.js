import Request from "@/api/common/request.js";

export default class Accounting {
  /**
   * Returns all the years that we have data for.
   */
  static getYears() {
    console.log("/api/v2/accounting/years/list called");
    return Request.getRequest('/api/v2/accounting/years/list');
  }

  /**
   * Returns all the expense types that there are.
   */
  static getExpenseTypes() {
    console.log("/api/v2/accounting/expense_types/list");
    return Request.getRequest('/api/v2/accounting/expense_types/list');
  }

  /**
   * Returns all the income types that there are.
   */
  static getIncomeTypes() {
    console.log("/api/v2/accounting/income_types/list called");
    return Request.getRequest('/api/v2/accounting/income_types/list');
  }

  /**
   * Returns the statistics for a given year. Year can be 'any' or
   * a specific number
   */
  static getStatistics(year) {
    console.log("/api/v2/accounting/statistics/get called with year:", year);
    if (year == null) {
      year = 0;
    }
    const requestData = {
      year: parseInt(year),
    };

    return Request.postRequest("/api/v2/accounting/statistics/get", requestData);
  }

  /**
   * Get all transactions for a given year. Year can be 'any' or
   * a specific number.
   * @param {*} year
   */
  static getTransactions(year) {
    console.log("/api/v2/accounting/transactions/get called with:", year);
    if (year == null) {
      year = 0;
    }

    const requestData = {
      year: parseInt(year),
    };

    return Request.postRequest("/api/v2/accounting/transactions/get", requestData);
  }

  /**
   * Delete a specific transaction.
   * @param {*} transaction
   */
  static deleteTransaction(transaction) {
    console.log("/api/v2/accounting/transactions/delete called with:", transaction);
    const requestData = {
      table_id: parseInt(transaction.tbl),
      row_id: parseInt(transaction.id),
    };
    return Request.postRequest("/api/v2/accounting/transactions/delete", requestData);
  }

  /**
   * Adds an income transaction
   * @param {*} incomeTransaction containing the following fields:
   * {
   *  amount:  ...
   *  typeId:  ...         income type ID
   *  date:    ...
   *  userId:  ...         ID of the user associated with this transaction
   *  comment: ...
   * }
   */
  static addIncome(incomeTransaction) {
    console.log("/api/v2/accounting/income/add called with:", incomeTransaction);
    const requestData = {
      amount: incomeTransaction.amount,
      type_id: incomeTransaction.typeId,
      date: incomeTransaction.date,
      user_id: incomeTransaction.userId,
      comment: incomeTransaction.comment,
    };
    return Request.postRequest("/api/v2/accounting/income/add", requestData);
  }

  /**
   * Adds an expense transaction
   * @param {*} expenseTransaction containing the following fields:
   * {
   *  amount:  ...
   *  typeId:  ...         expense type ID
   *  date:    ...
   *  userId:  ...         ID of the user associated with this transaction
   *  comment: ...
   * }
   */
  static addExpense(expenseTransaction) {
    console.log("/api/v2/accounting/expense/add called with:", expenseTransaction);
    const requestData = {
      amount: expenseTransaction.amount,
      type_id: expenseTransaction.typeId,
      date: expenseTransaction.date,
      user_id: expenseTransaction.userId,
      comment: expenseTransaction.comment,
    };
    return Request.postRequest("/api/v2/accounting/expense/add", requestData);
  }
}
