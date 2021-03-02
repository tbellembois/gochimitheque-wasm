package ajax

import "sync"

// SafeQueryFilter is a concurent
// query filter used in Ajax requests.
type SafeQueryFilter struct {
	lock   sync.Mutex
	locked bool
	QueryFilter
}

func (s *SafeQueryFilter) Lock() {
	s.locked = true
	s.lock.Lock()
}

func (s *SafeQueryFilter) Unlock() {
	if s.locked {
		s.locked = false
		s.lock.Unlock()
	}
}

func (s *SafeQueryFilter) Clean() {
	s.Lock()
	s.QueryFilter = QueryFilter{}
	s.Unlock()
}

func (s *SafeQueryFilter) CleanExceptProduct() {

	backupProduct := s.QueryFilter.Product
	backupProductFilterLabel := s.QueryFilter.ProductFilterLabel
	s.Clean()
	s.QueryFilter.Product = backupProduct
	s.QueryFilter.ProductFilterLabel = backupProductFilterLabel

}
