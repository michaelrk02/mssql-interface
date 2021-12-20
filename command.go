package main

import (
    "errors"
    "fmt"
    "strings"
)

type cmdArg struct {
    name string
    optional bool
}

type cmdHelper struct {
    syntax string
    description string
    handler func(map[string]string) error
}

type cmdRoute struct {
    desc *cmdDesc
    handler func(map[string]string) error
}

//
// cmdDesc structure
//

type cmdDesc struct {
    name string
    args []cmdArg
}

func parseCmdDesc(str string) (*cmdDesc, error) {
    var err error

    desc := new(cmdDesc)

    tokens := parse(str)
    desc.args = make([]cmdArg, len(tokens) - 1)

    for i := range tokens {
        if i == 0 {
            desc.name = tokens[i]
        } else {
            if (tokens[i][0] == '<') && (tokens[i][len(tokens[i]) - 1] == '>') {
                if (i <= 1) || ((i > 1) && !desc.args[i - 2].optional) {
                    desc.args[i - 1].name = tokens[i][1:(len(tokens[i]) - 1)]
                    desc.args[i - 1].optional = false
                } else {
                    err = errors.New("no required parameters are allowed after an optional parameter")
                }
            } else if (tokens[i][0] == '[') && (tokens[i][len(tokens[i]) - 1] == ']') {
                desc.args[i - 1].name = tokens[i][1:(len(tokens[i]) - 1)]
                desc.args[i - 1].optional = true
            } else {
                err = errors.New("invalid format")
            }
        }
    }

    return desc, err
}

func (desc *cmdDesc) String() string {
    var sb strings.Builder
    sb.WriteString(desc.name)
    for i := range desc.args {
        if !desc.args[i].optional {
            sb.WriteString(fmt.Sprintf(" <%s>", desc.args[i].name))
        } else {
            sb.WriteString(fmt.Sprintf(" [%s]", desc.args[i].name))
        }
    }
    return sb.String()
}

func (desc *cmdDesc) parse(tokens []string) (map[string]string, error) {
    var err error

    values := make(map[string]string)
    for i := range desc.args {
        if i + 1 < len(tokens) {
            values[desc.args[i].name] = tokens[i + 1]
        } else {
            if !desc.args[i].optional {
                err = errors.New("unsatisfied parameter requirements")
            }
            break
        }
    }

    return values, err
}

