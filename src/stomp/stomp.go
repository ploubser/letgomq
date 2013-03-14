package stomp

/*
  Basic stomp message structure. These can be either Client
  or Server messages.

  Client commands : SEND, SUBSCRIBE, UNSUBSCRIBE, BEGIN, COMMIT,
                    ABORT, ACK,  DISCONNECT, CONNECT

  Server commands : CONNECTED, MESSAGE, RECEIPT, ERROR

  All stomp messages are implicitly terminated with a \000
  character. This needs to be placed in any result string
  being sent over the wire.
*/

import(
  "strings"
  "regexp"
)

const(
  Version = "1.0" //Stomp protocol version
  Headers = "SEND|SUBSCRIBE|UNSUBSCRIBE|BEGIN|COMMIT|ABORT|ACK|DISCONNECT|CONNECT|CONNECTED|MESSAGE|RECEIPT|ERROR" //Stomp Headers
)


type StompMessage struct{
  Command string
  Headers map[string] string
  Body string
}

/*
   Creates and returns a new stomp messasge
*/
func NewStompMessage(command string, headers map[string]string, body string) *StompMessage{
  return &StompMessage{Command: command, Headers: headers, Body: body}
}


/*
  Tries to convert a string into a stomp message struct
  Returns error code 0 on success, 1 on failure
*/
func ToStompMessage(msg string)(*StompMessage, uint){
  if len(msg) < 1{
    return &StompMessage{}, 1
  }

  // Get rid of the trailing null character in the stomp string
  msg = strings.TrimRight(msg, "\000")
  msg_slice := strings.Split(msg,"\n")

  var command string

  if m, _ := regexp.MatchString(Headers ,msg_slice[0]); m{
    command = msg_slice[0]
  } else{
    return &StompMessage{}, 1
  }

  headers := make(map[string] string)
  var body string

  for i:=1; i < len(msg_slice); i++{

    if m, _ := regexp.MatchString(".+:.+", msg_slice[i]); m{
      k := strings.Split(msg_slice[i], ":")
      headers[k[0]] = k[1]
    }else if msg_slice[i] != ""{
      if len(headers) == 0{
        return &StompMessage{}, 1
      }

      body = strings.Join(msg_slice[i:len(msg_slice)], "\n")
      break
    }
  }

  return &StompMessage{Command: command, Headers: headers, Body: body}, 0
}

/*
  Returns a string representation of the stomp message ready
  for transmission over the wire
*/
func (s StompMessage) ToString() string{
  var headerstring []string
  for k, _ := range s.Headers{
    headerstring = append(headerstring, k + ": " + s.Headers[k])
  }
  return s.Command + "\n" +
         strings.Join(headerstring, "\n") + "\n\n" +
         s.Body + "\000"
}

/*
  Compares the given stomp message version with the supported version
*/
func SupportsVersion(version string) bool{
  return Version == version
}
