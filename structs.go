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

var nodesLocal struct {
	Nodes interface {
	} `json:"nodes"`
}

var clusterState struct {
	MasterNode string `json:"master_node"`
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