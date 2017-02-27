package core

type Service interface {
	SetId(id uint)
	Id() uint
}
