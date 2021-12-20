package main

import (
    "database/sql"
    "errors"
    "fmt"
    "os"
    "strings"
)

func helpHandler(args map[string]string) error {
    fmt.Printf("=== COMMANDS LIST ===\n")
    fmt.Printf("Parameters in [] are optional\n")
    for e := cmds.Front(); e != nil; e = e.Next() {
        cmd := e.Value.(cmdHelper)

        fmt.Printf(" - %s -- %s\n", cmd.syntax, cmd.description)
    }
    return nil
}

func listHandler(args map[string]string) error {
    var err error

    var query string
    if args["type"] == "tables" {
        query = "EXEC sp_tables @table_owner = 'dbo', @table_type = \"'TABLE'\""
    } else if args["type"] == "views" {
        query = "EXEC sp_tables @table_owner = 'dbo', @table_type = \"'VIEW'\""
    } else if args["type"] == "procedures" {
        query = "EXEC sp_stored_procedures @sp_owner = 'dbo'"
    } else {
        return errors.New("invalid type")
    }

    var rs *sql.Rows
    rs, err = db.Query(query)
    if err != nil {
        return err
    }

    dt := generateDataTable(rs, outputWidth())
    dt.print(os.Stdout)

    return nil
}

func showHandler(args map[string]string) error {
    var err error

    var rs *sql.Rows
    rs, err = db.Query(fmt.Sprintf("SELECT * FROM %s", args["object"]))
    if err != nil {
        return err
    }

    dt := generateDataTable(rs, outputWidth())
    dt.print(os.Stdout)

    return nil
}

func describeHandler(args map[string]string) error {
    var err error

    var rs *sql.Rows
    rs, err = db.Query(fmt.Sprintf("EXEC sp_helptext '%s'", args["object"]))
    if err != nil {
        return err
    }

    var sb strings.Builder
    for rs.Next() {
        var buf string
        rs.Scan(&buf)
        sb.WriteString(buf)
    }
    fmt.Printf("Definition of %s:\n\n", args["object"])
    fmt.Printf("%s\n\n", sb.String())

    return nil
}

func executeHandler(args map[string]string) error {
    var err error

    var query string
    fmt.Printf("Press EOF when finished\n")
    fmt.Printf("=== BEGIN SQL QUERY ===\n")
    query = textInput(true)
    fmt.Printf("=== END SQL QUERY ===\n")

    var rs *sql.Rows
    rs, err = db.Query(query)
    if err != nil {
        return err
    }

    var dt *dataTable

    if _, ok := args["file"]; ok {
        cw := 15

        file := args["file"]
        if _, ok := args["cw"]; ok {
            fmt.Sscanf(args["cw"], "%d", &cw)
        }

        var f *os.File
        f, err = os.Create(file)
        if err != nil {
            return err
        }
        defer f.Close()

        cols, _ := rs.Columns()
        dt = generateDataTable(rs, 1 + len(cols) * (cw + 1))
        dt.print(f)
    } else {
        dt = generateDataTable(rs, outputWidth())
        dt.print(os.Stdout)
    }

    return nil
}

func execHandler(args map[string]string) error {
    var err error

    query := args["statement"]

    var rs *sql.Rows
    rs, err = db.Query(query)
    if err != nil {
        return err
    }

    var dt *dataTable

    if _, ok := args["file"]; ok {
        cw := 15

        file := args["file"]
        if _, ok := args["cw"]; ok {
            fmt.Sscanf(args["cw"], "%d", &cw)
        }

        var f *os.File
        f, err = os.Create(file)
        if err != nil {
            return err
        }
        defer f.Close()

        cols, _ := rs.Columns()
        dt = generateDataTable(rs, 1 + len(cols) * (cw + 1))
        dt.print(f)
    } else {
        dt = generateDataTable(rs, outputWidth())
        dt.print(os.Stdout)
    }

    return nil
}

func exitHandler(args map[string]string) error {
    running = false
    return nil
}

