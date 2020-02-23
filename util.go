package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/paulbellamy/ratecounter"
)

var counter = ratecounter.NewRateCounter(1 * time.Second)

func makeDir(path string) {
	fullPath, _ := filepath.Abs(path)

	os.MkdirAll(fullPath, os.ModePerm)
}

func writeFile(path string, jsonData []byte) {
	ioutil.WriteFile(path, jsonData, os.ModePerm)
	counter.Incr(1)
}

func makeHash(rawData string) string {
	data := []byte(rawData)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
