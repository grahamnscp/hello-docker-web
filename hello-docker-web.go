package main

import (
  "fmt"
  "net/http"
  "os"
  "time"
)

var Hostname, oserr = os.Hostname()

func main() {
  if oserr != nil {
    fmt.Println("io.hostname error: ", oserr, " hostname value:", Hostname)
  }

  http.HandleFunc("/", index)
  http.HandleFunc("/hello", helloHandler)
  http.HandleFunc("/hello/", helloHandler)
  http.HandleFunc("/docker", dockerHandler)
  http.HandleFunc("/docker/", dockerHandler)
  http.HandleFunc("/whale", whaleHandler)
  http.HandleFunc("/whale/", whaleHandler)
  err := http.ListenAndServe(":8080", nil)

  if err != nil {
    fmt.Println("Serve Http:", err)
  }
}

func index(w http.ResponseWriter, r *http.Request) {

  fmt.Fprint(w,"\n")
  fmt.Fprint(w,"                              ##\n")
  fmt.Fprint(w,"                        ## ## ##        ==\n")
  fmt.Fprint(w,"                     ## ## ## ## ##    ===\n")
  fmt.Fprint(w,"                 /`````````````````\\___/ ===\n")
  fmt.Fprint(w,"            ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~~/~ === ~~~\n")
  fmt.Fprint(w,"                 \\______ o           __/\n")
  fmt.Fprint(w,"                   \\    \\         __/\n")
  fmt.Fprint(w,"                    \\____\\_______/\n")
  fmt.Fprint(w," _           _    _                _            _\n")
  fmt.Fprint(w,"| |     ___ | |  | |    ___     __| | ___   ___| | _____ _ __\n")
  fmt.Fprint(w,"| |___ / _ \\| |  | |   / _ \\   / _  |/ _ \\ / __| |/ / _ \\ '__|\n")
  fmt.Fprint(w,"|  _  |  __/| |__| |__| (_) | | (_| | (_) | (__|   <  __/ |\n")
  fmt.Fprint(w,"|_| |_|\\___/ \\___|\\___|\\___/   \\____|\\___/ \\___|_|\\_\\___|_|\n")
  fmt.Fprint(w,"\n")
  dt := time.Now()
  fmt.Fprint(w,"[", dt.Format("01-02-2006 15:04:05.00"), "] ", "Container hostname: ", Hostname, "\n")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Fprint(w,"\n")
  fmt.Fprint(w," _           _    _\n")
  fmt.Fprint(w,"| |     ___ | |  | |    ___\n")
  fmt.Fprint(w,"| |___ / _ \\| |  | |   / _ \\\n")
  fmt.Fprint(w,"|  _  |  __/| |__| |__| (_) |\n")
  fmt.Fprint(w,"|_| |_|\\___/ \\___|\\___|\\___/\n")
  fmt.Fprint(w,"\n")
  dt := time.Now()
  fmt.Fprint(w,"[", dt.Format("01-02-2006 15:04:05.00"), "] ", "Container hostname: ", Hostname, "\n")
}

func dockerHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Fprint(w,"\n")
  fmt.Fprint(w,"     _            _\n")
  fmt.Fprint(w,"  __| | ___   ___| | _____ _ __\n")
  fmt.Fprint(w," / _  |/ _ \\ / __| |/ / _ \\ '__|\n")
  fmt.Fprint(w,"| (_| | (_) | (__|   <  __/ |\n")
  fmt.Fprint(w," \\____|\\___/ \\___|_|\\_\\___|_|\n")
  fmt.Fprint(w,"\n")
  dt := time.Now()
  fmt.Fprint(w,"[", dt.Format("01-02-2006 15:04:05.00"), "] ", "Container hostname: ", Hostname, "\n")
}

func whaleHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Fprint(w,"\n")
  fmt.Fprint(w,"                   ##\n")
  fmt.Fprint(w,"             ## ## ##        ==\n")
  fmt.Fprint(w,"          ## ## ## ## ##    ===\n")
  fmt.Fprint(w,"      /`````````````````\\___/ ===\n")
  fmt.Fprint(w," ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~~/~ === ~~~\n")
  fmt.Fprint(w,"      \\______ o           __/\n")
  fmt.Fprint(w,"        \\    \\         __/\n")
  fmt.Fprint(w,"         \\____\\_______/\n")
  fmt.Fprint(w,"\n")
  dt := time.Now()
  fmt.Fprint(w,"[", dt.Format("01-02-2006 15:04:05.00"), "] ", "Container hostname: ", Hostname, "\n")
}
