import Accounting from "@/api/accounting";

const state = () => ({
  selectedYear: 'any',
  years: [],                      // available years
  expenseTypes: [],               // types of different expenses
  incomeTypes: [],                // types of different incomes
  transactions: [],
  totalPayments: 0,
  totalPaymentsAllTime: 0,
  totalSessionPayments: 0,
  totalSessionPaymentsAllTime: 0,
  totalExpenditures: 0,
  totalExpendituresAllTime: 0,
  balance: 0,
  totalSessionUsed: 0,
  totalSessionUsedAllTime: 0,
  sessionProfit: 0,
  sessionProfitAllTime: 0
});

const getters = {
  getYears: (state) => {
    return state.years;
  },
  getExpenseTypes: (state) => {
    return state.expenseTypes;
  },
  getIncomeTypes: (state) => {
    return state.incomeTypes;
  },
  getTotalPayments: (state) => {
    return Math.round(Number(state.totalPayments)*100)/100;
  },
  getTotalPaymentsAllTime: (state) => {
    return Math.round(Number(state.totalPaymentsAllTime)*100)/100;
  },
  getTotalSessionPayments: (state) => {
    return Math.round(Number(state.totalSessionPayments)*100)/100;
  },
  getTotalSessionPaymentsAllTime: (state) => {
    return Math.round(Number(state.totalSessionPaymentsAllTime)*100)/100;
  },
  getTotalExpenditures: (state) => {
    return Math.round(Number(state.totalExpenditures)*100)/100;
  },
  getTotalExpendituresAllTime: (state) => {
    return Math.round(Number(state.totalExpendituresAllTime)*100)/100;
  },
  getBalance: (state) => {
    return Math.round(Number(state.balance)*100)/100;
  },
  getTotalSessionUsed: (state) => {
    return Math.round(Number(state.totalSessionUsed)*100)/100;
  },
  getTotalSessionUsedAllTime: (state) => {
    return Math.round(Number(state.totalExpendituresAllTime)*100)/100;
  },
  getSessionProfit: (state) => {
    return Math.round(Number(state.sessionProfit)*100)/100;
  },
  getSessionProfitAllTime: (state) => {
    return Math.round(Number(state.sessionProfitAllTime)*100)/100;
  },
  getSessionsBalance: (state) => {
    return Math.round(Number(state.totalSessionPaymentsAllTime - state.totalSessionUsedAllTime)*100)/100;
  },
  getTransactions: (state) => {
    return state.transactions;
  }
};

const actions = {
  queryYears ({ commit }) {
    console.log("Action triggered: queryYears");
    return new Promise((resolve, reject) => {
      Accounting.getYears()
      .then((years) => {
        commit('setYears', years);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  queryExpenseTypes ({ commit }) {
    console.log("Action triggered: queryExpenseTypes");
    return new Promise((resolve, reject) => {
      Accounting.getExpenseTypes()
      .then((types) => {
        commit('setExpenseTypes', types);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  queryIncomeTypes ({ commit }) {
    console.log("Action triggered: queryIncomeTypes");
    return new Promise((resolve, reject) => {
      Accounting.getIncomeTypes()
      .then((types) => {
        commit('setIncomeTypes', types);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  queryStatistics ({ commit }, year) {
    console.log("Action triggered: queryStatistics");
    return new Promise((resolve, reject) => {
      Accounting.getStatistics(year)
      .then((statistics) => {
        commit('setSelectedYear', year);
        commit('setStatistics', statistics);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  queryTransactions ({ commit }, year) {
    console.log("Action triggered: queryTransactions");
    return new Promise((resolve, reject) => {
      Accounting.getTransactions(year)
      .then((statistics) => {
        commit('setSelectedYear', year);
        commit('setTransactions', statistics);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  deleteTransaction ({ dispatch, state }, transaction) {
    console.log("Action triggered: deleteTransaction");
    return new Promise((resolve, reject) => {
      Accounting.deleteTransaction(transaction)
      .then(() => {
        dispatch('queryTransactions', state.selectedYear);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  addIncome ({ dispatch, state }, income) {
    console.log("Action triggered: addIncome");
    return new Promise((resolve, reject) => {
      Accounting.addIncome(income)
      .then(() => {
        dispatch('queryTransactions', state.selectedYear);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  },
  addExpense ({ dispatch, state }, expense) {
    console.log("Action triggered: addExpense");
    return new Promise((resolve, reject) => {
      Accounting.addExpense(expense)
      .then(() => {
        dispatch('queryTransactions', state.selectedYear);
        resolve();
      })
      .catch((errors) => {
        reject(errors);
      })
    })
  }
};

const mutations = {
  setYears (state, years){
    state.years = years;
  },
  setExpenseTypes (state, types){
    state.expenseTypes = types;
  },
  setIncomeTypes (state, types){
    state.incomeTypes = types
  },
  setStatistics (state, statistics){
    console.log('TODO', statistics);

    // Sessions
    state.totalSessionPayments = statistics.total_session_payment_selected_year;
    state.totalSessionPaymentsAllTime = statistics.total_session_payment;

    state.totalSessionUsed = statistics.total_used_selected_year;
    state.totalSessionUsedAllTime = statistics.total_used;

    // Total payments
    state.totalPayments = statistics.total_payment_selected_year;
    state.totalPaymentsAllTime = statistics.total_payment;

    // Total expenditure
    state.totalExpenditures = statistics.total_expenditure_selected_year;
    state.totalExpendituresAllTime = statistics.total_expenditure;
    // total_expenditure_no_refunds: "180100.54"
    // total_expenditure_selected_year_no_refunds: "7.00"
    
    // Balance
    state.balance = statistics.current_balance;

    // Profit/Loss
    state.sessionProfit = statistics.current_session_profit_selected_year;
    state.sessionProfitAllTime = statistics.current_profit_selected_year;
    
    // admin_minutes: 18584
    // admin_minutes_selected_year: 0
    // currency: "CHF"
    // guest_minutes: 0
    // guest_minutes_selected_year: 0
    // member_minutes: 59441
    // member_minutes_selected_year: 0
    // total_open: 1966.4599999999919
  },
  setTransactions (state, transactions){
    state.transactions = transactions;
  },
  setSelectedYear (state, year){
    state.selectedYear = year;
  }
};

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}

