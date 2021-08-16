package mongo

import "github.com/globalsign/mgo"

func New(connectionUrl string) (*mgo.Database, error) {
	session, err := mgo.Dial(connectionUrl)
	if err != nil {
		return nil, err
	}

	return session.DB(""), nil
}