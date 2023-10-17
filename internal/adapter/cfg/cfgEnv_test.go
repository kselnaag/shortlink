package adapterCfg_test

import (
	"os"
	"path/filepath"
	adapterCfg "shortlink/internal/adapter/cfg"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCfgEnv(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("CfgEnv", func(t *testing.T) {
		cfg := &T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_MODE:  "fprintf",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_MODE: "gin",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_MODE:   "mock",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   "",
			SL_DB_LOGIN:  "login",
			SL_DB_PASS:   "password",
			SL_DB_DBNAME: "shortlink",
		}
		log := adapterLog.NewLogFprintf(cfg)

		wd, _ := os.Getwd()
		for !strings.HasSuffix(wd, "shortlink") {
			wd = filepath.Dir(wd)
		}
		filename := filepath.Join(wd, "test/test.env")

		adapterCfg.NewCfgEnvFile(filename, cfg, log)
		newcfg := &T.CfgEnv{
			SL_APP_NAME:  "aaa",
			SL_LOG_MODE:  "bbb",
			SL_LOG_LEVEL: "ccc",
			SL_HTTP_MODE: "ddd",
			SL_HTTP_IP:   "eee",
			SL_HTTP_PORT: "fff",
			SL_DB_MODE:   "ggg",
			SL_DB_IP:     "hhh",
			SL_DB_PORT:   "iii",
			SL_DB_LOGIN:  "jjj",
			SL_DB_PASS:   "kkk",
			SL_DB_DBNAME: "mmm",
		}
		asrt.Equal(newcfg, cfg)
	})
}
