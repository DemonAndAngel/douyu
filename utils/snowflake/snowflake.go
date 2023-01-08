package snowflake

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1584979200000 // 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp >= now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := int64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func GetWorkerId() (int64, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return int64(ip[3]), nil
}

var worker *Worker

func init() {
	if worker == nil {
		workerId, _ := GetWorkerId()
		worker, _ = NewWorker(workerId)
	}
}

func GenId() int64 {
	return worker.GetId()
}

func GenStringId() string {
	return fmt.Sprint(GenId())
}

// 生成订单号
func GenOrderNo() string {
	cstZone := time.FixedZone("CST", 8*3600) // 东八
	now := time.Now().In(cstZone)
	return now.Format("20060102150405") + GenStringId()
}

// 微信支付订单生成
func GenOrderNoWithWechatPay() string {
	cstZone := time.FixedZone("CST", 8*3600) // 东八
	now := time.Now().In(cstZone)

	snowString := GenStringId()
	nowStr := now.Format("200601021504")

	if len(snowString)+len(nowStr) == 30 {
		return nowStr + "00" + snowString
	}

	if len(snowString)+len(nowStr) == 31 {
		return nowStr + "0" + snowString
	}

	return now.Format("20060102150405") + GenStringId()
}
