package parser

import (
	"testing"
	"time"
)

type CmdArgs struct {
	ConfFile     string        `required:"false" name:"config" default:"/etc/daemon.conf" description:"Конфигурационный файл"`
	Daemon       bool          `required:"false" name:"daemon" default:"true" description:"Запуск приложения в режиме daemon"`
	Pool         uint64        `required:"false" name:"pool" default:"5" description:"кол-во пул потоков"`
	TimeOut      float64       `required:"false" name:"timeout" default:"2.5" description:"time out"`
	DurationTime time.Duration `required:"false" name:"duration" default:"3ms" description:"duration time"`
}

func TestArgParse(t *testing.T) {
	arg := &CmdArgs{}
	if err := GetArguments(arg); err != nil {
		t.Error("flag parse error", err)
	}

	if arg.ConfFile != "/etc/daemon.conf" {
		t.Errorf("exected %q got %q", "/etc/daemon.conf", arg.ConfFile)
	}

	if arg.Daemon != true {
		t.Errorf("exected %t got %t", true, arg.Daemon)
	}

	if arg.Pool != 5 {
		t.Errorf("exected %d got %d", true, arg.Pool)
	}

	if arg.TimeOut != 2.5 {
		t.Errorf("exected %f got %f", true, arg.TimeOut)
	}

	t.Logf("arg: %#v \n", arg)
}
