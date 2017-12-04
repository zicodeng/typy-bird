package models

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ErrTypieBirdNotFound is returned when the typie bird can't be found
var ErrTypieBirdNotFound = errors.New("typie bird not found")

//MongoStore implements Store for MongoDB
type MongoStore struct {
	session *mgo.Session
	dbname  string
	colname string
}

//NewMongoStore constructs a new MongoStore
func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		panic("nil pointer passed for session")
	}
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
	}
}

func (s *MongoStore) GetTopScores() ([]*TypieBird, error)  {
	topScores := make([]*TypieBird, 10)
	col := s.session.DB(s.dbname).C(s.colname)
	err := col.Find(nil).Limit(10).Sort("Record").All(topScores)
	if err != nil {
		return nil, fmt.Errorf("error retrieving top 10 typie birds")
	}
	return topScores, nil
}

//GetByID returns the User with the given ID
func (s *MongoStore) GetByID(id bson.ObjectId) (*TypieBird, error) {
	typieBird := &TypieBird{}
	col := s.session.DB(s.dbname).C(s.colname)
	err := col.FindId(id).One(typieBird)
	if err != nil {
		return nil, ErrTypieBirdNotFound
	}
	return typieBird, nil
}

//GetByUserName returns the User with the given Username
func (s *MongoStore) GetByUserName(username string) (*TypieBird, error) {
	typieBird := &TypieBird{}
	col := s.session.DB(s.dbname).C(s.colname)
	err := col.Find(bson.M{"username": username}).One(typieBird)
	if err != nil {
		return nil, ErrTypieBirdNotFound
	}
	return typieBird, nil
}

//InsertTypieBird inserts a new typie bird into the mongo store and returns the typie bird
func (s *MongoStore) InsertTypieBird(newTypie *TypieBird) (*TypieBird, error) {
	col := s.session.DB(s.dbname).C(s.colname)
	if err := col.Insert(newTypie); err != nil {
		return nil, fmt.Errorf("error inserting typie bird: %v", err)
	}
	return newTypie, nil
}

//InsertWords inserts a new word from a dictionary into the mongo store
func (s *MongoStore) InsertWords(word string) (string, error) {
	col := s.session.DB(s.dbname).C(s.colname)
	if err := col.Insert(word); err != nil {
		return "", fmt.Errorf("error inserting word: %v", err)
	}
	return word, nil
}

//Update applies UserUpdates to the given user ID
func (s *MongoStore) Update(typieBirdID bson.ObjectId, updates *Updates) error {
	typie := &TypieBird{}
	change := mgo.Change{
		Update: bson.M{"$set": updates},
    ReturnNew: true,
	}
	col := s.session.DB(s.dbname).C(s.colname)
	if _, err := col.FindId(typieBirdID).Apply(change, typie); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

//Delete deletes the typie bird with the given ID
func (s *MongoStore) Delete(typieBirdID bson.ObjectId) error {
	col := s.session.DB(s.dbname).C(s.colname)
	if err := col.RemoveId(typieBirdID); err != nil {
		return fmt.Errorf("error removing user: %v", err)
	}
	return nil
}
