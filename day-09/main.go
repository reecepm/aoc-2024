package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

type Block struct {
	ID  *int
	Pos int
}

type Disk struct {
	blocks []*int
}

func NewDisk(blocks []*int) *Disk {
	return &Disk{blocks: blocks}
}

func (d *Disk) Clone() *Disk {
	newBlocks := make([]*int, len(d.blocks))
	copy(newBlocks, d.blocks)
	return NewDisk(newBlocks)
}

func (d *Disk) Checksum() int {
	sum := 0
	for pos, block := range d.blocks {
		if block != nil {
			sum += pos * (*block)
		}
	}
	return sum
}

type FileInfo struct {
	ID        int
	Positions []int
	Size      int
}

func (d *Disk) GetFileInfo() []FileInfo {
	fileMap := make(map[int][]int)
	for pos, block := range d.blocks {
		if block != nil {
			fileMap[*block] = append(fileMap[*block], pos)
		}
	}

	var files []FileInfo
	for id, positions := range fileMap {
		files = append(files, FileInfo{
			ID:        id,
			Positions: positions,
			Size:      len(positions),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ID > files[j].ID
	})

	return files
}

func (d *Disk) CompactIndividualBlocks() {
	finalized := make([]bool, len(d.blocks))

	for i := range d.blocks {
		if d.blocks[i] == nil {
			for j := len(d.blocks) - 1; j >= 0; j-- {
				if d.blocks[j] != nil && !finalized[j] {
					d.blocks[i] = d.blocks[j]
					d.blocks[j] = nil
					finalized[i] = true
					break
				}
			}
			continue
		}
		finalized[i] = true
	}
}

func (d *Disk) CompactWholeFiles() {
	files := d.GetFileInfo()

	for _, file := range files {
		if gap := d.findSuitableGap(file.Size, file.Positions[0]); gap >= 0 {
			d.moveFile(file, gap)
		}
	}
}

func (d *Disk) findSuitableGap(size int, currentPos int) int {
	gapStart := -1
	gapSize := 0

	for i := 0; i < currentPos; i++ {
		if d.blocks[i] == nil {
			if gapStart == -1 {
				gapStart = i
			}
			gapSize++
			if gapSize >= size {
				return gapStart
			}
		} else {
			gapStart = -1
			gapSize = 0
		}
	}
	return -1
}

func (d *Disk) moveFile(file FileInfo, newPos int) {
	for _, pos := range file.Positions {
		d.blocks[pos] = nil
	}

	for i := 0; i < file.Size; i++ {
		value := file.ID
		d.blocks[newPos+i] = &value
	}
}

func main() {
	blocks, err := parseInput("input.txt")
	if err != nil {
		panic(err)
	}

	p1Input := make([]*int, len(blocks))
	copy(p1Input, blocks)
	p2Input := make([]*int, len(blocks))
	copy(p2Input, blocks)

	p1 := partOne(p1Input)
	p2 := partTwo(p2Input)

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", p2)
}

func partOne(blocks []*int) int {
	disk := NewDisk(blocks)
	disk.CompactIndividualBlocks()
	return disk.Checksum()
}

func partTwo(blocks []*int) int {
	disk := NewDisk(blocks)
	disk.CompactWholeFiles()
	return disk.Checksum()
}

func parseInput(inputPath string) ([]*int, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []*int
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		for i := 0; i < len(line); i++ {
			parsedChar, err := strconv.Atoi(string(line[i]))
			if err != nil {
				return nil, err
			}

			var block *int

			if i%2 == 0 {
				adjusted := i / 2
				block = &adjusted
			}

			for j := 0; j < parsedChar; j++ {
				blocks = append(blocks, block)
			}
		}
	}

	return blocks, nil
}
