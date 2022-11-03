# Snowflake News

Run a server that presents recent Snowflake news stories.

## Getting Started

Run:
```
$ go run main.go
```

Next, make a GET request to `localhost:8282/news`, and you'll receive recent news stories about Snowflake!

Example:
```
$ curl localhost:8282/news
Recent News Stories
===============================

Why Salesforce, Snowflake, and CrowdStrike All Cratered on Wednesday -- November, 2 2022 15:41PM
-------------------------------
Summary: A broad cross section of the stock market tumbled on Wednesday, as investors focused on the macro economy and the Federal Reserve Bank's ongoing campaign to battle runaway inflation, which has remained stubbornly near 40-year highs.  The latest Fed rate hike and the corresponding commentary did little to calm jittery investors.  With that as a backdrop, shares of Salesforce (NYSE: CRM) slipped 6.1%, Snowflake (NYSE: SNOW) stock was down 7.4%, and CrowdStrike (NASDAQ: CRWD) slipped as much as 7.8%.
Read more: https://finnhub.io/api/news?id=2fd7e72f1340b695827504d11a6014e699cc6dadefb00613c182aa017661b305

===============================

An Intrinsic Calculation For Snowflake Inc. (NYSE:SNOW) Suggests It's 37% Undervalued -- November, 2 2022 09:49AM
-------------------------------
Summary: In this article we are going to estimate the intrinsic value of Snowflake Inc. ( NYSE:SNOW ) by taking the expected...
Read more: https://finnhub.io/api/news?id=8d17f5ff3855177c4b93b8198cf65e07661cfd5c605514bcb03a679c2f120e14

===============================

Street Wrap: Today's Top 15 Upgrades, Downgrades, Initiations -- November, 2 2022 02:09AM
-------------------------------
Summary: Abiomed upgrade, Twilio downgrade and Microsoft initiation among today's top calls on Wall Street ABMD, JNJ, KKR, ARCT, CSLLY, DNLI, PINC, TWLO, NDAQ, AVA, BPMC, ECL, MSFT, SNOW, ZS, WDAY, PANW, OKTA, NET, DT, DDOG, CHKP, HUBS, CRWD, LEGN, HLMN, SRNE
Read more: https://finnhub.io/api/news?id=20075d5e407661aea1905b05da108b0e8fd17e3bc7177d1f11e7d4d8516a80cc
```

## Server

The server exposes the 3 most recent news stories mentioning SNOW published in the past 3 days.

Details:
- Hosted at `localhost:8282` by default, and exposes a single endpoint, `/news`.
- Pulls data from finnhub.io's API. Specifically, it hits the `/api/v1/company-news` endpoint. See https://finnhub.io/docs/api/company-news.
- Serves cached data. The cache is updated once every minute.

## Improvements

Potential future improvements to this server:
- Allow for querying news from any company ticker (using query parameter)
- Allow for querying other types of information (stock price, market news)

Potential future code improvements:
- Add unit testing
- Add integration test
- Create Collector interface for different types of collectors (News, StockPrice, MarketNews, etc.)
