package manager

import (
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
	return m.collection.Insert(game)
}

func (m *GameManager) UpdateGame(id string, game *protos.Game) error {
	offerMap := make(map[protos.Shop]*protos.Offer, len(game.Offers))
	if existed, _ := m.Find(id); existed != nil {
		for _, offer := range existed.Offers {
			offerMap[offer.Shop] = offer
		}
	}

	for _, offer := range game.Offers {
		offerMap[offer.Shop] = offer
	}

	game.Offers = make([]*protos.Offer, 0, len(offerMap))
	for _, offer := range offerMap {
		game.Offers = append(game.Offers, offer)
	}

	return m.collection.Update(bson.M{"_id": id}, game)
}

func (m *GameManager) Find(id string) (*protos.Game, error) {
	var game protos.Game
	if err := m.collection.Find(bson.M{"_id": id}).One(&game); err != nil {
		return nil, err
	}
	return &game, nil
}
