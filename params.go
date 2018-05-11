package main

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func check(e error, m string) {
	if e != nil {
		log.Error("[Error]: ", m, e)
	}
}

type Params struct {
	Url     string
	Prefix  string
	Logfile string
	Files   []string
	Refresh int
	Self    bool
	Debug   bool
}

func (p *Params) init() {
	var filePath string

	flag.StringVar(&p.Url, "url", "http://rancher-metadata.rancher.internal", "Rancher metadata url. RANCHER_TEMPLATE_URL")
	flag.StringVar(&p.Prefix, "prefix", "2016-07-29", "Rancher metadata prefix. RANCHER_TEMPLATE_PREFIX")
	flag.StringVar(&p.Logfile, "logfile", "/opt/tools/rancher-template/log/rancher-template.log", "Rancher template log fie. RANCHER_TEMPLATE_LOGFILE")
	flag.StringVar(&filePath, "templates", "/opt/tools/rancher-template/etc/*.yml", "Templates config files, wildcard allowed between quotes. RANCHER_TEMPLATE_FILES")
	flag.IntVar(&p.Refresh, "refresh", 300, "Rancher metadata refresh time in seconds. RANCHER_TEMPLATE_REFRESH")
	flag.BoolVar(&p.Self, "self", false, "Get self stack data or all. RANCHER_TEMPLATE_SELF")
	flag.BoolVar(&p.Debug, "debug", false, "Run in debug mode. RANCHER_TEMPLATE_DEBUG")

	flag.Parse()

	p.setEnvVar(&filePath)

	p.getFiles(filePath)
}

func (p *Params) setEnvVar(filePath *string) {
	var err error
	var auxInt int
	var auxBool bool

	url := os.Getenv("RANCHER_TEMPLATE_URL")
	if len(url) > 0 {
		p.Url = url
	}

	prefix := os.Getenv("RANCHER_TEMPLATE_PREFIX")
	if len(prefix) > 0 {
		p.Prefix = prefix
	}

	logfile := os.Getenv("RANCHER_TEMPLATE_LOGFILE")
	if len(logfile) > 0 {
		p.Logfile = logfile
	}

	files := os.Getenv("RANCHER_TEMPLATE_FILES")
	if len(files) > 0 {
		*filePath = files
	}

	refresh := os.Getenv("RANCHER_TEMPLATE_REFRESH")
	if len(refresh) > 0 {
		auxInt, err = strconv.Atoi(refresh)
		if err == nil {
			p.Refresh = auxInt
		}
	}

	self := os.Getenv("RANCHER_TEMPLATE_SELF")
	if len(self) > 0 {
		auxBool, err = strconv.ParseBool(self)
		if err == nil {
			p.Self = auxBool
		}
	}

	debug := os.Getenv("RANCHER_TEMPLATE_DEBUG")
	if len(debug) > 0 {
		auxBool, err = strconv.ParseBool(debug)
		if err == nil {
			p.Debug = auxBool
		}
	}
}

func (p *Params) getFiles(f string) {
	var err error

	p.Files, err = filepath.Glob(f)
	if err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
