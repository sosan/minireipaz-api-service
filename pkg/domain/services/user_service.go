package services

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/brokerclient"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
)

type UserServiceInterface interface {
	SynUser(user *models.Users) (created, exist bool)
}

var _ UserServiceInterface = (*UserService)(nil)

type UserService struct {
	repoHTTP   *httpclient.UserRepository
	repoRedis  *redisclient.UserRedisRepository
	repoBroker *brokerclient.UserKafkaRepository
}

func NewUserService(newRepoHTTP *httpclient.UserRepository, newRepoRedis *redisclient.UserRedisRepository, newRepoBroker *brokerclient.UserKafkaRepository) *UserService {
	return &UserService{
		repoHTTP:   newRepoHTTP,
		repoRedis:  newRepoRedis,
		repoBroker: newRepoBroker,
	}
}

func (u *UserService) SynUser(user *models.Users) (created, exist bool) {
	exist, err := u.repoRedis.CheckUserExist(user)
	if err != nil {
		log.Printf("ERROR | Cannot access to repo redis %v", err)
		return false, false
	}
	if exist {
		return false, true
	}

	// new user
	exist, err = u.repoRedis.CheckLockExist(user)
	if err != nil {
		log.Printf("ERROR | Cannot access to repo redis %v", err)
		return false, false
	}
	if exist {
		return false, false
	}
	// not exist generate lock for about 20 seconds
	// and 10 retries
	locked, err := u.repoRedis.AddLock(user)
	if err != nil {
		// TODO: Dead letter
		log.Printf("ERROR | Needs to Added to dead letter %s", user.Sub)
	}

	if !locked {
		log.Printf("WARN | Cannot created lock for user %s", user.Sub)
	}

	setDefaults(user) // default roleID
	sended := u.repoBroker.CreateUser(user)

	// in case cannot get autoremoved
	u.repoRedis.RemoveLock(user)

	return sended, false
}

func setDefaults(user *models.Users) {
	if user.Status == 0 {
		user.Status = models.StatusActive
	}

	if user.RoleID == 0 {
		user.RoleID = generateDefaultUserRoleID()
	}
}

func generateDefaultUserRoleID() models.UserRoleID {
	return models.RoleUser
}
