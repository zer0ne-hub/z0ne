package core

import (
	"z0ne/internal/recon"

	"github.com/fatih/color"
)

func RunRecon(target string) {
	TargetType := detectTargetType(target)

	switch TargetType {
	case IP:
		color.Green("IP address detected: %s", target)
		recon.RunNaabu(target, "", "")
		recon.RunSubfinder(target)
		recon.RunDnsX(target)
		recon.RunHttpX(target)
		recon.RunKatana(target)
		recon.RunNuclei(target)
	case DOMAIN:
		color.Green("Domain detected: %s", target)
		recon.RunNaabu(target, "", "")
		recon.RunSubfinder(target)
		recon.RunDnsX(target)
		recon.RunHttpX(target)
		recon.RunKatana(target)
		recon.RunNuclei(target)
	case URL:
		color.Green("URL detected: %s", target)
	case FILE:
		color.Green("FilePath detected: %s", target)
		color.Red("File path is not supported yet")
	default:
		color.Red("Unknown target type: %s", target)
		color.Red("Please specify a valid target. Supported types are: IP, Domain, URL, FilePath")
	}
}
