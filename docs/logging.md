# Logging in kyverno-json

kyverno-json includes a built-in logging system that can be controlled using the `--verbosity` flag.

## Usage

You can set the log level using the `--verbosity` flag when running any kyverno-json command:

```bash
kyverno-json scan --verbosity debug --policy policy.yaml --payload payload.json
```

## Available Log Levels

The following log levels are available, in order of increasing verbosity:

1. `error` - Only log errors
2. `warn` - Log warnings and errors
3. `info` - Log informational messages, warnings, and errors (default)
4. `debug` - Log debug messages, informational messages, warnings, and errors

## Examples

### Using the default log level (info)

```bash
kyverno-json scan --policy policy.yaml --payload payload.json
```

### Enabling debug logs for troubleshooting

```bash
kyverno-json scan --verbosity debug --policy policy.yaml --payload payload.json
```

### Reducing log output to only errors

```bash
kyverno-json scan --verbosity error --policy policy.yaml --payload payload.json
```

## Log Format

Logs are output in JSON format, which makes them easier to parse and analyze with tools like `jq`. Each log entry includes:

- Timestamp
- Log level
- Message
- Additional contextual fields

Example:

```json
{"level":"info","ts":"2023-04-01T12:34:56.789Z","msg":"Evaluating resources","resourceCount":1,"policyCount":2}
```

## Integration with Other Systems

Since the logs are output in a structured format, they can be easily integrated with log management systems like ELK, Loki, or CloudWatch.

## Future Improvements

In the future, we could enhance the logging system to include:
- Additional output formats
- Log file destination configuration
- More granular log level control per component 