package adbkit

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	OKAY = "OKAY"
	FAIL = "FAIL"
	DATA = "DATA"
	SEND = "SEND"
	DONE = "DONE"
	RECV = "RECV"
	STAT = "STAT"
	LIST = "LIST"
	DENT = "DENT"
)

type Device struct {
	Serial      string `json:"serial"`
	State       string `json:"state"`
	Product     string `json:"product"`
	Model       string `json:"model"`
	Device      string `json:"device"`
	TransportId int    `json:"transport_id"`
}

// Devices 获取一个节点设备列表
func (c Client) Devices() ([]Device, error) {
	resp, err := c.Command("host:devices")
	if err != nil {
		return nil, err
	}
	if string(resp[0:4]) == OKAY {
		var devices []Device
		for _, line := range strings.Split(string(resp[8:]), "\n") {
			device := strings.Split(line, "\t")
			if len(device) > 1 {
				devices = append(devices, Device{Serial: device[0], State: device[1]})
			}
		}
		return devices, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb response: Fail")
	}
	return nil, errors.New("error response: " + string(resp))
}

// Devices 获取一个详细节点设备列表
func (c Client) Lists() ([]Device, error) {
	resp, err := c.Command("host:devices-l")
	if err != nil {
		return nil, err
	}
	if string(resp[0:4]) == OKAY {
		var devices []Device
		for _, line := range strings.Split(string(resp[8:]), "\n") {
			var device Device
			line := strings.Fields(strings.TrimSpace(line))
			if len(line) < 2 {
				continue
			}
			for i, item := range line {
				switch true {
				case i == 0:
					device.Serial = line[0]
				case i == 1:
					device.State = line[1]
				case strings.Contains(item, "product:"):
					product := strings.Split(item, ":")
					device.Product = product[1]
				case strings.Contains(item, "model:"):
					model := strings.Split(item, ":")
					device.Model = model[1]
				case strings.Contains(item, "device:"):
					dev := strings.Split(item, ":")
					device.Device = dev[1]
				case strings.Contains(item, "transport_id:"):
					transportId := strings.Split(item, ":")
					device.TransportId, _ = strconv.Atoi(transportId[1])
				}
			}
			devices = append(devices, device)
		}
		return devices, nil
	} else if string(resp[0:4]) == FAIL {
		return nil, errors.New("adb response: Fail")
	}
	return nil, errors.New("error response: " + string(resp))
}

// Connect 从一个节点连接一个设备
func (c Client) Connect(ip string, port int) (bool, error) {
	resp, err := c.Command(fmt.Sprintf("host:connect:#%s:#%d", ip, port))
	if err != nil {
		return false, err
	}
	if string(resp[0:4]) == OKAY {
		length, _ := strconv.Atoi(string(resp[4:8]))

		var res = strings.Trim(string(resp[8:8+length]), "\n")
		if strings.Contains(res, "failed to connect") || strings.Contains(res, "unable to connect to") {
			return false, errors.New("failed to connect device")
		}
		if strings.Contains(res, "already connected to") || strings.Contains(res, "connected to") {
			return true, nil
		}
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb response: Fail")
	}
	return false, errors.New("error response: " + string(resp))
}

// Disconnect 断开设备
func (c Client) Disconnect(serial string) (bool, error) {
	resp, err := c.Command(fmt.Sprintf("host:disconnect:#%s", serial))
	if err != nil {
		return false, err
	}
	if string(resp[0:4]) == OKAY {
		length, _ := strconv.Atoi(string(resp[4:8]))
		var res = strings.Trim(string(resp[8:8+length]), "\n")
		if strings.Contains(res, "No such device") {
			return false, errors.New("no such device")
		}
		return true, nil
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb response: Fail")
	}
	return false, errors.New("error response: " + string(resp))
}

// Kill 断开一个节点得所有设备
func (c Client) Kill() (bool, error) {
	resp, err := c.Command("host:kill")
	if err != nil {
		return false, err
	}
	if string(resp[0:4]) == OKAY {
		return true, nil
	} else if string(resp[0:4]) == FAIL {
		return false, errors.New("adb response: Fail")
	}
	return false, errors.New("error response: " + string(resp))
}

type TrackDevice func(devices []Device, err error)

// TrackDevices  回调 设备连接变化
//  err := adbkit.New("127.0.0.1", 5037).TrackDevices(func(devices []adbkit.Device, err error) {
//		fmt.Println(err)
//		fmt.Println(devices)
//	})
func (c Client) TrackDevices(callback TrackDevice) error {
	return c.Callback("host:track-devices", func(buf []byte, err error) {
		if string(buf) == OKAY {
			return
		}
		if string(buf) == FAIL {
			callback(nil, errors.New("adb response: Fail"))
		}
		var devices []Device
		for _, line := range strings.Split(string(buf), "\n") {
			device := strings.Split(line, "\t")
			if len(device) > 1 {
				devices = append(devices, Device{Serial: device[0], State: device[1]})
			}
		}
		callback(devices, nil)
	})
}
