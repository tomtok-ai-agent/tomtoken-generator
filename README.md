# TomToken Generator

This repository contains two implementations of a TomToken stream generator in Go. TomToken is a flexible token format that allows modeling text with different types of content.

## Implementations

### 1. Without Channels (`tomtoken_without_channels.go`)

This implementation generates tokens directly in the main execution flow, writing them to the output as they are generated. It uses a traditional approach where all processing happens in a single function.

Key features:
- Direct token generation and serialization
- Manual JSON formatting with line length control
- Command-line parameters for output destination and token count
- Efficient memory usage through immediate serialization

### 2. With Channels (`tomtoken_with_channels.go`)

This implementation uses Go channels to separate token generation from serialization. The generator runs in a separate goroutine and sends tokens through a channel to the main routine which handles serialization.

Key features:
- Separation of concerns: generation and serialization are decoupled
- Streaming approach using Go channels
- Buffered channel for performance optimization
- Same command-line interface as the first implementation
- Memory-efficient for large token streams

## Usage

Both programs support the same command-line parameters:

```
-output string
    Output: 'stdout' or path to .json file (default "stdout")
-tokens int
    Number of TomToken tokens to generate (default varies by implementation)
```

### Examples

Generate 10,000 tokens and output to stdout:
```
./tomtoken_without_channels -tokens=10000
```

Generate 1 million tokens and save to a file:
```
./tomtoken_with_channels -tokens=1000000 -output=output.json
```

## TomToken Format

TomToken is a flexible token format that can represent:

1. String tokens - Regular text strings
2. Numeric references - References to entries in the referenceMap (e.g., spaces, tabs)
3. Structured tokens - Objects with type, mode, count, and value properties

The output is a JSON structure with:
- `metadata`: Information about the generated content
- `referenceMap`: A mapping of numeric IDs to their string values
- `content`: An array of tokens (strings, numbers, or objects)

## Performance Considerations

The channel-based implementation is more memory-efficient for large token streams as it doesn't need to store all tokens in memory before serialization. This makes it suitable for generating very large token streams.

## License

MIT
