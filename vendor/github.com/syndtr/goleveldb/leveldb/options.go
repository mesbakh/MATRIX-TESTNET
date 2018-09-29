// Copyright 2018 The MATRIX Authors as well as Copyright 2014-2017 The go-ethereum Authors
// This file is consisted of the MATRIX library and part of the go-ethereum library.
//
// The MATRIX-ethereum library is free software: you can redistribute it and/or modify it under the terms of the MIT License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, 
//and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject tothe following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, 
//WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISINGFROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// Copyright (c) 2012, Suryandaru Triandana <syndtr@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func dupOptions(o *opt.Options) *opt.Options {
	newo := &opt.Options{}
	if o != nil {
		*newo = *o
	}
	if newo.Strict == 0 {
		newo.Strict = opt.DefaultStrict
	}
	return newo
}

func (s *session) setOptions(o *opt.Options) {
	no := dupOptions(o)
	// Alternative filters.
	if filters := o.GetAltFilters(); len(filters) > 0 {
		no.AltFilters = make([]filter.Filter, len(filters))
		for i, filter := range filters {
			no.AltFilters[i] = &iFilter{filter}
		}
	}
	// Comparer.
	s.icmp = &iComparer{o.GetComparer()}
	no.Comparer = s.icmp
	// Filter.
	if filter := o.GetFilter(); filter != nil {
		no.Filter = &iFilter{filter}
	}

	s.o = &cachedOptions{Options: no}
	s.o.cache()
}

const optCachedLevel = 7

type cachedOptions struct {
	*opt.Options

	compactionExpandLimit []int
	compactionGPOverlaps  []int
	compactionSourceLimit []int
	compactionTableSize   []int
	compactionTotalSize   []int64
}

func (co *cachedOptions) cache() {
	co.compactionExpandLimit = make([]int, optCachedLevel)
	co.compactionGPOverlaps = make([]int, optCachedLevel)
	co.compactionSourceLimit = make([]int, optCachedLevel)
	co.compactionTableSize = make([]int, optCachedLevel)
	co.compactionTotalSize = make([]int64, optCachedLevel)

	for level := 0; level < optCachedLevel; level++ {
		co.compactionExpandLimit[level] = co.Options.GetCompactionExpandLimit(level)
		co.compactionGPOverlaps[level] = co.Options.GetCompactionGPOverlaps(level)
		co.compactionSourceLimit[level] = co.Options.GetCompactionSourceLimit(level)
		co.compactionTableSize[level] = co.Options.GetCompactionTableSize(level)
		co.compactionTotalSize[level] = co.Options.GetCompactionTotalSize(level)
	}
}

func (co *cachedOptions) GetCompactionExpandLimit(level int) int {
	if level < optCachedLevel {
		return co.compactionExpandLimit[level]
	}
	return co.Options.GetCompactionExpandLimit(level)
}

func (co *cachedOptions) GetCompactionGPOverlaps(level int) int {
	if level < optCachedLevel {
		return co.compactionGPOverlaps[level]
	}
	return co.Options.GetCompactionGPOverlaps(level)
}

func (co *cachedOptions) GetCompactionSourceLimit(level int) int {
	if level < optCachedLevel {
		return co.compactionSourceLimit[level]
	}
	return co.Options.GetCompactionSourceLimit(level)
}

func (co *cachedOptions) GetCompactionTableSize(level int) int {
	if level < optCachedLevel {
		return co.compactionTableSize[level]
	}
	return co.Options.GetCompactionTableSize(level)
}

func (co *cachedOptions) GetCompactionTotalSize(level int) int64 {
	if level < optCachedLevel {
		return co.compactionTotalSize[level]
	}
	return co.Options.GetCompactionTotalSize(level)
}