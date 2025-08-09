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
        :key="calculateRowKey(r)"
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

<script setup>
import { ref, watch, computed } from "vue";
import { sortBy, reverse, isEqual } from "lodash";
import hash from "object-hash";

const props = defineProps({
  customTableClass: {
    type: String,
    required: false,
  },
  size: {
    type: String,
    required: false,
  },
  columns: {
    type: Array,
    required: true,
  },
  rows: {
    type: Array,
    required: true,
  },
  rowClassFunction: {
    type: Function,
    required: false,
  },
  rowClass: {
    type: String,
    required: false,
  },
  selectable: {
    type: Boolean,
    required: false,
  },
  selectMode: {
    type: String,
    required: false,
  },
});

const emit = defineEmits(["rowclick", "select-row"]);

const sortCol = ref(null);
const sortOrder = ref(0);
const tColumns = ref([]);
const tRows = ref([]);
const selectedRows = ref([]);

const tableClass = computed(() => {
  if (props.customTableClass != null) {
    return (
      "table " +
      getTableSizing() +
      " table-hover table-flex " +
      props.customTableClass
    );
  }
  return "table " + getTableSizing() + " table-hover table-flex";
});

watch(
  () => props.rows,
  (newRows) => {
    console.log("rows changed to:", newRows);
    selectedRows.value = [];
    emit("select-row", []);

    if (newRows == null || newRows.length == 0) {
      tRows.value = [];
      return;
    }

    tRows.value = newRows;
    sort();
  },
);

watch(
  () => props.columns,
  (newCols) => {
    console.log("columns changed to:", newCols);
    sortCol.value = null;

    if (newCols == null || newCols.length == 0) {
      tColumns.value = [];
      return;
    }

    tColumns.value = newCols;
    return;
  },
);

const calculateRowKey = (row) => {
  return hash(row);
};

const sort = () => {
  sortByCol(sortCol.value, false);
};

const sortByCol = (col, changeSorting = true) => {
  if (col == null || col.sortable == null || !col.sortable) {
    return;
  }

  if (
    changeSorting &&
    sortCol.value != null &&
    sortCol.value.label == col.label
  ) {
    sortOrder.value = (sortOrder.value + 1) % 2;
  } else if (changeSorting && sortCol.value != null) {
    sortOrder.value = 0;
  }
  sortCol.value = col;
  var tempRows = sortBy(tRows.value, [
    function (r) {
      return r[col.key];
    },
  ]);
  if (sortOrder.value == 0) {
    tRows.value = reverse(tempRows);
  } else {
    tRows.value = tempRows;
  }
};

const getTableSizing = () => {
  if (props.size == null) {
    return "";
  } else if (props.size == "small") {
    return "table-sm";
  }
  return "";
};

const getColumnClass = (col) => {
  // TODO handle sortable or not with underline
  if (col.class != null) {
    return "fw-lighter " + col.class;
  }
  return "fw-lighter text-start";
};

const getRowClass = (row) => {
  let rowClass = "";
  if (props.rowClass != null) {
    rowClass = props.rowClass;
  } else if (props.rowClassFunction != null) {
    rowClass = props.rowClassFunction(row);
  }
  if (
    props.selectable == true &&
    props.selectMode == "single" &&
    selectedRows.value.length == 1 &&
    isEqual(row, selectedRows.value[0])
  ) {
    rowClass += " selected";
  }
  return rowClass;
};

const rowClickHandler = (row) => {
  console.log("row clicked", row);
  if (props.selectable == true && props.selectMode == "single") {
    if (selectedRows.value.length == 1 && isEqual(row, selectedRows.value[0])) {
      selectedRows.value = [];
      emit("select-row", []);
    } else {
      selectedRows.value = [row];
      emit("select-row", selectedRows.value);
    }
  }
  emit("rowclick", row);
};

const getCellSlotName = (columnName) => {
  return "cell(" + columnName + ")";
};

if (props.columns != null && props.columns.length > 0) {
  tColumns.value = props.columns;
  // set default sort column
  sortCol.value = tColumns.value[0];
  console.log("TableModule: columns set to", props.columns);
}
if (props.rows != null && props.rows.length > 0) {
  tRows.value = props.rows;
  sort();
}
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

.table > :not(caption) > * > * {
  /* to remove any bootstrap background styling */
  background-color: rgba(0, 0, 0, 0);
}

.selected {
  background-color: rgb(127, 218, 224);
}
</style>
