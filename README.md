# Snowflake News

Run a server that presents recent Snowflake news stories.

## Getting Started

Simply run:
```
$ go run main.go
```

Next, simply make a GET request to `localhost:8282/news`, and you'll receive recent news stories about Snowflake!

## Server

The server exposes the most recent news stories mentioning SNOW published in the past 3 days.

The server is hosted at `localhost:8282` by default, and exposes a single endpoint, `/news`

The server pulls data from finnhub.io's API. Specifically, it hits the `/api/v1/company-news` endpoint. See https://finnhub.io/docs/api/company-news.

The server caches data from the API endpoint, and updates this cache once every minute. 

## Improvements

Potential future improvements to this server:
- Allow for querying news from any company ticker (using query parameter)
- Allow for querying other types of information (stock price, market news)

Potential future code improvements:
- Add unit testing
- Add integration test
- Create Collector interface for different types of collectors (News, StockPrice, MarketNews, etc.)
