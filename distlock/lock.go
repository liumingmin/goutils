package distlock

type DistLock interface {
	Lock(timeout int) bool
	Unlock()
}

//var gLocakKey = func() string {
//	localKey := os.Args[0]
//
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		return localKey
//	}
//
//	for _, address := range addrs {
//		// 检查ip地址判断是否回环地址
//		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				localKey += "~" + ipnet.IP.String()
//			}
//		}
//	}
//	return localKey
//}()
