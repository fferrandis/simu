package hdio

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	hdcfg "scality/hyperdrive/config"
	. "scality/hyperdrive/hdcluster"
	"strconv"
	"sync"
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
	nr     int
	curr   int
	closed bool
	msg    []chan HDIoChanMsg
	ret    []chan HDIoChanRsp
	lock   sync.Mutex
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
		ret, load := HDClusterSrvPut(nextjob.datalen)
		if ret != true {
			hdio.ret[id] <- HDIoChanRsp{ret, load, 500}
		} else {
			hdio.ret[id] <- HDIoChanRsp{ret, load, 200}
		}
	}
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
		loadstr := strconv.FormatUint(data.load, 10)
		resp.Header().Add("X-IO-Load", loadstr)
		body_str := "{\"scal-response-time\" : " + loadstr + "}\n"
		resp.Header().Add("Content-Length", strconv.Itoa(len(body_str)))
		resp.WriteHeader(data.code)
		io.WriteString(resp, body_str)

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
	io.WriteString(resp, "realtime statistics not available yet")
	resp.WriteHeader(501)
}

func addsrv(resp http.ResponseWriter, req *http.Request) {
	var nrdisk int
	var capacity uint64

	for hdr, value := range req.Header {
		fmt.Println(hdr)
		switch hdr {
		case "X-Sim-Nrdisk":
			nrdisk, _ = strconv.Atoi(value[0])
		case "X-Sim-Diskcapa":
			fmt.Println(value)
			capacity, _ = strconv.ParseUint(value[0], 10, 64)
		}
	}
	HDClusterSrvAdd(nrdisk, capacity)
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
	if found == false {
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
	for _, srv := range hdcfg.HDCFG.Hdservers {
		HDClusterSrvAdd(srv.Nr_disk, srv.Capacity)
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}
