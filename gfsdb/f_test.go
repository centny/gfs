package gfsdb

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/ffcm"
	_ "github.com/Centny/gfs/test"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
	"regexp"
	"runtime"
	"testing"
	"time"
)

func TestF(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	mgo.C(CN_F).RemoveAll(nil)
	var do_f = func(i int) int {
		var rt = &F{
			Path: "xxx",
			SHA:  "abc",
			MD5:  "xyz",
			Pub:  "/s",
			Mark: []string{"jkk"},
		}
		var updated, err = FOI_F(rt)
		if err != nil {
			t.Error(err.Error())
			return 0
		}
		return updated
	}
	var updated = do_f(0)
	if updated < 1 {
		t.Error("error")
		return
	}
	used, _ := tutil.DoPerf(100, "", func(i int) {
		do_f(i)
	})
	tc, err := CountF()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if tc != 1 {
		t.Error("error")
		return
	}
	fmt.Printf("done with used(%vms),per(%vms)\n", used, used/100)
	rt, err := FindHashF("abc", "xyz")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = FindF(rt.Id)
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = FindMarkF(rt.Mark[0])
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = FindHashF("", "")
	if err == nil {
		t.Error("error")
		return
	}
	_, err = AddMarkF(rt.Id, []string{"kjj"})
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = FindMarkF("kjj")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = FindPubF("/s")
	if err != nil {
		t.Error(err.Error())
		return
	}
	//
	_, err = FOI_F(&F{})
	if err == nil {
		t.Error("error")
		return
	}
	//
	_, err = FOI_F(&F{
		Path: "sdfd",
	})
	if err == nil {
		t.Error("error")
		return
	}
	//
	tmgo.Mock = true
	tmgo.SetMckC("Query-Apply", 0)
	_, err = FOI_F(&F{
		Path: "sfkdf",
		SHA:  "abc",
		MD5:  "xyz",
	})
	if err == nil {
		t.Error("error")
		return
	}
	tmgo.ClearMock()
	//
}

func TestFFCM(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	mgo.C(CN_F).RemoveAll(nil)
	ffcm.StartTest("../../ffcm/ffcm_s.properties", "../../ffcm/ffcm_c.properties", dtm.MemDbc, NewFFCM_H())
	time.Sleep(3 * time.Second)
	fmt.Println(ffcm.SRV)
	var rt = &F{
		Path: "xx.mp4",
		SHA:  "abc",
		MD5:  "xyz",
		Mark: []string{"jkk"},
	}
	ffcm.SRV.Db.(*dtm.MemH).Errs["Add"] = util.Err("mock error")
	var _, err = FOI_F(&F{
		Path: "xxkjk.mp4",
		SHA:  "abcsd",
		MD5:  "xyzfd",
		Mark: []string{"jkk"},
	})
	if err != nil {
		t.Error("error")
		return
	}
	ffcm.SRV.Db.(*dtm.MemH).Errs["Add"] = nil
	_, err = FOI_F(rt)
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = FOI_F(&F{
		Path: "XXXX",
		SHA:  "abcx",
		MD5:  "xyzx",
		Mark: []string{"jkk"},
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	// if true {
	// 	return
	// }
	for {
		rt, err = FindF(rt.Id)
		if err != nil {
			t.Error(err.Error())
			return
		}
		fmt.Println("waiting result...")
		if len(rt.Info) > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("result->", util.S2Json(rt.Info))
	fmt.Println(rt.Id)
}

func TestReg(t *testing.T) {
	fmt.Println(regexp.MustCompile("^[^X]+[^K]+.*$").MatchString("XXX"))
}

func TestFFCM_H_err(t *testing.T) {
	var ffcm = NewFFCM_H()
	//
	var err = ffcm.OnDone(nil, &dtm.Task{
		Proc: map[string]*dtm.Proc{
			"xx": &dtm.Proc{},
		},
	})
	if err == nil {
		t.Error("error")
		return
	}
	//
	err = ffcm.OnDone(nil, &dtm.Task{
		Proc: map[string]*dtm.Proc{
			"xx": &dtm.Proc{
				Res: "sss",
			},
		},
	})
	if err == nil {
		t.Error("error")
		return
	}
	//
	err = ffcm.OnDone(nil, &dtm.Task{
		Proc: map[string]*dtm.Proc{
			"xx": &dtm.Proc{
				Res: util.Map{},
			},
		},
	})
	if err == nil {
		t.Error("error")
		return
	}
}
