package common

import (
	"crypto/md5"
	"encoding/binary"
	"net"
)

type Event uint16

const (
	EVENT_CLINET_1 Event = 0
	EVENT_SERVER_2 Event = 1
	EVENT_SERVER_3 Event = 2
	EVENT_SERVER_4 Event = 3
	EVENT_CLIENT_5 Event = 4
	EVENT_CLIENT_6 Event = 5
	EVENT_SERVER_7 Event = 6
	EVENT_SERVER_8 Event = 7
)

var EventComment = map[Event]string{
	EVENT_CLINET_1: "1.1 client -->> i1p1",         //
	EVENT_SERVER_2: "Test1.0 1.2 client <<-- i1p1", // 检测是否能ping通
	EVENT_SERVER_3: "Test2.0 1.3 client <<-- i2p2", // 检测是否是公网IP
	EVENT_SERVER_4: "Test3.0 1.4 client <<-- i1p3", // 检测是否是IP限制型
	EVENT_CLIENT_5: "2.5 client -->> i1p1",         //
	EVENT_CLIENT_6: "2.6 client -->> i2p2",         //
	EVENT_SERVER_7: "Test4.1 2.7 client <<-- i1p1", //
	EVENT_SERVER_8: "Test4.2 2.8 client <<-- i2p2", // 根据Test4的两次回包是否一致，判断是否是端口限制型 或者是否是可预测对称型
}

type Request struct {
	OpID      int64
	EventID   Event
	TimeStamp int64
}

type Response struct {
	OpID              int64
	EventID           Event
	ClientPublicIP    uint32
	ClientPublicPort  int32
	ServerOtherIP     uint32
	ServerOtherPort   int32
	RequestTimeStamp  int64
	ResponseTimeStmap int64
	Md5               uint64
}

var (
	md5Salt = []byte{0x00, 0x21, 0x16, 0x01, 0x02, 0x05, 0x06, 0x18}
)

func ResponseBinaryEncode(res Response) []byte {
	b := make([]byte, 50)
	binary.LittleEndian.PutUint64(b[0:], uint64(res.OpID))
	binary.LittleEndian.PutUint16(b[8:], uint16(res.EventID))
	binary.LittleEndian.PutUint32(b[10:], uint32(res.ClientPublicIP))
	binary.LittleEndian.PutUint32(b[14:], uint32(res.ClientPublicPort))
	binary.LittleEndian.PutUint32(b[18:], uint32(res.ServerOtherIP))
	binary.LittleEndian.PutUint32(b[22:], uint32(res.ServerOtherPort))
	binary.LittleEndian.PutUint64(b[26:], uint64(res.RequestTimeStamp))
	binary.LittleEndian.PutUint64(b[34:], uint64(res.ResponseTimeStmap))
	binary.LittleEndian.PutUint64(b[42:], uint64(res.Md5))
	return b
}
func ResponseBinaryDecode(b []byte) Response {

	var res Response
	res.OpID = int64(binary.LittleEndian.Uint64(b[0:]))
	res.EventID = Event(binary.LittleEndian.Uint16(b[8:]))
	res.ClientPublicIP = uint32(binary.LittleEndian.Uint32(b[10:]))
	res.ClientPublicPort = int32(binary.LittleEndian.Uint32(b[14:]))
	res.ServerOtherIP = uint32(binary.LittleEndian.Uint32(b[18:]))
	res.ServerOtherPort = int32(binary.LittleEndian.Uint32(b[22:]))
	res.RequestTimeStamp = int64(binary.LittleEndian.Uint64(b[26:]))
	res.ResponseTimeStmap = int64(binary.LittleEndian.Uint64(b[34:]))
	res.Md5 = uint64(binary.LittleEndian.Uint64(b[42:]))
	return res
}

func (res *Response) CalMd5() {
	b := ResponseBinaryEncode(*res)
	md5Bytes := md5.Sum(append(b[:42], md5Salt...))
	res.Md5 = binary.LittleEndian.Uint64(md5Bytes[:8])
}

func (res Response) CheckMd5isAvaible() bool {
	b := ResponseBinaryEncode(res)
	md5Bytes := md5.Sum(append(b[:42], md5Salt...))
	md5uint64 := binary.LittleEndian.Uint64(md5Bytes[:8])
	return md5uint64 == res.Md5
}

func Ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

