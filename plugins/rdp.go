package plugins

//
//import (
//	"errors"
//	"fmt"
//	"github.com/tomatome/grdp/core"
//	"github.com/tomatome/grdp/glog"
//	"github.com/tomatome/grdp/protocol/nla"
//	"github.com/tomatome/grdp/protocol/pdu"
//	"github.com/tomatome/grdp/protocol/rfb"
//	"github.com/tomatome/grdp/protocol/sec"
//	"github.com/tomatome/grdp/protocol/t125"
//	"github.com/tomatome/grdp/protocol/tpkt"
//	"github.com/tomatome/grdp/protocol/x224"
//	"log"
//	"net"
//	"os"
//	"sync"
//	"time"
//	"x-crack/models"
//	"x-crack/vars"
//)
//
//func RdpConn(ip, domain, user, password string, port int, timeout time.Duration) (bool, error) {
//	target := fmt.Sprintf("%s:%d", ip, port)
//	g := NewClient(target, glog.NONE)
//	err := g.Login(domain, user, password, timeout)
//
//	if err == nil {
//		return true, nil
//	}
//
//	return false, err
//}
//
//type Client struct {
//	Host string // ip:port
//	tpkt *tpkt.TPKT
//	x224 *x224.X224
//	mcs  *t125.MCSClient
//	sec  *sec.Client
//	pdu  *pdu.Client
//	vnc  *rfb.RFB
//}
//
//func NewClient(host string, logLevel glog.LEVEL) *Client {
//	glog.SetLevel(logLevel)
//	logger := log.New(os.Stdout, "", 0)
//	glog.SetLogger(logger)
//	return &Client{
//		Host: host,
//	}
//}
//
//func (g *Client) Login(domain, user, pwd string, timeout time.Duration) error {
//	conn, err := net.DialTimeout("tcp", g.Host, timeout)
//	defer func() {
//		if conn != nil {
//			conn.Close()
//		}
//	}()
//	if err != nil {
//		return fmt.Errorf("[dial err] %v", err)
//	}
//	glog.Info(conn.LocalAddr().String())
//
//	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
//	g.x224 = x224.New(g.tpkt)
//	g.mcs = t125.NewMCSClient(g.x224)
//	g.sec = sec.NewClient(g.mcs)
//	g.pdu = pdu.NewClient(g.sec)
//
//	g.sec.SetUser(user)
//	g.sec.SetPwd(pwd)
//	g.sec.SetDomain(domain)
//	//g.sec.SetClientAutoReconnect()
//
//	g.tpkt.SetFastPathListener(g.sec)
//	g.sec.SetFastPathListener(g.pdu)
//	g.pdu.SetFastPathSender(g.tpkt)
//
//	//开启下面的任一选项,
//	//g.x224.SetRequestedProtocol(x224.PROTOCOL_SSL)
//	//g.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)
//
//	err = g.x224.Connect()
//	if err != nil {
//		return fmt.Errorf("[x224 connect err] %v", err)
//	}
//	glog.Info("wait connect ok")
//	wg := &sync.WaitGroup{}
//	breakFlag := false
//	wg.Add(1)
//
//	g.pdu.On("error", func(e error) {
//		err = e
//		glog.Error("error", e)
//		g.pdu.Emit("done")
//	})
//	g.pdu.On("close", func() {
//		err = errors.New("close")
//		glog.Info("on close")
//		g.pdu.Emit("done")
//	})
//	g.pdu.On("success", func() {
//		err = nil
//		glog.Info("on success")
//		g.pdu.Emit("done")
//	})
//	g.pdu.On("ready", func() {
//		glog.Info("on ready")
//		g.pdu.Emit("done")
//	})
//	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
//		glog.Info("on update:", rectangles)
//	})
//	g.pdu.On("done", func() {
//		if breakFlag == false {
//			breakFlag = true
//			wg.Done()
//		}
//	})
//	wg.Wait()
//	return err
//}
//
//func ScanRDP(service models.Service) (err error, result models.ScanResult) {
//	result.Service = service
//	IpPort := fmt.Sprintf("rdp://%v:%d", service.Ip, service.Port)
//	fmt.Println(IpPort)
//	flag, err := RdpConn(service.Ip, "DESKTOP-Q1Test", service.Username, service.Password, service.Port, vars.TimeOut)
//	if flag == true && err == nil {
//		result.Result = true
//
//	} else {
//		fmt.Println(fmt.Errorf("[x224 connect err] %v", err))
//	}
//	return err, result
//}
