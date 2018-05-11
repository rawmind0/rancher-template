package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type rancherTemplate struct {
	Name        string `description:"Template name"`
	Hash        string `description:"Template data hash"`
	Destination string `description:"Template destination file" yaml:"destination,omitempty"`
	Source      string `description:"Template source file" yaml:"source,omitempty"`
	Action      string `description:"Template action if change" yaml:"action,omitempty"`
}

func (r *rancherTemplate) getConfig(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithFields(log.Fields{"file": file, "error": err}).Error("Failed reading config yaml.")
		return err
	}

	err = yaml.Unmarshal(content, &r)
	if err != nil {
		log.WithFields(log.Fields{"file": file, "error": err}).Error("Failed unmarshaling config yaml.")
	}

	_, r.Name = filepath.Split(r.Source)

	r.Hash = r.getDestinationHash()

	return err
}

func (r *rancherTemplate) getDestinationHash() string {
	content, err := ioutil.ReadFile(r.Destination)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", sha256.Sum256(content))
}

func (r *rancherTemplate) getDataHash(w []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(w))
}

func (r *rancherTemplate) getHash() string {
	return r.Hash
}

func (r *rancherTemplate) updateHash(h string) bool {
	if r.hasChanged(h) {
		log.WithFields(log.Fields{"Old": r.Hash, "New": h}).Debug("Updating hash.")
		r.Hash = h
		return true
	}
	return false
}

func (r *rancherTemplate) hasChanged(h string) bool {
	return h != r.Hash
}

func (r *rancherTemplate) doAction() {
	if r.Action != "" {
		log.WithField("action", r.Action).Info("Executing Action.")
		err := exec.Command("sh", "-c", r.Action).Run()
		if err != nil {
			log.WithFields(log.Fields{"action": r.Action, "error": err}).Error("Failed executing action.")
		}
	}
}

func (r *rancherTemplate) getTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
		"replace": func(s, old, new string) string {
			return strings.Replace(s, old, new, -1)
		},
		"tolower": func(s string) string {
			return strings.ToLower(s)
		},
		"contains": func(s, c string) bool {
			return strings.Contains(s, c)
		},
		"ishealthy": func(s string) bool {
			return strings.Contains(s, "healthy")
		},
		"isrunning": func(s string) bool {
			return strings.Contains(s, "running")
		},
		"getStringValue":      getStringValue,
		"getBoolValue":        getBoolValue,
		"getIntValue":         getIntValue,
		"getInt64Value":       getInt64Value,
		"getSliceStringValue": getSliceStringValue,
		"hasLabel":            has,
		"hasLabelPrefix":      hasPrefix,
	}
}

func (r *rancherTemplate) execute(data interface{}) {
	log.WithField("file", r.Source).Debug("Executing template.")

	t, err := template.New(r.Name).Funcs(r.getTemplateFunc()).ParseFiles(r.Source)
	if err != nil {
		log.WithFields(log.Fields{"file": r.Source, "error": err}).Error("Failed parsing template.")
		return
	}

	var destBuf bytes.Buffer
	err = t.Execute(&destBuf, data)
	if err != nil {
		log.WithFields(log.Fields{"file": r.Source, "error": err}).Error("Failed executing template.")
		return
	}

	destBytes := destBuf.Bytes()
	if r.updateHash(r.getDataHash(destBytes)) {
		err := ioutil.WriteFile(r.Destination, destBytes, 0644)
		if err != nil {
			log.WithFields(log.Fields{"file": r.Destination, "error": err}).Error("Failed writing file.")
			return
		}

		log.WithField("file", r.Destination).Info("Template has been updated")

		r.doAction()
	}
}

type rancherTemplates struct {
	rancherTemplates []*rancherTemplate
}

func newRancherTemplates(files []string) *rancherTemplates {
	var temp = &rancherTemplates{}

	err := temp.getConfig(files)
	if err != nil {
		log.WithField("error", err).Error("Failed creating rancherTemplates.")
		return nil
	}

	return temp
}

func (r *rancherTemplates) execute(data interface{}) {
	for _, tmpl := range r.rancherTemplates {
		tmpl.execute(data)
	}
}

func (r *rancherTemplates) getConfig(files []string) error {
	var err error
	for _, file := range files {
		var temp = &rancherTemplate{}

		err = temp.getConfig(file)
		if err == nil {
			r.rancherTemplates = append(r.rancherTemplates, temp)
		}
	}

	if len(files) == 0 {
		log.Fatal("No template config file found.")
	}

	return err
}
