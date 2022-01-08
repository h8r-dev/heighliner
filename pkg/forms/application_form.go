package forms

type ApplicationForm struct {
	StackID    string                 `form:"stack_id" json:"stack_id" binding:"required,min=3,max=100"`
	Parameters map[string]interface{} `form:"parameters" json:"parameters" binding:"-"`
}
