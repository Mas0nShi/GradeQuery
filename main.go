package main

import (
	"encoding/json"
	"github.com/Mas0nShi/MHttp"
	"github.com/Mas0nShi/goConsole/console"
	"github.com/crufter/goquery"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CourseInfo struct {
	AcadYear     string `json:"学年"`
	Term         string `json:"学期"`
	CourseCode   string `json:"课程代码"`
	CourseName   string `json:"课程名称"`
	CourseNature string `json:"课程性质"`
	CourseAttr   string `json:"课程归属"`
	Credit       string `json:"学分"`
	GradePoint   string `json:"绩点"`
	Grade        string `json:"成绩"`
	MinorMark    string `json:"辅修标记"`
	RetestMark   string `json:"补考成绩"`
	RetakeGrades string `json:"重修成绩"`
	CollegeName  string `json:"学院名称"`
	Remarks      string `json:"备注"`
	ReworkMark   string `json:"重修标记"`
	CourseEnName string `json:"课程英文名称"`
}
type total []interface{}
type resMsg struct {
	Success int    `json:"success"`
	Error   string `json:"error"`
	Msg     total  `json:"msg"`
	Average string `json:"average"`
	Count   int    `json:"count"`
}

func throwErrorMsg(msg string) string {
	r, _ := json.Marshal(resMsg{Msg: total{}, Success: 0, Error: msg})
	return MHttp.Bytes2str(r)
}

func getTextMid(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)
	}
	str = MHttp.Bytes2str([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = MHttp.Bytes2str([]byte(str)[:m])
	return str
}
func parseCourseInfo(dom goquery.Nodes) string {
	nodes := dom.Find("#Datagrid1 tbody tr")
	text := dom.Find("#pjxfjd").Text()

	if nodes.Length() == 0 {
		return throwErrorMsg("Your session may expire.")
	}

	arr := strings.Split(text, "平均学分绩点：")
	if len(arr) == 0 {
		log.Fatalln("error in get #pjxfjd")
	}

	to2 := make(total, nodes.Length()-1)
	cont := 0
	nodes.Each(func(index int, element *goquery.Node) {
		if index > 0 {
			reg := regexp.MustCompile(`>(.+?)</td>`)
			mageData := reg.FindAllStringSubmatch(nodes.Eq(index).Html(), -1)

			if len(mageData) == 16 {
				to2[cont] = CourseInfo{
					AcadYear:     mageData[0][1],
					Term:         mageData[1][1],
					CourseCode:   mageData[2][1],
					CourseName:   mageData[3][1],
					CourseNature: mageData[4][1],
					CourseAttr:   strings.Trim(mageData[5][1], " "),
					Credit:       mageData[6][1],
					GradePoint:   strings.Trim(mageData[7][1], " "),
					Grade:        mageData[8][1],
					MinorMark:    mageData[9][1],
					RetestMark:   strings.Trim(mageData[10][1], " "),
					RetakeGrades: strings.Trim(mageData[11][1], " "),
					CollegeName:  mageData[12][1],
					Remarks:      strings.Trim(mageData[13][1], " "),
					ReworkMark:   mageData[14][1],
					CourseEnName: strings.Trim(mageData[15][1], " "),
				}
				cont++
			}
		}
	})
	comBytes, _ := json.Marshal(resMsg{Msg: to2, Count: cont, Success: 1, Error: "", Average: arr[1]})
	return MHttp.Bytes2str(comBytes)
}

// query
// types : 1-term / 2-years
// queryType : 1-CourseInfo
func query(host string, types int, queryType int, SessionId, stuID, name, queryID, queryName, acadYear, term string) string {
	var (
		http    = new(MHttp.MHttp)
		url     = ""
		ret     = ""
		data    = ""
		headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64",
		}
	)
	hostUrl := "http://" + host
	http.AutoHeaders(true)
	http.SetRequestHeaders(headers)
	http.SetCookie("ASP.NET_SessionId", SessionId)

	url = hostUrl + "/xscj_gc2.aspx?xh=" + stuID + "&xm=" + name + "&gnmkdm=N121611"
	http.Open("GET", url)
	http.Send(nil)
	ret = http.GetResponseText()
	__VIEWSTATE := url2.QueryEscape(getTextMid(ret, "__VIEWSTATE\" value=\"", "\" />"))
	__VIEWSTATEGENERATOR := getTextMid(ret, "__VIEWSTATEGENERATOR\" value=\"", "\" />")

	switch types {
	case 1:
		data = "__VIEWSTATE=" + __VIEWSTATE + "&__VIEWSTATEGENERATOR=" + __VIEWSTATEGENERATOR + "&ddlXN=" + acadYear + "&ddlXQ=" + term + "&Button1=%E6%8C%89%E5%AD%A6%E6%9C%9F%E6%9F%A5%E8%AF%A2"
	case 2:
		data = "__VIEWSTATE=" + __VIEWSTATE + "&__VIEWSTATEGENERATOR=" + __VIEWSTATEGENERATOR + "&ddlXN=" + acadYear + "&ddlXQ=" + term + "&Button5=%E6%8C%89%E5%AD%A6%E5%B9%B4%E6%9F%A5%E8%AF%A2"
	default:
		return throwErrorMsg("Param: type error.")
	}

	url = hostUrl + "/xscj_gc2.aspx?xh=" + queryID + "&xm=" + queryName + "&gnmkdm=N121611"
	http.Open("POST", url)
	http.Send(data)
	ret = http.GetResponseText()
	dom, err := goquery.ParseString(ret)
	if err != nil {
		log.Fatalln(err)
	}

	refdata := ""
	switch queryType {
	case 1:
		refdata = parseCourseInfo(dom)
	default:
		refdata = throwErrorMsg("Param: queryTypes error.")
	}
	return refdata
}

