API for bit.dj

## Getting started:
  You will need postgres as well as redis pre-installed with no special configurations

  Install Deps:
  `$ ./script/bootstrap`

  Provision database:
  `$ ./script/recycle`

  Test it:
  `$ ./script/test`

  Build it:
  `$ go build`

  Run it:
  `./Exgo`


## Resources
`POST /user`
```json
{
  "userName": "Bob",
  "email": "BobAndAllThingsBob@Bob.com",
  "password": "Boberson"
}

```

`GET /socket`
`Pretty self explanitory. Requires auth atm.`

`POST /session`
```json
{
  "userName": "Bob",
  "password": "Boberson"
}

```

Returns:
```json
{"token" : "<sessionToken>"}
```
