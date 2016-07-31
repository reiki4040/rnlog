rnlog
===

simple logging library wrapper.
output json format log to logger.

## rnlog do

- rnlog can change log level dynamically.

### log item

output below items with JSON format
- time (RFC3339)
- level (TRACE/DEBUG/INFO/ERROR, NOTICE/FATAL)
- message (string)
- items (map[string]interface{})

sample formated with jq
```
{
    "time": "2016-07-31T01:02:25+09:00",
    "level": "NOTICE",
    "message": "starting rnlog...",
    "items": {}
}
```

## TODO

- benchmark and improve
- consider interface
- testing
