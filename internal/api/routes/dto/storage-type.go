package dto

type FindStorageTypeDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StorageTypeFilterDTO struct {
	Paging
	Name string `query:"name"`
}

func (f *StorageTypeFilterDTO) GetLimit() int {
	return f.Limit
}

func (f *StorageTypeFilterDTO) SetLimit(limit int) {
	f.Limit = limit
}

func (f *StorageTypeFilterDTO) GetOffset() int {
	return f.Offset
}

func (f *StorageTypeFilterDTO) SetOffset(offset int) {
	f.Offset = offset
}

type UpdateStorageTypeDTO struct {
	Name string `json:"name"`
}

type CreateStorageTypeDTO struct {
	Name string `json:"name" validate:"required"`
}