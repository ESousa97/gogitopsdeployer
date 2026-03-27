package ssh

import (
	"fmt"
	"os"

	"gogitopsdeployer/internal/config"

	"golang.org/x/crypto/ssh"
)

// Service gerencia a execucao de comandos remotos via SSH.
type Service struct {
	cfg *config.Config
}

// NewService cria uma nova instancia do servico SSH.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// RunCommands conecta a VPS e executa a lista de comandos configurada.
func (s *Service) RunCommands() (string, error) {
	if s.cfg.SSHHost == "" {
		return "", nil // SSH nao configurado
	}

	fmt.Printf("[SSH] Conectando a %s@%s...\n", s.cfg.SSHUser, s.cfg.SSHHost)

	// 1. Carrega a chave privada
	key, err := os.ReadFile(s.cfg.SSHKeyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	// 2. Configura o cliente
	config := &ssh.ClientConfig{
		User: s.cfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Nota: Em producao, use verificacao de HostKey real
	}

	// 3. Conecta ao host
	client, err := ssh.Dial("tcp", s.cfg.SSHHost+":22", config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	var combinedOutput string

	// 4. Executa os comandos
	for _, cmd := range s.cfg.SSHCommands {
		fmt.Printf("[SSH] Executando: %s\n", cmd)
		
		session, err := client.NewSession()
		if err != nil {
			return combinedOutput, fmt.Errorf("failed to create session: %v", err)
		}

		output, err := session.CombinedOutput(cmd)
		combinedOutput += string(output) + "\n"
		
		if err != nil {
			fmt.Printf("[SSH] Erro no comando: %v\n", err)
			fmt.Printf("[SSH] Output de Erro:\n%s\n", string(output))
			session.Close()
			return combinedOutput, err
		}

		fmt.Printf("[SSH] Output:\n%s\n", string(output))
		session.Close()
	}

	fmt.Println("[SSH] Todos os comandos executados com sucesso.")
	return combinedOutput, nil
}

// RunRollback executa o comando de rollback configurado na VPS.
func (s *Service) RunRollback() (string, error) {
	if s.cfg.SSHHost == "" || s.cfg.RollbackCommand == "" {
		return "", nil
	}

	fmt.Printf("[SSH] Iniciando ROLLBACK em %s...\n", s.cfg.SSHHost)

	// Carrega a chave privada
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
