package catalog

import "strings"

type Coin struct {
	ID        string
	Name      string
	Ticker    string
	IconAsset string
}

var trackedCoins = []Coin{
	{ID: "bitcoin", Name: "Bitcoin", Ticker: "BTC", IconAsset: "coins/bitcoin.png"},
	{ID: "ethereum", Name: "Ethereum", Ticker: "ETH", IconAsset: "coins/ethereum.png"},
	{ID: "the-open-network", Name: "TON Coin", Ticker: "TON", IconAsset: "coins/the-open-network.png"},
	{ID: "solana", Name: "Solana", Ticker: "SOL", IconAsset: "coins/solana.png"},
	{ID: "dogecoin", Name: "Dogecoin", Ticker: "DOGE", IconAsset: "coins/dogecoin.png"},
	{ID: "ripple", Name: "Ripple", Ticker: "XRP", IconAsset: "coins/ripple.png"},
	{ID: "litecoin", Name: "Litecoin", Ticker: "LTC", IconAsset: "coins/litecoin.png"},
}

var aliases = map[string][]string{
	"bitcoin":          {"btc", "btc-bitcoin"},
	"ethereum":         {"eth", "eth-ethereum"},
	"the-open-network": {"ton", "toncoin", "ton-toncoin", "toncoin-toncoin"},
	"solana":           {"sol", "sol-solana"},
	"dogecoin":         {"doge", "doge-dogecoin"},
	"ripple":           {"xrp", "xrp-xrp"},
	"litecoin":         {"ltc", "ltc-litecoin"},
}

var (
	byID       map[string]Coin
	aliasToID  map[string]string
	trackedIDs []string
)

func init() {
	byID = make(map[string]Coin, len(trackedCoins))
	aliasToID = make(map[string]string, len(trackedCoins)*4)
	trackedIDs = make([]string, 0, len(trackedCoins))

	for _, coin := range trackedCoins {
		byID[coin.ID] = coin
		trackedIDs = append(trackedIDs, coin.ID)

		registerAlias(coin.ID, coin.ID)
		registerAlias(coin.ID, coin.Ticker)
		for _, alias := range aliases[coin.ID] {
			registerAlias(coin.ID, alias)
		}
	}
}

func Tracked() []Coin {
	coins := make([]Coin, len(trackedCoins))
	copy(coins, trackedCoins)
	return coins
}

func IDs() []string {
	ids := make([]string, len(trackedIDs))
	copy(ids, trackedIDs)
	return ids
}

func IDsCSV() string {
	return strings.Join(trackedIDs, ",")
}

func Lookup(id string) (Coin, bool) {
	coin, ok := byID[normalize(id)]
	return coin, ok
}

func CanonicalID(values ...string) string {
	for _, value := range values {
		id, ok := aliasToID[normalize(value)]
		if ok {
			return id
		}
	}
	return ""
}

func registerAlias(id, alias string) {
	normalized := normalize(alias)
	if normalized == "" {
		return
	}
	aliasToID[normalized] = id
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
