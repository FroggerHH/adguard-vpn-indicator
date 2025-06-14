package cli

import "strings"

// VpnStatus хранит разобранную информацию о состоянии AdGuard VPN.
type VpnStatus struct {
	IsConnected bool
	Location    string
}

// ParseStatus анализирует текстовый вывод команды `adguardvpn-cli status`
// и возвращает структурированное состояние.
func ParseStatus(output string) VpnStatus {
	defaultStatus := VpnStatus{IsConnected: false, Location: ""}
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return defaultStatus
	}

	firstLine := lines[0]

	if strings.HasPrefix(firstLine, "Connected to ") {
		endOfLocation := strings.Index(firstLine, " in ")
		if endOfLocation == -1 {
			return defaultStatus
		}
		location := firstLine[13:endOfLocation]
		return VpnStatus{IsConnected: true, Location: location}
	}

	if strings.HasPrefix(firstLine, "VPN is disconnected") {
		return defaultStatus
	}

	return defaultStatus
}
