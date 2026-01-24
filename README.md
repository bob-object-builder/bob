![bob-banner](https://github.com/user-attachments/assets/1b786ab4-99f7-4d29-a257-e44759189938)

# BOB

**Bob** is a lightweight, declarative transpiler that allows you to define database schemas and queries using a simple, human-readable DSL (Domain Specific Language) called `bob`. Bob then converts these definitions into SQL code for various database engines, such as SQLite, MariaDB, and PostgreSQL.

## ❓ What is Bob?

Bob is designed to simplify the process of modeling databases and writing queries. Instead of writing verbose SQL, you describe your tables, relationships, and queries in a concise format. Bob parses this format and generates the corresponding SQL code automatically.

## ✨ Key Features

- 📝 **Declarative Syntax:** Define tables, fields, relationships, and queries in a clear and readable way.
- 🌐 **Multi-Database Support:** Generate SQL for SQLite, MariaDB, and PostgreSQL.
- 🔗 **Automatic Relationships:** Easily define foreign keys and joins.
- 📊 **Aggregations and Filters:** Use simple syntax for counts, filters, and groupings.
- 🧩 **Extensible:** The internal architecture allows for adding new drivers and custom logic.

## ⚙️ How Does It Work?

1. 📄 **Write a `.bob` file:**
   Describe your database schema and queries using Bob's DSL.

2. 🏃 **Run Bob:**
   Use the command line to transpile your `.bob` file into SQL:

   ```
   go run main.go -i input.bob -d sqlite
   ```

   Replace `sqlite` with `mariadb` or `postgresql` as needed.

## Flags

| Flag                  | Description                                                                | Required / Optional |
| --------------------- | -------------------------------------------------------------------------- | ------------------- |
| `-i` <br> `--input`   | Path to the input `.bob` file                                              | Optional            |
| `-q` <br> `--query`   | Direct query string input instead of file                                  | Optional            |
| `-d` <br> `--driver`  | Database driver: `mariadb`, `postgres`, or `sqlite`                      | Required            |
| `-o` <br> `--output`  | Output file path for saving the generated SQL (default output is terminal) | Optional            |
| `-v` <br> `--version` | Show Bob version and exit                                                  | Optional            |

You must include some kind of flag to provide "the query", either `-i` or `-q`.

3. 💾 **Get SQL Output:**
   By default, the generated SQL is printed to the terminal. You can specify an output file with the `-o {dirName}` to save your SQL output.
   The output is generated inside a folder.
   If no folder is specified, one named **sql** will be created automatically.

   Inside this folder you will find:
   - A **tables.sql** file containing the database definition and creation.
   - An **actions** folder with all the other operations (**INSERT**, **UPDATE**, **DELETE**, etc.).

     It is important to note that the action names are inherited directly from the source.

## 🧪 Example

```
table Users {
    id id
    name string
    email string unique index
    created_at createdAt
}

get Users {
    id
    name
    email
    if email != ""
}
```

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT UNIQUE,
    created_at DATETIME
);

SELECT id, name, email FROM users WHERE email != "";
```

## 📚 Supported Syntax

- `table TableName { ... }` — Define tables and columns.
- `get TableName { ... }` — Define queries.
- `unique`, `index` — Column modifiers.
- `-> RelatedTable foreign_key { ... }` — Joins.
- `count(*)`, `group field`, `if condition` — Aggregations, grouping, and filtering.
- 💬 Comments with `//`.

See the documentation for a complete syntax reference.

## 🏗️ Internal Architecture

- ⚙️ **Lexer/Parser:** Reads and interprets the `.bob` file.
- 🧠 **Models:** Represents tables, actions, and queries.
- 🔀 **Transpiler:** Converts the parsed model into SQL for the selected database engine.
- 🛠️ **Drivers:** Abstracts differences between SQL dialects.

## 🙌 Why Use Bob?

- ⏱️ Rapidly prototype and evolve your database schema.
- ✂️ Write less boilerplate SQL.
- 📁 Maintain database logic in a single, versionable file.
- 🔄 Easily switch between different SQL engines.

---
