package main

type board struct {
	cells         [][]bool
	nextGenCells         [][]bool
	ch chan bool
	generation    int
	width, height int
}

func (b *board) ifAlive(x, y int) int {
	if x < 0 || x >= b.width {
		return 0
	}

	if y < 0 || y >= b.height {
		return 0
	}

	if b.cells[y][x] {
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
	b.nextGenCells, b.cells = b.cells, b.nextGenCells
}

func (b *board) computeNextGen() {
	width, height := b.width, b.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			go func(x,y int) {
				n := b.countNeighbours(x, y)
				if b.cells[y][x] {
					b.nextGenCells[y][x] = n == 2 || n == 3
				} else {
					b.nextGenCells[y][x] = n == 3
				}
				b.ch<-true
			}(x,y)
		}
	}
	count := height * width
	for i := 0; i < count; i++ {
		<-b.ch
	}
}

func (b *board) createGrid(w, h int) {
	b.cells = make([][]bool, h)
	b.nextGenCells =  make([][]bool, h)
	for y := 0; y < h; y++ {
		b.cells[y] = make([]bool, w)
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
		for i := range b.cells {
			b.cells[i] = append(b.cells[i], make([]bool, xDelta)...)
			b.nextGenCells[i] = append(b.nextGenCells[i], make([]bool, xDelta)...)
		}
	}

	if yDelta > 0 {
		// add empty rows
		b.cells = append(b.cells, make([][]bool, yDelta)...)
		b.nextGenCells = append(b.nextGenCells, make([][]bool, yDelta)...)
		for y := b.height; y < h; y++ {
			b.cells[y] = make([]bool, w)
			b.nextGenCells[y] = make([]bool, w)
		}
	}

	b.width = w
	b.height = h
}

func (b *board) load() {
	// gun
	b.cells[5][1] = true
	b.cells[5][2] = true
	b.cells[6][1] = true
	b.cells[6][2] = true

	b.cells[3][13] = true
	b.cells[3][14] = true
	b.cells[4][12] = true
	b.cells[4][16] = true
	b.cells[5][11] = true
	b.cells[5][17] = true
	b.cells[6][11] = true
	b.cells[6][15] = true
	b.cells[6][17] = true
	b.cells[6][18] = true
	b.cells[7][11] = true
	b.cells[7][17] = true
	b.cells[8][12] = true
	b.cells[8][16] = true
	b.cells[9][13] = true
	b.cells[9][14] = true

	b.cells[1][25] = true
	b.cells[2][23] = true
	b.cells[2][25] = true
	b.cells[3][21] = true
	b.cells[3][22] = true
	b.cells[4][21] = true
	b.cells[4][22] = true
	b.cells[5][21] = true
	b.cells[5][22] = true
	b.cells[6][23] = true
	b.cells[6][25] = true
	b.cells[7][25] = true

	b.cells[3][35] = true
	b.cells[3][36] = true
	b.cells[4][35] = true
	b.cells[4][36] = true

	// spaceship
	b.cells[34][2] = true
	b.cells[34][3] = true
	b.cells[34][4] = true
	b.cells[34][5] = true
	b.cells[35][1] = true
	b.cells[35][5] = true
	b.cells[36][5] = true
	b.cells[37][1] = true
	b.cells[37][4] = true
}

func newBoard(minX, minY int) *board {
	b := &board{ch: make(chan bool)}
	b.createGrid(minX, minY)

	return b
}
