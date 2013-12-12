package main

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "log"
  "strings"
  "strconv"
)

type pt struct {
  sent float64
  recv float64
}

var data map[string]pt

func adddata(s string, p pt) {
  if _,ok := data[s]; ok {
    p2 := pt{data[s].sent + p.sent, data[s].recv+p.recv}
    data[s] = p2
  } else {
    data[s] = pt{0.0,0.0}
  }
}

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  data = map[string]pt {}
  flag.Parse()
  filename := flag.Args()[0]
  file, err := os.Open(filename)
  check(err)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    parseline(scanner.Text())
  }
  if err = scanner.Err(); err != nil {
    log.Fatal(err)
  }
  prettyprint()
}

func parseline(line string) {
  l := strings.Fields(line)
  if len(l) < 3 {
    return
  }
  recv, err := strconv.ParseFloat(l[len(l)-1],64)
  if err != nil { return }
  sent, err := strconv.ParseFloat(l[len(l)-2],64)
  if err != nil { return }
  processCol := strings.Join(l[0:len(l)-2],"_")
  processdata := strings.Split(processCol, "/")
  if len(processdata) < 3 { return }
  pname := strings.Join(processdata[0:len(processdata)-2], "/")
  pname = strings.Split(pname, ":")[0]
  adddata(pname, pt{ sent, recv })
}


func prettyprint() {
  for proc, _ := range data {
    fmt.Printf("%s\t%f\t%f\n", proc, data[proc].sent, data[proc].recv)
  }
}
