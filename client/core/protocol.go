package core

import (
	"net"
	"encoding/binary"
	"bytes"
	"github.com/golang/protobuf/proto"
	"fmt"
)
const HeaderLen = 10
const MaxBodyLen = 2048000


func ReadPacket(tcpConn *net.TCPConn) *Message {
	//读取包头的字节
	header := make([]byte, HeaderLen)
	_, err := tcpConn.Read(header)
	if err != nil {
		fmt.Println("包头读取失败", err , tcpConn.RemoteAddr(), header)
		return nil
	}
	//提取包头中的信息
	length, seq, cmd, version := extractHeader(header)
	if length < 0 {
		fmt.Println("包体长度为负数", err, tcpConn.RemoteAddr(), header)
		return	nil
	}
	if length > MaxBodyLen {
		fmt.Println("包体长度超大", err, tcpConn.RemoteAddr(), header)
		return nil
	}
	if cmd == 0  {
		fmt.Println("包头cmd为空", tcpConn.RemoteAddr(), header)
		return nil
	}
	if cmd < 0 || cmd > 127 {
		fmt.Println("包头cmd异常", tcpConn.RemoteAddr(), header)
		return nil
	}
	var dataBytes []byte
	if length > 0 {
		dataBytes,err = readFullData(tcpConn, length)
		if cmd == SyncOfflineMsg {
			fmt.Println(dataBytes)
		}
		if err != nil {
			fmt.Println("包体读取异常", tcpConn.RemoteAddr(), header)
		}
	}
	msg := new(Message)
	msg.Cmd = cmd
	msg.Seq = seq
	msg.Version = version
	msg.formatData(dataBytes)
	return msg
}


func readFullData(tcpConn *net.TCPConn, needLen int) ([]byte, error) {
	dataBytes := make([]byte, needLen)
	readPos := 0
	for {
		readLen, err := tcpConn.Read(dataBytes[readPos:])
		if err != nil {
			fmt.Println("包体读取失败", err, tcpConn.RemoteAddr(), "body", dataBytes[0:200])
			return dataBytes, err
		}
		readPos += readLen
		if readLen == 0 || readPos == needLen {
			break
		}
	}
	return dataBytes, nil
}


func ReadRoutePacket(tcpConn *net.TCPConn) *RouteMessage {
	//读取包头的字节
	header := make([]byte, HeaderLen)
	_, err := tcpConn.Read(header)

	if err != nil {
		fmt.Println("包头读取失败", err)
		return nil
	}
	//提取包头中的信息
	length, _, cmd, _ := extractHeader(header)
	if cmd == 0  {
		fmt.Println("包头cmd为空")
		return nil
	}
	var dataBytes []byte
	if length > 0 {
		dataBytes,err = readFullData(tcpConn, length)
		if err != nil {
			fmt.Println("包体读取异常", tcpConn.RemoteAddr(), header)
		}
	}
	msg := new(RouteMessage)
	msg.Cmd = cmd
	msg.formatData(dataBytes)
	return msg
}


func extractHeader(header []byte) (int, int, int8, int8) {
	var length int32
	var seq int32
	buffer := bytes.NewBuffer(header)
	binary.Read(buffer, binary.BigEndian, &length)
	binary.Read(buffer, binary.BigEndian, &seq)
	cmd, _ := buffer.ReadByte()
	version, _ := buffer.ReadByte()
	return int(length), int(seq), int8(cmd), int8(version)
}


func MakePacket(msg *Message) []byte {
	data := msg.Data
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("消息msg编码失败", err, msg)
	}
	length := int32(len(dataBytes))
	seq := int32(1)
	version := byte(1)
	cmd := byte(msg.Cmd)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, length)
	binary.Write(buffer, binary.BigEndian, seq)
	buffer.Write([]byte{cmd, version})
	buffer.Write(dataBytes)
	return buffer.Bytes()
}


func MakeRoutePacket(msg *RouteMessage) []byte {
	routeData := msg.RouteData
	routeDataBytes, err := proto.Marshal(routeData)
	if err != nil {
		fmt.Println("消息msg编码失败", err, msg)
	}
	length := int32(len(routeDataBytes))
	seq := int32(1)
	version := byte(1)
	cmd := byte(msg.Cmd)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, length)
	binary.Write(buffer, binary.BigEndian, seq)
	buffer.Write([]byte{cmd, version})
	buffer.Write(routeDataBytes)
	return buffer.Bytes()
}
