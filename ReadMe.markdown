# useful.types

Branch: master — [![master|Build Status](https://travis-ci.org/nathanaelle/useful.types.svg?branch=master)](https://travis-ci.org/nathanaelle/useful.types)


## What is useful.types ?

useful.types is a collection of types for golang
there are essentialy for the marshaling and unmarshaling of data with a minimum of data verification

## Implemented interfaces

  * [fmt.Stringer](https://golang.org/pkg/fmt/#Stringer)
  * [flag.Value](https://golang.org/pkg/flag/#Value)
  * [json.Unmarshaler](https://golang.org/pkg/encoding/json/#Unmarshaler)
  * [json.Marshaler](https://golang.org/pkg/encoding/json/#Marshaler)
  * [toml.Unmarshaler](https://godoc.org/github.com/naoina/toml#Unmarshaler)


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
  * add tests
