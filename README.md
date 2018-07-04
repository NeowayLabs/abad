# abad

[![GoDoc](https://godoc.org/github.com/NeowayLabs/abad?status.svg)](https://godoc.org/github.com/NeowayLabs/abad)
[![Build Status](https://travis-ci.org/NeowayLabs/abad.svg?branch=master)](https://travis-ci.org/NeowayLabs/abad)
[![Go Report Card](https://goreportcard.com/badge/github.com/NeowayLabs/abad)](https://goreportcard.com/report/github.com/NeowayLabs/abad)

Abad stands for [Abaddon](https://en.wikipedia.org/wiki/Abaddon) the destroyer and torturer of men.

Why this name ? Because developing this will be torture, we have been forsaken by God and
left on the hands of a torturing angel (we must have done something pretty awful).

## End to End Tests

To run the end to end tests you can run:

```make
make test-e2e
```

These tests assume that you have both d8 and abad installed on the host that is
running the tests. d8 is a repl/debugger from the
[V8](https://developers.google.com/v8/) engine, which is our reference implementation
used to compare if we are interpreting JavaScript code correctly.

As you can imagine, compiling V8 is not exactly easy...or fun. So we provide a docker
image for this. To run the tests inside the docker dev environment you can run:

```make
make dev-test-e2e
```

Adding new tests is very easy, just add new code samples to **tests/e2e/testdata** and
they will imediately be tested when you run the end to end tests.