package goapplog

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"os"
)

var (
	DefaultName            = "go-app"
	DefaultIP              = "0.0.0.0"
	DefaultPort            = 0
	DefaultNode            = 0
	DefaultLevel           = logrus.InfoLevel
	DefaultFilePath        = "app.log"
	DefaultTimestampFormat = "2006-01-02 15:04:05.000"
	// elasticsearch
	DefaultElasticsearchIndex = "applog"
)

type Options struct {
	Name            string
	IP              string
	Port            int
	Node            int
	Level           logrus.Level
	ReportCaller    bool
	FileEnable      bool
	FilePath        string
	TimestampFormat string
	// elasticsearch
	Elasticsearch      *elastic.Client
	ElasticsearchIndex string
}

type Option func(o *Options)

func Name(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func IP(s string) Option {
	return func(o *Options) {
		o.IP = s
	}
}

func Port(i int) Option {
	return func(o *Options) {
		o.Port = i
	}
}

func Node(i int) Option {
	return func(o *Options) {
		o.Node = i
	}
}

func Level(level logrus.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func ReportCaller(b bool) Option {
	return func(o *Options) {
		o.ReportCaller = b
	}
}

func FilePath(s string) Option {
	return func(o *Options) {
		o.FilePath = s
	}
}

func FileEnable(b bool) Option {
	return func(o *Options) {
		o.FileEnable = b
	}
}

func TimestampFormat(s string) Option {
	return func(o *Options) {
		o.TimestampFormat = s
	}
}

func Elasticsearch(c *elastic.Client) Option {
	return func(o *Options) {
		o.Elasticsearch = c
	}
}

func ElasticsearchIndex(s string) Option {
	return func(o *Options) {
		o.ElasticsearchIndex = s
	}
}

func newOptions(options ...Option) Options {
	opt := Options{
		Name:               DefaultName,
		IP:                 DefaultIP,
		Port:               DefaultPort,
		Node:               DefaultNode,
		Level:              DefaultLevel,
		ReportCaller:       true,
		FileEnable:         false,
		FilePath:           DefaultFilePath,
		TimestampFormat:    DefaultTimestampFormat,
		Elasticsearch:      nil,
		ElasticsearchIndex: DefaultElasticsearchIndex,
	}

	for _, o := range options {
		o(&opt)
	}

	return opt
}

func NewLogger(options ...Option) *logrus.Entry {
	return newLogger(options...)
}

func newLogger(options ...Option) *logrus.Entry {
	o := newOptions(options...)
	e := logrus.WithFields(logrus.Fields{
		"name": o.Name,
		"ip":   o.IP,
		"port": o.Port,
		"node": o.Node,
	})
	e.Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: o.TimestampFormat})
	e.Logger.SetLevel(o.Level)
	e.Logger.SetReportCaller(o.ReportCaller)
	e.Logger.SetOutput(os.Stdout)
	// file
	if o.FileEnable {
		f, err := os.OpenFile(o.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		e.Logger.SetOutput(f)
	}
	// Elasticsearch
	if o.Elasticsearch != nil {
		hook, err := elogrus.NewAsyncElasticHook(o.Elasticsearch, o.IP, o.Level, o.ElasticsearchIndex)
		if err != nil {
			panic(err)
		}
		e.Logger.Hooks.Add(hook)
	}
	return e
}
