package pdkshell

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// PdkCmd struct
type PdkCmd struct {
	exe string
}

// New create new session
func New() *PdkCmd {
	exe, _ := exec.LookPath("cmd.exe")
	return &PdkCmd{
		exe: exe,
	}
}

func (p *PdkCmd) Execute(args []string) {
	var pdkInstallDir string
	if runtime.GOOS == "windows" {
		pdkInstallDir = "C:/PROGRA~1/PUPPET~1/DEVELO~1"
	} else {
		pdkInstallDir = "/opt/puppetlabs/pdk"
	}
	pdkRubyVer := "2.4.10"

	rubyexe := filepath.Join(pdkInstallDir, "private", "ruby", pdkRubyVer, "bin", "ruby")
	pdkexe := filepath.Join(pdkInstallDir, "private", "ruby", pdkRubyVer, "bin", "pdk")
	certDir := filepath.Join(pdkInstallDir, "ssl", "certs")
	certPem := filepath.Join(pdkInstallDir, "ssl", "cert.pem")

	if runtime.GOOS == "windows" {
		args = append([]string{"/c", rubyexe, "-S", "--", pdkexe}, args...)
	} else {
		args = append([]string{"/bin/sh", rubyexe, "-S", "--", pdkexe}, args...)
		log.Printf("argS: %v", args)
	}

	cmd := exec.Command(p.exe, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("SSL_CERT_DIR=%s", certDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("SSL_CERT_FILE=%s", certPem))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

}
