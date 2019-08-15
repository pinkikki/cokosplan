package port

import (
	"context"
	"reflect"

	"github.com/olivere/elastic/v7"
	"github.com/pinkikki/cokosplan/es"
)

type Expense struct {
	Event     string `json:"event"`
	MaxAmount int64  `json:"max_amount"`
	MinAmount int64  `json:"min_amount"`
	// TODO 日付型に変更したい
	PaymentDate string `json:"payment_date"`
}

type ExpenseCriteria struct {
	Event string `json:"event"`
	// TODO 日付型に変更したい
	FromPaymentDate string `json:"payment_date"`
	ToPaymentDate   string `json:"payment_date"`
}

func List(c ExpenseCriteria) ([]Expense, error) {
	client := es.NewClient()
	s := client.Search().Index("expense")
	withExpenseCriteria(s, c)

	result, err := s.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return translate(result)
}

func Save(e Expense) error {
	client := es.NewClient()

	_, err := client.Index().Index("expense").BodyJson(e).Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}

func translate(r *elastic.SearchResult) ([]Expense, error) {
	expenses := []Expense{}
	var expense Expense
	for _, item := range r.Each(reflect.TypeOf(expense)) {
		e := item.(Expense)
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func withExpenseCriteria(s *elastic.SearchService, c ExpenseCriteria) {
	withEventCriteria(s, c)
	withFromPaymentDateCriteria(s, c)
	withToPaymentDateCriteria(s, c)
}

func withEventCriteria(s *elastic.SearchService, c ExpenseCriteria) *elastic.SearchService {
	if c.Event != "" {
		return s.Query(elastic.NewMatchQuery(getMetaValue(c, 0), c.Event))
	}
	return s
}

func withFromPaymentDateCriteria(s *elastic.SearchService, c ExpenseCriteria) *elastic.SearchService {
	if c.FromPaymentDate != "" {
		return s.Query(elastic.NewRangeQuery(getMetaValue(c, 1)).Gte(c.FromPaymentDate))
	}
	return s
}

func withToPaymentDateCriteria(s *elastic.SearchService, c ExpenseCriteria) *elastic.SearchService {
	if c.ToPaymentDate != "" {
		return s.Query(elastic.NewRangeQuery(getMetaValue(c, 2)).Lte(c.ToPaymentDate))
	}
	return s
}

func getMetaValue(c ExpenseCriteria, i int) string {
  return reflect.TypeOf(c).Field(i).Tag.Get("json")
}
