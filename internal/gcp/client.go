package gcp

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type GoogleAdminService struct {
	service *admin.Service
}

func Service(ServiceAccount string, subjectEmail string) (*GoogleAdminService, error) {

	creds, err := os.ReadFile(ServiceAccount)
	if err != nil {
		logrus.WithField("service", "gcp").Error("Falha ao ler o arquivo de credenciais", err)
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(creds, admin.AdminDirectoryGroupMemberScope)
	config.Subject = subjectEmail
	if err != nil {
		logrus.WithField("service", "gcp").Error("Falha ao criar configuração OAuth2", err)
		return nil, err
	}

	ctx := context.Background()
	service, err := admin.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		logrus.WithField("service", "gcp").Error("Não foi possível recuperar o diretório Cliente", err)
		return nil, err
	}

	return &GoogleAdminService{
		service: service,
	}, nil
}

func (s *GoogleAdminService) GetMembers(groupKey string) ([]string, error) {

	members, err := s.service.Members.List(groupKey).Do()
	if err != nil {
		logrus.WithField("service", "gcp").Error("Erro ao listar os membros", err)
		return nil, err
	}

	users := []string{}
	for _, member := range members.Members {
		users = append(users, member.Email)
		logrus.WithField("service", "gcp").Info(groupKey, ": ", member.Email)
	}
	return users, nil
}

func (s *GoogleAdminService) InsertMember(groupKey string, user string) {

	result, err := s.service.Members.Insert(groupKey, &admin.Member{Email: user}).Do()
	if err != nil {
		logrus.WithField("service", "gcp").Errorf("Falha ao adicionar %v", err)
	}
	logrus.WithField("service", "gcp").Info("Adicionado: ", result.Email)
}

func (s *GoogleAdminService) DeleteMember(groupKey string, user string) {

	err := s.service.Members.Delete(groupKey, user).Do()
	if err != nil {
		logrus.WithField("service", "gcp").Errorf("Falha ao remover %v", err)
	}
	logrus.WithField("service", "gcp").Info("Removido: ", user)
}
