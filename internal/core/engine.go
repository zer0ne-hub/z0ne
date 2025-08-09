package core

import (
	"fmt"
	"z0ne/internal/recon"
)

func RunRecon(target string) {
	TargetType := detectTargetType(target)

	switch TargetType {
	case IP:
		fmt.Println("IP address detected:", target)
		recon.RunNaabu(target, "", "")
		recon.RunSubfinder(target)
		recon.RunDnsX(target)
		recon.RunHttpX(target)
		recon.RunKatana(target)
		recon.RunNuclei(target)
	case DOMAIN:
		fmt.Println("Domain detected:", target)
		recon.RunNaabu(target, "", "")
		recon.RunSubfinder(target)
		recon.RunDnsX(target)
		recon.RunHttpX(target)
		recon.RunKatana(target)
		recon.RunNuclei(target)
	case URL:
		fmt.Println("URL detected:", target)
	case FILE:
		fmt.Println("FilePath detected:", target)
	default:
		fmt.Println("Unknown target type:", target)
		fmt.Println("Please specify a valid target. Supported types are: IP, Domain, URL, FilePath")
	}
}
