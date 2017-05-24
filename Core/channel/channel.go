// Copyright 2017 GoCrab neeke@php.net All Rights Reserved.
//
// Licensed under the GNU General Public License, Version 3.
//
// Powered By CloudWise

/*
demo
package tasks

import (
	"github.com/CloudWise-OpenSource/GoCrab/Example/SendProxy/enums"
	"github.com/CloudWise-OpenSource/GoCrab/Example/SendProxy/models"
	"github.com/CloudWise-OpenSource/GoCrab/Core/channel"
	"fmt"
	"time"
)

var (
	ChannelSharded channel.Sharded
)

func writer() bool {
	if models.GetCount() >= 10 {
		for i := 0; i < 10; i++ {
			a, err := models.Pop()
			if err != nil {
				return false
			}
			fmt.Printf("a ---  %v\n", a.Content)

			time.Sleep(time.Second * 1)
		}
	}

	fmt.Printf("len ---  %d\n\n", models.GetCount())

	return true
}

func reader() bool {
	//fmt.Printf("reader  %v\n", models.GetCount())
	return true
}

func TaskInit() {
	channel.SetChannelCount(enums.ENUM_CHANNEL_COUNT)
	ChannelSharded = channel.Sharded{make(chan int), make(chan int)}
	channel.ShardedWatching(ChannelSharded, writer, reader)
}

func WriteWatcher() {
	ChannelSharded.Writer <- 1
}

func ReadWatcher() {
	<-ChannelSharded.Reader
}
*/

package channel

import (
	"fmt"
)

type Sharded struct {
	Reader chan int
	Writer chan int
}

type Watcher func() bool

var (
	channelCount int = 10
)

func SetChannelCount(count int) {
	channelCount = count
}

func ShardedWatching(sh Sharded, writer Watcher, reader Watcher) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	//fmt.Printf("ShardedWatching channelCount --- %d\n", channelCount)

	for i := 0; i < channelCount; i++ {
		go func() {
			var value int = 0
			for {
				select {
				case value = <-sh.Writer:
					writer()
					break
				case sh.Reader <- value:
					reader()
					break
				}
			}
		}()
	}

}
