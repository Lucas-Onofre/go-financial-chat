package marketdataprovider

type MarketDataProviderPort interface {
	GetMarketData(stockCommand string) (string, error)
}
