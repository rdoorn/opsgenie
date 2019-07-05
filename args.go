package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// attempt to find time in a string, and attempts to find string parameters
func parseArgs(args ...string) (timeString string, rest string) {
	//timeString = "24h"
	str := []string{}
	strLock := false
	strLock2 := false
	for _, a := range args {

		// end of string "
		if a[len(a)-1:] == "\"" && strLock == true {
			str = append(str, a[:len(a)-1])
			strLock = false
			continue
		}

		// end of string ''
		if a[len(a)-1:] == "'" && strLock2 == true {
			str = append(str, a[:len(a)-1])
			strLock2 = false
			continue
		}
		// lock of "
		if strLock {
			str = append(str, a)
			continue
		}
		// lock of '
		if strLock2 {
			str = append(str, a)
			continue
		}

		// start of string "
		if a[:1] == "\"" && strLock == false {
			str = append(str, a[1:])
			strLock = true
			continue
		}

		// start of string '
		if a[:1] == "'" && strLock2 == false {
			str = append(str, a[1:])
			strLock2 = true
			continue
		}

		if isTime(a) {
			timeString = a
		} else {
			str = append(str, a)
		}
	}

	rest = strings.Join(str, " ")

	return
}

func isTime(str string) bool {
	r, err := regexp.Compile(`^(\d+[mhds])+$`)
	if err != nil {
		panic(err)
	}
	return r.MatchString(str)
}

func addTime(t time.Time, s string) time.Time {
	return modTime(t, s, 1)
}

func subTime(t time.Time, s string) time.Time {
	return modTime(t, s, -1)
}

func modTime(t time.Time, s string, mod int) time.Time {
	r, _ := regexp.Compile(`(\d+)([dmhs])`)
	res := r.FindAllStringSubmatch(s, -1)
	for _, re := range res {
		val, _ := strconv.Atoi(re[1])
		dur := time.Duration(val * mod)
		switch re[2] {
		case "h":
			t = t.Add(dur * time.Hour)
		case "m":
			t = t.Add(dur * time.Minute)
		case "s":
			t = t.Add(dur * time.Second)
		case "d":
			t = t.Add(dur * 24 * time.Hour)
		}
	}
	return t
}

func arrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
