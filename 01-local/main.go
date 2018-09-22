//    Copyright 2018 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"os"
	"runtime/trace"
	"sync"
)

const Limit = 100

func main() {
	file, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := trace.Start(file); err != nil {
		log.Fatal(err)
	}
	defer trace.Stop()

	ctx := context.Background()
	ctx, task := trace.NewTask(ctx, "printFibonacci")
	defer task.End()
	ch := genFibonacci(ctx, 10)
	for n := range ch {
		log.Println(n.String())
	}
}

func genFibonacci(ctx context.Context, length int) <-chan big.Int {
	ctx, task := trace.NewTask(ctx, "genFibonacci")
	defer task.End()
	ch := make(chan big.Int)
	var wg sync.WaitGroup
	go func(ch chan<- big.Int) {
		for i := 0; i < length; i++ {
			wg.Add(1)
			defer trace.StartRegion(ctx, "genFibonacci").End()
			defer wg.Done()
			n := rand.Intn(Limit)
			fibonacci(ctx, n, ch)
		}
		close(ch)
	}(ch)
	return ch
}

func fibonacci(ctx context.Context, n int, ch chan<- big.Int) {
	a, b := big.NewInt(0), big.NewInt(1)
	for i := 0; i < n; i++ {
		a.Add(a, b)
		a, b = b, a
	}
	ch <- *a
}
