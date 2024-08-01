package storage

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/mdigger/translit"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}
type db struct {
	sync.RWMutex
	logger     logger
	collection *mongo.Collection
}

func generateIntegerFromObjectId(id primitive.ObjectID) int {
	lastBytes := id[len(id)-3:]
	return int(lastBytes[0])<<16 | int(lastBytes[1])<<8 | int(lastBytes[2])
}
func generateTranslitName(name string, id primitive.ObjectID) string {
	rxSpaces := regexp.MustCompile(`\s+`)
	reg := regexp.MustCompile(`[^\p{L}\s]`)
	return fmt.Sprintf("%s-%d", strings.ReplaceAll(strings.TrimSpace(rxSpaces.ReplaceAllString(translit.Ru(reg.ReplaceAllString(strings.ToLower(strings.ReplaceAll(name, "-", " ")), "")), " ")), " ", "-"), generateIntegerFromObjectId(id))
}
func NewStorage(storage *mongo.Database, collection string, logger logger) *db {
	db := &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
	return db
}
