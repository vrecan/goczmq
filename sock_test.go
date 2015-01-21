package goczmq

import (
	"bytes"
	"testing"
)

func TestSendFrame(t *testing.T) {
	pushSock := NewSock(PUSH)
	defer pushSock.Destroy()

	pullSock := NewSock(PULL)
	defer pullSock.Destroy()

	_, err := pullSock.Bind("inproc://test-sock")
	if err != nil {
		t.Errorf("repSock.Bind failed: %s", err)
	}

	err = pushSock.Connect("inproc://test-sock")
	if err != nil {
		t.Errorf("reqSock.Connect failed: %s", err)
	}

	pushSock.SendFrame([]byte("Hello"), 0)
	msg, flag, err := pullSock.RecvFrame()
	if bytes.Compare(msg, []byte("Hello")) != 0 {
		t.Errorf("expected 'Hello' received '%s'", msg)
	}

	if flag != 0 {
		t.Errorf("flag shouled have been 0, is '%d'", flag)
	}
}

func TestSendMessage(t *testing.T) {
	pushSock := NewSock(PUSH)
	defer pushSock.Destroy()

	pullSock := NewSock(PULL)
	defer pullSock.Destroy()

	_, err := pullSock.Bind("inproc://test-sock")
	if err != nil {
		t.Errorf("repSock.Bind failed: %s", err)
	}

	err = pushSock.Connect("inproc://test-sock")
	if err != nil {
		t.Errorf("reqSock.Connect failed: %s", err)
	}

	pushSock.SendMessage([][]byte{[]byte("Hello")})
	msg, err := pullSock.RecvMessage()
	if err != nil {
		t.Errorf("pullsock.RecvMessage() failed: %s", err)
	}

	if bytes.Compare(msg[0], []byte("Hello")) != 0 {
		t.Errorf("expected 'Hello' received '%s'", msg)
	}
}

func TestPUBSUB(t *testing.T) {
	_, err := NewPUB("bogus://bogus")
	if err == nil {
		t.Error("NewPUB should have returned error and did not")
	}

	_, err = NewSUB("bogus://bogus", "")
	if err == nil {
		t.Error("NewSUB should have returned error and did not")
	}

	pub, err := NewPUB("inproc://pub1,inproc://pub2")
	if err != nil {
		t.Errorf("NewPUB failed: %s", err)
	}
	defer pub.Destroy()

	sub, err := NewSUB("inproc://pub1,inproc://pub2", "")
	if err != nil {
		t.Errorf("NewSUB failed: %s", err)
	}
	defer sub.Destroy()

	err = pub.SendFrame([]byte("test pub sub"), 0)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	frame, _, err := sub.RecvFrame()
	if err != nil {
		t.Errorf("RecvFrame failed: %s", err)
	}

	if string(frame) != "test pub sub" {
		t.Errorf("Expected 'test pub sub', received %s", frame)
	}
}

func TestREQREP(t *testing.T) {
	_, err := NewREQ("bogus://bogus")
	if err == nil {
		t.Error("NewREQ should have returned error and did not")
	}

	_, err = NewREP("bogus://bogus")
	if err == nil {
		t.Error("NewREP should have returned error and did not")
	}

	rep, err := NewREP("inproc://rep1,inproc://rep2")
	if err != nil {
		t.Errorf("NewREP failed: %s", err)
	}
	defer rep.Destroy()

	req, err := NewREQ("inproc://rep1,inproc://rep2")
	if err != nil {
		t.Errorf("NewREQ failed: %s", err)
	}
	defer req.Destroy()

	err = req.SendFrame([]byte("Hello"), 0)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	reqframe, _, err := rep.RecvFrame()
	if err != nil {
		t.Errorf("RecvFrame failed: %s", err)
	}

	if string(reqframe) != "Hello" {
		t.Errorf("Expected 'Hello', received '%s", string(reqframe))
	}

	err = rep.SendFrame([]byte("World"), 0)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	repframe, _, err := req.RecvFrame()
	if err != nil {
		t.Errorf("RecvFrame failed: %s", err)
	}

	if string(repframe) != "World" {
		t.Errorf("Expected 'World', received '%s", string(repframe))
	}

}

func TestPUSHPULL(t *testing.T) {
	_, err := NewPUSH("bogus://bogus")
	if err == nil {
		t.Error("NewPUSH should have returned error and did not")
	}

	_, err = NewPULL("bogus://bogus")
	if err == nil {
		t.Error("NewPULL should have returned error and did not")
	}

	push, err := NewPUSH("inproc://push1,inproc://push2")
	if err != nil {
		t.Errorf("NewPUSH failed: %s", err)
	}
	defer push.Destroy()

	pull, err := NewPULL("inproc://push1,inproc://push2")
	if err != nil {
		t.Errorf("NewPULL failed: %s", err)
	}
	defer pull.Destroy()

	err = push.SendFrame([]byte("Hello"), 1)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	err = push.SendFrame([]byte("World"), 0)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	msg, err := pull.RecvMessage()
	if err != nil {
		t.Errorf("RecvMessage failed: %s", err)
	}

	if string(msg[0]) != "Hello" {
		t.Errorf("Expected 'Hello', received '%s", string(msg[0]))
	}

	if string(msg[1]) != "World" {
		t.Errorf("Expected 'World', received '%s", string(msg[0]))
	}
}

