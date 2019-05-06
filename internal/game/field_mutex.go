package game

// setMatrixValue set a value to matrix
func (field *Field) setMatrixValue(x, y, v int) {
	field.matrixM.Lock()
	field.Matrix[x][y] = v
	field.matrixM.Unlock()
}

// setMatrixValue set a value to matrix
func (field *Field) incrementMatrixValue(x, y int) {
	field.matrixM.Lock()
	field.Matrix[x][y]++
	field.matrixM.Unlock()
}

// setMatrixValue set a value to matrix
func (field *Field) matrixFree() {
	field.matrixM.Lock()
	field.Matrix = nil
	field.matrixM.Unlock()
}

// setMatrixValue set a value to matrix
func (field *Field) historyFree() {
	field.historyM.Lock()
	field.History = nil
	field.historyM.Unlock()
}

// setMatrixValue set a value to matrix
func (field *Field) lessThenMine(x, y int) bool {
	field.matrixM.RLock()
	v := field.Matrix[x][y] < CellMine
	field.matrixM.RUnlock()
	return v
}

// getMatrixValue get a value from matrix
func (field *Field) getMatrixValue(x, y int) int {
	field.matrixM.RLock()
	v := field.Matrix[x][y]
	field.matrixM.RUnlock()
	return v
}

// setMatrixValue set a value to matrix
func (field *Field) setToHistory(cell Cell) {
	field.historyM.Lock()
	defer field.historyM.Unlock()
	field.History = append(field.History, cell)
}

// setMatrixValue set a value to matrix
func (field *Field) decrementCellsLeft() {
	field.cellsLeftM.Lock()
	field.CellsLeft--
	field.cellsLeftM.Unlock()
}

// getMatrixValue get a value from matrix
func (field *Field) getCellsLeft() int {
	field.cellsLeftM.RLock()
	v := field.CellsLeft
	field.cellsLeftM.RUnlock()
	return v
}

// setMatrixValue set a value to matrix
func (field *Field) setDone() {
	field.doneM.Lock()
	field.done = true
	field.doneM.Unlock()
}

// getMatrixValue get a value from matrix
func (field *Field) getDone() bool {
	field.doneM.RLock()
	v := field.done
	field.doneM.RUnlock()
	return v
}
