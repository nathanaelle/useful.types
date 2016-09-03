# useful.types

![License](http://img.shields.io/badge/license-Simplified_BSD-blue.svg?style=flat) [![Go Doc](http://img.shields.io/badge/godoc-useful.types-blue.svg?style=flat)](http://godoc.org/github.com/nathanaelle/useful.types) [![Build Status](https://travis-ci.org/nathanaelle/useful.types.svg?branch=master)](https://travis-ci.org/nathanaelle/useful.types)


## What is useful.types ?

useful.types is a collection of types for golang
There are essentialy for the marshaling and unmarshaling of data with a minimum of data verification
useful.types provides also some useful functions.

## Implemented interfaces

  * [x] [fmt.Stringer](https://golang.org/pkg/fmt/#Stringer)
  * [x] [flag.Value](https://golang.org/pkg/flag/#Value)
  * [x] [json.Unmarshaler](https://golang.org/pkg/encoding/json/#Unmarshaler)
  * [x] [json.Marshaler](https://golang.org/pkg/encoding/json/#Marshaler)
  * [x] [encoding.TextMarshaler](https://golang.org/pkg/encoding/#TextMarshaler)
  * [x] [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler)
  * [x] [toml.Unmarshaler](https://godoc.org/github.com/naoina/toml#Unmarshaler)


## Provided types

  * [ ] BitSet
  * [x] Duration
  * [x] FQDN
  * [x] IpAddr
  * [x] path
  * [x] URL
  * [x] StoreSize
  * [x] UUID


## License
2-Clause BSD

## Todo / wish-list

  * write comments
  * clean some ugly stuff
  * add more tests