func TestROUTERDEALER(t *testing.T) {
	_, err := NewDEALER("bogus://bogus")
	if err == nil {
		t.Error("NewDEALER should have returned error and did not")
	}

	_, err = NewROUTER("bogus://bogus")
	if err == nil {
		t.Error("NewROUTER should have returned error and did not")
	}

	dealer, err := NewDEALER("inproc://router1,inproc://router2")
	if err != nil {
		t.Errorf("NewDEALER failed: %s", err)
	}
	defer dealer.Destroy()

	router, err := NewROUTER("inproc://router1,inproc://router2")
	if err != nil {
		t.Errorf("NewROUTER failed: %s", err)
	}
	defer router.Destroy()

	err = dealer.SendFrame([]byte("Hello"), 0)
	if err != nil {
		t.Errorf("SendMessage failed: %s", err)
	}

	msg, err := router.RecvMessage()
	if err != nil {
		t.Errorf("RecvMessage failed: %s", err)
	}
	if len(msg) != 2 {
		t.Error("message should have 2 frames")
	}

	if string(msg[1]) != "Hello" {
		t.Errorf("Expected 'Hello', received '%s", string(msg[0]))
	}

	msg[1] = []byte("World")

	err = router.SendMessage(msg)
	if err != nil {
		t.Errorf("SendMessage failed: %s", err)
	}

	msg, err = dealer.RecvMessage()
	if err != nil {
		t.Errorf("RecvMessage failed: %s", err)
	}

	if len(msg) != 1 {
		t.Error("message should have 1 frames")
	}

	if string(msg[0]) != "World" {
		t.Errorf("Expected 'World', received '%s", string(msg[0]))
	}
}

func TestXSUBXPUB(t *testing.T) {
	_, err := NewXPUB("bogus://bogus")
	if err == nil {
		t.Error("NewXPUB should have returned error and did not")
	}

	_, err = NewXSUB("bogus://bogus")
	if err == nil {
		t.Error("NewXSUB should have returned error and did not")
	}

	xpub, err := NewXPUB("inproc://xpub1,inproc://xpub2")
	if err != nil {
		t.Errorf("NewXPUB failed: %s", err)
	}
	defer xpub.Destroy()

	xsub, err := NewXSUB("inproc://xpub1,inproc://xpub2")
	if err != nil {
		t.Errorf("NewXSUB failed: %s", err)
	}
	defer xsub.Destroy()
}

func TestPAIR(t *testing.T) {
	_, err := NewPAIR("bogus://bogus")
	if err == nil {
		t.Error("NewPAIR should have returned error and did not")
	}

	pair1, err := NewPAIR(">inproc://pair")
	if err != nil {
		t.Errorf("NewPAIR failed: %s", err)
	}
	defer pair1.Destroy()

	pair2, err := NewPAIR("@inproc://pair")
	if err != nil {
		t.Errorf("NewPAIR failed: %s", err)
	}
	defer pair2.Destroy()
}

func TestSTREAM(t *testing.T) {
	_, err := NewSTREAM("bogus://bogus")
	if err == nil {
		t.Error("NewSTREAM should have returned error and did not")
	}

	stream1, err := NewSTREAM(">inproc://stream")
	if err != nil {
		t.Errorf("NewSTREAM failed: %s", err)
	}
	defer stream1.Destroy()

	stream2, err := NewSTREAM("@inproc://stream")
	if err != nil {
		t.Errorf("NewSTREAM failed: %s", err)
	}
	defer stream2.Destroy()

}

func TestPollin(t *testing.T) {
	push, err := NewPUSH("inproc://pollin")
	if err != nil {
		t.Errorf("NewPUSH failed: %s", err)
	}
	defer push.Destroy()

	pull, err := NewPULL("inproc://pollin")
	if err != nil {
		t.Errorf("NewPULL failed: %s", err)
	}
	defer pull.Destroy()

	if pull.Pollin() {
		t.Errorf("Pollin returned true should be false")
	}

	err = push.SendFrame([]byte("Hello World"), 0)
	if err != nil {
		t.Errorf("SendFrame failed: %s", err)
	}

	if !pull.Pollin() {
		t.Errorf("Pollin returned false should be true")
	}
}

func TestPollout(t *testing.T) {
	push := NewSock(PUSH)
	_, err := push.Bind("inproc://pollout")
	if err != nil {
		t.Errorf("failed binding test socket: %s", err)
	}
	defer push.Destroy()

	if push.Pollout() {
		t.Errorf("Pollout returned true should be false")
	}

	pull := NewSock(PULL)
	defer pull.Destroy()

	err = pull.Connect("inproc://pollout")
	if err != nil {
		t.Errorf("failed connecting test socket: %s", err)
	}

	if !push.Pollout() {
		t.Errorf("Pollout returned false should be true")
	}
}
