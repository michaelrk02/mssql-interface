package main

import (
    "container/list"
    "bufio"
    "fmt"
    "io"
    "os"
    "runtime"
    "strings"
    "unicode"

    "golang.org/x/term"
)

func textInput(eof bool) string {
    var sb strings.Builder

    rd := bufio.NewReader(os.Stdin)
    for true {
        var ch byte
        var err error

        ch, err = rd.ReadByte()
        if err == io.EOF {
            return sb.String()
        }

        if !eof {
            if (ch == '\n') || (ch == '\r') {
                if runtime.GOOS == "windows" {
                    rd.Discard(1)
                }
                return sb.String()
            }
        }

        sb.WriteByte(ch)
    }

    return ""
}

func wrap(str string, width, index int) string {
    start := index * width;
    if start < len(str) {
        end := (index + 1) * width;
        if end > len(str) {
            end = len(str)
        }
        return str[start:end]
    }
    return ""
}

func ellipsis(str string, maxLength int) string {
    end := maxLength
    if end > len(str) {
        end = len(str)
        return fmt.Sprintf("%s...", str[0:end])
    }
    return str
}

func parse(txt string) []string {
    tokensLs := list.New()

    quoted := false
    quotedEsc := false
    var tokenSb strings.Builder
    for i := range txt {
        ch := txt[i]
        push := false
        insert := false

        if quoted {
            if quotedEsc {
                if (ch == '"') || (ch == '\\') {
                    push = true
                    quotedEsc = false
                }
            } else {
                if ch == '"' {
                    insert = true
                    quoted = false
                } else if ch == '\\' {
                    quotedEsc = true
                } else {
                    push = true
                }
            }
        } else {
            if unicode.IsSpace(rune(ch)) {
                insert = true
            } else if ch == '"' {
                quoted = true
                insert = true
            } else {
                push = true
                if i == len(txt) - 1 {
                    insert = true
                }
            }
        }

        if push {
            tokenSb.WriteByte(ch)
        }
        if insert {
            if tokenSb.Len() > 0 {
                tokensLs.PushBack(tokenSb.String())
                tokenSb.Reset()
            }
        }
    }

    tokens := make([]string, tokensLs.Len())
    for i := 0; tokensLs.Front() != nil; i++ {
        tokens[i] = tokensLs.Front().Value.(string)
        tokensLs.Remove(tokensLs.Front())
    }

    return tokens
}

func outputWidth() int {
    var w int
    w, _, _ = term.GetSize(int(os.Stdout.Fd()))
    return w
}

