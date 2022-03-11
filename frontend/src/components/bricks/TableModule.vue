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
        <div @click="function() { sortByCol(c) }">
          {{ c.label }}
          <!-- select the right sorting icon -->
          <i
            v-if="
              c.sortable != null &&
              c.sortable == true &&
              sortCol != null &&
              c.label == sortCol.label &&
              sortOrder == 0
            "
            class="bi bi-caret-down-fill"
          ></i>
          <i
            v-if="
              c.sortable != null &&
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
    <tbody v-if="tRows != null && tRows.length > 0" role="rowgroup">
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
          <slot
            :name="getCellSlotName(tColumns[c].key)"
            :row="r"
            :cell="r[tColumns[c].key]"
          >
            <!-- in case we do not have a custom template, we render either
            the value directly or use the custom text formatter -->
            <!-- use either the value directly -->
            <div
              v-if="
                tColumns[c].formatter == null &&
                tColumns[c].htmlFormatter == null
              "
            >
              {{ r[tColumns[c].key] }}
            </div>
            <!-- the custom text formatter -->
            <div v-if="tColumns[c].formatter != null">
              {{
                tColumns[c].formatter(r[tColumns[c].key], tColumns[c].key, r)
              }}
            </div>
          </slot>
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
import { sortBy, reverse, isEqual } from "lodash";
export default {
  name: "TableModule",
  props: [
    "customTableClass",
    "size",
    "columns",
    "rows",
    "rowClassFunction",
    "rowClass",
    "selectable",
    "selectMode",
  ],
  emits: ["rowclick", "select-row"],
  data() {
    return {
      sortCol: null,
      sortOrder: 0,
      tColumns: [],
      tRows: [],
      selectedRows: [],
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
      console.log("rows changed to:", newRows);
      this.selectedRows = [];
      this.$emit("select-row", []);

      if (newRows == null || newRows.length == 0) {
        this.tRows = [];
        return;
      }

      this.tRows = newRows;
      this.sort();
    },
    columns: function (newCols) {
      console.log("columns changed to:", newCols);
      this.sortCol = null;

      if (newCols == null || newCols.length == 0) {
        this.tColumns = [];
        return;
      }

      this.tColumns = newCols;
      return;
    }
  },
  methods: {
    sort: function () {
      this.sortByCol(this.sortCol, false);
    },
    sortByCol: function (col, changeSorting = true) {
      if (col == null || col.sortable == null || !col.sortable) {
        return;
      }

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
      // TODO handle sortable or not with underline
      if (col.class != null) {
        return "fw-lighter " + col.class;
      }
      return "fw-lighter text-start";
    },
    getRowClass: function (row) {
      let rowClass = "";
      if (this.rowClass != null) {
        rowClass = this.rowClass;
      }
      else if (this.rowClassFunction != null) {
        rowClass = this.rowClassFunction(row);
      }
      if (
        this.selectable == true &&
        this.selectMode == "single" &&
        this.selectedRows.length == 1 &&
        isEqual(row, this.selectedRows[0])
      ) {
        rowClass += " selected";
      }
      return rowClass;
    },
    rowClickHandler: function (row) {
      console.log("row clicked", row);
      if (this.selectable == true && this.selectMode == "single") {
        if (
          this.selectedRows.length == 1 &&
          isEqual(row, this.selectedRows[0])
        ) {
          this.selectedRows = [];
          this.$emit("select-row", []);
        } else {
          this.selectedRows = [row];
          this.$emit("select-row", this.selectedRows);
        }
        this.$forceUpdate();
      }
      this.$emit("rowclick", row);
    },
    getCellSlotName: function (columnName) {
      return "cell(" + columnName + ")";
    },
  },
  created() {
    if(this.columns != null && this.columns.length > 0) {
      this.tColumns = this.columns;
      // set default sort column
      this.sortCol = this.tColumns[0];
      console.log("TableModule: columns set to", this.columns);
    }
    if(this.rows != null && this.rows.lenth > 0) {
      this.tRows = this.rows;
      this.sort();
    }
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

.selected {
  background-color: rgb(127, 218, 224);
}
</style>
