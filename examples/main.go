package main

import (
	"reflect"

	"github.com/sirupsen/logrus"
	. "gitlab.com/jdbellamy/kundalini"
)

func main() {

	buf := []int{}
	ptr := reflect.ValueOf(&buf)

	types := []reflect.Type{}
	typesPtr := reflect.ValueOf(&types)

	v := []int{0, 1, 2, 3, 4}

	k, err := Wrap(v).
		Filter(even()).
		Map(double()).
		Export(ptr).
		Filter(firstN(2)).
		Concat(Wrap(buf).
			Filter(firstN(1)).
			ReleaseOrPanic()).
		Reduce(8, sum()).
		Push().
		Types().
		Export(typesPtr).
		Pop().
		Release()

	if err != nil {
		logrus.Errorln(err)
	}

	logrus.Infoln("     k:", k)
	logrus.Infoln("   buf:", buf)
	logrus.Infoln(" types:", types)
}

func even() Predicate {
	return func(x interface{}) bool {
		return x.(int)%2 == 0
	}
}

func double() Fn {
	return func(x interface{}) interface{} {
		return x.(int) * 2
	}
}

func sum() Transform {
	return func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}
}

func firstN(n int) Predicate {
	count := 0
	return func(interface{}) bool {
		var r = false
		if count < n {
			r = true
		}
		count = count + 1
		return r
	}
}

func init() {
	textFmt := new(logrus.TextFormatter)
	textFmt.DisableTimestamp = true
	logrus.SetFormatter(textFmt)
	logrus.SetLevel(logrus.DebugLevel)
}
