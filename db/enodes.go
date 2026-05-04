package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethpandaops/dora/utils"
	"net"
	"net/url"
	"os"
	"strconv"
)

type Enodes struct {
	EnodesMap   map[uint64]*EnodeEntry
	EnodesMapIP map[string][]*EnodeEntry
	HostMap     map[string]*HostEntry
}

type EnodeEntry struct {
	index uint64
	Value string
	IP    string
}

type HostEntry struct {
	IP   string
	Host string
}

func InitEnodes() *Enodes {
	d := &Enodes{
		EnodesMap:   map[uint64]*EnodeEntry{},
		EnodesMapIP: map[string][]*EnodeEntry{},
		HostMap:     map[string]*HostEntry{},
	}
	d.HostMap["162.220.9.10"] = &HostEntry{Host: "Oak", IP: "162.220.9.10"}
	d.HostMap["66.45.229.2"] = &HostEntry{Host: "Pine", IP: "66.45.229.2"}
	d.HostMap["162.220.9.74"] = &HostEntry{Host: "Aspen", IP: "162.220.9.74"}
	d.HostMap["173.214.174.218"] = &HostEntry{Host: "Spruce", IP: "173.214.174.218"}
	d.HostMap["206.72.206.158"] = &HostEntry{Host: "Fir", IP: "206.72.206.158"}
	d.HostMap["64.20.55.118"] = &HostEntry{Host: "Cedar", IP: "64.20.55.118"}
	d.HostMap["161.129.71.58"] = &HostEntry{Host: "Cypress", IP: "161.129.71.58"}
	d.HostMap["66.45.234.170"] = &HostEntry{Host: "Trichet", IP: "66.45.234.170"}
	d.HostMap["64.20.61.126"] = &HostEntry{Host: "Draghi", IP: "64.20.61.126"}
	d.HostMap["205.209.106.50"] = &HostEntry{Host: "Lagarde", IP: "205.209.106.50"}
	d.HostMap["205.209.120.202"] = &HostEntry{Host: "Yellen", IP: "205.209.120.202"}
	d.HostMap["67.217.48.138"] = &HostEntry{Host: "Powell", IP: "67.217.48.138"}
	d.HostMap["205.209.110.94"] = &HostEntry{Host: "Maribor", IP: "205.209.110.94"}
	d.HostMap["67.217.61.66"] = &HostEntry{Host: "Run", IP: "67.217.61.66"}
	d.HostMap["74.50.75.74"] = &HostEntry{Host: "Greenspan", IP: "74.50.75.74"}
	d.HostMap["216.158.236.2"] = &HostEntry{Host: "Izola", IP: "216.158.236.2"}
	d.HostMap["64.20.46.82"] = &HostEntry{Host: "Koper", IP: "64.20.46.82"}
	d.HostMap["205.209.113.226"] = &HostEntry{Host: "Dubai", IP: "205.209.113.226"}
	d.HostMap["204.13.238.122"] = &HostEntry{Host: "Novomesto", IP: "204.13.238.122"}
	d.HostMap["216.219.80.14"] = &HostEntry{Host: "Mouse", IP: "216.219.80.14"}
	d.HostMap["74.50.76.102"] = &HostEntry{Host: "Bear", IP: "74.50.76.102"}
	d.HostMap["205.209.112.50"] = &HostEntry{Host: "Fox", IP: "205.209.112.50"}
	d.HostMap["198.96.94.94"] = &HostEntry{Host: "Deer", IP: "198.96.94.94"}
	d.HostMap["67.217.57.66"] = &HostEntry{Host: "Polar", IP: "67.217.57.66"}

	/*d.HostMap["67.217.61.66"] = &HostEntry{Host: "run.merapi.io", IP: "67.217.61.66"}
	d.HostMap["66.45.234.170"] = &HostEntry{Host: "trichet.merapi.io", IP: "66.45.234.170"}
	d.HostMap["74.50.76.102"] = &HostEntry{Host: "bear.ethpar.net", IP: "74.50.76.102"}
	d.HostMap["67.217.48.138"] = &HostEntry{Host: "powell.merapi.io", IP: "67.217.48.138"}
	d.HostMap["216.158.236.2"] = &HostEntry{Host: "izola.merapi.io", IP: "216.158.236.2"}
	d.HostMap["216.219.80.14"] = &HostEntry{Host: "mouse.ethpar.net", IP: "216.219.80.14"}
	d.HostMap["206.72.206.158"] = &HostEntry{Host: "fir.ethpar.net", IP: "206.72.206.158"}
	d.HostMap["162.220.9.10"] = &HostEntry{Host: "oak.ethpar.net", IP: "162.220.9.10"}
	d.HostMap["66.45.229.2"] = &HostEntry{Host: "pine.ethpar.net", IP: "66.45.229.2"}
	d.HostMap["205.209.113.226"] = &HostEntry{Host: "dubai.merapi.io", IP: "205.209.113.226"}
	d.HostMap["64.20.55.118"] = &HostEntry{Host: "pine.ethpar.net", IP: "64.20.55.118"}*/

	d.InitFromFile()
	return d
}

func (en Enodes) GetENode(index uint64) string {
	if en.EnodesMap[index] != nil {
		return en.EnodesMap[index].IP
	}
	return ""
}

func (en Enodes) InitFromFile() {

	isBCW := false
	if utils.Config.Monitor.Subject == "mainnet" {
		isBCW = true
	}

	data, err := os.ReadFile(utils.Config.Alert.EnodesDataDir)
	if err != nil {
		return
	}

	var rawMap map[string]string
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return
	}

	var enodes []EnodeEntry
	for id, val := range rawMap {
		u, err := url.Parse(val)
		if err != nil {
			continue
		}

		ip, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			continue
		}
		fmt.Printf("ID: %s\n", ip)

		ind, err := strconv.ParseUint(id, 10, 64)

		if err != nil {
			continue
		}

		if isBCW && (ind > 1566 && ind < 1688) {
			ip = "BCW"
		}

		if en.HostMap[ip] != nil {
			ip = en.HostMap[ip].Host
		}
		enode := EnodeEntry{
			index: ind,
			Value: val,
			IP:    ip,
		}
		en.EnodesMap[ind] = &enode
		if en.EnodesMapIP[enode.IP] == nil {
			en.EnodesMapIP[enode.IP] = []*EnodeEntry{}
		}
		en.EnodesMapIP[enode.IP] = append(en.EnodesMapIP[enode.IP], &enode)
		enodes = append(enodes, EnodeEntry{
			index: ind,
			Value: val,
			IP:    ip,
		})
	}

	/*for _, entry := range enodes {
		fmt.Printf("ID: %s, Enode: %s\n", entry.index, entry.Value)
	}*/
}
