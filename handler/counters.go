package handler

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"../models"
	"github.com/gin-gonic/gin"
)

type Bets struct {
	Text             struct{}     `json:"text"`
	Bet              Bettype      `json:"bet"`
	BallNum          interface{}  `json:"ballNum"`
	MobileBetText    Bettype      `json:"mobileBetText"`
	BetType          string       `json:"betType"`
	Selection        string       `json:"selection"`
	CounterID        int64        `json:"counterId"`
	CounterName      string       `json:"counterName"`
	DrawNo           int64        `json:"drawNo"`
	Selections       []Selections `json:"selections"`
	Stake            interface{}  `json:"stake"`
	IsShowChipBox    bool         `json:"isShowChipBox"`
	Uivalid          bool         `json:"uivalid"`
	EstWinning       int          `json:"estWinning"`
	DetectSoftKeyPad struct{}     `json:"detectSoftKeyPad"`
	HideChipBox      interface{}  `json:"hideChipBox"`
	ShowValid        bool         `form:"showValid" json:"showValid"`
}

type Bettype struct {
	BetType   string      `json:"betType"`
	Selection interface{} `json:"selection"`
}

type Selections struct {
	ID     int     `json:"id"`
	Odds   float64 `json:"odds"`
	MinBet float64 `json:"minBet"`
	MaxBet float64 `json:"maxBet"`
}

//下注
func PlaceBet(c *gin.Context) {
	var bet []Bets
	// if err := c.ShouldBindJSON(&bet); err == nil {
	// 	fmt.Println(&bet)
	// } else {
	// 	fmt.Println(err.Error())
	// }

	if err := c.BindJSON(&bet); err == nil {
		fmt.Println(&bet)
	} else {
		fmt.Println(err.Error())
	}
}

//单个数据抓取
func CounterID(c *gin.Context) {

	//fmt.Println(c.Param("number"))

	formate := "2006-01-02T15:04:05+08:00"
	cid, err := strconv.ParseInt(c.Param("number"), 10, 64)

	if err != nil {
		fmt.Println("cid is :", err.Error())
	}

	//标注上一期
	var drawno int64

	b := bytes.Buffer{}
	b.WriteString(`{"isSuccess": true,`)
	b.WriteString(`"data":{ `)
	if v, err := models.GetCounter(cid); err == nil {

		b.WriteString(`"id":` + strconv.FormatInt(v.Id, 10) + `,`)
		b.WriteString(`"name": "` + v.Name + `",`)
		b.WriteString(`"official": "` + v.Official + `",`)
		b.WriteString(`"status": ` + strconv.Itoa(v.Status) + `,`)

		// //draw
		if info, err := models.GetDraw(time.Now(), v.Id); info != nil && err == nil {

			drawno = info.Drawno - 1

			b.WriteString(`"draw": {`)
			b.WriteString(`"drawNo": ` + strconv.FormatInt(info.Drawno, 10) + `,`)
			b.WriteString(`"drawStatus": ` + strconv.Itoa(info.Drawstatus) + `,`)
			//b.WriteString(`"status": ` + strconv.Itoa(v.Status) + `,`)
			b.WriteString(`"startTime": "` + info.Starttime.Format(formate) + `",`)
			b.WriteString(`"endTime": "` + info.Endtime.Format(formate) + `",`)
			b.WriteString(`"betClosedMMSS": "` + info.Betclosedmmss + `",`)
			b.WriteString(`"isCloseManually": ` + strconv.FormatBool(info.Isclosemanually) + ``)
			b.WriteString(`},`)
		}
		// //draw

		//selections
		b.WriteString(`"selections": {`)
		selections, _ := models.GetSelections(v.Id)
		for x, s := range selections {

			b.WriteString(`"` + s.Name + `": {`)
			b.WriteString(`"id": ` + strconv.FormatInt(s.Id, 10) + `,`)
			b.WriteString(`"odds": ` + Float64toStr(s.Odds) + `,`)
			b.WriteString(`"minBet": ` + Float64toStr(s.Minbet) + `,`)
			b.WriteString(`"maxBet": ` + Float64toStr(s.Maxbet) + ``)

			//fmt.Println("#######################", x, len(selections))
			if x >= len(selections)-1 {
				b.WriteString(`}`)
			} else {
				b.WriteString(`},`)
			}

		}
		b.WriteString(`},`)
		//selections

		//gameResult
		b.WriteString(`"gameResult": {`)

		if info, err := models.GetDrawno(drawno); err == nil {
			b.WriteString(`"counterId": ` + strconv.FormatInt(v.Id, 10) + `,`)
			b.WriteString(`"drawNo": ` + strconv.FormatInt(v.Id, 10) + `,`)
			b.WriteString(`"drawTime":  "` + info.Starttime.Format(formate) + `",`)
			b.WriteString(`"drawStatus": ` + strconv.Itoa(info.Drawstatus) + `,`)
			b.WriteString(`"voidReason": ` + strconv.Itoa(info.Voidreason) + `,`)
			b.WriteString(`"resultBalls": ` + info.Resultballs + ``)
		}
		b.WriteString(`},`)
		//end gameResult

		b.WriteString(`"seq": ` + strconv.Itoa(v.Seq) + `,`)
		b.WriteString(`"ballCount": ` + strconv.Itoa(v.Ballcount) + `,`)
		b.WriteString(`"resultOpenIntervalSecond": ` + strconv.Itoa(v.Resultopenintervalsecond) + `,`)
		b.WriteString(`"resultWaitingIntervalSecond": ` + strconv.Itoa(v.Resultwaitingintervalsecond) + ``)

		b.WriteString("}")
	}

	b.WriteString("}")

	c.Request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	c.Writer.WriteString(b.String())

}

