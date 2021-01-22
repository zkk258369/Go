package server

import (
	"encoding/json"
	"fmt"
	reuse "github.com/libp2p/go-reuseport"
	"go.uber.org/zap"
	"net"
	"strconv"
	"time"

	"natdetect/common"
)

var serverLog *zap.SugaredLogger = nil

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr,
	opID int64, eventID common.Event,
	conn2 *net.UDPConn, requestTimeStamp int64,
) {
	serverOtherIP, serverOtherPortStr, err := net.SplitHostPort(conn2.LocalAddr().String())
	if err != nil {
		serverLog.Info("[server] [sendResponse] SplitHostPort filed, connAddr: ",
			conn.LocalAddr().String(), " ERROR: ", err)
		return
	}

	serverOtherPort, err := strconv.Atoi(serverOtherPortStr)
	if err != nil {
		serverLog.Info("[server] [sendResponse] Atoi file, serverOtherPortStr",
			serverOtherPortStr, " ERROR: ", err)
		return
	}

	for i := 0; i < 2; i++ {

		res := common.Response{
			opID,
			eventID,
			common.Ip2int(addr.IP),
			int32(addr.Port),
			common.Ip2int(net.ParseIP(serverOtherIP)),
			int32(serverOtherPort),
			requestTimeStamp,
			time.Now().Unix(),
			0,
		}
		res.CalMd5()

		buf, err := json.Marshal(res)
		if err != nil {
			serverLog.Info("[server] [sendResponse] Marshal Response filed, Response = ",
				res," ERROR: ", err)
		}

		n, err := conn.WriteToUDP(buf, addr)
		if err != nil {
			serverLog.Error("[server] [sendResponse] WriteToUDP filed, write ",
				n, "bytes"," ERROR:", err)
		}
	}
}

func serverRead(conn11, conn12, conn21, conn22 *net.UDPConn) {
	for {
		buf := make([]byte, 1024)
		req := &common.Request{}
		readn, remoteAddr, err := conn11.ReadFromUDP(buf)
		if err != nil {
			serverLog.Error("[server] [serverRead] ReadFromUDP filed, read ",
				readn, "bytes "," ERROR: ",err)
			continue
		}
		err = json.Unmarshal(buf[:readn], req)
		if err != nil {
			serverLog.Error("[server] [serverRead] Unmarshal Request filed, buf = ",
				buf," ERROR: ",err)
			continue
		}

		serverLog.Info(req.EventID, "   server:", conn11.LocalAddr().String(), "<-client:", remoteAddr.String())
		if req.EventID == common.EVENT_CLINET_1 {
			go sendResponse(conn11, remoteAddr, req.OpID, common.EVENT_SERVER_2, conn22, req.TimeStamp)
			go sendResponse(conn21, remoteAddr, req.OpID, common.EVENT_SERVER_3, conn22, req.TimeStamp)
			go sendResponse(conn12, remoteAddr, req.OpID, common.EVENT_SERVER_4, conn22, req.TimeStamp)
		} else if req.EventID == common.EVENT_CLIENT_5 {
			go sendResponse(conn11, remoteAddr, req.OpID, common.EVENT_SERVER_7, conn22, req.TimeStamp)
		} else if req.EventID == common.EVENT_CLIENT_6 {
			go sendResponse(conn11, remoteAddr, req.OpID, common.EVENT_SERVER_8, conn22, req.TimeStamp)
		}
		serverLog.Info(req)
	}
}
func WarpRecoverForever(conn11, conn12, conn21, conn22 *net.UDPConn) {
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					serverLog.Info("ERR", err)
				}
			}()
			serverRead(conn11, conn12, conn21, conn22)
		}()
	}
}

func ServerRun(ip1, ip2, ip3, ip4, ip5, ip6 string, port1, port2, port3, port4 int) error {
	serverLog = common.LogInit("./nat_detect.log", 7, 500)
	serverLog.Info("--------------------------------------------------")
	serverLog.Info("--------------new run Server----------------------")
	serverLog.Info("--------------------------------------------------")

	var ipgroup [][]string = [][]string{{ip1, ip2}, {ip3, ip4}, {ip5, ip6}, {ip2, ip1}, {ip4, ip3}, {ip6, ip5}}
	for i := 0; i < 6; i++ {
		if ipgroup[i][0] == "" || ipgroup[i][1] == "" {
			continue
		}

		addrList := []string{
			ipgroup[i][0] + ":" + strconv.Itoa(port1),
			ipgroup[i][0] + ":" + strconv.Itoa(port2),
			ipgroup[i][1] + ":" + strconv.Itoa(port3),
			ipgroup[i][1] + ":" + strconv.Itoa(port4),
		}
		conns := make([]*net.UDPConn, 4)
		for i, addr := range addrList {
			fmt.Println(addr)
			conn, err := reuse.ListenPacket("udp", addr)
			if err != nil {
				return err
			}
			u, _ := conn.(*net.UDPConn)
			conns[i] = u
			defer conns[i].Close()
		}

		WarpRecoverForever(conns[0], conns[1], conns[2], conns[3])
		WarpRecoverForever(conns[1], conns[0], conns[2], conns[3])
		WarpRecoverForever(conns[2], conns[3], conns[0], conns[1])
		WarpRecoverForever(conns[3], conns[2], conns[0], conns[1])

	}
	for {
		time.Sleep(time.Second * 60)
	}
	return nil
}


/*
func main() {
	ip1 := "127.0.0.1"
	//ip2 := "10.9.130.62"
	ip2 := "10.17.1.150"
	ip1port1 := 9000
	ip1port2 := 9001
	ip2port1 := 9002
	ip2port2 := 90031
	ServerRun(ip1, ip2, "", "","","",ip1port1, ip1port2, ip2port1, ip2port2)

	fmt.Printf("%v\n",ip1)
}
*/
