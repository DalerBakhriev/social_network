package store

import "github.com/DalerBakhriev/social_network/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(int) (*model.User, error)
	Update(*model.User) error
	GetTopUsers(int) ([]*model.User, error)
	GetFriendsList(int) ([]*model.User, error)
	SendFriendRequest(int, int) error
	AcceptFriendRequest(int, int) error
}
