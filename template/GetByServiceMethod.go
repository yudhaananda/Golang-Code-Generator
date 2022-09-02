func (s *[name]Service) Get[nameUpper]By[itemUpper]([itemParam]) (entity.[nameUpper], error) {

	[name], err := s.[name]Repository.FindBy[itemUpper]([item])

	if err != nil {
		return [name], err
	}

	if [name].Id == 0 {
		return [name], errors.New("[name] not found")
	}

	return [name], nil
}