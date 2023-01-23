
### Manual context propagation example
This is a simple example of how to manually extract the trace id context for propagation of the trace
id to a separate system. This is useful if you have multiple items each with their own trace id and
batch send them to a separate system.

### Output
Output of this example program
```json
{
  "items": [
    {
      "integer": 1,
      "result": 0,
      "trace_carrier": {
        "traceparent": "00-fb2178782d2b4b955b4dc4eff38c79ec-b01f373da8dd63f7-01"
      }
    },
    {
      "integer": 1,
      "result": 0,
      "trace_carrier": {
        "traceparent": "00-5bf062ce28ea663132e74808cdca7322-4dddaba5da7447d7-01"
      }
    },
    {
      "integer": 1,
      "result": 0,
      "trace_carrier": {
        "traceparent": "00-9a78ec4a109bf9b31f9091dbe419dc52-98a18b277d2f9b53-01"
      }
    }
  ]
}
```
