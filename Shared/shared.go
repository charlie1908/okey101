package Shared

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Byte []byte

func (b Byte) byteAppendComma(text string) []byte {
	for _, num := range strings.Split(text, ",") {
		intNum, _ := strconv.Atoi(strings.TrimSpace(num))
		b = append(b, byte(intNum))
	}
	return b
}

type configuration struct {
	SQLURL                 string
	SQLTIMEOUT             time.Duration
	SECRETKEY              string
	REDISPASSWORD          string
	TOKENEXPIRETIME        time.Duration
	REFRESHTOKENEXPIRETIME time.Duration
	REDISURL               string
	BYTES                  []byte
	ELASTICUSER            string
	ELASTICPASSWORD        string
	ELASTICURL             string
	ELASTICLOGININDEX      string
	ELASTICAUDITINDEX      string
	ELASTICERRORINDEX      string
	RABBITMQHOST           string
	RABBITMQUSER           string
	RABBITMQPASSWORD       string
	RABBITMQPORT           int
	LANG                   string
	MONGOURL               string
}

var byteStrings = getConfigPrm("encryption", "bytes").(string)
var bytes = make(Byte, 0).byteAppendComma(byteStrings)
var Config = configuration{
	SQLURL:                 getConfigPrm("database", "sqlurl").(string),
	SQLTIMEOUT:             time.Second * time.Duration(getConfigPrm("database", "sqltimeout").(int)),
	SECRETKEY:              "2909012565820034",
	REDISPASSWORD:          getConfigPrm("redis", "password").(string),
	REDISURL:               getConfigPrm("redis", "url").(string),
	TOKENEXPIRETIME:        time.Minute * time.Duration(getConfigPrm("redis", "tokenExpireTime").(int)),
	REFRESHTOKENEXPIRETIME: time.Minute * time.Duration(getConfigPrm("redis", "refreshTokenExpireTime").(int)),
	BYTES:                  bytes,
	ELASTICUSER:            getConfigPrm("elastic", "username").(string),
	ELASTICPASSWORD:        getConfigPrm("elastic", "password").(string),
	ELASTICURL:             getConfigPrm("elastic", "url").(string),
	ELASTICLOGININDEX:      getConfigPrm("elastic", "elasticLoginIndex").(string),
	ELASTICAUDITINDEX:      getConfigPrm("elastic", "elasticAuditIndex").(string),
	ELASTICERRORINDEX:      getConfigPrm("elastic", "elasticErrorIndex").(string),
	RABBITMQHOST:           getConfigPrm("rabbitmq", "host").(string),
	RABBITMQUSER:           getConfigPrm("rabbitmq", "username").(string),
	RABBITMQPASSWORD:       getConfigPrm("rabbitmq", "password").(string),
	RABBITMQPORT:           getConfigPrm("rabbitmq", "port").(int),
	LANG:                   getConfigPrm("", "lang").(string),
	MONGOURL:               getConfigPrm("mongo", "url").(string),
}

func getConfigPrm(root string, prm string) interface{} {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")

	// Projenin kök dizinini bul
	exePath, _ := os.Getwd() // Current working directory
	projectRoot := findProjectRoot(exePath, "okey101")

	sharedPath := filepath.Join(projectRoot, "Shared")
	v.AddConfigPath(sharedPath)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %w", err))
	}

	if root == "" {
		return v.Get(prm)
	}
	return v.Get(root + "." + prm)
}

// "okey101" klasörünü bulana kadar yukarı çıkar
func findProjectRoot(start string, projectDirName string) string {
	for {
		if filepath.Base(start) == projectDirName {
			return start
		}
		parent := filepath.Dir(start)
		if parent == start {
			break
		}
		start = parent
	}
	panic("project root not found")
}
