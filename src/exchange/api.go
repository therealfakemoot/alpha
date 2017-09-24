package exchange

import (
	"encoding/json"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PriceQueryStats struct {
	Count int

	Duration *time.Duration

	High float64
	Low  float64

	Open  float64
	Close float64

	VolumeFrom float64
	VolumeTo   float64
}

// PricePoint represents a price-point-in-time provided by an exchange.
type PricePoint struct {
	Time       int64   `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
	From       string
	To         string
}

// Key returns a string suitable for use in BuntDB.
func (p *PricePoint) Key() string {
	return p.From + ":" + p.To + ":" + string(p.Time)
}

// Embed returns a Rich Embed representation of a PricePoint.
//
// The only parameter is to allow for Embed() methods to return a "minimal" version, i.e.
// without graphics or abbreviated messages.
func (p *PricePoint) Embed(min bool) *dgo.MessageEmbed {
	e := &dgo.MessageEmbed{}
	e.Title = p.From + ":" + p.To
	ts := time.Unix(p.Time, 0).Format(time.RFC3339)
	e.Timestamp = ts
	e.Description = "Exchange rate for " + p.From + "at " + ts

	e.Provider = &dgo.MessageEmbedProvider{URL: "https://min-api.cryptocompare.com/", Name: "CryptoCompare"}
	e.Author = &dgo.MessageEmbedAuthor{URL: "https://min-api.cryptocompare.com/", Name: "CryptoCompare"}

	highField := &dgo.MessageEmbedField{}
	highField.Name = "High"
	highField.Inline = true
	highField.Value = fmt.Sprintf("%f", p.High)

	lowField := &dgo.MessageEmbedField{}
	lowField.Name = "Low"
	lowField.Inline = true
	lowField.Value = fmt.Sprintf("%f", p.Low)

	e.Fields = []*dgo.MessageEmbedField{highField, lowField}

	e.Footer = &dgo.MessageEmbedFooter{Text: "Caveat emptor."}

	return e
}

// PriceQuery defines the exchange response for the HistoMinute endpoint
type PriceQuery struct {
	From              string
	To                string
	Response          string        `json:"Response"`
	Type              int           `json:"Type"`
	Aggregated        bool          `json:"Aggregated"`
	Prices            []*PricePoint `json:"Data"`
	TimeTo            int64         `json:"TimeTo"`
	TimeFrom          int64         `json:"TimeFrom"`
	FirstValueInArray bool          `json:"FirstValueInArray"`
	ConversionType    struct {
		Type             string `json:"type"`
		ConversionSymbol string `json:"conversionSymbol"`
	} `json:"ConversionType"`
}

func (pq *PriceQuery) Stats() *PriceQueryStats {
	s := &PriceQueryStats{}

	high := 0.0
	low := 0.0
	open := 0.0
	close := 0.0
	volumefrom := 0.0
	volumeto := 0.0

	for _, p := range pq.Prices {
		high = high + p.High
		low = low + p.Low
		open = open + p.Open
		close = close + p.Close
		volumefrom = volumeto + p.VolumeFrom
		volumeto = volumeto + p.VolumeTo
	}
	s.High = high
	s.Low = low
	s.Open = open
	s.Close = close
	s.VolumeTo = volumeto
	s.VolumeFrom = volumefrom

	s.Count = len(pq.Prices)

	d := time.Unix(pq.TimeTo, 0).Sub(time.Unix(pq.TimeFrom, 0))
	s.Duration = &d

	return s

}

// Embed returns a Discord Rich Embed representation of a series of PricePoints
func (pq *PriceQuery) Embed(min bool) *dgo.MessageEmbed {
	e := &dgo.MessageEmbed{}
	e.Title = pq.From + ":" + pq.To

	d := time.Unix(pq.TimeTo, 0).Sub(time.Unix(pq.TimeFrom, 0))

	//ts := time.Unix(p.Time, 0).Format(time.RFC3339)
	//e.Timestamp = ts
	e.Description = "Stats for the past **" + d.String() + "**"

	stats := pq.Stats()

	high := &dgo.MessageEmbedField{}
	high.Name = "Average High"
	high.Value = fmt.Sprintf("[%f]", stats.High)
	high.Inline = true

	low := &dgo.MessageEmbedField{}
	low.Name = "Average Low"
	low.Value = fmt.Sprintf("[%f]", stats.Low)
	low.Inline = true

	timestamp := &dgo.MessageEmbedField{}
	timestamp.Name = "Self destruct timer"
	timestamp.Value = "5"
	timestamp.Inline = false

	e.Fields = []*dgo.MessageEmbedField{high, low, timestamp}

	if pq.Prices[0].High > pq.Prices[len(pq.Prices)-1].High {
		e.Color = 0x00FF00
	} else {
		e.Color = 0xFF0000
	}

	return e
}

// NewPricePoint whatever
func NewPricePoint(hp *PricePoint, fsym string, tsym string) *PricePoint {

	hp.To = fsym
	hp.From = tsym

	return hp
}

// NewPriceQuery creates a PriceQuery value with 0 or more PricePoints
//
// l, dictates how many results are desired. It looks like
//  the CryptoCompare API returns a minimum of 2 responses no matter what the limit
//  is set to.
// fsym is the 'origin' currency
// tsym is the 'target' currency
func NewPriceQuery(pq *PriceQuery, fsym string, tsym string) *PriceQuery {
	pq.From = fsym
	pq.To = tsym

	return pq
}

const apiBase = "https://min-api.cryptocompare.com/data/"

var (
	apiMinute = "histominute"
)

// HistoMinute shut up
func HistoMinute(l int, fsym string, tsym string) *PriceQuery {
	endpoint := apiBase + apiMinute

	v := url.Values{}

	v.Set("fsym", fsym)
	v.Set("tsym", tsym)
	v.Set("limit", fmt.Sprintf("%d", l))
	// v.Set("limit", string(l))
	v.Set("aggregate", "1")

	resp, err := http.Get(endpoint + "?" + v.Encode())
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
	}

	respBytes, _ := ioutil.ReadAll(resp.Body)

	r := PriceQuery{}
	json.Unmarshal(respBytes, &r)

	pq := NewPriceQuery(&r, tsym, fsym)

	for _, p := range pq.Prices {
		p.From = fsym
		p.To = tsym
	}

	return pq
}
