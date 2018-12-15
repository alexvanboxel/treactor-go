package execute

import (
	"testing"
)

func TestSuccess(t *testing.T) {
	for _, test := range []struct {
		in string
		expected string
		//want logging.LogEntry
	}{
		//{"[Ur]", ""},
		//{"[Ur],log:1", ""},
		//{"5s[Ur]", ""},
		//{"5[Ur,log:1]", ""},
		//{"5[Ur,log:1,xyz:4]", ""},
		//{"5[Ur,log:1,xyz:4]+5[Ur,log:1,xyz:4]", ""},
		{"2[5[Ur,log:1,xyz:4]+5[Ur,log:1,xyz:4]],x:1,y:2", ""},
	} {
		t.Logf(test.in)
		_, err := Parse(test.in)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		//if test.expected != err {
		//	t.FailNow()
		//}

	}

}

//func TestFail(t *testing.T) {
//
//	for _, test := range []struct {
//		in string
//		expected string
//		//want logging.LogEntry
//	}{
//		{"5x[Ur]", ""},
//		{"[Ur", ""},
//	} {
//		t.Logf(test.in)
//		_, err := Parse(test.in)
//		if err == nil {
//			t.FailNow()
//		}
//		//if test.expected != err {
//		//	t.FailNow()
//		//}
//
//	}
//
//}
