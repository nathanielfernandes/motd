package slp

import (
	"errors"
	"fmt"
	"net"
)

const (
	NEXTSTATE_STATUS   int32 = 1
	NEXTSTATE_LOGIN    int32 = 2
	NEXTSTATE_TRANSFER int32 = 3
)

type Server struct {
	onStatus func(conn net.Conn) (StatusResponse, error)
	onLogin  func(conn net.Conn, loginStart LoginStart) error
	onError  func(conn net.Conn, err error)
}

func NewServer() *Server {
	return &Server{
		onStatus: func(conn net.Conn) (StatusResponse, error) { return StatusResponse{}, nil },
		onLogin:  func(conn net.Conn, loginStart LoginStart) error { return nil },
		onError:  func(conn net.Conn, err error) { fmt.Println(err) },
	}
}

func (s *Server) OnStatus(f func(conn net.Conn) (StatusResponse, error)) *Server {
	s.onStatus = f
	return s
}

func (s *Server) OnLogin(f func(conn net.Conn, loginStart LoginStart) error) *Server {
	s.onLogin = f
	return s
}

func (s *Server) OnError(f func(conn net.Conn, err error)) *Server {
	s.onError = f
	return s
}

func (s *Server) HandleConnection(conn net.Conn) (err error) {
	defer conn.Close()

	var handshake Handshake
	if err := handshake.Read(conn); err != nil {
		return errors.New("Error reading handshake: " + err.Error())
	}

	switch handshake.NextState {
	case NEXTSTATE_STATUS:
		if err = ReadStatusRequest(conn); err != nil {
			return errors.New("Error reading status request: " + err.Error())
		}

		var status StatusResponse
		if status, err = s.onStatus(conn); err != nil {
			return errors.New("Error handling status: " + err.Error())
		}

		packet := NewPacket(0x00)
		if _, err := status.Write(&packet); err != nil {
			return errors.New("Error writing status response: " + err.Error())
		}

		if _, err := packet.WritePacket(conn); err != nil {
			return errors.New("Error writing status response: " + err.Error())
		}

		if err := ListenAndPong(conn); err != nil {
			return errors.New("Error sending ping: " + err.Error())
		}
	case NEXTSTATE_LOGIN:
		var loginStart LoginStart
		if err := loginStart.Read(conn); err != nil {
			return errors.New("Error reading login start: " + err.Error())
		}

		if err := loginStart.WriteSuccess(conn); err != nil {
			return errors.New("Error writing login success: " + err.Error())
		}

		if err := ReadStatusRequest(conn); err != nil {
			return errors.New("Error reading login acknowledged: " + err.Error())
		}

		if err := s.onLogin(conn, loginStart); err != nil {
			return errors.New("Error handling login: " + err.Error())
		}

	}

	return nil
}

func (s *Server) ListenAndServe(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic("Error listening: " + err.Error())
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	if err := s.HandleConnection(conn); err != nil {
		s.onError(conn, err)
	}
}
