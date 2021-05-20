package items

import (
	"encoding/json"
	"errors"

	"github.com/chernyshev-alex/bookstore_items-api/client/es"
	"github.com/chernyshev-alex/bookstore_items-api/domain/queries"
	"github.com/chernyshev-alex/bookstore_utils_go/rest_errors"
)

const (
	itemName = "items"
	typeItem = "_doc"
)

func (it *Item) Save() rest_errors.RestErr {
	result, err := es.Client.Index(itemName, typeItem, it)
	if err != nil {
		rest_errors.NewInternalServerError("save item error", err)
	}

	it.Id = result.Id
	return nil
}

func (it *Item) Get() rest_errors.RestErr {
	result, err := es.Client.Get(itemName, typeItem, it.Id)
	if err != nil {
		rest_errors.NewInternalServerError("save item error", err)
	}

	it.Id = result.Id
	return nil
}

func (it Item) Search(q queries.EsQuery) ([]Item, rest_errors.RestErr) {
	result, err := es.Client.Search(itemName, q.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerError("save item error", errors.New("ELK error"))
	}

	items := make([]Item, result.TotalHits())
	for idx, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()

		var item Item
		if err = json.Unmarshal(bytes, &item); err != nil {
			return nil, rest_errors.NewInternalServerError("elk parse search response error", errors.New("ELK error"))
		}

		item.Id = hit.Id
		items[idx] = item
	}

	if len(items) == 0 {
		return nil, rest_errors.NewNotFoundError("not found")
	}

	return items, nil
}
