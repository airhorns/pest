package main

import (
  "bufio"
  "fmt"
  "flag"
  "github.com/hornairs/go-xmpp"
  "github.com/sloonz/go-iconv"
  "log"
  "os"
  "strings"
)

var server   = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password = flag.String("password", "", "password")

func fromUTF8(s string) string {
  ic, err := iconv.Open("char", "UTF-8")
  if err != nil {
    return s
  }
  defer ic.Close()
  ret, _ := ic.Conv(s)
  return ret
}

func toUTF8(s string) string {
  ic, err := iconv.Open("UTF-8", "char")
  if err != nil {
    return s
  }
  defer ic.Close()
  ret, _ := ic.Conv(s)
  return ret
}

func main() {
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "usage: example [options]\n")
    flag.PrintDefaults()
    os.Exit(2)
  }
  flag.Parse()
  if *username == "" || *password == "" {
    flag.Usage()
  }

  talk, err := xmpp.NewClient(*server, *username, *password)
  if err != nil {
    log.Fatal(err)
  }
  log.Print("Booted")
  log.Print(talk.Roster)

  go func() {
    for {
      chat, err := talk.Recv()
      if err != nil {
        log.Fatal(err)
      }
      fmt.Println(chat.Remote, fromUTF8(chat.Text))
    }
  }()
  for {
    in := bufio.NewReader(os.Stdin)
    line, err := in.ReadString('\n')
    if err != nil {
      continue
    }
    line = strings.TrimRight(line, "\n")

    tokens := strings.SplitN(line, " ", 2)
    if len(tokens) == 2 {
      talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: toUTF8(tokens[1])})
    }
  }
}
