package stomp

import(
  "testing"
)

func TestNewStompMessage(t *testing.T){
  command := "CONNECT"
  headers := []string {"header1:foo", "header2:bar"}
  body := "stomp message body"

  msg := NewStompMessage(command, headers, body)

  if msg.Command != command{
    t.Errorf("Stomp message Commands do not match")
  }

  for i := 0; i < len(msg.Headers); i++{
    if msg.Headers[i] != headers[i]{
      t.Errorf("Stomp message headers do not match")
    }
  }

  if msg.Body != body{
    t.Errorf("Stomp message body does not match")
  }
}

func TestToStompMessage(t *testing.T){
  stomp_string := "CONNECT\nheader1:foo\nheader2:bar\n\nfoobar\nbarfoo\000"
  msg, code := ToStompMessage(stomp_string)

  if code != 0{
    t.Errorf("Could not convert string into stomp message")
  }

  if msg.Command != "CONNECT"{
    t.Errorf("Message command does not match")
  }

  if msg.Headers[0] != "header1:foo" && msg.Headers[1] != "header2:bar"{
    t.Errorf("Message headers do not match")
  }

  if msg.Body != "foobar\nbarfoo"{
    t.Errorf("Message body does not match")
  }

  msg, code = ToStompMessage("")
  if code != 1{
    t.Errorf("Tried to create stomp message from empty string")
  }

  msg, code = ToStompMessage("NOTVALID")
  if code != 1{
    t.Errorf("Tried to create stomp message with invalid command")
  }

  msg, code = ToStompMessage("CONNECT\n\nfoobars\000")
  if code !=1 {
    t.Errorf("Tried to create stomp message without headers")
  }
}

func TestToString(t *testing.T){
  headers := []string {"header1:foo", "header2:bar"}
  msg := NewStompMessage("CONNECT", headers, "foobar\nbarfoo")

  if msg.ToString() != "CONNECT\nheader1:foo\nheader2:bar\n\nfoobar\nbarfoo\000"{
    t.Errorf("Message structure does not translate to the correct string")
  }
}

func TestSupportsVersion(t *testing.T){
  if !SupportsVersion("1.0"){
    t.Errorf("Incorrectly compared version 1.0")
  }

  if SupportsVersion("1.1"){
    t.Errorf("Incorrectly compared version 1.1")
  }
}
