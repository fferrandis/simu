package hdio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/fferrandis/simu/scality/hyperdrive/config"
	"github.com/fferrandis/simu/scality/hyperdrive/hdcluster"
)

type HDIoChanMsg struct {
	datalen uint64
}

type HDIoChanRsp struct {
	ret  bool
	load uint64
	code int
}

type HDIoSrv struct {
	nr      int
	curr    int
	closed  bool
	msg     []chan HDIoChanMsg
	ret     []chan HDIoChanRsp
	lock    sync.Mutex
	datalen uint64
	load    uint64
}

var hdio HDIoSrv

func getnextchan() int {
	i := 0
	hdio.lock.Lock()
	{
		if hdio.curr >= hdio.nr {
			hdio.curr = 0
		}
		i = hdio.curr
		hdio.curr++
	}
	hdio.lock.Unlock()
	return i
}

func worker(id int) {
	for {
		nextjob, ok := <-hdio.msg[id]

		if ok == false {
			fmt.Println("killing worker ", id, ok)
			return
		}

		fmt.Println("worker ", id, "inject ", nextjob, "bytes")
		ret, load := hdcluster.HDClusterSrvPut(nextjob.datalen)
		if ret != true {
			hdio.ret[id] <- HDIoChanRsp{ret, load, 500}
		} else {
			hdio.ret[id] <- HDIoChanRsp{ret, load, 200}
		}
	}
}

type Answer struct {
	Scal_response_time    uint64
	Scal_accumulated_len  uint64
	Scal_accumulated_load uint64
}

func PutData(datalen uint64, resp http.ResponseWriter) (bool, string) {
	if hdio.closed == true {
		resp.WriteHeader(503)
		return false, "already closed"
	}

	if datalen != 0 {
		id := getnextchan()
		hdio.msg[id] <- HDIoChanMsg{datalen}

		data := <-hdio.ret[id]
		hdio.datalen += datalen
		hdio.load += data.load
		tp := Answer{data.load, hdio.datalen, hdio.load}
		bodyStr, err := json.Marshal(tp)
		resp.WriteHeader(data.code)
		if err == nil {
			resp.Write([]byte(bodyStr))
		}

	} else {
		for i := 0; i < hdio.nr; i++ {
			close(hdio.msg[i])
		}
		hdio.closed = true
		resp.WriteHeader(200)
	}
	return true, "ok"
}

func createworkerpool(nr int) {
	hdio.msg = make([]chan HDIoChanMsg, nr)
	hdio.ret = make([]chan HDIoChanRsp, nr)
	hdio.nr = nr

	for i := 0; i < nr; i++ {
		hdio.msg[i] = make(chan HDIoChanMsg, 1)
		hdio.ret[i] = make(chan HDIoChanRsp, 1)
		go worker(i)
	}
}

func root(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(501)
	resp.Write([]byte("realtime statistics not available yet"))
}

func addsrv(resp http.ResponseWriter, req *http.Request) {
	//var nrdisk int
	//var capacity uint64

	//for hdr, value := range req.Header {
	//	fmt.Println(hdr)
	//	switch hdr {
	//	case "X-Sim-Nrdisk":
	//		nrdisk, _ = strconv.Atoi(value[0])
	//	case "X-Sim-Diskcapa":
	//		fmt.Println(value)
	//		capacity, _ = strconv.ParseUint(value[0], 10, 64)
	//	}
	//}
	//HDClusterSrvAdd(nrdisk, capacity)
	resp.WriteHeader(http.StatusOK)
}

func delsrv(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(501)
}

func onput(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		fmt.Println("bad method : ", req.Method)
		return
	}
	m, _ := url.ParseQuery(req.URL.RawQuery)
	val, found := m["datalen"]
	if !found {
		fmt.Println("datalen not specified ")
		return
	}
	datalen, err2 := strconv.ParseUint(val[0], 10, 64)
	if err2 != nil {
		fmt.Println("cannot parse datalen", val, err2)
		return
	}
	PutData(datalen, resp)
}

var handler_map = map[string]func(http.ResponseWriter, *http.Request){
	"/":        root,
	"/srv/add": addsrv,
	"/srv/del": delsrv,
	"/put":     onput,
}

func HDIoStart() {
	createworkerpool(10)
	for u, handler := range handler_map {
		http.HandleFunc(u, handler)
	}
	hdcluster.Init(config.HDCFG)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
