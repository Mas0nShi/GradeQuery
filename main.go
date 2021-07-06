package main

import (
	"encoding/json"
	"github.com/Mas0nShi/MHttp"
	"github.com/Mas0nShi/goConsole/console"
	"github.com/crufter/goquery"
	"log"
	url2 "net/url"
	"regexp"
	"strings"
)

type CourseInfo struct {
	AcadYear string `json:"学年"`
	Term string `json:"学期"`
	CourseCode string `json:"课程代码"`
	CourseName string `json:"课程名称"`
	CourseNature string `json:"课程性质"`
	CourseAttr string `json:"课程归属"`
	Credit string `json:"学分"`
	GradePoint string `json:"绩点"`
	Grade string `json:"成绩"`
	MinorMark string `json:"辅修标记"`
	RetestMark string `json:"补考成绩"`
	RetakeGrades string `json:"重修成绩"`
	CollegeName string `json:"学院名称"`
	Remarks string `json:"备注"`
	ReworkMark string `json:"重修标记"`
	CourseEnName string `json:"课程英文名称"`
}
type total []interface{}

func getTextMid(str, start, end string) string {

	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)  // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func query(types int, SessionId, stuID, name, queryID , queryName , acadYear , term string) string {
	var (
		http = new(MHttp.MHttp)
		url = ""
		ret = ""
		data = ""

		headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64",
		}
	)
	name = url2.QueryEscape(name)




	url = "http://jw.ypc.edu.cn/xs_main.aspx?xh=" + stuID
	http.Open("GET", url)
	http.SetCookie("ASP.NET_SessionId",SessionId)
	http.SetRequestHeaders(headers)

	http.Send(nil)

	url = "http://jw.ypc.edu.cn/xscj_gc2.aspx?xh=" + stuID + "&xm=" + name + "&gnmkdm=N121611"
	http.Open("GET", url)
	http.Send(nil)
	ret = http.GetResponseText()
	__VIEWSTATE := url2.QueryEscape(getTextMid(ret, "__VIEWSTATE\" value=\"","\" />"))
	__VIEWSTATEGENERATOR := getTextMid(ret, "__VIEWSTATEGENERATOR\" value=\"","\" />")



	switch types {
	case 1:
		data = "__VIEWSTATE=" + __VIEWSTATE +"&__VIEWSTATEGENERATOR=" + __VIEWSTATEGENERATOR + "&ddlXN=" + acadYear + "&ddlXQ="+ term +"&Button1=%E6%8C%89%E5%AD%A6%E6%9C%9F%E6%9F%A5%E8%AF%A2"
	case 2:
		data = "__VIEWSTATE=" + __VIEWSTATE +"&__VIEWSTATEGENERATOR=" + __VIEWSTATEGENERATOR + "&ddlXN=" + acadYear + "&ddlXQ="+ term +"&Button1=%E6%8C%89%E5%AD%A6%E6%9C%9F%E6%9F%A5%E8%AF%A2"
	default:
		panic("type error.")
	}

	url = "http://jw.ypc.edu.cn/xscj_gc2.aspx?xh=" + queryID + "&xm=" + queryName + "&gnmkdm=N121611"
	http.Open("POST", url)
	http.Send(data)

	ret = http.GetResponseText()
	dom ,err := goquery.ParseString(ret)
	if err != nil {
		log.Fatalln(err)
	}
	nodes := dom.Find("#Datagrid1 tbody tr")

	to2 := make(total, nodes.Length()-1)
	cont := 0
	nodes.Each(func(index int, element *goquery.Node) {
		if index > 0 {
			reg := regexp.MustCompile(`>(.+?)</td>`)
			mageData := reg.FindAllStringSubmatch(nodes.Eq(index).Html(), -1)
			if len(mageData) == 16 {
				stuctC := CourseInfo{
					AcadYear: mageData[0][1],
					Term: mageData[1][1],
					CourseCode: mageData[2][1],
					CourseName: mageData[3][1],
					CourseNature: mageData[4][1],
					CourseAttr: strings.Trim(mageData[5][1], " "),
					Credit: mageData[6][1],
					GradePoint: strings.Trim(mageData[7][1], " "),
					Grade: mageData[8][1],
					MinorMark: mageData[9][1],
					RetestMark: strings.Trim(mageData[10][1], " "),
					RetakeGrades: strings.Trim(mageData[11][1], " "),
					CollegeName: mageData[12][1],
					Remarks: strings.Trim(mageData[13][1], " "),
					ReworkMark: mageData[14][1],
					CourseEnName: strings.Trim(mageData[15][1], " "),
				}
				to2[cont] = stuctC
				cont++
			}
		}
	})

	comBytes, _ := json.Marshal(to2)
	return MHttp.Bytes2str(comBytes)
}

func main() {
	sessionId := ""

	data := query(1,sessionId, "0417200322", "施滢琦","0417200322","施滢琦","2020-2021","2")
	console.Log(data)
}
