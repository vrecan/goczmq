package goczmq

import (
	"testing"
)

func TestZproxy(t *testing.T) {
	proxy := NewZproxy()
	proxy.SetFrontend(PULL, "inproc://proxy_front")
	proxy.SetBackend(PUSH, "inproc://proxy_back")
	proxy.SetCapture("inproc://proxy_capture")

	send := NewZsock(PUSH)
	err := send.Connect("inproc://proxy_front")
	if err != nil {
		t.Error(err)
	}

	recv := NewZsock(PULL)
	err = recv.Connect("inproc://proxy_back")
	if err != nil {
		t.Error(err)
	}

	cap := NewZsock(PULL)
	err = cap.Bind("inproc://proxy_capture")
	if err != nil {
		t.Error(err)
	}

	send.SendBytes([]byte("hello proxy"), 0)

	b, err := cap.RecvBytes()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "hello proxy" {
		t.Error("message is wrong")
	}

	b, err = recv.RecvBytes()
	if err != nil {
		t.Error(err)
	}
	if string(b) != "hello proxy" {
		t.Error("message is wrong")
	}
	proxy.Destroy()
}