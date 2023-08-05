package models

import (
	"sync"
)

type RWPortGroupMap struct {
	Mu sync.Mutex
	Data map[int] string
}

type User struct {
	Username 	string
	Role		string
}
