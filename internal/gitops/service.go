// Package gitops provides a simplified interface for Git operations,
// specifically for monitoring and Updating repositories.
package gitops

import (
	"fmt"
	"os"

	"github.com/ESousa97/gogitopsdeployer/internal/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Service manages the Git repository lifecycle and version detection.
// It wraps the go-git library to provide high-level operations.
type Service struct {
	cfg *config.Config
}

// NewService initializes a new [Service] with the provided [config.Config].
func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

// EnsureClone checks if the repository exists locally at the configured path
// and performs a clone if it is missing.
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

// CheckForUpdates fetches the latest references from the remote (origin)
// and compares the local HEAD with the remote target branch (master/main).
// It returns true and the new hash if an update is detected.
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

// UpdateLocal performs a pull operation to synchronize the local repository
// with the latest changes from the remote tracking branch.
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
