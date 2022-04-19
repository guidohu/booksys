import Request from "@/api/common/request.js";

export default class Accounting {
  /**
   * Returns all the years that we have data for.
   */
  static getYears() {
    console.log("accounting/getYears called");
    return Request.getRequest('/api/payment.php?action=get_years');
  }

  /**
   * Returns all the expense types that there are.
   */
  static getExpenseTypes() {
    console.log("accounting/getExpenseTypes called");
    return Request.getRequest('/api/payment.php?action=get_expenditure_types');
  }

  /**
   * Returns all the income types that there are.
   */
  static getIncomeTypes() {
    console.log("accounting/getIncomeTypes called");
    return Request.getRequest('/api/payment.php?action=get_payment_types');
  }

  /**
   * Returns the statistics for a given year. Year can be 'any' or
   * a specific number
   */
  static getStatistics(year) {
    console.log("accounting/getStatistics called with:", year);
    if (year == null) {
      year = "any";
    }

    const requestData = {
      year: year,
    };

    return Request.postRequest("/api/payment.php?action=get_statistics", requestData);
  }

  /**
   * Get all transactions for a given year. Year can be 'any' or
   * a specific number.
   * @param {*} year
   */
  static getTransactions(year) {
    console.log("accounting/getTransactions called with:", year);
    if (year == null) {
      year = "any";
    }

    const requestData = {
      year: year,
    };

    return Request.postRequest("/api/payment.php?action=get_transactions", requestData);
  }

  /**
   * Delete a specific transaction.
   * @param {*} transaction
   */
  static deleteTransaction(transaction) {
    console.log("accounting/deleteTransaction called with:", transaction);
    const requestData = {
      table_id: transaction.tbl,
      row_id: transaction.id,
    };
    return Request.postRequest("/api/payment.php?action=delete_transaction", requestData);
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
    console.log("accounting/addIncome called with:", incomeTransaction);
    const requestData = {
      amount: incomeTransaction.amount,
      type_id: incomeTransaction.typeId,
      date: incomeTransaction.date,
      user_id: incomeTransaction.userId,
      comment: incomeTransaction.comment,
    };
    return Request.postRequest("/api/payment.php?action=add_payment", requestData);
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
    console.log("accounting/addExpense called with:", expenseTransaction);
    const requestData = {
      amount: expenseTransaction.amount,
      type_id: expenseTransaction.typeId,
      date: expenseTransaction.date,
      user_id: expenseTransaction.userId,
      comment: expenseTransaction.comment,
    };
    return Request.postRequest("/api/payment.php?action=add_expenditure", requestData);
  }
}
