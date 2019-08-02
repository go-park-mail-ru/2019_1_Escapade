package game

// sets the value in the min matrix
func (field *Field) setMatrixValue(x, y, v int32) {
	field.matrixM.Lock()
	field._matrix[x][y] = v
	field.matrixM.Unlock()
}

// incrementMatrixValue increments the counter of mines in
// the neighborhood
func (field *Field) incrementMatrixValue(x, y int32) {
	field.matrixM.Lock()
	field._matrix[x][y]++
	field.matrixM.Unlock()
}

// matrixFree free the memory allocated for the matrix
func (field *Field) matrixFree() {
	field.matrixM.Lock()
	for i := 0; i < len(field._matrix); i++ {
		field._matrix[i] = nil
	}
	field._matrix = nil
	field.matrixM.Unlock()
}

// History get slice of opened cells
func (field *Field) History() []*Cell {
	field.historyM.Lock()
	v := field._history
	field.historyM.Unlock()
	return v
}

// setHistory set new slice of opened cells
func (field *Field) setHistory(history []*Cell) {
	field.historyM.Lock()
	field._history = history
	field.historyM.Unlock()
}

// historyFree free the memory allocated for the slice
// of opened cells
func (field *Field) historyFree() {
	field.historyM.Lock()
	field._history = nil
	field.historyM.Unlock()
}

// lessThenMine returns true if there is a min counter
// in the cell located in the coordinates 'x','y'
func (field *Field) lessThenMine(x, y int32) bool {
	field.matrixM.RLock()
	v := field._matrix[x][y] < CellMine
	field.matrixM.RUnlock()
	return v
}

// matrixValue get the value from the min matrix
func (field *Field) matrixValue(x, y int32) int32 {
	field.matrixM.RLock()
	v := field._matrix[x][y]
	field.matrixM.RUnlock()
	return v
}

// setToHistory set the cell to the slice of opened cells
func (field *Field) setToHistory(cell *Cell) {
	field.historyM.Lock()
	defer field.historyM.Unlock()
	field._history = append(field._history, cell)
}

// decrementCellsLeft decrements the number of remaining cells
func (field *Field) decrementCellsLeft() {
	field.cellsLeftM.Lock()
	field._cellsLeft--
	field.cellsLeftM.Unlock()
}

// cellsLeft get the number of remaining cells
func (field *Field) cellsLeft() int32 {
	field.cellsLeftM.RLock()
	v := field._cellsLeft
	field.cellsLeftM.RUnlock()
	return v
}

// setCellsLeft set new number of remaining cells
func (field *Field) setCellsLeft(cellsLeft int32) {
	field.cellsLeftM.Lock()
	field._cellsLeft = cellsLeft
	field.cellsLeftM.Unlock()
}

// setDone sets the done flag to true, thus indicating that
// field is preparing to free memory
func (field *Field) setDone() {
	field.doneM.Lock()
	field._done = true
	field.doneM.Unlock()
}

// Done returns true if the field is preparing to free memory
func (field *Field) Done() bool {
	field.doneM.RLock()
	v := field._done
	field.doneM.RUnlock()
	return v
}
