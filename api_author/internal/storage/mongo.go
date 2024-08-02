package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mdigger/translit"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_author/pkg/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	defer db.seedAuthors()
	return db
}
func (d *db) seedAuthors() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if count, _ := d.collection.CountDocuments(ctx, bson.D{}); count > 0 {
		d.logger.Infof("there are %d authors in base, seeding canceled", count)
		return
	}
	authors := []Author{
		{
			Name: "Сергей Есенин",
			About: `Русский поэт и писатель. Одна из крупнейших личностей Серебряного века. Представитель новокрестьянской поэзии и лирики, а в более позднем периоде творчества — имажинизма.

В разные периоды творчества в его стихотворениях находили отражение социал-демократические идеи, образы революции и Родины, деревни и природы, любви и поиска счастья.`,
			ProfilePicture: "https://priazovskoe.ru/media/resized/c2WD9HOjUq_QkATEUfUD9jPbDGQFlq84FK_a-GlEHYA/rs:fit:1024:768/aHR0cHM6Ly9wcmlh/em92c2tvZS5ydS9t/ZWRpYS9wcm9qZWN0/X21vXzM5OC8yZi9i/Yi81NS9mZC8yMy9i/Yi8yODcyNTYxMS5q/cGc.jpg",
		},
		{
			Name: "Александр Пушкин",
			About: `Русский поэт, драматург и прозаик, заложивший основы русского реалистического направления, литературный критик и теоретик литературы, историк, публицист, журналист, редактор и издатель.

Один из самых авторитетных литературных деятелей первой трети XIX века. Ещё при жизни Пушкина сложилась его репутация величайшего национального русского поэта. Пушкин рассматривается как основоположник современного русского литературного языка.`,
			ProfilePicture: "https://www.prlib.ru/sites/default/files/book_preview/6cbf6f02-880a-4569-8859-9198cf2909eb/234829_doc1_6D69F2AA-93BE-47D2-ACDD-D327150EBF5A.jpg",
		},
		{
			Name: "Лев Толстой",
			About: `Один из наиболее известных русских писателей и мыслителей, один из величайших в мире писателей-романистов.

Участник обороны Севастополя. Просветитель, публицист, религиозный мыслитель, его авторитетное мнение послужило причиной возникновения нового религиозно-нравственного течения — толстовства. За свои взгляды был отлучён от РПЦ. Член-корреспондент Императорской Академии наук (1873), почётный академик по разряду изящной словесности (1900). Был номинирован на Нобелевскую премию по литературе (1902, 1903, 1904, 1905). Впоследствии отказался от дальнейших номинаций. Классик мировой литературы.

Писатель, ещё при жизни признанный главой русской литературы. Творчество Льва Толстого ознаменовало новый этап в русском и мировом реализме, выступив мостом между классическим романом XIX века и литературой XX века. Лев Толстой оказал сильное влияние на эволюцию европейского гуманизма, а также на развитие реалистических традиций в мировой литературе. Произведения Льва Толстого многократно экранизировались и инсценировались; его пьесы ставились на сценах всего мира. Лев Толстой был самым издаваемым в СССР писателем за 1918—1986 годы: общий тираж 3199 изданий составил 436,261 млн экземпляров.`,
			ProfilePicture: "https://globus-nsk.ru/upload/_thumbs/STRONG_284x426_f1bbe82dc2283b772f566a6a19c71f59.jpg",
		},
	}

	for _, a := range authors {
		if _, err := d.CreateAuthor(ctx, &a); err != nil {
			d.logger.Fatalf("error seeding author %v: %v", a, err)
		}
	}
	d.logger.Infof("seeded %d authors", len(authors))
}
func (d *db) CreateAuthor(ctx context.Context, author *Author) (*Author, error) {
	author.Id = primitive.NewObjectID()
	author.TranslitName = generateTranslitName(author.Name, author.Id)

	result, err := d.collection.InsertOne(ctx, author)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok || id != author.Id {
		return nil, status.Error(codes.Internal, "error retrieving author id")
	}

	return author, nil
}
func (d *db) GetAuthors(ctx context.Context, id []primitive.ObjectID, translit []string) ([]*Author, error) {
	authors := make([]*Author, 0)

	if len(id) > 0 {
		result, err := d.collection.Find(ctx, bson.M{"_id": bson.D{{Key: "$in", Value: id}}})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		var temp []*Author
		err = result.All(ctx, &temp)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		authors = append(authors, temp...)
	}
	if len(translit) > 0 {
		result, err := d.collection.Find(ctx, bson.M{"translit": bson.D{{Key: "$in", Value: translit}}})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		var temp []*Author
		err = result.All(ctx, &temp)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		authors = append(authors, temp...)
	}
	if len(id) == 0 && len(translit) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no id or translit name argument presented")
	}
	if len(authors) == 0 {
		var str string
		for _, i := range id {
			str += i.Hex() + ","
		}
		status, _ := status.New(codes.NotFound, "no authors found").WithDetails(&shared_pb.ErrorDetail{
			Field:       "id",
			Struct:      "authors_pb.GetAuthorsRequest",
			Actualvalue: strings.Trim(str, ","),
		}, &shared_pb.ErrorDetail{
			Field:       "translit",
			Struct:      "authors_pb.GetAuthorsRequest",
			Actualvalue: strings.Join(translit, ","),
		})
		return nil, status.Err()
	}

	return authors, nil
}
