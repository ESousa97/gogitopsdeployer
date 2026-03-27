// Package ssh provides capabilities for remote command execution
// using the SSH protocol, including support for primary and rollback commands.
package ssh

import (
	"fmt"
	"os"

	"gogitopsdeployer/internal/config"

	"golang.org/x/crypto/ssh"
)

// Service encapsulates the SSH client logic, handling authentication
// and command dispatch to the target remote host.
type Service struct {
	cfg *config.Config
}

// NewService returns a new [Service] initialized with the provided [config.Config].
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// RunCommands connects to the remote host using the configured SSH credentials
// and executes the global set of deployment commands sequentially.
// It returns the combined output of all commands or an error.
func (s *Service) RunCommands() (string, error) {
	if s.cfg.SSHHost == "" {
		return "", nil // SSH not configured
	}

	fmt.Printf("[SSH] Connecting to %s@%s...\n", s.cfg.SSHUser, s.cfg.SSHHost)

	// 1. Load private key
	key, err := os.ReadFile(s.cfg.SSHKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	// 2. Client setup
	config := &ssh.ClientConfig{
		User: s.cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, use real HostKey verification
	}

	// 3. Connect to host
	client, err := ssh.Dial("tcp", s.cfg.SSHHost+":22", config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	var combinedOutput string

	// 4. Run commands
	for _, cmd := range s.cfg.SSHCommands {
		fmt.Printf("[SSH] Executing: %s\n", cmd)
		
		session, err := client.NewSession()
		if err != nil {
			return combinedOutput, fmt.Errorf("failed to create session: %v", err)
		}

		output, err := session.CombinedOutput(cmd)
		combinedOutput += string(output) + "\n"
		
		if err != nil {
			fmt.Printf("[SSH] Command error: %v\n", err)
			fmt.Printf("[SSH] Error Output:\n%s\n", string(output))
			session.Close()
			return combinedOutput, err
		}

		fmt.Printf("[SSH] Output:\n%s\n", string(output))
		session.Close()
	}

	fmt.Println("[SSH] All commands executed successfully.")
	return combinedOutput, nil
}

// RunRollback connects to the remote host and executes the recovery command
// defined in the configuration (e.g., git checkout HEAD^).
func (s *Service) RunRollback() (string, error) {
	if s.cfg.SSHHost == "" || s.cfg.RollbackCommand == "" {
		return "", nil
	}

	fmt.Printf("[SSH] Starting ROLLBACK on %s...\n", s.cfg.SSHHost)

	// Load private key
	key, err := os.ReadFile(s.cfg.SSHKeyPath)
	if err != nil {
		return "", err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: s.cfg.SSHUser,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", s.cfg.SSHHost+":22", config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(s.cfg.RollbackCommand)
	return string(output), err
}
