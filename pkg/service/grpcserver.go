package service

import (
	"context"
	"gRPC-Chat/pkg/api"
	"log"
	"sync"
)

type Connection struct {
	stream api.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
}

func (s *Server) CreateStream(ctx *api.Connect, api api.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: api,
		id:     ctx.User.Id,
		active: true,
		error:  make(chan error),
	}
	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *api.Message) (*api.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *api.Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				log.Fatal("Sending message to: ", conn.stream)
				if err != nil {
					log.Fatalf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}

		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &api.Close{}, nil
}
