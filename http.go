package main
import (
 "log"
 "net/http"
 "net"
 "fmt"
 "os"
 "time"
 "io/ioutil"
 "io"
 "strconv"
 "path/filepath"
 "flag"

)

var (
	Config      JSON
	Out_log     *os.File
	Config_time int64
)

type JSON struct {
    LocalPort  string `json: "localport"`
    HostName   string `json: "hostname"`
    Storage    string `json: "storage"`
    LogFile    string `json: "logfile"`
}

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func LogToFile(str string) {
	os.Remove(Config.LogFile)
	loger, _ := os.OpenFile(Config.LogFile, os.O_CREATE, 0777)
	loger.Write([]byte(str))
	loger.Close()
}

// Or if you want to do some more things like logging:
type logServer struct {
    hdl http.Handler
}

func (l *logServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if (r.URL.Path != "/favicon.ico") { log.Println(r.RemoteAddr + r.URL.Path) }
    l.hdl.ServeHTTP(w, r)
}
// /Or if you want to do some more things like logging:

func searchFiles(dir string) { // dir is the parent directory you what to search
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Fatal(err)
    }
    for _, file := range files {
	size := strconv.FormatInt(file.Size(),10)
        log.Println("find file: "+file.Name()+", size: "+size)
    }
}


func visit(path string, f os.FileInfo, err error) error {
  size := strconv.FormatInt(f.Size(),10)
  log.Println("find file: "+f.Name()+", size: "+size)
  return nil
} 


func LookFiles() {
  root := Config.Storage
  filepath.Walk(root, visit)
}


func main() {

	RELEASE_DATE := time.Now().Format("2006/01/02 15:04:05")
        fmt.Println("Stand-alone small http server from https://github.com/fly304625, compile at", RELEASE_DATE)
	Out_log, _ = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE, 0755)
	log.SetOutput(io.MultiWriter(Out_log, os.Stdout))
	
	err := ParseConf()
	if err != nil {
		log.Println("Error open config file\r\n")
		panic(err)
	}

        fmt.Println("run -todo=scan - for scan files")
	fmt.Println("WebPath:", "http://"+Config.HostName +":"+ Config.LocalPort)
	fmt.Println("Storage:", Config.Storage)
	fmt.Println("LogFile:", Config.LogFile)

 wordPtr := flag.String("todo", "web", "a string")
 flag.Parse()
 fmt.Println("get flag from user:", *wordPtr)

 if (*wordPtr == "web") {
  log.Println("Listening...")
  panic(http.ListenAndServe(":"+Config.LocalPort, &logServer{
        hdl: http.FileServer(http.Dir(Config.Storage)),
  }))
 }

 if (*wordPtr == "scan") { LookFiles() }

}

