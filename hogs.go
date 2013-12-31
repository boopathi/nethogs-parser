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
  "time"
  "runtime"
  "runtime/pprof"
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

type DATA struct {
  val map[string]pt
  filename string
}

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

//flags
var (
  datatable string
  datatableclass string
  ptype string
  cpuprofile string
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  start := time.Now()
  flag.StringVar(&datatable, "datatable", "", "Datatable server details")
  flag.StringVar(&datatableclass, "class", "nethogs", "Classname for Datatable")
  flag.StringVar(&ptype, "type", "", "How to print the data to STDOUT")
  flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU Profile to file")
  flag.Parse()
  if cpuprofile != "" {
    f, err := os.Create(cpuprofile)
    if err != nil { log.Fatal(err) }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }
  collection := make([]DATA, 1)
  // A channel for sending and receiving values
  dchan := make(chan DATA, flag.NArg()+1)
  for i:=0; i<flag.NArg(); i++ {
    filename := flag.Args()[i]
    go parsefile(filename,dchan)
  }
  for i:=0; i<flag.NArg(); i++ {
    data := <-dchan
    collection = append(collection, data)
    if ptype == "pretty" {
      data.prettyprint()
    } else if ptype == "csv" {
      fmt.Print(data.getcsv())
    } else {
      log.Printf("[DONE] [%s] %s", time.Since(start), data.filename)
    }
  }
  if datatable != "" {
    send_to_datatable(collection)
  }
}

// This will be called from a go routine
func parsefile(filename string, d chan DATA) {
  data := DATA{ map[string]pt{}, filename }
  file, err := os.Open(filename)
  check(err)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    data.parseline(scanner.Text())
  }
  if err = scanner.Err(); err != nil {
    log.Print(err)
    d <- data
  }
  d <- data
}

func (d DATA) parseline(line string) {
  data := d.val
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

func (d DATA) prettyprint() {
  fmt.Printf("Output for file = %s\n\n", d.filename)
  data := d.val
  for proc, _ := range data {
    fmt.Printf("%40s\t%10.2f\t%10.2f\t%40s\n", proc, data[proc].sent, data[proc].recv,
      strings.Join(data[proc].users.get(), ",") )
  }
  fmt.Printf("\n\n")
}

func (d DATA) getcsv() string{
  data := d.val
  csv := ""
  hostname,_ := os.Hostname()
  for proc, _ := range data {
    csv = csv + fmt.Sprintf("%s,%.2f,%.2f,%s,%s,%s\n",
      proc, data[proc].sent, data[proc].recv,
      strings.Join(data[proc].users.get(), " "),
      d.filename,hostname)
  }
  return csv
}

func send_to_datatable(collection []DATA) {
  hostname, _ := os.Hostname()
  params := url.Values{}
  params.Set("class", datatableclass)
  params.Set("host", hostname)
  csvdata := ""
  for i := range collection {
    csvdata = csvdata + collection[i].getcsv()
  }
  params.Set("data", csvdata)
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
