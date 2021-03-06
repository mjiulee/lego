package utils

import (
	"fmt"
	"time"
)

const (
	KC_TIME_GREP_MIN  = 0 // 分钟
	KC_TIME_GREP_HOUR = 1 // 小时
	KC_TIME_GREP_DAY  = 2 // 小时
)

func TimeBefor(t *time.Time, gtpye int, grep time.Duration) (f *time.Time) {
	ts := ""
	if gtpye == KC_TIME_GREP_MIN {
		ts = "-1m"
	} else if gtpye == KC_TIME_GREP_HOUR {
		ts = "-1h"
	} else {
		ts = "-24h"
	}

	m, _ := time.ParseDuration(ts)
	m1 := t.Add(grep * m)
	return &m1
}

func DateBefor(t *time.Time, ndate int) (f *time.Time) {
	yesTime := t.AddDate(0, 0, -ndate)
	return &yesTime
}

func TimeAfter(t *time.Time, gtpye int, grep time.Duration) (f *time.Time) {
	ts := ""
	if gtpye == KC_TIME_GREP_MIN {
		ts = "1m"
	} else if gtpye == KC_TIME_GREP_HOUR {
		ts = "1h"
	} else {
		ts = "24h"
	}

	m, _ := time.ParseDuration(ts)
	m1 := t.Add(grep * m)
	return &m1
}


func GetTimeStamp() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d", timestamp)
}
