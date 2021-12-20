package main

import (
    "container/list"
    "database/sql"
)

func generateDataTable(rs *sql.Rows, maxWidth int) *dataTable {
    dt := new(dataTable)

    var cols []string
    cols, _ = rs.Columns()

    ls := list.New()
    for rs.Next() {
        t := make([]string, len(cols))
        tp := make([]interface{}, len(cols))
        for i := range tp {
            tp[i] = &t[i]
        }
        rs.Scan(tp...)
        ls.PushBack(t)
    }

    dt.init(maxWidth, len(cols), ls.Len())
    for i := range cols {
        dt.columns[i].name = cols[i]
        dt.columns[i].computeHeight()
    }
    for i := 0; ls.Front() != nil; i++ {
        t := ls.Front().Value.([]string)

        dt.rows[i].values = make([]string, len(cols))
        for j := range cols {
            dt.rows[i].values[j] = t[j]
            dt.rows[i].computeHeight(dt.columns[0].width)
        }

        ls.Remove(ls.Front())
    }

    return dt
}

