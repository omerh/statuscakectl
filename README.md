# statuscakectl

[Statuscake](https://www.statuscake.com/statuscake-long-page/?a_aid=5d6fc4349afd6&a_bid=af013c39) is a Website Uptime & Performance Monitoring.

They offer a free plan and paid plans.

This is a small binary written in Golnag which allow you to control statuscake via API

Currently allow you to create/list/delete uptime tests and ssl monitoring.

## Configuration

Set the following environment variables from statuscake in your machine

```bash
export STATUSCAKE_USER=your_statuscake_user
export STATUSCAKE_KEY=your_key
```

## Run with Docker

You can also use it with docker:

```bash
docker run -e STATUSCAKE_USER=your_statuscake_user -e STATUSCAKE_KEY=your_key omerha/statuscakectl:latest statuscakectl list ssl
```

## Examples

```bash
# listing
statuscakectl list ssl
statuscakectl list uptime

# create
statuscakectl create ssl -d domain.com
statuscakectl create uptime --domain https://www.domain.com --checkrate 30 --type HTTP

# delete
statuscakectl delete ssl -d domain.com
statuscakectl delete ssl --id 1111111
statuscakectl delete uptime -d https://www.domain.com
statuscakectl delete uptime --id 1111111
```

If you'ed like to test them out I would appriciate it if you do it via this affiliation [link](https://www.statuscake.com/statuscake-long-page/?a_aid=5d6fc4349afd6&a_bid=af013c39) to help support my time working on this cool tool.

## build

```bash
# Use go modules add the following env var GO111MODULE=on
go build -ldflags "-s -w"
```

### Contributing

Fork, implement, add tests, pull request, get my everlasting thanks and a respectable place here :).

### Copyright

Copyright (c) 2020 Omer Haim, [@omerhaim](http://twitter.com/omerhaim).
See [LICENSE](LICENSE) for further details.