func Mobile(c *gin.Context) {

	formate := "2006-01-02T15:04:05+08:00"
	//now := time.Now().Format(formate)

	b := bytes.Buffer{}
	b.WriteString(`{"isSuccess": true,`)
	b.WriteString(`"data":{ `)
	b.WriteString(`"announcements": [],`)

	//openBets
	b.WriteString(`"openBets": {`)
	b.WriteString(`"totalCount": 0,`)
	b.WriteString(`"totalStake": 0.0,`)
	b.WriteString(`"totalReturnAmount": 0.0,`)
	b.WriteString(`"wagers": []`)
	b.WriteString(`},`)
	//openBets

	//counters
	b.WriteString(`"counters": [`)

	//标注上一期
	var drawno int64
	now := time.Now()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	lists, _ := models.GetCounters("seq")
	for i, v := range lists {

		b.WriteString(`{`)
		b.WriteString(`"id":` + strconv.FormatInt(v.Id, 10) + `,`)
		b.WriteString(`"name": "` + v.Name + `",`)
		b.WriteString(`"official": "` + v.Official + `",`)
		b.WriteString(`"status": ` + strconv.Itoa(v.Status) + `,`)

		// //draw
		if info, err := models.GetDraw(time.Now(), v.Id); info != nil && err == nil {

			drawno = info.Drawno - 1

			b.WriteString(`"draw": {`)
			b.WriteString(`"drawNo": ` + strconv.FormatInt(info.Drawno, 10) + `,`)
			b.WriteString(`"drawStatus": ` + strconv.Itoa(info.Drawstatus) + `,`)
			//b.WriteString(`"status": ` + strconv.Itoa(v.Status) + `,`)
			b.WriteString(`"startTime": "` + info.Starttime.Format(formate) + `",`)
			b.WriteString(`"endTime": "` + info.Endtime.Format(formate) + `",`)
			b.WriteString(`"betClosedMMSS": "` + info.Betclosedmmss + `",`)
			b.WriteString(`"isCloseManually": ` + strconv.FormatBool(info.Isclosemanually) + ``)
			b.WriteString(`},`)
		}
		// //draw

		//selections
		b.WriteString(`"selections": {`)
		selections, _ := models.GetSelections(v.Id)
		for x, s := range selections {

			b.WriteString(`"` + s.Name + `": {`)
			b.WriteString(`"id": ` + strconv.FormatInt(s.Id, 10) + `,`)
			b.WriteString(`"odds": ` + Float64toStr(s.Odds) + `,`)
			b.WriteString(`"minBet": ` + Float64toStr(s.Minbet) + `,`)
			b.WriteString(`"maxBet": ` + Float64toStr(s.Maxbet) + ``)

			//fmt.Println("#######################", x, len(selections))
			if x >= len(selections)-1 {
				b.WriteString(`}`)
			} else {
				b.WriteString(`},`)
			}

		}
		b.WriteString(`},`)
		//selections

		//gameResult
		b.WriteString(`"gameResult": {`)

		if info, err := models.GetDrawno(drawno); err == nil {
			b.WriteString(`"counterId": ` + strconv.FormatInt(v.Id, 10) + `,`)
			b.WriteString(`"drawNo": ` + strconv.FormatInt(v.Id, 10) + `,`)
			b.WriteString(`"drawTime":  "` + info.Starttime.Format(formate) + `",`)
			b.WriteString(`"drawStatus": ` + strconv.Itoa(info.Drawstatus) + `,`)
			b.WriteString(`"voidReason": ` + strconv.Itoa(info.Voidreason) + `,`)
			b.WriteString(`"resultBalls": ` + info.Resultballs + ``)
		}
		b.WriteString(`},`)
		//end gameResult

		b.WriteString(`"seq": ` + strconv.Itoa(v.Seq) + `,`)
		b.WriteString(`"ballCount": ` + strconv.Itoa(v.Ballcount) + `,`)
		b.WriteString(`"resultOpenIntervalSecond": ` + strconv.Itoa(v.Resultopenintervalsecond) + `,`)
		b.WriteString(`"resultWaitingIntervalSecond": ` + strconv.Itoa(v.Resultwaitingintervalsecond) + ``)

		if i >= len(lists)-1 {
			b.WriteString(`}`)
		} else {
			b.WriteString(`},`)
		}

	}

	b.WriteString(`],`)
	//end counters

	//trendsList
	b.WriteString(`"trendsList": [`)
	for i, v := range lists {
		b.WriteString(`{`)
		b.WriteString(`"counterId": ` + strconv.FormatInt(v.Id, 10) + `,`)

		b.WriteString(`"trends": [`)

		if dlists, err := models.GetDraws(time.Now(), beginTime, v.Id); err == nil {
			for i, v := range dlists {

				strNo := strconv.FormatInt(v.Drawno, 10)
				str := strNo[8:len(strNo)]

				intNo, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("strconv Atoi error")
				}

				b.WriteString("{")
				b.WriteString(`"counterId": 0,`)
				b.WriteString(`"drawNo": ` + strconv.Itoa(intNo) + `,`)
				b.WriteString(`"drawTime":  "` + v.Drawtime.Format(formate) + `",`)
				b.WriteString(`"drawStatus": ` + strconv.Itoa(v.Drawstatus) + `,`)
				b.WriteString(`"voidReason": ` + strconv.Itoa(v.Voidreason) + `,`)
				b.WriteString(`"resultBalls": ` + v.Resultballs + ``)

				if i >= len(dlists)-1 {
					b.WriteString(`}`)
				} else {
					b.WriteString(`},`)
				}
			}
		}

		b.WriteString(`]`)

		if i >= len(lists)-1 {
			b.WriteString(`}`)
		} else {
			b.WriteString(`},`)
		}

	}
	b.WriteString(`],`)
	//trendsList

	//betslip
	b.WriteString(`"betSlip": "[{\"counterId\":320,\"drawNo\":61,\"bet\":{\"betType\":\"num\",\"selection\":\"66\"},\"ballNum\":\"53\"}]",`)

	//member begin
	b.WriteString(`"member": {`)
	b.WriteString(`"MemberId": 268042,`)
	b.WriteString(`"MemberName": "conku188",`)
	b.WriteString(`"CurrencyCode": "RMB",`)
	b.WriteString(`"Balance": "1341.2319"`)
	b.WriteString(`}`)
	//member end

	b.WriteString(`}`)
	b.WriteString("}")

	c.Request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	c.Writer.WriteString(b.String())
}

