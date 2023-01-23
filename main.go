package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ItemBatch is a batch of items we want to move from one system to another
// without each item losing the trace information that goes with it.
type ItemBatch struct {
	Items []Item `json:"items"`
}

type Item struct {
	// Integer is just some value that I made up for this example
	Integer int `json:"integer"`

	// Result is just some result that I made up for this example
	Result int `json:"result"`

	// propagation.MapCarrier is simply a map[string]string that implements
	// the propagation.TextMapCarrier interfaces. This allows us to pass
	// TraceCarrier to the propagator to set the correct context for this item
	TraceCarrier propagation.MapCarrier `json:"trace_carrier"`
}

func produceItem(i int, tp *sdktrace.TracerProvider) Item {
	prop := propagation.TraceContext{}
	tr := tp.Tracer("producer")
	ctx, span := tr.Start(context.Background(), "produce item")
	defer span.End()

	item := Item{Integer: 1, TraceCarrier: make(propagation.MapCarrier)}
	prop.Inject(ctx, item.TraceCarrier)
	return item
}

func producer(tp *sdktrace.TracerProvider) ItemBatch {
	var batch ItemBatch
	for i := 0; i < 3; i++ {
		batch.Items = append(batch.Items, produceItem(i, tp))
	}
	return batch
}

func consumer(tp *sdktrace.TracerProvider, batch ItemBatch) {
	prop := propagation.TraceContext{}
	tr := tp.Tracer("consumer")

	for _, item := range batch.Items {
		extractedCtx := prop.Extract(context.Background(), item.TraceCarrier)
		ctx, span := tr.Start(extractedCtx, "consume item")
		item := addResult(ctx, item)
		span.End()
		fmt.Printf("Item Result: %d\n", item.Result)
	}
}

func addResult(ctx context.Context, item Item) Item {
	_, span := otel.Tracer("consumer").Start(ctx, "addResult")
	item.Result = item.Integer + 1
	defer span.End()
	return item
}

func main() {

	// ==============================
	// Simple Initialization
	// ==============================
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	// ==============================

	// Producer simulates a system that creates the trace for each item
	// and then needs to send that batch of items to another system
	batch := producer(tp)

	// Pretend we marshalled this batch of items and sent it to a
	// different system here.
	b, err := json.MarshalIndent(batch, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("== Serialized Items ===\n")
	fmt.Printf("%s", string(b))
	fmt.Printf("\n=======================\n")

	consumer(tp, batch)
}
