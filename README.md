# NovaDB

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**NovaDB** is a lightweight, in-memory relational database management system (RDBMS) written in Go. It provides a simple SQL-like interface through a **REPL** (Read-Eval-Print Loop) environment. NovaDB supports:

- Creating tables with multiple column types
- Adding columns dynamically
- CRUD operations (Create, Read, Update, Delete)
- Listing all tables and describing table structure
- SQL query parsing, planning, and execution

This project was inspired by the **Pesapal RDBMS Challenge**, designed to demonstrate building a simple RDBMS from scratch with a practical REPL interface.

---

## Features

1. **REPL Mode**
   - Interactive terminal to execute SQL queries.
   - Syntax highlighting for basic SQL keywords.
   - Multi-line SQL queries supported (end queries with `;`).

2. **Table Management**
   - Create new tables: `CREATE TABLE table_name;`
   - List all tables: `SHOW TABLES;`
   - Describe a table's structure: `DESCRIBE table_name;`
   - Add columns to existing tables: `ALTER TABLE table_name ADD COLUMN column_name TYPE;`

3. **CRUD Operations**
   - Insert rows: `INSERT INTO table_name (columns) VALUES (values);`
   - Select rows: `SELECT * FROM table_name;`
   - Update rows: `UPDATE table_name SET column=value WHERE id=...;`
   - Delete rows: `DELETE FROM table_name WHERE id=...;`

4. **Supported Column Types**
   - `INT` → integer numbers
   - `TEXT` → strings
   - `FLOAT` → floating-point numbers
   - `BOOL` → true/false
   - `DATE` → date values

5. **In-memory Storage**
   - No external database required.
   - Data exists only during runtime of the REPL.

---

## Demo

Below is a 1-minute demo showing NovaDB in action via the REPL:

<video src="novadb_demo_1.mp4" width="800" autoplay loop muted></video>

- Create a table `users`
- Add columns `id`, `names`, `age`
- Insert rows
- Query rows with `SELECT`
- Update and delete rows dynamically
- Describe table structure
- List all tables

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://golang.org/doc/install) installed
- Git installed

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/NovaDB.git
cd NovaDB
Running the REPL
bash
Copy code
go run ./cmd/novadb --mode=repl
You will see:

pgsql
Copy code
NovaDB REPL. Type 'exit;' to quit. End SQL with ';'
> 
Example session:

sql
Copy code
> CREATE TABLE users;
Table 'users' created successfully

> ALTER TABLE users ADD COLUMN id INT;
Column 'id' added to 'users'

> ALTER TABLE users ADD COLUMN names TEXT;
Column 'names' added to 'users'

> INSERT INTO users (id, names) VALUES (1, 'Alice');
Query executed successfully

> SELECT * FROM users;
id  names
--  -----
1   Alice
Exiting
sql
Copy code
> EXIT;
Bye!
Project Structure
graphql
Copy code
NovaDB/
├─ cmd/          # Entry point for REPL
├─ internal/
│  ├─ engine/    # Execution engine
│  ├─ parser/    # SQL parser
│  ├─ planner/   # Query planner
│  └─ storage/   # In-memory storage structures
├─ assets/       # Demo video and GIF for README
├─ README.md
└─ go.mod
Contributing
NovaDB is an educational project and open to contributions.
To contribute:

Fork the repository.

Create a branch: git checkout -b feature/my-feature

Commit your changes: git commit -m "Add feature"

Push: git push origin feature/my-feature

Open a Pull Request

License
This project is licensed under the MIT License - see the LICENSE file for details.

Notes
This project is focused on learning and experimentation with database internals.

All data is stored in-memory; restarting the REPL clears all tables and rows.

The UI mode is excluded for simplicity; the focus is on understanding the REPL and database engine.

Author: Martin Murithi
Inspired by the Pesapal RDBMS Challenge.
