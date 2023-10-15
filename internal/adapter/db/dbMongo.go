package adapterDB

import (
	"context"
	T "shortlink/internal/apptype"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _ T.IDB = (*DBMongo)(nil)

type DBMongo struct {
	cfg  *T.CfgEnv
	log  T.ILog
	conn *mongo.Client
}

func NewDBMongo(cfg *T.CfgEnv, log T.ILog) DBMongo {
	return DBMongo{
		cfg: cfg,
		log: log,
	}
}

func (m *DBMongo) SaveLinkPair(links T.DBlinksDTO) bool {
	ctx := context.Background()
	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := m.conn.Ping(ctxPing, readpref.Primary()); err != nil {
		m.log.LogError(err, "(DBMongo).Migration(): unable to Ping()")
		return false
	}

	coll := m.conn.Database(m.cfg.SL_DB_DBNAME).Collection(m.cfg.SL_DB_DBNAME)
	res, err := coll.InsertOne(ctx, bson.D{{Key: "slink", Value: links.Short}, {Key: "llink", Value: links.Long}})
	if err != nil {
		m.log.LogError(err, "(DBMongo).SaveLinkPair(): unable to InsertOne()")
		return false
	}
	m.log.LogDebug("(DBMongo).SaveLinkPair(): InsertOne() %s", res.InsertedID)
	return true
}

func (m *DBMongo) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	ctx := context.Background()
	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := m.conn.Ping(ctxPing, readpref.Primary()); err != nil {
		m.log.LogError(err, "(DBMongo).LoadLinkPair(): unable to Ping()")
		return T.DBlinksDTO{}
	}

	var res struct {
		ID    primitive.ObjectID `bson:"_id"` //nolint:tagliatelle // mongodb internal label
		Slink string             `bson:"slink"`
		Llink string             `bson:"llink"`
	}
	coll := m.conn.Database(m.cfg.SL_DB_DBNAME).Collection(m.cfg.SL_DB_DBNAME)
	if err := coll.FindOne(ctx, bson.D{{Key: "slink", Value: links.Short}}).Decode(&res); err != nil {
		m.log.LogError(err, "(DBMongo).LoadLinkPair(): unable to FindOne()")
		return T.DBlinksDTO{}
	}
	m.log.LogDebug("(DBMongo).LoadLinkPair(): FindOne() %s", res)
	return T.DBlinksDTO{Short: res.Slink, Long: res.Llink}
}

func (m *DBMongo) LoadAllLinkPairs() []T.DBlinksDTO {
	res := make([]T.DBlinksDTO, 0, 8)
	ctx := context.Background()
	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()
	if err := m.conn.Ping(ctxPing, readpref.Primary()); err != nil {
		m.log.LogError(err, "(DBMongo).LoadAllLinkPairs(): unable to Ping()")
		return res
	}

	coll := m.conn.Database(m.cfg.SL_DB_DBNAME).Collection(m.cfg.SL_DB_DBNAME)
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		m.log.LogError(err, "(DBMongo).LoadAllLinkPairs(): unable to Find()")
	}
	defer cur.Close(ctx)
	var cursRes struct {
		ID    primitive.ObjectID `bson:"_id"` //nolint:tagliatelle // mongodb internal label
		Slink string             `bson:"slink"`
		Llink string             `bson:"llink"`
	}
	for cur.Next(ctx) {
		if err := cur.Decode(&cursRes); err != nil {
			m.log.LogError(err, "(DBMongo).LoadAllLinkPairs(): unable to Decode() cursor")
		}
		res = append(res, T.DBlinksDTO{Short: cursRes.Slink, Long: cursRes.Llink})
	}
	m.log.LogDebug("(DBMongo).LoadAllLinkPairs(): Find() %s", res)
	return res
}

func (m *DBMongo) Migration() {
	ctx := context.Background()
	ctxPing, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	if err := m.conn.Ping(ctxPing, readpref.Primary()); err != nil {
		m.log.LogError(err, "(DBMongo).Migration(): unable to Ping()")
		return
	}
	db := m.conn.Database(m.cfg.SL_DB_DBNAME)

	if err := db.Collection(m.cfg.SL_DB_DBNAME).Drop(ctx); err != nil {
		m.log.LogError(err, "(DBMongo).Migration(): unable to Drop() collection")
		return
	}

	if err := db.CreateCollection(ctx, m.cfg.SL_DB_DBNAME); err != nil {
		m.log.LogError(err, "(DBMongo).Migration(): unable to Create() collection")
		return
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "slink", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	coll := db.Collection(m.cfg.SL_DB_DBNAME)
	idxName, idxErr := coll.Indexes().CreateOne(ctx, indexModel)
	if idxErr != nil {
		m.log.LogWarn("(DBMongo).Migration(): unable to Create() index")
	}
	m.log.LogDebug("(DBMongo).Migration(): Index %s", idxName)

	if res, err := coll.InsertOne(ctx, bson.D{{Key: "slink", Value: "5clp60"}, {Key: "llink", Value: "http://lib.ru"}}); err == nil {
		m.log.LogDebug("(DBMongo).Migration(): InsertOne() %s", res.InsertedID)
	}
	if res, err := coll.InsertOne(ctx, bson.D{{Key: "slink", Value: "dhiu79"}, {Key: "llink", Value: "http://google.ru"}}); err == nil {
		m.log.LogDebug("(DBMongo).Migration(): InsertOne() %s", res.InsertedID)
	}
}

func (m *DBMongo) Connect() func(e error) {
	if m.cfg.SL_DB_PORT == "" {
		m.cfg.SL_DB_PORT = ":27017"
	}
	// mgURI := "mongodb://user:password@host:27017"
	mgURI := "mongodb://" + m.cfg.SL_DB_LOGIN + ":" + m.cfg.SL_DB_PASS + "@" + m.cfg.SL_DB_IP + m.cfg.SL_DB_PORT
	opt := options.Client().ApplyURI(mgURI)
	ctx := context.Background()
	conn, err := mongo.Connect(ctx, opt)
	if err != nil {
		m.log.LogError(err, "(DBMongo).Connect(): unable to connect to mongodb: "+mgURI)
		return func(e error) {}
	} else {
		m.conn = conn
		m.log.LogInfo("mongodb connected: " + mgURI)
	}
	m.Migration()
	return func(e error) {
		ctxSHD, cancelSHD := context.WithTimeout(ctx, 5*time.Second)
		defer cancelSHD()
		if err := m.conn.Disconnect(ctxSHD); err != nil {
			m.log.LogError(err, "(DBMongo).Connect(): mongodb disconnection error")
		}
		if e != nil {
			m.log.LogError(e, "(DBMongo).Connect(): mongodb shutdown with error")
		}
		m.conn = nil
		m.log.LogInfo("mongodb disconnected")
	}
}
