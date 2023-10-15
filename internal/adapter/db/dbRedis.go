package adapterDB

import (
	"context"
	T "shortlink/internal/apptype"

	"github.com/redis/go-redis/v9"
)

var _ T.IDB = (*DBRedis)(nil)

type DBRedis struct {
	cfg  *T.CfgEnv
	log  T.ILog
	conn *redis.Client
}

func NewDBRedis(cfg *T.CfgEnv, log T.ILog) DBRedis {
	return DBRedis{
		cfg: cfg,
		log: log,
	}
}

func (r *DBRedis) SaveLinkPair(links T.DBlinksDTO) bool {
	ctx := context.Background()
	if err := r.conn.Set(ctx, links.Short, links.Long, 0).Err(); err != nil {
		r.log.LogError(err, "(DBRedis).SaveLinkPair(): redis db Set() error")
		return false
	}
	r.log.LogDebug("(DBRedis).SaveLinkPair(): %s, %s", links.Short, links.Long)
	return true
}

func (r *DBRedis) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	ctx := context.Background()
	val, err := r.conn.Get(ctx, links.Short).Result()
	if err != nil {
		r.log.LogError(err, "(DBRedis).LoadLinkPair(): redis db Get() error")
		return T.DBlinksDTO{}
	}
	r.log.LogDebug("(DBRedis).LoadLinkPair(): %s", val)
	return T.DBlinksDTO{Short: links.Short, Long: val}
}

func (r *DBRedis) LoadAllLinkPairs() []T.DBlinksDTO {
	res := make([]T.DBlinksDTO, 0, 8)
	ctx := context.Background()
	iter := r.conn.Scan(ctx, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.conn.Get(ctx, iter.Val()).Result()
		if err != nil {
			r.log.LogError(err, "(DBRedis).LoadAllLinkPairs(): redis db Get() error")
		}
		res = append(res, T.DBlinksDTO{Short: iter.Val(), Long: val})
	}
	if err := iter.Err(); err != nil {
		r.log.LogError(err, "(DBRedis).LoadAllLinkPairs(): redis db Iter() error")
	}
	r.log.LogDebug("(DBRedis).LoadAllLinkPairs(): %s", res)
	return res
}

func (r *DBRedis) Migration() {
	ctx := context.Background()
	if err := r.conn.FlushDB(ctx).Err(); err != nil {
		r.log.LogError(err, "(DBRedis).Migration(): redis db FlushDB() error")
	}
	if err := r.conn.Set(ctx, "5clp60", "http://lib.ru", 0).Err(); err != nil {
		r.log.LogError(err, "(DBRedis).Migration(): redis db Set() error")
	}
	if err := r.conn.Set(ctx, "dhiu79", "http://google.ru", 0).Err(); err != nil {
		r.log.LogError(err, "(DBRedis).Migration(): redis db Set() error")
	}
}

func (r *DBRedis) Connect() func(e error) {
	if r.cfg.SL_DB_PORT == "" {
		r.cfg.SL_DB_PORT = ":6378"
	}
	// rdURI := "redis://default:password@localhost:6379/db-number?protocol=3"
	rdURI := "redis://" + "default" + ":" + r.cfg.SL_DB_PASS + "@" + r.cfg.SL_DB_IP + r.cfg.SL_DB_PORT + "/0?protocol=3"
	opts, err := redis.ParseURL(rdURI)
	if err != nil {
		r.log.LogError(err, "(DBRedis).Connect(): unable to connect to redis db: "+rdURI)
		return func(e error) {}
	}
	conn := redis.NewClient(opts)
	status := conn.Ping(context.Background())
	_, err = status.Result()
	if err != nil {
		r.log.LogError(err, "(DBRedis).Connect(): unable to Ping() to redis db: "+rdURI)
		return func(e error) {}
	} else {
		r.conn = conn
		r.log.LogInfo("redis db connected: " + rdURI)
	}
	r.Migration()
	return func(e error) {
		if err := r.conn.Close(); err != nil {
			r.log.LogError(e, "(DBRedis).Connect(): redis db disconnection error")
		}
		if e != nil {
			r.log.LogError(e, "(DBRedis).Connect(): redis db shutdown with error")
		}
		r.conn = nil
		r.log.LogInfo("redis db disconnected")
	}
}
