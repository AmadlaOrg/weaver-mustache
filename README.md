# weaver-mustache

Amadla Weaver plugin for the Mustache template engine.

## Usage

```bash
# Show plugin info
weaver-mustache info

# Render a template with JSON input from stdin
echo '{"name": "nginx"}' | weaver-mustache render -t config.mustache

# Render with YAML file input
weaver-mustache render -t config.mustache -f data.yaml

# Render to output file
weaver-mustache render -t config.mustache -f data.yaml -o output.conf
```

## License

Copyright (c) Amadla. All rights reserved.
