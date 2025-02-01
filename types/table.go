package types

type TableColumn struct {
	Name    string
	Type    string
	Primary bool
	Unique  bool
	NotNull bool
	Default string
}
