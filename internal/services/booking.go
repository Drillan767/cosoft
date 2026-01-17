package services

func (s *Service) UpdateCredits() (*float64, error) {
	return s.store.UpdateCredits()
}
