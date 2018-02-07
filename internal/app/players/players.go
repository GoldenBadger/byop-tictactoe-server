package players

import (
	"sync"
)

type Players struct {
	ids    map[uint]bool
	mutex  sync.RWMutex
	nextID uint
}

func New() *Players {
	p := new(Players)
	p.ids = make(map[uint]bool)
	return p
}

func (p *Players) NewPlayer() uint {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.ids[p.nextID] = true
	ret := p.nextID
	p.nextID++
	return ret
}

func (p *Players) PlayerExists(player uint) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	_, ok := p.ids[player]
	return ok
}
