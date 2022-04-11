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
    return new Promise((resolve, reject) => {
      if (year == null) {
        year = "any";
      }

      const requestData = {
        year: year,
      };

      fetch("/api/payment.php?action=get_statistics", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("accounting/getStatistics response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.log(
                  "accounting/getStatistics: Cannot retrieve statistics, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "accounting/getStatistics: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("accounting/getStatistics", error);
          reject([error]);
        });
    });
  }

  /**
   * Get all transactions for a given year. Year can be 'any' or
   * a specific number.
   * @param {*} year
   */
  static getTransactions(year) {
    console.log("accounting/getTransactions called with:", year);
    return new Promise((resolve, reject) => {
      if (year == null) {
        year = "any";
      }

      const requestData = {
        year: year,
      };

      fetch("/api/payment.php?action=get_transactions", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("accounting/getTransactions response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.log(
                  "accounting/getTransactions: Cannot retrieve transactions, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "accounting/getTransactions: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("accounting/getTransactions", error);
          reject([error]);
        });
    });
  }

  /**
   * Delete a specific transaction.
   * @param {*} transaction
   */
  static deleteTransaction(transaction) {
    console.log("accounting/deleteTransaction called with:", transaction);
    return new Promise((resolve, reject) => {
      const requestData = {
        table_id: transaction.tbl,
        row_id: transaction.id,
      };

      fetch("/api/payment.php?action=delete_transaction", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("accounting/deleteTransaction response data:", data);
              if (data.ok) {
                resolve();
              } else {
                console.log(
                  "accounting/deleteTransaction: Cannot delete transactions, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "accounting/deleteTransaction: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("accounting/deleteTransaction", error);
          reject([error]);
        });
    });
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
    return new Promise((resolve, reject) => {
      const requestData = {
        amount: incomeTransaction.amount,
        type_id: incomeTransaction.typeId,
        date: incomeTransaction.date,
        user_id: incomeTransaction.userId,
        comment: incomeTransaction.comment,
      };

      fetch("/api/payment.php?action=add_payment", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("accounting/addIncome response data:", data);
              if (data.ok) {
                resolve();
              } else {
                console.log(
                  "accounting/addIncome: Cannot add income transactions, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "accounting/addIncome: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("accounting/addIncome", error);
          reject([error]);
        });
    });
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
    return new Promise((resolve, reject) => {
      const requestData = {
        amount: expenseTransaction.amount,
        type_id: expenseTransaction.typeId,
        date: expenseTransaction.date,
        user_id: expenseTransaction.userId,
        comment: expenseTransaction.comment,
      };

      fetch("/api/payment.php?action=add_expenditure", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("accounting/addExpense response data:", data);
              if (data.ok) {
                resolve();
              } else {
                console.log(
                  "accounting/addExpense: Cannot add income transactions, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "accounting/addExpense: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("accounting/addExpense", error);
          reject([error]);
        });
    });
  }
}
