export default class Accounting {

  /**
   * Returns all the years that we have data for.
   */
  static getYears() {
    console.log("accounting/getYears called");
    return new Promise((resolve, reject) => {
      fetch('/api/payment.php?action=get_years', {
        method: "GET",
        cache: "no-cache",
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("accounting/getYears response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("accounting/getYears: Cannot retrieve years, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("accounting/getYears: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("accounting/getYears", error);
        reject([error]);
      })
    })
  }

  /**
   * Returns the statistics for a given year. Year can be 'any' or
   * a specific number
   */
  static getStatistics(year) {
    console.log("accounting/getStatistics called with:", year);
    return new Promise((resolve, reject) => {
      if(year == null){
        year = 'any';
      }

      const requestData = {
        year: year
      };

      fetch('/api/payment.php?action=get_statistics', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("accounting/getStatistics response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("accounting/getStatistics: Cannot retrieve statistics, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("accounting/getStatistics: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("accounting/getStatistics", error);
        reject([error]);
      })
    })
  }

  /**
   * Get all transactions for a given year. Year can be 'any' or
   * a specific number.
   * @param {*} year 
   */
  static getTransactions(year) {
    console.log("accounting/getTransactions called with:", year);
    return new Promise((resolve, reject) => {
      if(year == null){
        year = 'any';
      }

      const requestData = {
        year: year
      };

      fetch('/api/payment.php?action=get_transactions', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("accounting/getTransactions response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("accounting/getTransactions: Cannot retrieve transactions, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("accounting/getTransactions: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("accounting/getTransactions", error);
        reject([error]);
      })
    })
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
        row_id:   transaction.id
      };

      fetch('/api/payment.php?action=delete_transaction', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("accounting/deleteTransaction response data:", data);
            if(data.ok){
              resolve();
            }else{
              console.log("accounting/deleteTransaction: Cannot delete transactions, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("accounting/deleteTransaction: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("accounting/deleteTransaction", error);
        reject([error]);
      })
    })
  }
}