package game

import "github.com/Odyssey-Classic/server/internal/services/network"

type Player struct {
	client *network.Client
}

func (p *Player) Send(msg any) {

}
