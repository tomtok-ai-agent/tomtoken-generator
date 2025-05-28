# TomToken Generator

This repository contains two implementations of TomToken generation with JSON serialization in Go:

1. **Without Channels**: A straightforward implementation that generates all tokens at once and stores them in memory.
2. **With Channels**: An efficient implementation using Go channels for stream processing, which minimizes memory usage.

## Features

- Random word generation with configurable token count
- Various token types including words, spaces, newlines, and repeats
- JSON serialization with proper formatting
- Command-line parameters for the channel-based implementation

## Implementation Comparison

### Without Channels
- Simpler implementation
- Stores all tokens in memory
- Suitable for smaller datasets

### With Channels
- Memory-efficient stream processing
- Clear separation between generation and processing
- Highly scalable for large datasets
- Supports command-line parameters

## Usage

### Without Channels

```bash
go run tomtoken_without_channels.go
```

### With Channels

```bash
# Output to stdout
go run tomtoken_with_channels.go

# Output to file with custom token count
go run tomtoken_with_channels.go -output output.json -tokens 50000
```

## Performance

The channel-based implementation is particularly efficient for large datasets, similar to the `yield` approach in Python and C#, and works effectively even with millions of elements.
