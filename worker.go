package drudge

type Worker interface {
	Insert() error
	Load() error
	Update() error
	Delete() error
}
