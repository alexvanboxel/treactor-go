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
		{"[Ur]", "1s[Ur]"},
		{"[Ur],log:1", "1s[Ur],log:1"},
		{"5s[Ur]", "5s[Ur]"},
		{"5[Ur,log:1]", "5s[Ur,log:1]"},
		{"5[Ur,log:1,xyz:4]", "5s[Ur,log:1,xyz:4]"},
		{"5[Ur,log:1,xyz:4]^5[Ur,log:1,xyz:4]", "5s[Ur,log:1,xyz:4]^5s[Ur,log:1,xyz:4]"},
		{"2[5[Ur,log:1,xyz:4]^5[Ur,log:1,xyz:4]],x:1,y:2", "2s[5[Ur,log:1,xyz:4]^5[Ur,log:1,xyz:4]],x:1,y:2"},
	} {
		t.Logf(test.in)
		plan, err := Parse(test.in)
		t.Logf(plan.String())
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if test.expected != plan.String() {
			t.FailNow()
		}

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
