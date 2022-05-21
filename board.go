package main

type board struct {
	currentGenCells [][]bool
	nextGenCells    [][]bool
	generation      int
	width, height   int
}

func (b *board) ifAlive(x, y int) int {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return 0
	}

	if b.currentGenCells[y][x] {
		return 1
	}

	return 0
}

func (b *board) countNeighbours(x, y int) int {
	sum := 0

	sum += b.ifAlive(x-1, y-1)
	sum += b.ifAlive(x, y-1)
	sum += b.ifAlive(x+1, y-1)

	sum += b.ifAlive(x-1, y)
	sum += b.ifAlive(x+1, y)

	sum += b.ifAlive(x-1, y+1)
	sum += b.ifAlive(x, y+1)
	sum += b.ifAlive(x+1, y+1)

	return sum
}

func (b *board) nextGen() {
	b.computeNextGen()
	b.generation++

	b.currentGenCells, b.nextGenCells = b.nextGenCells, b.currentGenCells
}

func (b *board) computeNextGen() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			n := b.countNeighbours(x, y)

			if b.currentGenCells[y][x] {
				b.nextGenCells[y][x] = n == 2 || n == 3
			} else {
				b.nextGenCells[y][x] = n == 3
			}
		}
	}
}

func (b *board) createGrid(w, h int) {
	b.currentGenCells = make([][]bool, h)
	b.nextGenCells = make([][]bool, h)
	for y := 0; y < h; y++ {
		b.currentGenCells[y] = make([]bool, w)
		b.nextGenCells[y] = make([]bool, w)
	}

	b.width = w
	b.height = h
}

func (b *board) ensureGridSize(w, h int) {
	if w <= 0 || h <= 0 { // for some reason we can be packed in below minsize on mobile - fyne#718
		return
	}
	yDelta := h - b.height
	xDelta := w - b.width

	if xDelta > 0 {
		// extend existing rows
		for y := 0; y < b.height; y++ {
			b.currentGenCells[y] = append(b.currentGenCells[y], make([]bool, xDelta)...)
			b.nextGenCells[y] = append(b.nextGenCells[y], make([]bool, xDelta)...)
		}
	}

	if yDelta > 0 {
		// add empty rows
		b.currentGenCells = append(b.currentGenCells, make([][]bool, yDelta)...)
		b.nextGenCells = append(b.nextGenCells, make([][]bool, yDelta)...)
		for y := b.height; y < h; y++ {
			b.currentGenCells[y] = make([]bool, w)
			b.nextGenCells[y] = make([]bool, w)
		}
	}

	b.width = w
	b.height = h
}

func (b *board) load() {
	// gun
	b.currentGenCells[5][1] = true
	b.currentGenCells[5][2] = true
	b.currentGenCells[6][1] = true
	b.currentGenCells[6][2] = true

	b.currentGenCells[3][13] = true
	b.currentGenCells[3][14] = true
	b.currentGenCells[4][12] = true
	b.currentGenCells[4][16] = true
	b.currentGenCells[5][11] = true
	b.currentGenCells[5][17] = true
	b.currentGenCells[6][11] = true
	b.currentGenCells[6][15] = true
	b.currentGenCells[6][17] = true
	b.currentGenCells[6][18] = true
	b.currentGenCells[7][11] = true
	b.currentGenCells[7][17] = true
	b.currentGenCells[8][12] = true
	b.currentGenCells[8][16] = true
	b.currentGenCells[9][13] = true
	b.currentGenCells[9][14] = true

	b.currentGenCells[1][25] = true
	b.currentGenCells[2][23] = true
	b.currentGenCells[2][25] = true
	b.currentGenCells[3][21] = true
	b.currentGenCells[3][22] = true
	b.currentGenCells[4][21] = true
	b.currentGenCells[4][22] = true
	b.currentGenCells[5][21] = true
	b.currentGenCells[5][22] = true
	b.currentGenCells[6][23] = true
	b.currentGenCells[6][25] = true
	b.currentGenCells[7][25] = true

	b.currentGenCells[3][35] = true
	b.currentGenCells[3][36] = true
	b.currentGenCells[4][35] = true
	b.currentGenCells[4][36] = true

	// spaceship
	b.currentGenCells[34][2] = true
	b.currentGenCells[34][3] = true
	b.currentGenCells[34][4] = true
	b.currentGenCells[34][5] = true
	b.currentGenCells[35][1] = true
	b.currentGenCells[35][5] = true
	b.currentGenCells[36][5] = true
	b.currentGenCells[37][1] = true
	b.currentGenCells[37][4] = true
}

func newBoard(minX, minY int) *board {
	b := &board{}
	b.createGrid(minX, minY)

	return b
}
