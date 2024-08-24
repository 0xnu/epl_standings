## EPL Standings

Acquire [BBC Sport](https://www.bbc.co.uk/sport/football/tables) English Premier League data with [Go](https://go.dev/).

> For educational purposes only! 😎✌🏽

### Getting Started

Export environment variables:

```sh
export SCRAPER_API_KEY=
export FOOTBALL_HOST=localhost
export FOOTBALL_USER=root
export FOOTBALL_PASS=rootroot
export FOOTBALL_DBNAME=epl_data
```

Run the commands below in your terminal: 

```sh
go mod init epl_standing
go mod tidy
go run epl_standing.go

## Optional
go build
./epl_standing
```

### License

This project is licensed under the [BSD 3-Clause](LICENSE) License.

### Copyright

(c) 2020 - 2024 [Finbarrs Oketunji](https://finbarrs.eu).