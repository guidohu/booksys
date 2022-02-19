<template>
  <table :class="tableClass" role="table">
    <!-- table header -->
    <thead role="rowgroup">
      <th
        v-for="c in tColumns"
        :key="c.key"
        role="columnheader"
        :class="getColumnClass(c)"
      >
        <div @click="sortByCol(c)">
          {{ c.label }}
          <!-- select the right sorting icon -->
          <i
            v-if="
              c.sortable == true &&
              sortCol != null &&
              c.label == sortCol.label &&
              sortOrder == 0
            "
            class="bi bi-caret-down-fill"
          ></i>
          <i
            v-if="
              c.sortable == true &&
              sortCol != null &&
              c.label == sortCol.label &&
              sortOrder == 1
            "
            class="bi bi-caret-up-fill"
          ></i>
        </div>
      </th>
    </thead>
    <!-- table body -->
    <tbody v-if="tRows.length > 0" role="rowgroup">
      <tr
        v-for="r in tRows"
        :key="r.id"
        role="row"
        :class="getRowClass(r)"
        @click="rowClickHandler(r)"
      >
        <td
          v-for="c in Array(tColumns.length).keys()"
          :key="c"
          :class="getColumnClass(tColumns[c])"
        >
          <!-- use either the value directly or the custom formatter -->
          <div v-if="tColumns[c].formatter == null">
            {{ r[tColumns[c].key] }}
          </div>
          <div v-else>
            {{ tColumns[c].formatter(r[tColumns[c].key], tColumns[c].key, r) }}
          </div>
        </td>
      </tr>
    </tbody>
    <tbody v-else>
      <tr class="text-center">
        <td :colspan="tColumns.length">no entries</td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import { sortBy, reverse } from "lodash";
export default {
  name: "TableModule",
  props: ["customTableClass", "size", "columns", "rows", "rowClassFunction"],
  emits: ["rowclick"],
  data() {
    return {
      sortCol: null,
      sortOrder: 0,
      tColumns: [],
      tRows: [],
    };
  },
  computed: {
    tableClass: function () {
      if (this.customTableClass != null) {
        return (
          "table " +
          this.getTableSizing() +
          " table-hover table-flex " +
          this.customTableClass
        );
      }
      return "table " + this.getTableSizing() + " table-hover table-flex";
    },
  },
  watch: {
    rows: function (newRows) {
      // sort by col
      if (newRows == null || newRows.length == 0) {
        this.tRows = [];
      }

      this.tRows = newRows;
      this.sort();
    },
  },
  methods: {
    sort: function () {
      this.sortByCol(this.sortCol, false);
    },
    sortByCol: function (col, changeSorting = true) {
      if (
        changeSorting &&
        this.sortCol != null &&
        this.sortCol.label == col.label
      ) {
        this.sortOrder = (this.sortOrder + 1) % 2;
      } else if (changeSorting && this.sortCol != null) {
        this.sortOrder = 0;
      }
      this.sortCol = col;
      var tempRows = sortBy(this.tRows, [
        function (r) {
          return r[col.key];
        },
      ]);
      if (this.sortOrder == 0) {
        this.tRows = reverse(tempRows);
      } else {
        this.tRows = tempRows;
      }
    },
    getTableSizing: function () {
      if (this.size == null) {
        return "";
      } else if (this.size == "small") {
        return "table-sm";
      }
      return "";
    },
    getColumnClass: function (col) {
      if (col.class != null) {
        return "fw-lighter " + col.class;
      }
      return "fw-lighter text-start";
    },
    getRowClass: function (row) {
      if (this.rowClassFunction != null) {
        return this.rowClassFunction(row);
      }
      return "";
    },
    rowClickHandler: function (row) {
      this.$emit("rowclick", row);
    },
  },
  created() {
    this.tColumns = this.columns;
    this.tRows = this.rows;
    this.sortCol = this.tColumns[0];
    this.sort();
  },
};
</script>

<style>
.table-size {
  overflow-y: scroll;
  width: 100%;
  height: 100%;
}

.table-flex {
  border-collapse: collapse;
  width: 100%;
  overflow: scroll;
}

thead th {
  background-color: rgb(255, 255, 255);
  position: sticky;
  top: 0;
}

.table > :not(:first-child) {
  border-top: 1.5px solid currentColor;
}
</style>
