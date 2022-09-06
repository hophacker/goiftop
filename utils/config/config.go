package config

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/fs714/goiftop/utils/log"
	"github.com/google/gopacket/pcap"
)

var IfaceListString string
var GroupListString string
var Engine string
var IsDecodeL4 bool
var PrintEnable bool
var PrintInterval int64
var WebHookEnable bool
var WebHookUrl string
var WebHookInterval int64
var WebHookPostTimeout int
var WebHookNodeId string
var WebHookNodeOamAddr string
var IsEnableHttpSrv bool
var HttpSrvAddr string
var HttpSrvPort string
var IsProfiling bool
var IsShowVersion bool

type NfLogConfig struct {
	IfaceName string
	GroupId   int
	Direction pcap.Direction
}

var IfaceList []string
var NflogConfigList []NfLogConfig

// GetOutboundInterface returns the outbound interface
func GetOutboundInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Name == "eth0" || strings.HasPrefix(iface.Name, "enp") || strings.HasPrefix(iface.Name, "ens") {
			return iface.Name, nil
		}
	}
	return "", errors.New("failed to find outbound interface")
}

func ParseIfaces() {
	for _, iface := range strings.Split(IfaceListString, ",") {
		IfaceList = append(IfaceList, strings.TrimSpace(iface))
	}
}

func ParseNflogConfig() (err error) {
	for _, gpString := range strings.Split(GroupListString, ",") {
		gp := strings.Split(strings.TrimSpace(gpString), ":")

		if len(gp) != 3 {
			err = errors.New("invalid interface, group id and direction list: " + GroupListString)
			log.Errorf(err.Error())
			return
		}

		iface := strings.TrimSpace(gp[0])
		var groupId int
		groupId, err = strconv.Atoi(strings.TrimSpace(gp[1]))
		if err != nil {
			err = errors.New("invalid interface, group id and direction list: " + GroupListString)
			log.Errorf(err.Error())
			return
		}

		var direction pcap.Direction
		if strings.ToLower(strings.TrimSpace(gp[2])) == "in" {
			direction = pcap.DirectionIn
		} else if strings.ToLower(strings.TrimSpace(gp[2])) == "out" {
			direction = pcap.DirectionOut
		} else {
			err = errors.New("invalid interface, group id and direction list: " + GroupListString)
			log.Errorf(err.Error())
			return
		}

		nflogConf := NfLogConfig{
			IfaceName: iface,
			GroupId:   groupId,
			Direction: direction,
		}

		IfaceList = append(IfaceList, iface)
		NflogConfigList = append(NflogConfigList, nflogConf)
	}

	return
}
