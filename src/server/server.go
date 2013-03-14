package server

import(
  "stomp"
  "net"
  "fmt"
  "os"
  "bufio"
)

var channels = make(map[string] chan string)
var connections = make(map[string] []net.Conn)
var hosts = make([]string,0)

func startQueue(name string){
  for{
    select{
      case body := <- channels[name]:
        stomp_message := stomp.NewStompMessage("MESSAGE", map[string]string {"destination":name}, body)

        for i := 0; i < len(connections[name]); i++{
          connections[name][i].Write([]byte(stomp_message.ToString()))
        }
    }
  }
}

func handleMessage(msg *stomp.StompMessage, conn net.Conn, ){
  switch msg.Command{
    case "CONNECT":
      //Respond with a connected frame
      //Add hostname to the connection pool
      hosts = append(hosts, conn.RemoteAddr().String())
      response := stomp.NewStompMessage("CONNECTED", map[string]string{"session": "0"}, "").ToString()
      conn.Write([]byte(response))

    case "SUBSCRIBE":
        // Check if destination exists
        // create the destination if it doesnt
        // Add connection to list
        // TODO Deal with ack headers
        // TODO Deal with multiple subscribes to the same topic
      if _, present := msg.Headers["destination"]; present{

        if _, present = channels[msg.Headers["destination"]]; !present{
          fmt.Println("Creating Queue : " + msg.Headers["destination"])
          channels[msg.Headers["destination"]] = make(chan string)

          go startQueue(msg.Headers["destination"])
        }

        if _, present = connections[msg.Headers["destination"]]; !present{
          connections[msg.Headers["destination"]] = make([]net.Conn,0)
          connections[msg.Headers["destination"]] = append(connections[msg.Headers["destination"]], conn)
        }else{
          connections[msg.Headers["destination"]] = append(connections[msg.Headers["destination"]], conn)
        }
      }

    case "UNSUBSCRIBE":
      // Remove the connection from the connection pool
      if _, present := connections[msg.Headers["destination"]]; present{
        for i:= 0; i < len(connections[msg.Headers["destination"]]); i++{
          if connections[msg.Headers["destination"]][i].RemoteAddr().String() == conn.RemoteAddr().String(){
            //TODO FIX THIS
            //copy(connections[msg.Headers["destination"]][i-1:], connections[msg.Headers["destination"]][i+1:])
            break
          }
        }
      }

    case "SEND":
      if _, present := channels[msg.Headers["destination"]]; present{
        channels[msg.Headers["destination"]] <- msg.Body
      }

    case "DISCONNECT":
      if _, present := connections[msg.Headers["destination"]]; present{
        for i:= 0; i < len(connections[msg.Headers["destination"]]); i++{
          if connections[msg.Headers["destination"]][i].RemoteAddr().String() == conn.RemoteAddr().String(){
            connections[msg.Headers["destination"]][i].Close()
            //TODO FIX THIS
            //connections[msg.Headers["destination"]][i] = append(connections[msg.Headers["destination"]][i-1:], connections[msg.Headers["destination"]][i+1:])
          }
        }
      }

    case "BEGIN":
    //TODO: Implement

    case "COMMIT":
    //TODO: Implement

    case "ABORT":
    //TODO: Implement

    case "ACK":
    //TODO: Implement
  }
}

func handleConnection(conn net.Conn){
  reader := bufio.NewReader(conn)

  for{
    msg, _ := reader.ReadString('\000')
    stomp_message, err := stomp.ToStompMessage(msg)

    if err != 0{
      return
    }else{
      handleMessage(stomp_message, conn)
    }
  }
}

/*
  We're going to need a way to get stats back
*/

/*
  Server starts and listens on the port
  on each new connection, start a go routine
  to deal with its messages, create any missing queues
  and topics. It then does what it needs to do with the
  message and dies.
*/
func Start(port string){
  ln, err := net.Listen("tcp", ":" + port)

  if err != nil{
    fmt.Println("Could not listen on tcp port " + port)
    os.Exit(1)
  }
  fmt.Println("Starting letgomq server on port " + port)
  /*
    This is the meat of the server function.
    We loop, waiting for connections. Connections
    are established and queues are created.
  */
  for{
    conn, err := ln.Accept()

    if err != nil{

      fmt.Println("An error occurred while trying to esablish connection with " + conn.RemoteAddr().String() + " .")
    }

    go handleConnection(conn)
  }
}
