// The MIT License (MIT)
//
// Copyright (c) 2015 Jamie Alquiza
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	endpoints = [][]string{
		[]string{"cluster-health", "_cluster/health"},
		[]string{"cluster-stats", "_cluster/stats"},
	}

	nodeIp         string
	nodePort       string
	updateInterval int

	stats   = make(map[string][]byte)
	metrics = make(map[string]int)
)

func init() {
	flag.StringVar(&nodeIp, "ip", "127.0.0.1", "ElasticSearch IP address")
	flag.StringVar(&nodePort, "port", "9200", "ElasticSearch port")
	flag.IntVar(&updateInterval, "interval", 30, "update interval")
	flag.Parse()
}

func getNodeName() (string, error) {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/_nodes/_local/name")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(contents, &NodeName)

	var name string
	for k, _ := range NodeName.Nodes.(map[string]interface{}) {
		name = k
	}
	return name, nil
}

func queryEndpoint(endpoint string) ([]byte, error) {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func fetchMetrics() ([]byte, error) {
	for i := range endpoints {
		key, endpoint := endpoints[i][0], endpoints[i][1]

		resp, err := queryEndpoint(endpoint)
		if err != nil {
			return nil, err
		}

		stats[key] = resp
	}

	json.Unmarshal(stats["cluster-stats"], &clusterStats)
	json.Unmarshal(stats["cluster-health"], &clusterHealth)

	metrics["es-stats.state.red"] = 0
	metrics["es-stats.state.yellow"] = 0
	metrics["es-stats.state.green"] = 0
	// Flip value according to read state.
	metrics["es-stats.state."+clusterHealth.Status] = 1

	metrics["es-stats.shards.active_primary_shards"] = clusterHealth.ActivePrimaryShards
	metrics["es-stats.shards.active_shards"] = clusterHealth.ActiveShards
	metrics["es-stats.shards.relocating_shards"] = clusterHealth.RelocatingShards
	metrics["es-stats.shards.initializing_shards"] = clusterHealth.InitializingShards
	metrics["es-stats.shards.unassigned_shards"] = clusterHealth.UnassignedShards

	metrics["es-stats.indices"] = clusterStats.Indices.Count
	metrics["es-stats.docs"] = clusterStats.Indices.Docs.Count
	metrics["es-stats.cluster_cpu_cores"] = clusterStats.Nodes.Os.AvailableProcessors
	metrics["es-stats.cluster_memory"] = clusterStats.Nodes.Os.Mem.TotalInBytes

	metrics["es-stats.nodes.master"] = clusterStats.Nodes.Count.MasterOnly
	metrics["es-stats.nodes.data"] = clusterStats.Nodes.Count.DataOnly
	metrics["es-stats.nodes.master_data"] = clusterStats.Nodes.Count.MasterData
	metrics["es-stats.nodes.client"] = clusterStats.Nodes.Count.Client

	metrics["es-stats.fs.total"] = clusterStats.Nodes.Fs.TotalInBytes
	metrics["es-stats.fs.available"] = clusterStats.Nodes.Fs.AvailableInBytes
	storageUsed := metrics["es-stats.fs.total"] - metrics["es-stats.fs.available"]
	metrics["es-stats.fs.used"] = storageUsed

	metrics["es-stats.mem.jvm.heap_used_in_bytes"] = clusterStats.Nodes.Jvm.Mem.HeapUsedInBytes
	metrics["es-stats.mem.jvm.heap_max_in_bytes"] = clusterStats.Nodes.Jvm.Mem.HeapMaxInBytes
	metrics["es-stats.mem.store.size_in_bytes"] = clusterStats.Indices.Store.SizeInBytes
	metrics["es-stats.mem.store.throttle_time_in_millis"] = clusterStats.Indices.Store.ThrottleTimeInMillis
	metrics["es-stats.mem.fielddata.memory_size_in_bytes"] = clusterStats.Indices.Fielddata.MemorySizeInBytes
	metrics["es-stats.mem.fielddata.evictions"] = clusterStats.Indices.Fielddata.Evictions
	metrics["es-stats.mem.filter_cache.memory_size_in_bytes"] = clusterStats.Indices.FilterCache.MemorySizeInBytes
	metrics["es-stats.mem.filter_cache.evictions"] = clusterStats.Indices.FilterCache.Evictions
	metrics["es-stats.mem.id_cache.memory_size_in_bytes"] = clusterStats.Indices.IdCache.MemorySizeInBytes
	metrics["es-stats.mem.completion.size_in_bytes"] = clusterStats.Indices.Completion.SizeInBytes
	metrics["es-stats.mem.segments.count"] = clusterStats.Indices.Segments.Count
	metrics["es-stats.mem.segments.memory_in_bytes"] = clusterStats.Indices.Segments.MememoryInBytes
	metrics["es-stats.mem.segments.index_writer_memory_in_bytes"] = clusterStats.Indices.Segments.IndexWriterMemoryInBytes
	metrics["es-stats.mem.segments.index_writer_max_memory_in_bytes"] = clusterStats.Indices.Segments.IndexWriterMaxMemoryInBytes
	metrics["es-stats.mem.segments.version_map_memory_in_bytes"] = clusterStats.Indices.Segments.VersionMapMemoryInBytes
	metrics["es-stats.mem.segments.fixed_bit_set_memory_in_bytes"] = clusterStats.Indices.Segments.FixedBitSetMemoryInBytes

	metricsJson, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return nil, err
	}

	return metricsJson, nil
}

