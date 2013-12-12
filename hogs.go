package main

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "log"
  "io/ioutil"
  "strings"
  "strconv"
  "net/http"
  "net/url"
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

//flags
var (
  filename,datatable string
  pprint bool
)

func main() {
  data = map[string]pt {}
  flag.StringVar(&datatable, "datatable", "", "Datatable server details")
  flag.BoolVar(&pprint, "prettyprint", false, "Pretty Print")
  flag.Parse()
  filename = flag.Args()[0]
  file, err := os.Open(filename)
  check(err)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    parseline(scanner.Text())
  }
  if err = scanner.Err(); err != nil {
    log.Fatal(err)
  }
  if pprint {
    prettyprint()
  }
  if datatable != "" {
    send_to_datatable()
  }
}

func parseline(line string) {
  l := strings.Fields(line)
  if len(l) < 3 { return }
  recv, err := strconv.ParseFloat(l[len(l)-1],64)
  if err != nil { return }
  sent, err := strconv.ParseFloat(l[len(l)-2],64)
  if err != nil { return }
  processCol := strings.Join(l[0:len(l)-2],"_")
  processdata := strings.Split(processCol, "/")
  if len(processdata) < 3 { return }
  pname := strings.Join(processdata[0:len(processdata)-2], "/")
  if strings.Index(pname, ":") != -1 && strings.Index(pname, "-") != -1 {
    pname = strings.Split(pname, "-")[0]
  }
  adddata(pname, pt{ sent, recv })
}

func prettyprint() {
  for proc, _ := range data {
    fmt.Printf("%40s\t%10.2f\t%10.2f\n", proc, data[proc].sent, data[proc].recv)
  }
}

func getcsv() string{
  csv := ""
  hostname,_ := os.Hostname()
  for proc, _ := range data {
    csv = csv + fmt.Sprintf("%s,%.2f,%.2f,%s,%s\n",
      proc, data[proc].sent, data[proc].recv,filename,hostname)
  }
  return csv
}

func send_to_datatable() {
  classname := "nethogs"
  hostname, _ := os.Hostname()
  params := url.Values{}
  params.Set("class", classname)
  params.Set("host", hostname)
  params.Set("data", getcsv())
  if strings.Index(datatable, "http://") != 0 {
    datatable = "http://" + datatable
  }
  resp, err := http.PostForm(datatable + "/api/put", params)
  check(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  check(err)
  fmt.Println(string(body))
}
