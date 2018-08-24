package core

import (
	"im/core/log"
	"net"
)

const HeaderLen = 10
const MaxBodyLen = 102400
const BaseBodyLen = 1024


type PacketData struct {
	PacketH		[]byte
	PacketB		[]byte
}

func ReadPacket(tcpConn *net.TCPConn, readPacket *PacketData) (*Message, error) {
	//读取包头的字节
	_, err := tcpConn.Read(readPacket.PacketH)
	if err != nil {
		return nil, err
	}
	//提取包头中的信息
	length, seq, cmd, version := extractHeader(readPacket.PacketH)
	if length < 0 {
		log.Infoln("包体长度为负数", err, tcpConn.RemoteAddr(), readPacket.PacketH)
		return nil, err
	}
	if length > MaxBodyLen {
		log.Infoln("包体长度超大", err, tcpConn.RemoteAddr(), readPacket.PacketH)
		return nil, err
	}
	if cmd == 0 {
		log.Infoln("包头cmd为空", tcpConn.RemoteAddr(), readPacket.PacketH)
		return nil, err
	}
	if cmd < 0 || cmd > 127 {
		log.Infoln("包头cmd异常", tcpConn.RemoteAddr(), readPacket.PacketH)
		return nil, err
	}
	msg := new(Message)
	msg.Cmd = cmd
	msg.Seq = seq
	msg.Version = version

	if length > 0 {
		if length > cap(readPacket.PacketB) {
			readPacket.PacketB = make([]byte, length)
		}
		_, err = readFullData(tcpConn, readPacket.PacketB, length)
		if err != nil {
			log.Infoln("包体读取异常", tcpConn.RemoteAddr(), readPacket.PacketH)
			return nil, err
		}
		msg.formatData(readPacket.PacketB[0:length])
	} else {
		msg.formatData(nil)
	}
	return msg, nil
}

func readFullData(tcpConn *net.TCPConn, dataBytes []byte, needLen int) ([]byte, error) {
	readPos := 0
	for {
		readLen, err := tcpConn.Read(dataBytes[readPos:needLen])
		if err != nil {
			return dataBytes, err
		}
		readPos += readLen
		if readLen == 0 || readPos == needLen {
			break
		}
	}
	return dataBytes, nil
}

func ReadRoutePacket(tcpConn *net.TCPConn, readPacket *PacketData) (*RouteMessage, error) {
	//读取包头的字节
	_, err := tcpConn.Read(readPacket.PacketH)
	if err != nil {
		return nil, err
	}
	//提取包头中的信息
	length, _, cmd, _ := extractHeader(readPacket.PacketH)
	if cmd == 0 {
		return nil, err
	}
	msg := new(RouteMessage)
	msg.Cmd = cmd

	if length > 0 {
		if length > cap(readPacket.PacketB) {
			readPacket.PacketB = make([]byte, length)
		}
		_, err = readFullData(tcpConn, readPacket.PacketB, length)
		if err != nil {
			log.Infoln("包体读取异常", tcpConn.RemoteAddr(), readPacket.PacketH)
			return nil, err
		}
		msg.formatData(readPacket.PacketB[0:length])
	} else {
		msg.formatData(nil)
	}
	return msg, nil
}

func extractHeader(header []byte) (int, int, int8, int8) {
	length := int(header[3]) | int(header[2])<<8 | int(header[1])<<16 | int(header[0])<<24
	seq := int(header[7]) | int(header[6])<<8 | int(header[5])<<16 | int(header[4])<<24
	cmd := int8(header[8])
	version := int8(header[9])
	return length, seq, cmd, version
}


/*	var length int32
	var seq int32
	buffer := bytes.NewBuffer(header)
	binary.Read(buffer, binary.BigEndian, &length)
	binary.Read(buffer, binary.BigEndian, &seq)
	cmd, _ := buffer.ReadByte()
	version, _ := buffer.ReadByte()*/
/*	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, length)
	binary.Write(buffer, binary.BigEndian, seq)
	buffer.Write([]byte{cmd, version})
	buffer.Write(routeDataBytes)*/

func SetPacketHeader(cmd byte, version byte, seq int, bodyLen int, packetData []byte) {
	packetData[0] = byte(bodyLen >> 24)
	packetData[1] = byte(bodyLen >> 16)
	packetData[2] = byte(bodyLen >> 8)
	packetData[3] = byte(bodyLen)
	packetData[4] = byte(seq >> 24)
	packetData[5] = byte(seq >> 16)
	packetData[6] = byte(seq >> 8)
	packetData[7] = byte(seq)
	packetData[8] = cmd
	packetData[9] = version
}

