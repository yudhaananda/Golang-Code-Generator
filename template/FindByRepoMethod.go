
func (r *[name]Repository) FindBy[itemUpper]([itemParam]) (entity.[nameUpper], error) {
	var [name] entity.[nameUpper]

	err := r.db.Where("[item_] = ? AND deleted_by = ?", [item], "").Find(&[name]).Error

	if err != nil {
		return [name], err
	}

	return [name], nil
}