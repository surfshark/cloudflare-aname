//
// Copyright 2023 Laurynas Četyrkinas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"gopkg.in/yaml.v3"

	"github.com/surfshark/cloudflare-aname/pkg/cfaname"
)

type configCloudflare struct {
	ApiToken string `yaml:"api-token"`
	ZoneID string `yaml:"zone-id"`
}

type configRecord struct {
	Name string `yaml:"name"`
	Target string `yaml:"target"`
	TTL int `yaml:"ttl"`
}

type config struct {
	Cloudflare configCloudflare `yaml:"cloudflare"`
	Record configRecord `yaml:"record"`
}

var conf config

func init() {
	p := flag.String("conf", "config.yaml", "Configuration file")
	flag.Parse()
	b, err := os.ReadFile(*p)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if conf.Cloudflare.ApiToken == "" || conf.Cloudflare.ZoneID == "" || conf.Record.Name == "" || conf.Record.Target == "" {
		log.Fatal("mandatory configuration parameters missing")
	}
	if conf.Record.TTL == 0 {
		conf.Record.TTL = 60
	}
}

func main() {
	api, err := cloudflare.NewWithAPIToken(conf.Cloudflare.ApiToken)
	if err != nil {
		log.Fatal(err)
	}
	record := cfaname.New(api, conf.Cloudflare.ZoneID, conf.Record.Name, conf.Record.TTL,
		conf.Record.Target)
	err = record.Update(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
