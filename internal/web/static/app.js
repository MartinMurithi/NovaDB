let currentTable = null;

async function fetchTables() {
  const res = await fetch("/tables");
  const data = await res.json();
  const list = document.getElementById("tablesList");
  list.innerHTML = "";
  data.tables.forEach((t) => {
    const li = document.createElement("li");
    li.textContent = t;
    li.onclick = () => loadTable(t);
    list.appendChild(li);
  });
}

async function loadTable(name) {
  currentTable = name;
  document.getElementById("tableTitle").textContent = name;

  // Describe table
  const descRes = await fetch(`/table/${name}/describe`);
  const desc = await descRes.json();
  const colDiv = document.getElementById("columnsInfo");
  colDiv.innerHTML =
    "<b>Columns:</b> " +
    desc.columns.map((c) => `${c.name}(${c.type})`).join(", ");

  // Fetch rows
  const rowsRes = await fetch(`/table/${name}`);
  const rows = await rowsRes.json();
  renderRows(rows, desc.columns);
}

function renderRows(rows, columns) {
  const tableDiv = document.getElementById("tableData");
  tableDiv.innerHTML = "";

  if (rows.length === 0) {
    tableDiv.textContent = "(no rows)";
    renderAddRowForm(columns);
    return;
  }

  const table = document.createElement("table");
  const header = document.createElement("tr");
  columns.forEach((c) => {
    const th = document.createElement("th");
    th.textContent = c.name;
    header.appendChild(th);
  });
  const thActions = document.createElement("th");
  thActions.textContent = "Actions";
  header.appendChild(thActions);
  table.appendChild(header);

  rows.forEach((r) => {
    const tr = document.createElement("tr");
    columns.forEach((c) => {
      const td = document.createElement("td");
      td.textContent = r[c.name] ?? "";
      tr.appendChild(td);
    });
    // Actions
    const td = document.createElement("td");
    const delBtn = document.createElement("button");
    delBtn.textContent = "Delete";
    delBtn.onclick = () => deleteRow(r.id);
    td.appendChild(delBtn);
    tr.appendChild(td);

    table.appendChild(tr);
  });

  tableDiv.appendChild(table);
  renderAddRowForm(columns);
}

function renderAddRowForm(columns) {
  const formDiv = document.getElementById("addRowForm");
  formDiv.innerHTML = "";
  columns.forEach((c) => {
    const inp = document.createElement("input");
    inp.placeholder = c.name;
    inp.id = "row_" + c.name;
    formDiv.appendChild(inp);
  });
  const btn = document.createElement("button");
  btn.textContent = "Add Row";
  btn.onclick = addRow;
  formDiv.appendChild(btn);
}

async function addRow() {
  if (!currentTable) return;
  const inputs = document.querySelectorAll("#addRowForm input");
  const body = {};
  inputs.forEach((i) => {
    body[i.placeholder] = i.value;
  });

  await fetch(`/table/${currentTable}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  loadTable(currentTable);
}

async function deleteRow(id) {
  if (!currentTable) return;
  await fetch(`/table/${currentTable}/${id}`, { method: "DELETE" });
  loadTable(currentTable);
}

async function createTable() {
  const name = document.getElementById("newTableName").value;
  if (!name) return;
  await fetch("/table", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  });
  fetchTables();
}

async function addColumn() {
  if (!currentTable) return;
  const colName = document.getElementById("colName").value;
  const colType = document.getElementById("colType").value;
  if (!colName) return;

  await fetch(`/table/${currentTable}/column`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name: colName, type: colType }),
  });
  loadTable(currentTable);
}

// Initial load
fetchTables();
