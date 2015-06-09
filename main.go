// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fedorpatlin/sapccmsget/sapccms"
)

var flags struct {
	mtename *string
	host    *string
	sapsid  *string
	instNr  *string
	//	mode    *string // avg1, avg5, avg15
}

func prepareUsage() {
	hostname, _ := os.Hostname()
	flags.sapsid = flag.String("s", "", "SAP system name")
	flags.host = flag.String("h", hostname, "SAP host name")
	flags.instNr = flag.String("n", "00", "Instance number")
	//	flags.mode = flag.String("m", "avg1", "Value type of avg1, avg5, avg15" )
	flags.mtename = flag.String("e", "", "CCMS monitoring tree element full name")
}

func usage() {

}

func createMteRequest() *sapccms.MsgMtGetTidByNameRequest {
	req := sapccms.NewMsgMtGetTidByNameRequest()
	req.SoapRequest.Item = make([]sapccms.ALMTNAMEL, 1)
	if flags.mtename == nil {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	req.SoapRequest.Item[0].SetCompleteName(*flags.mtename)
	return req
}

func createMteResponse() *sapccms.MsgMtGetTidByNameResponse {
	res := sapccms.NewMsgMtGetTidByNameResponse()
	return res
}

func getMteByName(srv sapccms.SAPCCMS, name string) *sapccms.ALGTIDLNRC {
	var mteRequest = createMteRequest()
	var mteResponse = createMteResponse()

	if err := srv.MtGetTidByName(mteRequest, mteResponse); err != nil {
		log.Fatal(err.Error())
	}
	items := mteResponse.GetTidTable().Item
	if len(items) > 0 {
		if mteResponse.GetTidTable().GetItem()[0].Rc != "0" {

			log.Fatalf("MtGetTidByName: RC=%s for element %s\n", items[0].Rc, *flags.mtename)
			return nil

		} else {
			return &items[0]
		}
	} else {
		return nil
	}
}

func preparePerfRequest(tid sapccms.ALGLOBTID) *sapccms.MsgPerfReadRequest {
	msg := sapccms.NewMsgPerfReadRequest()
	msg.SoapRequest.Item = make([]sapccms.ALGLOBTID, 1)
	msg.SoapRequest.Item[0] = tid
	return msg
}

func getPerfByTid(srv sapccms.SAPCCMS, tid sapccms.ALGLOBTID) *sapccms.MsgPerfReadResponse {
	sendMsg := preparePerfRequest(tid)
	rcvMsg := sapccms.NewMsgPerfReadResponse()
	if err := srv.PerfRead(sendMsg, rcvMsg); err != nil {
		log.Fatal(err.Error())
	}
	return rcvMsg
}

func printToStdout(perf *sapccms.MsgPerfReadResponse) {
	items := perf.PerfReadResponse.TidTable.Item
	if len(items) > 0 {
		perfValue := items[0].PerfValue
		if perfValue.Avg01CountValue == "0" && perfValue.Avg01CountValue == "0" {
			fmt.Println(perfValue.AlertRelevantValue)
		} else {
			fmt.Println(perfValue.Avg01PerfValue)
		}
	} else {
		log.Fatalf("No performance values found for %s\n", *flags.mtename)
	}
}

func main() {
	prepareUsage()
	flag.Parse()
	srv := sapccms.NewSAPCCMS(fmt.Sprintf("http://%s:5%s13/SAPCCMS.cgi", *flags.host, *flags.instNr))
	tid := getMteByName(srv, *flags.mtename)
	if tid == nil {
		log.Fatalln("Element not found by name " + *flags.mtename)
	}
	perf := getPerfByTid(srv, tid.Tid)
	// output to zabbix-agent
	printToStdout(perf)
}
