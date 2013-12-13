package main

import (
  "fmt"
  "flag"
  "bufio"
  "os"
  "log"
  "io/ioutil"
  "os/user"
  "strings"
  "strconv"
  "net/http"
  "net/url"
)

type set map[string]struct{}

func (s set) add(x string) set {
  if s == nil {
    s = make(set)
  }
  s[x] = struct{}{}
  return s
}

func (s set) get() []string {
  x:=[]string{}
  for i := range s{
    x = append(x,i)
  }
  return x
}

type pt struct {
  sent float64
  recv float64
  users set
}

var data map[string]pt

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
  flag.BoolVar(&pprint, "pp", false, "Pretty Print")
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
  ud, err := user.LookupId(processdata[len(processdata)-1])
  if err != nil {
    ud = &user.User { "-1", "-1", "dontknow", "dontknow", "/home/dontknow" }
  }
  data[pname] = pt{
    data[pname].sent + sent,
    data[pname].recv + recv,
    data[pname].users.add(ud.Username),
  }
}

func prettyprint() {
  for proc, _ := range data {
    fmt.Printf("%40s\t%10.2f\t%10.2f\t%40s\n", proc, data[proc].sent, data[proc].recv,
      strings.Join(data[proc].users.get(), ",") )
  }
}

func getcsv() string{
  csv := ""
  hostname,_ := os.Hostname()
  for proc, _ := range data {
    csv = csv + fmt.Sprintf("%s,%.2f,%.2f,%s,%s,%s\n",
      proc, data[proc].sent, data[proc].recv,
      strings.Join(data[proc].users.get(), " "),
      filename,hostname)
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
