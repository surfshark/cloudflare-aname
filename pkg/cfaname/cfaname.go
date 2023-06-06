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

package cfaname

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

// CFANAME holds the configuration for the current ANAME/ALIAS record synchronization.
type CFANAME struct {
	api *cloudflare.API
	zoneID *cloudflare.ResourceContainer
	recordName string
	recordTTL int
	targetRecordName string
}

func New(api *cloudflare.API, zoneID, recordName string, recordTTL int, targetRecordName string) *CFANAME {
	return &CFANAME{
		api: api, zoneID: cloudflare.ZoneIdentifier(zoneID), recordName: recordName,
		recordTTL: recordTTL, targetRecordName: targetRecordName,
	}
}

func (cfaname *CFANAME) lookupCurrentRecords(ctx context.Context) ([]cloudflare.DNSRecord, error) {
	records, _, err := cfaname.api.ListDNSRecords(ctx, cfaname.zoneID,
		cloudflare.ListDNSRecordsParams{Name: cfaname.recordName, Type: "A"})
	if err != nil {
		return nil, err
	}
	records6, _, err := cfaname.api.ListDNSRecords(ctx, cfaname.zoneID,
		cloudflare.ListDNSRecordsParams{Name: cfaname.recordName, Type: "AAAA"})
	if err != nil {
		return nil, err
	}
	return append(records, records6...), nil
}

func recordTypeForIP(ip string) (string, error) {
	if strings.IndexByte(ip, '.') >= 0 {
		return "A", nil
	}
	if strings.IndexByte(ip, ':') >= 0 {
		return "AAAA", nil
	}
	return "", errors.New("could not determine record type for " + ip)
}

func recordsContainIP(records []cloudflare.DNSRecord, ip string) bool {
	for _, v := range records {
		if v.Content == ip {
			return true
		}
	}
	return false
}

func containsIP(ips []string, ip string) bool {
	for _, v := range ips {
		if v == ip {
			return true
		}
	}
	return false
}

func (cfaname *CFANAME) Update(ctx context.Context) error {
	targetIPs, err := net.LookupHost(cfaname.targetRecordName)
	if err != nil {
		return err
	}
	currentRecords, err := cfaname.lookupCurrentRecords(ctx)
	if err != nil {
		return err
	}
	keepIPs := []string{}
	// Add new target IPs
	for _, ip := range targetIPs {
		if !recordsContainIP(currentRecords, ip) {
			recordType, err := recordTypeForIP(ip)
			if err != nil {
				return err
			}
			_, err = cfaname.api.CreateDNSRecord(ctx, cfaname.zoneID,
				cloudflare.CreateDNSRecordParams{
					Type: recordType, Name: cfaname.recordName,
					Content: ip, TTL: cfaname.recordTTL,
				},
			)
			if err != nil {
				return err
			}
		} else {
			keepIPs = append(keepIPs, ip)
		}
	}
	// Remove old records
	for _, record := range currentRecords {
		if !containsIP(keepIPs, record.Content) {
			err = cfaname.api.DeleteDNSRecord(ctx, cfaname.zoneID, record.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