//开奖数据图形
func Trends(c *gin.Context) {
	formate := "2006-01-02T15:04:05+08:00"
	counterId, err := strconv.ParseInt(c.Param("number"), 10, 64)

	if err != nil {
		fmt.Println("counterId is :", err.Error())
	}

	now := time.Now()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	b := bytes.Buffer{}
	b.WriteString(`{"isSuccess": true,`)
	b.WriteString(`"data":[ `)

	if lists, err := models.GetDraws(time.Now(), beginTime, counterId); err == nil {
		for i, v := range lists {

			strNo := strconv.FormatInt(v.Drawno, 10)
			str := strNo[8:len(strNo)]

			intNo, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("strconv Atoi error")
			}

			b.WriteString("{")
			b.WriteString(`"counterId": 0,`)
			b.WriteString(`"drawNo": ` + strconv.Itoa(intNo) + `,`)
			b.WriteString(`"drawTime":  "` + v.Drawtime.Format(formate) + `",`)
			b.WriteString(`"drawStatus": 0,`)
			b.WriteString(`"voidReason": ` + strconv.Itoa(v.Voidreason) + `,`)
			b.WriteString(`"resultBalls": ` + v.Resultballs + ``)

			if i >= len(lists)-1 {
				b.WriteString(`}`)
			} else {
				b.WriteString(`},`)
			}
		}
	}

	b.WriteString(`]`)
	b.WriteString("}")

	c.Request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	c.Writer.WriteString(b.String())

}

func Float64toStr(v float64) string {
	string := strconv.FormatFloat(v, 'E', -1, 64)
	return string
}

func Float32toStr(v float64) string {
	string := strconv.FormatFloat(v, 'E', -1, 32)
	return string
}