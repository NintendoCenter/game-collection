package manager

import (
	"context"
	"sync"

	"NintendoCenter/game-collection/internal/protos"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameManager struct {
	collection *mongo.Collection
	mu         sync.RWMutex
}

const collectionName = "games"

func NewGameManager(db *mongo.Database) (*GameManager, error) {
	collection := db.Collection(collectionName)

	return &GameManager{
		collection: collection,
		mu: sync.RWMutex{},
	}, nil
}

func (m *GameManager) SaveGame(ctx context.Context, game *protos.Game) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := m.collection.InsertOne(ctx, game)
	return err
}

func (m *GameManager) UpdateGame(ctx context.Context, id string, game *protos.Game) error {
	offerMap := make(map[protos.Shop]*protos.Offer, len(game.Offers))

	m.mu.RLock()
	if existed, _ := m.Find(ctx, id); existed != nil {
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
	_, err := m.collection.UpdateOne(ctx, bson.M{"id": id}, game)
	return err
}

func (m *GameManager) Find(ctx context.Context, id string) (*protos.Game, error) {
	var game protos.Game
	m.mu.RLock()
	defer m.mu.RUnlock()
	if err := m.collection.FindOne(ctx, bson.M{"id": id}).Decode(&game); err != nil {
		return nil, err
	}
	return &game, nil
}

func (m *GameManager) SearchGames(ctx context.Context, filter *protos.FindGameRequest) (result []*protos.Game, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	query := bson.M{"title": bson.RegEx{Pattern: filter.Title, Options: "i"}}
	cursor, err := m.collection.Find(ctx, query)
	if err != nil {
		return
	}
	err = cursor.All(ctx, result)
	return
}