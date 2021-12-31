package vo

type Supplier struct {
	ID int64 `json:"id,omitempty"`
	// CreatedTime holds the value of the "created_time" field.
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Address holds the value of the "address" field.
	Address         string      `json:"address,omitempty"`
	Constacts       []*Constact `json:"Constacts,omitempty"`
	RemoveConstacts []*Constact `json:"removeConstacts,omitempty"`
}

type Constact struct {
	ID int64 `json:"id,omitempty"`
	// CreatedTime holds the value of the "created_time" field.
	// SupplierID holds the value of the "supplier_id" field.
	SupplierID int64 `json:"supplier_id,omitempty"`
	// Firstname holds the value of the "firstname" field.
	Firstname string `json:"firstname,omitempty"`
	// Lastname holds the value of the "lastname" field.
	Lastname string `json:"lastname,omitempty"`
	// Title holds the value of the "title" field.
	Title           string      `json:"title,omitempty"`
	Relations       []*Relation `json:"relations,omitempty"`
	RemoveRelations []*Relation `json:"removeRelations,omitempty"`
}

type Relation struct {
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// CreatedTime holds the value of the "created_time" field.
	// ContactID holds the value of the "contact_id" field.
	ContactID int64 `json:"contact_id,omitempty"`
	// Type holds the value of the "type" field.
	Type int8 `json:"type,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
}
