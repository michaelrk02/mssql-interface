package main

import (
    "errors"
    "fmt"
    "io"
)

type column struct {
    name string
    width int
    height int
}

func (c *column) computeHeight() {
    c.height = len(c.name) / c.width
    if len(c.name) % c.width != 0 {
        c.height++
    }
}

type row struct {
    height int
    values []string
}

func (r *row) computeHeight(cw int) {
    maxHeight := 1

    for i := range r.values {
        height := len(r.values[i]) / cw
        if len(r.values[i]) % cw != 0 {
            height++
        }

        if height > maxHeight {
            maxHeight++
        }
    }

    r.height = maxHeight
}

type dataTable struct {
    maxWidth int
    columns []column
    rows []row
}

func (dt *dataTable) init(maxWidth, numColumns, numRows int) error {
    dt.maxWidth = maxWidth

    dt.columns = make([]column, numColumns)
    for i := range dt.columns {
        dt.columns[i].width = (dt.maxWidth - numColumns - 1) / numColumns
        if dt.columns[i].width == 0 {
            return errors.New("too many columns")
        }
    }

    dt.rows = make([]row, numRows)

    return nil
}

func (dt *dataTable) print(w io.Writer) {
    // Print column headers

    junction(w)
    for i := range dt.columns {
        hSeparator(w, dt.columns[i].width)
        junction(w)
    }
    fmt.Fprintf(w, "\n")

    headerHeight := 1
    for i := range dt.columns {
        height := len(dt.columns[i].name) / dt.columns[i].width
        if len(dt.columns[i].name) != 0 {
            height++
        }
        if height > headerHeight {
            headerHeight = height
        }
    }

    for i := 0; i < headerHeight; i++ {
        vSeparator(w)
        for j := range dt.columns {
            fmt.Fprintf(w, "%*s", -dt.columns[j].width, wrap(fmt.Sprintf("@%s", dt.columns[j].name), dt.columns[j].width, i))
            vSeparator(w)
        }
        fmt.Fprintf(w, "\n")
    }

    junction(w)
    for i := range dt.columns {
        hSeparator(w, dt.columns[i].width)
        junction(w)
    }
    fmt.Fprintf(w, "\n")

    for i := range dt.rows {
        for j := 0; j < dt.rows[i].height; j++ {
            vSeparator(w)
            for k := range dt.rows[i].values {
                fmt.Fprintf(w, "%*s", -dt.columns[k].width, wrap(dt.rows[i].values[k], dt.columns[k].width, j))
                vSeparator(w)
            }
            fmt.Fprintf(w, "\n")
        }

        junction(w)
        for i := range dt.columns {
            hSeparator(w, dt.columns[i].width)
            junction(w)
        }
        fmt.Fprintf(w, "\n")
    }
}

func junction(w io.Writer) {
    fmt.Fprintf(w, "+")
}

func vSeparator(w io.Writer) {
    fmt.Fprintf(w, "|")
}

func hSeparator(w io.Writer, length int) {
    for i := 0; i < length; i++ {
        fmt.Fprintf(w, "-")
    }
}

