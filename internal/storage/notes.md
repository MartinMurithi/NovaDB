[Step 1] Table Creation
-----------------------
Create table "users"
Table.Columns = []
Table.Rows = []
Table.PrimaryIndex = {}

          |
          v

[Step 2] Define Columns (Schema)
--------------------------------
Add columns to Table.Columns:
- Column 1: Name="id", Type=INT, IsPrimaryKey=true
- Column 2: Name="name", Type=TEXT
- Column 3: Name="email", Type=TEXT, IsUnique=true

Schema is now set.
          |
          v

[Step 3] Insert Row
-------------------
User provides row data:
Row.Data = {"id": 1, "name": "Alice", "email": "a@b.com"}

Insert workflow:
1. Validate row keys exist in schema
2. Validate row value types match ColumnType
3. Extract PK: pkValue = Row.Data["id"] → 1
4. Check PrimaryIndex[pkValue] → not exists → OK
5. (Optional) Check unique constraints (email)
6. Append row to Table.Rows
7. Update PrimaryIndex[pkValue] = index of row in Rows slice

          |
          v

[Step 4] Table State After Insert
---------------------------------
Table.Rows = [
    0: {"id":1, "name":"Alice", "email":"a@b.com"}
]

Table.PrimaryIndex = {
    1 → 0
}

Columns unchanged
