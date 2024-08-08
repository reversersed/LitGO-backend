package storage

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=mongo.go -destination=mocks/mongo.go

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

func NewStorage(storage *mongo.Database, collection string, logger logger) *db {
	db := &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
	defer db.seedAdminAccount()
	return db
}
func (d *db) seedAdminAccount() {
	d.Lock()
	defer d.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, bson.M{"login": "admin"})
	if err := result.Err(); err != nil {
		d.logger.Info("starting seeding admin account...")
		pass, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		seedAdmin.Password = pass
		response, err := d.collection.InsertOne(ctx, seedAdmin)
		if err != nil {
			d.logger.Fatalf("cannot seed admin account: %v", err)
		}
		id, ok := response.InsertedID.(primitive.ObjectID)
		if !ok {
			d.logger.Fatalf("can't create id for admin document")
		}
		d.logger.Infof("admin account seeded with id %v", id.Hex())
		return
	}
	d.logger.Info("admin account exists. seed not executed")
}

func (d *db) FindById(ctx context.Context, id string) (*User, error) {
	d.RLock()
	defer d.RUnlock()

	primitive_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": primitive_id}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		d.logger.Warnf("error while fetching user from db: %v", err)
		return nil, err
	}
	var u User
	if err := result.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *db) FindByLogin(ctx context.Context, login string) (*User, error) {
	d.RLock()
	defer d.RUnlock()

	filter := bson.M{"login": login}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		d.logger.Warnf("error while fetching user from db: %v", err)
		return nil, err
	}
	var u User
	if err := result.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *db) FindByEmail(ctx context.Context, email string) (*User, error) {
	d.RLock()
	defer d.RUnlock()

	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := d.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		d.logger.Warnf("error while fetching user from db: %v", err)
		return nil, err
	}
	var u User
	if err := result.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
func (d *db) CreateUser(ctx context.Context, model *User) (primitive.ObjectID, error) {
	d.Lock()
	defer d.Unlock()

	contx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := d.collection.InsertOne(contx, model)
	if err != nil {
		d.logger.Warnf("error while user creation: %v", err)
		return primitive.ObjectID{}, status.Error(codes.Internal, err.Error())
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid, nil
	}
	d.logger.Warnf("cant get created user id: %v (%v)", oid.Hex(), oid)
	return primitive.ObjectID{}, status.Error(codes.Internal, "cant resolve user id")
}
