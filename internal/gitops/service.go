package gitops

import (
	"fmt"
	"os"

	"gogitopsdeployer/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Service gerencia as operacoes Git.
type Service struct {
	cfg *config.Config
}

// NewService cria uma nova instancia do servico GitOps.
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// EnsureClone garante que o repositorio esta clonado no disco.
func (s *Service) EnsureClone() error {
	_, err := git.PlainOpen(s.cfg.LocalPath)
	if err == nil {
		return nil
	}

	if err != git.ErrRepositoryNotExists {
		return err
	}

	fmt.Printf("Clonando repositorio %s em %s...\n", s.cfg.RepoURL, s.cfg.LocalPath)
	_, err = git.PlainClone(s.cfg.LocalPath, false, &git.CloneOptions{
		URL:      s.cfg.RepoURL,
		Progress: os.Stdout,
	})
	return err
}

// CheckForUpdates verifica se ha novos commits no remoto.
func (s *Service) CheckForUpdates() (bool, string, error) {
	repo, err := git.PlainOpen(s.cfg.LocalPath)
	if err != nil {
		return false, "", err
	}

	// Fetch para atualizar referencias remotas
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return false, "", err
	}

	// Pega o Hash local atual (HEAD)
	head, err := repo.Head()
	if err != nil {
		return false, "", err
	}
	currentHash := head.Hash().String()

	// Pega a referencia remota (assumindo branch padrao master/main)
	// Em um sistema real, poderiamos configurar a branch.
	remoteRef, err := repo.Reference(plumbing.ReferenceName("refs/remotes/origin/master"), true)
	if err != nil {
		// Tenta 'main' se 'master' nao existir
		remoteRef, err = repo.Reference(plumbing.ReferenceName("refs/remotes/origin/main"), true)
		if err != nil {
			return false, "", fmt.Errorf("could not find remote branch: %v", err)
		}
	}

	remoteHash := remoteRef.Hash().String()

	if currentHash != remoteHash {
		// Se mudou, faz o pull (ou apenas move o HEAD se for stateless)
		// Para simplificar, vamos apenas imprimir a deteccao.
		return true, remoteHash, nil
	}

	return false, currentHash, nil
}

// UpdateLocal atualiza o repositorio local para o hash detectado.
func (s *Service) UpdateLocal() error {
	repo, err := git.PlainOpen(s.cfg.LocalPath)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
