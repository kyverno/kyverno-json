package table

type Table struct {
	Rows []Row
}

func (t *Table) AddFailed(rows ...Row) {
	for _, row := range rows {
		if row.IsFailure {
			t.Rows = append(t.Rows, row)
		}
	}
}

func (t *Table) Add(rows ...Row) {
	t.Rows = append(t.Rows, rows...)
}

func (t *Table) GetRows() interface{} {
	return t.Rows
}
