// Hak Cipta (C) 2021, Michael Raditya Krisnadhi
package main

import (
    "container/list"
    "database/sql"
    "encoding/json"
    "fmt"
    "net/url"
    "os"

    _ "github.com/denisenkom/go-mssqldb"
)

var running bool
var cmds *list.List
var db *sql.DB

type DBConfig struct {
    Host string `json:"host"`
    User string `json:"user"`
    Pass string `json:"pass"`
    Name string `json:"name"`
}

func initDatabase() error {
    var err error

    var f *os.File
    f, err = os.Open("db.json")
    if err != nil {
        return err
    }
    defer f.Close()

    dec := json.NewDecoder(f)

    var conf DBConfig
    err = dec.Decode(&conf)
    if err != nil {
        return err
    }

    conf.User = url.QueryEscape(conf.User)
    conf.Pass = url.QueryEscape(conf.Pass)
    conf.Name = url.QueryEscape(conf.Name)
    dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", conf.User, conf.Pass, conf.Host, conf.Name)

    db, err = sql.Open("sqlserver", dsn)
    if err != nil {
        return err
    }

    return nil
}

func initCommands() error {
    cmds = list.New()

    cmds.PushBack(cmdHelper{syntax: "help", description: "show list of commands", handler: helpHandler})
    cmds.PushBack(cmdHelper{syntax: "execute [file] [cw]", description: "executes sql and saves to [file] with column width [cw] (default cw=15)", handler: executeHandler})
    cmds.PushBack(cmdHelper{syntax: "exec <statement> [file] [cw]", description: "executes sql <statement> (inline) and saves to [file] with column width [cw] (default cw=15)", handler: execHandler})
    cmds.PushBack(cmdHelper{syntax: "list <type>", description: "shows list of <type>(s) where <type> is one of: tables, views, procedures", handler: listHandler})
    cmds.PushBack(cmdHelper{syntax: "show <object>", description: "queries the contents of <object> (table or view)", handler: showHandler})
    cmds.PushBack(cmdHelper{syntax: "describe <object>", description: "show the definition of <object> (view or procedure)", handler: describeHandler})
    cmds.PushBack(cmdHelper{syntax: "exit", description: "exits the program", handler: exitHandler})

    return nil
}

func main() {
    var err error

    err = initDatabase()
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = initCommands()
    if err != nil {
        panic(err)
    }
    

    cmdRouter := make(map[string]cmdRoute)
    for e := cmds.Front(); e != nil; e = e.Next() {
        cmd := e.Value.(cmdHelper)

        var route cmdRoute
        route.desc, err = parseCmdDesc(cmd.syntax)
        if err != nil {
            fmt.Printf("E: %s\n", err)
        }
        route.handler = cmd.handler

        tokens := parse(cmd.syntax)
        cmdRouter[tokens[0]] = route
    }

    fmt.Printf("Welcome to MS SQL interface!\n")
    fmt.Printf("Enter `help` for list of available commands\n")

    running = true
    for running {
        var cmd string
        fmt.Printf("$ ")
        cmd = textInput(false)

        tokens := parse(cmd)
        if _, ok := cmdRouter[tokens[0]]; ok {
            var args map[string]string
            args, err = cmdRouter[tokens[0]].desc.parse(tokens)
            if err == nil {
                err = cmdRouter[tokens[0]].handler(args)
                if err != nil {
                    fmt.Printf("E: %s\n", err)
                }
            } else {
                fmt.Printf("E: %s\n", err)
            }
        } else {
            fmt.Printf("command not found :(\n")
        }
        fmt.Printf("\n")
    }
    fmt.Printf("Good bye\n")
}