var NodeName struct {
	Nodes interface {
	} `json:"nodes"`
}

var clusterHealth struct {
	Status              string `json:"status"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassined_shards"`
}

var clusterStats struct {
	Indices struct {
		Count int
		Docs  struct {
			Count int `json:"count"`
		} `json:"docs"`
		Store struct {
			SizeInBytes          int `json:"size_in_bytes"`
			ThrottleTimeInMillis int `json:"throttle_time_in_millis"`
		} `json:"store"`
		Fielddata struct {
			MemorySizeInBytes int `json:"memory_size_in_bytes"`
			Evictions         int `json:"evictions"`
		} `json:"fielddata"`
		FilterCache struct {
			MemorySizeInBytes int `json:"memory_size_in_bytes"`
			Evictions         int `json:"evictions"`
		} `json:"filter_cache"`
		IdCache struct {
			MemorySizeInBytes int `json:"memory_size_in_bytes"`
		} `json:"id_cache"`
		Completion struct {
			SizeInBytes int `json:"size_in_bytes"`
		} `json:"completion"`
		Segments struct {
			Count                       int `json:"count"`
			MememoryInBytes             int `json:"memory_in_bytes"`
			IndexWriterMemoryInBytes    int `json:"index_writer_memory_in_bytes"`
			IndexWriterMaxMemoryInBytes int `json:"index_writer_max_memory_in_bytes"`
			VersionMapMemoryInBytes     int `json:"version_map_memory_in_bytes"`
			FixedBitSetMemoryInBytes    int `json:"fixed_bit_set_memory_in_bytes"`
		} `json:"segments"`
	} `json:"indices"`
	Nodes struct {
		Count struct {
			MasterOnly int `json:"master_only"`
			DataOnly   int `json:"data_only"`
			MasterData int `json:"master_data"`
			Client     int `json:"client"`
		} `json:"count"`
		Os struct {
			AvailableProcessors int `json:"available_processors"`
			Mem                 struct {
				TotalInBytes int `json:"total_in_bytes"`
			} `json:"mem"`
		} `json:"os"`
		Jvm struct {
			Mem struct {
				HeapUsedInBytes int `json:"heap_used_in_bytes"`
				HeapMaxInBytes  int `json:"heap_max_in_bytes"`
			} `json:"mem"`
		} `json:"jvm"`
		Fs struct {
			TotalInBytes     int `json:"total_in_bytes"`
			AvailableInBytes int `json:"available_in_bytes"`
		} `json:"fs"`
	} `json:"nodes"`
}

func main() {
	// localhost:9200/_cluster/state/master_node, master_node

	// Grab node name.
	var nodeName *string
	retry := time.Tick(time.Duration(updateInterval) * time.Second)

	for _ = range retry {
		name, err := getNodeName()
		if err != nil {
			log.Printf("ElasticSearch unreachable: %s", err)

		} else {
			nodeName = &name
			break
		}
	}

	fmt.Println(*nodeName)
	
	// Run.
	tick := time.Tick(time.Duration(updateInterval) * time.Second)
	for {
		select {
		case <-tick:
			m, err := fetchMetrics()
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println(string(m))
		}
	}
}