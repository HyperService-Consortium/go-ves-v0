package helper

import (
	"errors"
	"fmt"
	"net"
)

func DecodeIP(ip []byte) (string, error) {
	if len(ip) == 6 {
		return fmt.Sprintf("%v.%v.%v.%v:%v", ip[0], ip[1], ip[2], ip[3], (uint16(ip[4])<<8)|uint16(ip[5])), nil
	} else if len(ip) == 18 {
		return fmt.Sprintf("[%v]:%v", net.IP(ip[0:16]), (uint16(ip[16])<<8)|uint16(ip[17])), nil
	} else {
		return "", errors.New("invalid length")
	}
}
