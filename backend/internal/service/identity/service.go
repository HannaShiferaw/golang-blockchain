package identity

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"coffee-consortium/backend/internal/domain"
	"coffee-consortium/backend/internal/pki"
)

type Repo interface {
	GetRoot(ctx context.Context) (certPEM, keyPEM string, found bool, err error)
	PutRoot(ctx context.Context, certPEM, keyPEM string) error
	PutIdentity(ctx context.Context, it domain.Identity) error
	ListIdentities(ctx context.Context) ([]domain.Identity, error)
}

type Service struct {
	mu         sync.RWMutex
	ca         *pki.RootCA
	identities map[string]domain.Identity
	repo       Repo
}

func NewService(repo Repo) (*Service, error) {
	var ca *pki.RootCA
	if repo != nil {
		certPEM, keyPEM, found, err := repo.GetRoot(context.Background())
		if err != nil {
			return nil, err
		}
		if found {
			loaded, err := pki.LoadRootCA(certPEM, keyPEM)
			if err != nil {
				return nil, err
			}
			ca = loaded
		}
	}
	if ca == nil {
		created, err := pki.NewRootCA("coffee-consortium-root-ca")
		if err != nil {
			return nil, err
		}
		ca = created
		if repo != nil {
			certPEM, _ := ca.CertPEM()
			keyPEM, _ := ca.PrivateKeyPEM()
			if err := repo.PutRoot(context.Background(), certPEM, keyPEM); err != nil {
				return nil, err
			}
		}
	}

	s := &Service{
		ca:         ca,
		identities: map[string]domain.Identity{},
		repo:       repo,
	}

	if repo != nil {
		items, err := repo.ListIdentities(context.Background())
		if err != nil {
			return nil, err
		}
		for _, it := range items {
			s.identities[it.ID] = it
		}
	}

	return s, nil
}

func (s *Service) RootCertPEM() (string, error) {
	return s.ca.CertPEM()
}

func (s *Service) CreateIdentity(name string, role domain.Role) (domain.Identity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	issued, err := s.ca.IssueIdentity(id, name, role)
	if err != nil {
		return domain.Identity{}, err
	}

	s.identities[id] = issued.ID
	if s.repo != nil {
		if err := s.repo.PutIdentity(context.Background(), issued.ID); err != nil {
			return domain.Identity{}, err
		}
	}
	return issued.ID, nil
}

func (s *Service) GetIdentity(id string) (domain.Identity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	it, ok := s.identities[id]
	if !ok {
		return domain.Identity{}, fmt.Errorf("identity not found: %s", id)
	}
	return it, nil
}

func (s *Service) ListIdentities() []domain.Identity {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]domain.Identity, 0, len(s.identities))
	for _, v := range s.identities {
		cp := v
		cp.PrivateKeyPEM = ""
		out = append(out, cp)
	}
	return out
}

func (s *Service) SeedDefaults() error {
	s.mu.RLock()
	empty := len(s.identities) == 0
	s.mu.RUnlock()
	if !empty {
		return nil
	}

	if _, err := s.CreateIdentity("exporter1", domain.RoleExporter); err != nil {
		return err
	}
	if _, err := s.CreateIdentity("buyer1", domain.RoleBuyer); err != nil {
		return err
	}
	if _, err := s.CreateIdentity("cbe1", domain.RoleBank); err != nil {
		return err
	}
	if _, err := s.CreateIdentity("customs1", domain.RoleCustoms); err != nil {
		return err
	}
	if _, err := s.CreateIdentity("shipper1", domain.RoleShipment); err != nil {
		return err
	}
	return nil
}

