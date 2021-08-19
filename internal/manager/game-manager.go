package manager

import (
	"context"
	"sync"

	"NintendoCenter/game-collection/internal/protos"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type GameManager struct {
	collection *mgo.Collection
	mu         sync.RWMutex
}

const collectionName = "games"

func NewGameManager(db *mgo.Database) (*GameManager, error) {
	collection := db.C(collectionName)
	idxTable := mgo.Index{
		Key:        []string{"title"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := collection.EnsureIndex(idxTable); err != nil {
		return nil, err
	}

	return &GameManager{
		collection: collection,
		mu: sync.RWMutex{},
	}, nil
}

func (m *GameManager) SaveGame(game *protos.Game) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.collection.Insert(game)
}

func (m *GameManager) UpdateGame(id string, game *protos.Game) error {
	offerMap := make(map[protos.Shop]*protos.Offer, len(game.Offers))

	m.mu.RLock()
	if existed, _ := m.Find(id); existed != nil {
		for _, offer := range existed.Offers {
			offerMap[offer.Shop] = offer
		}
	}
	m.mu.RUnlock()

	for _, offer := range game.Offers {
		offerMap[offer.Shop] = offer
	}

	game.Offers = make([]*protos.Offer, 0, len(offerMap))
	for _, offer := range offerMap {
		game.Offers = append(game.Offers, offer)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.collection.Update(bson.M{"id": id}, game)
}

func (m *GameManager) Find(id string) (*protos.Game, error) {
	var game protos.Game
	m.mu.RLock()
	defer m.mu.RUnlock()
	if err := m.collection.Find(bson.M{"id": id}).One(&game); err != nil {
		return nil, err
	}
	return &game, nil
}

func (m *GameManager) SearchGames(ctx context.Context, filter *protos.FindGameRequest) (result []*protos.Game, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	query := bson.M{"title": bson.RegEx{Pattern: filter.Title, Options: "i"}}
	err = m.collection.Find(query).All(&result)
	if err != nil {
		return
	}

	return
}