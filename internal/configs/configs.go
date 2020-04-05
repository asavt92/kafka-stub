package configs

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type KafkaConfig struct {
	BootstrapServers []string
	InTopic          []string
	OutTopic         []string
	Group            string
}

type JsonConfig struct {
	InJsonSchemaPath       string
	OutJsonExamplesDirPath string
	IdFieldMapping         map[string]string
}

const (
	configPathDefault string = "configs/config.yml"
)

var (
	configPath string
)

func InitConfig() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	viper.SetDefault("loglevel", "debug")

	flag.StringVar(&configPath, "config-path", configPathDefault, "path to config file")
	flag.Parse()

	if configPath != "" {
		log.Infof("Parsing config: %s", configPath)
		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			log.Panicf("Unable to read config file: %s", err)
		}
	} else {
		log.Infof("Config file is not specified.")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	//logger
	logLevelString := viper.GetString("loglevel")
	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		log.Panicf("Unable to parse loglevel: %s", logLevelString)
	}

	log.SetLevel(logLevel)

	_printConfigs()
}

func _printConfigs() {
	log.Infoln("---===      CONFIGS     ===---")
	for _, k := range viper.AllKeys() {
		log.Infof("config[  %-40s  ] = %v", k, viper.Get(k))
	}
	log.Infoln("---===        END      ===---")
}

func InitKafkaConfig() *KafkaConfig {
	c := KafkaConfig{
		BootstrapServers: viper.GetStringSlice("kafka.bootstrap.servers"),
		Group:            viper.GetString("kafka.group"),
		InTopic:          viper.GetStringSlice("kafka.in.topic"),
		OutTopic:         viper.GetStringSlice("kafka.out.topic"),
	}
	return &c
}

func InitJsonConfig() *JsonConfig {
	c := JsonConfig{
		InJsonSchemaPath:       viper.GetString("in.json.schema.path"),
		OutJsonExamplesDirPath: viper.GetString("out.json.examples.dir.path"),
		IdFieldMapping:         viper.GetStringMapString("id.field.mapping"),
	}
	return &c
}
