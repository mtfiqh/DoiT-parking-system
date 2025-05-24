package queuex_test

import (
	"github.com/mtfiqh/DoiT-parking-system/pkg/queuex"
	"sync"
	"sync/atomic"
	"testing"
)

func TestQueuex(t *testing.T) {

	type QueueOP func(queue *queuex.Queue[int]) (int, bool)
	type TestCase struct {
		Op           QueueOP
		Expected     int
		ExpectedBool bool
	}

	enqueue := func(val int) QueueOP {
		return func(queue *queuex.Queue[int]) (int, bool) {
			queue.Enqueue(val)
			return val, true
		}
	}

	dequeue := func() QueueOP {
		return func(queue *queuex.Queue[int]) (int, bool) {
			return queue.Dequeue()
		}
	}

	tests := []struct {
		name      string
		testCases []TestCase // Use lowercase for field name
	}{
		{
			name: "Enqueue and Dequeue",
			testCases: []TestCase{
				{
					Op:           enqueue(1),
					Expected:     1,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     1,
					ExpectedBool: true,
				},
			},
		},
		{
			name: "multiple enqueue and dequeue",
			testCases: []TestCase{
				{
					Op:           enqueue(1),
					Expected:     1,
					ExpectedBool: true,
				},
				{
					Op:           enqueue(2),
					Expected:     2,
					ExpectedBool: true,
				},
				{
					Op:           enqueue(3),
					Expected:     3,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     1,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     2,
					ExpectedBool: true,
				},
				{
					Op:           enqueue(1),
					Expected:     1,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     3,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     1,
					ExpectedBool: true,
				},
				{
					Op:           dequeue(),
					Expected:     0,
					ExpectedBool: false,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := queuex.NewQueue[int]()
			for _, tc := range test.testCases {
				val, ok := tc.Op(queue)
				if val != tc.Expected || ok != tc.ExpectedBool {
					t.Errorf("Expected (%d, %t), got (%d, %t)", tc.Expected, tc.ExpectedBool, val, ok)
				}
			}
		})
	}
}

func TestQueuexConcurrency(t *testing.T) {
	queue := queuex.NewQueue[int]()
	const numGoroutines = 10
	const numOperations = 1000

	var wg sync.WaitGroup
	var dequeuedItems int32
	var enqueuedItems int32

	// Launch goroutines that both enqueue and dequeue
	wg.Add(numGoroutines * 2)

	// Producer goroutines (enqueue)
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				queue.Enqueue(j)
				atomic.AddInt32(&enqueuedItems, 1)
			}
		}(i)
	}

	// Consumer goroutines (dequeue)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if _, ok := queue.Dequeue(); ok {
					atomic.AddInt32(&dequeuedItems, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Verify final state
	remaining := queue.Size
	totalProcessed := int(dequeuedItems) + remaining

	t.Logf("Enqueued: %d, Dequeued: %d, Remaining: %d",
		enqueuedItems, dequeuedItems, remaining)

	if totalProcessed != int(enqueuedItems) {
		t.Errorf("Expected total processed items to be %d, got %d",
			enqueuedItems, totalProcessed)
	}

	// verify remaining items can be dequeued
	values := queue.Print()
	if len(values) != remaining {
		t.Errorf("Expected %d remaining values, got %d", remaining, len(values))
	}
}