func getFormatTimeStr() string {
	t := time.Now()
	nanoT := strconv.FormatInt(t.UnixNano(), 10)
	return t.Format("2006-01-02 15:04:05") + "." + nanoT[10:13]
}
func logRequestInfo(r *http.Request, reslen int64) {
	format := r.Method + " - \"" + r.URL.String() + "\" " + strconv.FormatInt(reslen, 10) + " \"" + r.Header.Get("User-Agent")
	console.Info(format)
	fs, _ := os.OpenFile("web.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	fs.WriteString(getFormatTimeStr() + " - " + format + "\n")
	fs.Close()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	paser, _ := url2.ParseQuery(r.URL.RawQuery)
	qtt, _ := strconv.ParseInt(paser.Get("queryType"), 10, 64)
	tt, _ := strconv.ParseInt(paser.Get("type"), 10, 64)
	bokie := struct {
		host string // require

		session string // require
		user    string // require
		name    string // Optional

		queryId   string // require
		queryType int    // require
		types     int    // require
		queryName string // Optional

		acadYears string // require
		term      string // require
	}{
		host:    paser.Get("host"),
		session: paser.Get("session"),
		user:    paser.Get("user"),
		name:    paser.Get("name"),

		queryId:   paser.Get("queryId"),
		queryType: int(qtt),
		types:     int(tt),
		queryName: paser.Get("queryName"),

		acadYears: paser.Get("acadYears"),
		term:      paser.Get("term"),
	}
	var res string

	if bokie.host == "" || bokie.session == "" || bokie.user == "" || bokie.queryId == "" || bokie.queryType == 0 || bokie.types == 0 || bokie.acadYears == "" || bokie.term == "" {
		res = throwErrorMsg("Missing require params")
	} else {
		res = query(bokie.host, bokie.types, bokie.queryType, bokie.session, bokie.user, bokie.name, bokie.queryId, bokie.queryName, bokie.acadYears, bokie.term)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(MHttp.Str2bytes(res))
	logRequestInfo(r, int64(len(res)))

}

func main() {
	http.HandleFunc("/api/v1", IndexHandler)
	err := http.ListenAndServe(":13442", nil)
	if err != nil {
		panic("ERROR IN ListenAndServe")
	}

}
