package sapccms

import "testing"
import "fmt"

//import "log"

//import "reflect"
var host = "sap-dia1"
var domain = "miac.lan"
var sid = "ERP"
var s = NewSAPCCMS(fmt.Sprintf("http://%s.%s:50013/SAPCCMS.cgi", host, domain))

func getMonitor(hostname, sid, instanceNr, context, class string, elementpath ...string) string {
	name := sid + "\\"
	if context != "" {
		name = name + context
	}
	if (context != "") && (hostname != "") {
		name = name + "_"
	}
	if hostname != "" {
		name = name + hostname + "_" + sid + "_" + instanceNr
	}
	if class != "" {
		name = name + "\\" + class
	}
	if len(elementpath) > 0 {
		str := ""
		for _, v := range elementpath {
			str = str + "\\" + v
		}
		name = name + str
	}
	fmt.Println(name)
	return name
}

func Test_request(t *testing.T) {

	var req = NewMsgMtGetTidByNameRequest()
	var res = NewMsgMtGetTidByNameResponse()
	req.SoapRequest.SetItem(make([]ALMTNAMEL, 10))
	req.SoapRequest.Item[0].SetCompleteName(getMonitor(host, sid, "00", "", "R3Services", []string{"Dialog", "LongRunners"}...))

	if err := s.MtGetTidByName(req, res); err != nil {
		t.Fail()
	}
	if res == nil {
		t.Fail()
	}
	if res.MtGetTidByNameResponse.GetTidTable().GetItem()[0].GetRc() != "0" {
		t.Fail()
	}
}

func TestSequence(t *testing.T) {

	var req1 = NewMsgMtGetTidByNameRequest()
	var res1 = NewMsgMtGetTidByNameResponse()
	var req2 = NewMsgUtilMtReadAllRequest()
	var res2 = NewMsgUtilMtReadAllResponse()
	req1.SoapRequest.SetItem(make([]ALMTNAMEL, 10))
	req1.SoapRequest.Item[0].SetCompleteName(getMonitor(host, sid, "00", "", "R3Services", []string{"ICM", "General", "NoOfConnections"}...))

	if err := s.MtGetTidByName(req1, res1); err != nil {
		t.Fail()
	}
	req2.SetGetPerfData("1")
	req2.SoapRequest.SetItem(make([]ALTIDMAXLV, 1))
	req2.SoapRequest.Item[0].Tid = (res1.GetTidTable().Item[0].Tid)
	if err := s.UtilMtReadAll(req2, res2); err != nil {
		t.Fail()
	}
	if len(res2.GetPerfInfo().Item) > 0 {
		fmt.Println(res2.GetPerfInfo().Item[0].GetPerfValue().Avg01PerfValue)
		fmt.Println(res2.GetPerfInfo().Item[0].GetPerfValue().Avg01CountValue)
		fmt.Println(res2.GetPerfInfo().Item[0].GetPerfValue().Avg01SumValue)
	} else {
		t.Fail()
	}
}
