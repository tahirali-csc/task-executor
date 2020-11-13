package services

import (
	"github.com/task-executor/pkg/api"
	"github.com/task-executor/pkg/api-server/dbstore"
)

type RepoService struct{}

func (rs RepoService) FindByNamespaceAndName(namespace, name string) (*api.Repo, error) {
	selectStmt := `SELECT r.*, at.name as auth_type_name, st.name as secret_type_name 
	FROM repo r
	INNER JOIN auth_type at ON r.auth_type = at.id 
	INNER JOIN secret_type st ON r.secret_type = st.id 
	WHERE r.namespace=$1 AND r.name=$2`
	row := dbstore.DataSource.QueryRow(selectStmt, namespace, name)

	res := &api.Repo{}
	err := row.Scan(&res.Id, &res.Namespace, &res.Name, &res.SSHUrl, &res.HttpUrl,
		&res.AuthType.Id, &res.SecretType.Id, &res.SecretName, &res.CreatedTs, &res.UpdatedTs,
		&res.AuthType.Name, &res.SecretType.Name)

	return res, err
}
