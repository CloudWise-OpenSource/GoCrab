// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

package queue

import (
	"container/list"
	"errors"
	"math/rand"
	"time"
)

var (
	ErrQueueCount = errors.New("queueCount can't be less than 1.")
	ErrNoData     = errors.New("no data")
)

type Queue struct {
	QueueCount int
	Queues     map[int]*list.List
}

func InitQueue(queueCount int) (Queue, error) {
	var q Queue
	if queueCount < 1 {
		return q, ErrQueueCount
	}

	q.QueueCount = queueCount

	q.Queues = make(map[int]*list.List)
	for i := 0; i < queueCount; i++ {
		lq := list.New()
		q.Queues[i] = lq
	}

	return q, nil
}

func (q *Queue) Push(data interface{}, queueNo ...int) bool {
	queueNoLen := len(queueNo)

	if queueNoLen < 1 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		qNo := r.Intn(q.QueueCount)
		q.Queues[qNo].PushBack(data)
		return true
	} else {
		for i := 0; i < queueNoLen; i++ {
			q.Queues[queueNo[i]].PushBack(data)
		}
	}

	return true
}

func (q *Queue) Lens() (Count int) {
	Count = 0
	for i := 0; i < q.QueueCount; i++ {
		Count += q.Queues[i].Len()
	}

	return Count
}

func (q *Queue) Len(queueNo int) (Count int) {
	return q.Queues[queueNo].Len()
}

func (q *Queue) Pop(queueNo int) (data *list.Element, err error) {
	if q.Queues[queueNo].Len() < 1 {
		return data, ErrNoData
	}

	data = q.Queues[queueNo].Front()
	q.Queues[queueNo].Remove(data)

	return data, nil
}
