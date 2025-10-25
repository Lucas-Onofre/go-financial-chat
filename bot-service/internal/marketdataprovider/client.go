package marketdataprovider

import "github.com/go-resty/resty/v2"

var (
	stooqUrl = "https://stooq.com/q/l/"
)

type Client struct {
	client *resty.Client
}

func New() MarketDataProviderPort {
	return &Client{
		client: resty.New(),
	}
}

func (c *Client) GetMarketData(stockCommand string) (string, error) {
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"s": stockCommand,
			"f": "sd2t2ohlcv",
			"e": "csv",
			// TODO: verify usage of 'h' param
		}).
		Get(stooqUrl)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
