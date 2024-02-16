package repository

import (
	"github.com/elabosak233/pgshub/internal/model"
	"gorm.io/gorm"
)

type INatRepository interface {
	Insert(nat model.Nat) (n model.Nat, err error)
	FindByInstanceID(instanceIDs []uint) (nats []model.Nat, err error)
}

type NatRepository struct {
	Db *gorm.DB
}

func NewNatRepository(Db *gorm.DB) INatRepository {
	return &NatRepository{Db: Db}
}

func (t *NatRepository) Insert(nat model.Nat) (n model.Nat, err error) {
	result := t.Db.Table("nats").Create(&nat)
	return nat, result.Error
}

func (t *NatRepository) FindByInstanceID(instanceIDs []uint) (nats []model.Nat, err error) {
	result := t.Db.Table("nats").Where("instance_id IN ?", instanceIDs).Find(&nats)
	return nats, result.Error
}