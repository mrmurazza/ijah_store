package impl

import (
	"ijah-store/domain/item"
	"ijah-store/domain/request"
)

type service struct {
	repo item.Repository
}

func NewService(repo item.Repository) item.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateItemIfNotAny(request request.RestockOrderRequest) {
	i := item.Item{
		SKU: request.SKU,
		Name: request.ItemName,
		// stock will be updated separately
	}

	isExist := s.repo.IsExist(i)
	if !isExist {
		s.repo.Persist(i)
	}
}

func (s *service) UpdateItemStock(sku string, quantity int) {
	i := s.repo.GetItem(sku)

	i.Stock += quantity
	s.repo.UpdateStock(i)
}

// GetRequestedItemMap Given a list of SKU, this will return a map containing SKU -> its appropriate Item struct
func (s *service) GetRequestedItemMap(requestedSkuList []string) map[string]*item.Item {
	var itemMap = make(map[string]*item.Item)
	items := s.repo.GetItems(requestedSkuList[:])
	for _, product := range items {
		itemMap[product.SKU] = product
	}
	return itemMap
}

func (s *service) GetAllItems() []*item.Item {
	return s.repo.GetAllItems()
}
