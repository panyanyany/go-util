package etherscan_util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/cihub/seelog"
)

type TxTableItem struct {
	Tx     string
	Method string
	From   string
	Value  float64
	Time   time.Time
}

func NewTxTableItemFromSelection(selection *goquery.Selection) (r *TxTableItem, err error) {
	r = new(TxTableItem)
	selection.Find("td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			r.Tx = strings.TrimSpace(s.Text())
			break
		case 2: // method
			span := s.Find("span")
			r.Method = span.AttrOr("title", "-")
			break
		case 5: // time
			span := s.Find("span")
			dtStr, exists := span.Attr("title")
			if !exists {
				seelog.Errorf("not found datetime: tx=%v, text=%v", r.Tx, s.Text())
				return
			}
			tm, err := dateparse.ParseStrict(dtStr)
			if err != nil {
				err = fmt.Errorf("dateparse.ParseStrict(tx time): %w", err)
				seelog.Error(err)
				return
			}
			r.Time = tm.In(time.Local)
			break
		case 6: // method
			span := s.Find("span")
			r.From = span.Text()
			break
		case 9: // value
			valStr := strings.Split(s.Text(), " ")[0]
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				err = fmt.Errorf("strconv.ParseFloat(eth value): %w", err)
				seelog.Error(err)
				return
			}
			r.Value = val
			break
		}
	})
	return
}
