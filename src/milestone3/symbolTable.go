package milestone3

import "fmt"

type ObjectClass int

const (
	ObjConstant ObjectClass = iota
	ObjVariable
	ObjType
	ObjProcedure
	ObjFunction
	ObjProgram
	ObjField
)

func (o ObjectClass) String() string {
	switch o {
	case ObjConstant:
		return "constant"
	case ObjVariable:
		return "variable"
	case ObjType:
		return "type"
	case ObjProcedure:
		return "procedure"
	case ObjFunction:
		return "function"
	case ObjProgram:
		return "program"
	case ObjField:
		return "field"
	default:
		return "unknown"
	}
}

type TypeKind int

const (
	TypeNone TypeKind = iota
	TypeInteger
	TypeBoolean
	TypeChar
	TypeReal
	TypeArray
	TypeRecord
)

func (t TypeKind) String() string {
	switch t {
	case TypeInteger:
		return "integer"
	case TypeBoolean:
		return "boolean"
	case TypeChar:
		return "char"
	case TypeReal:
		return "real"
	case TypeArray:
		return "array"
	case TypeRecord:
		return "record"
	case TypeNone:
		return "void"
	default:
		return "unknown"
	}
}

// Menyimpan informasi identifier (konstanta, variabel, tipe, prosedur, fungsi)
type TabEntry struct {
	Identifier string
	Link       int         // Pointer ke identifier sebelumnya dalam scope yang sama (linked list)
	Obj        ObjectClass // Kelas objek (constant, variable, type, procedure, function)
	Type       TypeKind    // Tipe dasar identifier
	Ref        int         // Pointer ke atab/btab untuk tipe komposit (-1 jika tidak ada)
	Nrm        int         // Normal variable (1) atau var parameter (0)
	Lev        int         // Lexical level (0=global, 1=prosedur level 1, dst)
	Adr        int         // Address/offset atau nilai konstanta
}

// Menyimpan informasi detail array
type AtabEntry struct {
	Xtyp int // Tipe indeks array
	Etyp int // Tipe elemen array
	Eref int // Pointer ke detail elemen jika elemen adalah tipe komposit (-1 jika tidak ada)
	Low  int // Batas bawah indeks array
	High int // Batas atas indeks array
	Elsz int // Ukuran satu elemen array (dalam byte/unit memori)
	Size int // Total ukuran array (Elsz * (High - Low + 1))
}

// Menyimpan informasi block (prosedur, fungsi, record type)
type BtabEntry struct {
	Last int // Pointer ke identifier terakhir yang dideklarasikan di block ini
	Lpar int // Pointer ke parameter terakhir (0 untuk record)
	Psze int // Total ukuran parameter
	Vsze int // Total ukuran variabel lokal
}

type SymbolTable struct {
	Tab  []TabEntry
	Btab []BtabEntry
	Atab []AtabEntry

	// Indeks pointer
	TabIndex  int // Current index in Tab (start from 29 for reserved words)
	BtabIndex int // Current index in Btab
	AtabIndex int // Current index in Atab

	// Block/Scope management
	CurrentLevel int   // Current lexical level (0 = global)
	CurrentBlock int   // Current block index in Btab
	Display      []int // Display register untuk akses variabel di scope berbeda

	// Reserved words offset
	ReservedWordsCount int
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		Tab:  make([]TabEntry, 0, 1000),
		Btab: make([]BtabEntry, 0, 100),
		Atab: make([]AtabEntry, 0, 100),

		TabIndex:           0,
		BtabIndex:          0,
		AtabIndex:          0,
		CurrentLevel:       0,
		CurrentBlock:       -1,
		Display:            make([]int, 10), // Max 10 nested levels
		ReservedWordsCount: 29,
	}

	// Initialize global block (block 0) at level 0
	blockIndex := st.enterBlock()
	st.Display[0] = blockIndex
	st.CurrentLevel = 0
	st.CurrentBlock = blockIndex

	st.initReservedWords()

	return st
}

func (st *SymbolTable) initReservedWords() {
	reservedWords := []string{
		"dan", "larik", "mulai", "kasus", "konstanta", "bagi", "turun_ke",
		"lakukan", "selain_itu", "selesai", "untuk", "fungsi", "jika",
		"mod", "tidak", "dari", "atau", "prosedur", "program", "rekaman",
		"ulangi", "string", "maka", "ke", "tipe", "sampai", "variabel", "selama", "padat",
	}

	for i, word := range reservedWords {
		st.Enter(word, ObjConstant, TypeNone, 0, 1, i)
	}

	st.TabIndex = len(reservedWords)
}

// Masuk ke block baru (prosedur, fungsi, atau program utama)
func (st *SymbolTable) enterBlock() int {
	blockIndex := st.BtabIndex

	st.Btab = append(st.Btab, BtabEntry{
		Last: -1, // Belum ada identifier di block ini
		Lpar: 0,
		Psze: 0,
		Vsze: 0,
	})

	st.BtabIndex++
	st.CurrentBlock = blockIndex
	// Don't set Display here - let caller decide

	return blockIndex
}

// Masuk ke nested level baru
func (st *SymbolTable) enterLevel() {
	st.CurrentLevel++
	if st.CurrentLevel >= len(st.Display) {
		// Extend display if needed
		st.Display = append(st.Display, make([]int, 5)...)
	}
	// Update CurrentBlock from Display for this level
	if st.CurrentLevel < len(st.Display) {
		st.CurrentBlock = st.Display[st.CurrentLevel]
	}
}

// Masuk ke nested level baru dengan block baru
func (st *SymbolTable) enterLevelWithBlock() int {
	st.CurrentLevel++
	if st.CurrentLevel >= len(st.Display) {
		// Extend display if needed
		st.Display = append(st.Display, make([]int, 5)...)
	}
	// Create new block for this level
	blockIndex := st.enterBlock()
	st.Display[st.CurrentLevel] = blockIndex
	return blockIndex
}

// Keluar dari nested level
func (st *SymbolTable) exitLevel() {
	if st.CurrentLevel > 0 {
		st.CurrentLevel--
		st.CurrentBlock = st.Display[st.CurrentLevel]
	}
}

// Menambahkan identifier baru ke symbol table
func (st *SymbolTable) Enter(
	identifier string,
	obj ObjectClass,
	typ TypeKind,
	ref int,
	nrm int,
	adr int,
) int {
	index := st.TabIndex

	link := -1
	if st.CurrentBlock >= 0 && st.CurrentBlock < len(st.Btab) {
		link = st.Btab[st.CurrentBlock].Last
		st.Btab[st.CurrentBlock].Last = index
	}

	entry := TabEntry{
		Identifier: identifier,
		Link:       link,
		Obj:        obj,
		Type:       typ,
		Ref:        ref,
		Nrm:        nrm,
		Lev:        st.CurrentLevel,
		Adr:        adr,
	}

	st.Tab = append(st.Tab, entry)
	st.TabIndex++

	return index
}

// Menambahkan entry ke array table
func (st *SymbolTable) EnterArray(xtyp, etyp, eref, low, high, elsz int) int {
	index := st.AtabIndex

	// Check for potential overflow
	if high < low {
		low, high = high, low
	}
	size := elsz * (high - low + 1)

	entry := AtabEntry{
		Xtyp: xtyp,
		Etyp: etyp,
		Eref: eref,
		Low:  low,
		High: high,
		Elsz: elsz,
		Size: size,
	}

	st.Atab = append(st.Atab, entry)
	st.AtabIndex++

	return index
}

// Menambahkan entry ke block table (sudah di-handle di enterBlock)
func (st *SymbolTable) EnterBlock() int {
	return st.enterBlock()
}

// update vsze
func (st *SymbolTable) AddVariableSize(amount int) {
	if st.CurrentBlock >= 0 && st.CurrentBlock < len(st.Btab) {
		st.Btab[st.CurrentBlock].Vsze += amount
	}
}

// update psze
func (st *SymbolTable) AddParameterSize(amount int) {
	if st.CurrentBlock >= 0 && st.CurrentBlock < len(st.Btab) {
		st.Btab[st.CurrentBlock].Psze += amount
	}
}

// update lpar
func (st *SymbolTable) UpdateBlockLastParam(blockIndex, lpar int) {
	if blockIndex >= 0 && blockIndex < len(st.Btab) {
		st.Btab[blockIndex].Lpar = lpar
	}
}

// Cari identifier di symbol table (search dari scope saat ini ke global)
func (st *SymbolTable) Lookup(identifier string) (int, bool) {
	// Search dari level saat ini ke level 0 (global)
	for level := st.CurrentLevel; level >= 0; level-- {
		blockIndex := st.Display[level]

		// Search linked list di block ini
		if blockIndex >= 0 && blockIndex < len(st.Btab) {
			idx := st.Btab[blockIndex].Last
			for idx >= 0 && idx < len(st.Tab) {
				if st.Tab[idx].Identifier == identifier {
					return idx, true
				}
				idx = st.Tab[idx].Link
			}
		}
	}

	return -1, false
}

// Cari identifier hanya di scope saat ini
func (st *SymbolTable) LookupInCurrentScope(identifier string) (int, bool) {
	if st.CurrentBlock < 0 || st.CurrentBlock >= len(st.Btab) {
		return -1, false
	}

	idx := st.Btab[st.CurrentBlock].Last
	for idx >= 0 && idx < len(st.Tab) {
		if st.Tab[idx].Identifier == identifier {
			return idx, true
		}
		idx = st.Tab[idx].Link
	}

	return -1, false
}

// Ambil entry dari symbol table
func (st *SymbolTable) GetEntry(index int) (*TabEntry, error) {
	if index < 0 || index >= len(st.Tab) {
		return nil, fmt.Errorf("invalid tab index: %d", index)
	}
	return &st.Tab[index], nil
}

// Ambil entry dari array table
func (st *SymbolTable) GetArrayEntry(index int) (*AtabEntry, error) {
	if index < 0 || index >= len(st.Atab) {
		return nil, fmt.Errorf("invalid atab index: %d", index)
	}
	return &st.Atab[index], nil
}

// Ambil entry dari block table
func (st *SymbolTable) GetBlockEntry(index int) (*BtabEntry, error) {
	if index < 0 || index >= len(st.Btab) {
		return nil, fmt.Errorf("invalid btab index: %d", index)
	}
	return &st.Btab[index], nil
}

// Dapatkan ukuran tipe dalam byte/unit memori
func (st *SymbolTable) getTypeSize(typ TypeKind, ref int) int {
	switch typ {
	case TypeInteger, TypeBoolean, TypeChar:
		return 1
	case TypeReal:
		return 8
	case TypeArray:
		if ref >= 0 && ref < len(st.Atab) {
			return st.Atab[ref].Size
		}
		return 0
	case TypeRecord:
		if ref >= 0 && ref < len(st.Btab) {
			return st.Btab[ref].Vsze
		}
		return 0
	default:
		return 0
	}
}

// Cek apakah identifier sudah dideklarasikan
func (st *SymbolTable) IsDeclared(identifier string) bool {
	_, found := st.Lookup(identifier)
	return found
}

// Cek apakah identifier sudah dideklarasikan di scope saat ini
func (st *SymbolTable) IsDeclaredInCurrentScope(identifier string) bool {
	_, found := st.LookupInCurrentScope(identifier)
	return found
}

// Print seluruh symbol table untuk debugging
func (st *SymbolTable) PrintSymbolTable() {
	fmt.Println("\n========== SYMBOL TABLE (TAB) ==========")
	fmt.Printf("%-5s %-20s %-12s %-6s %-6s %-6s %-4s %-6s %-6s\n",
		"idx", "id", "obj", "type", "ref", "nrm", "lev", "adr", "link")
	fmt.Println("--------------------------------------------------------------------------------")

	for i := st.ReservedWordsCount; i < len(st.Tab); i++ {
		entry := st.Tab[i]
		fmt.Printf("%-5d %-20s %-12s %-6d %-6d %-6d %-4d %-6d %-6d\n",
			i, entry.Identifier, entry.Obj, entry.Type,
			entry.Ref, entry.Nrm, entry.Lev, entry.Adr, entry.Link)
	}

	fmt.Println("\n========== BLOCK TABLE (BTAB) ==========")
	fmt.Printf("%-5s %-6s %-6s %-6s %-6s\n", "idx", "last", "lpar", "psze", "vsze")
	fmt.Println("----------------------------------------")

	for i := 0; i < len(st.Btab); i++ {
		entry := st.Btab[i]
		fmt.Printf("%-5d %-6d %-6d %-6d %-6d\n",
			i, entry.Last, entry.Lpar, entry.Psze, entry.Vsze)
	}

	fmt.Println("\n========== ARRAY TABLE (ATAB) ==========")
	fmt.Printf("%-5s %-6s %-6s %-6s %-6s %-6s %-6s %-6s\n",
		"idx", "xtyp", "etyp", "eref", "low", "high", "elsz", "size")
	fmt.Println("----------------------------------------------------------------")

	for i := 0; i < len(st.Atab); i++ {
		entry := st.Atab[i]
		fmt.Printf("%-5d %-6d %-6d %-6d %-6d %-6d %-6d %-6d\n",
			i, entry.Xtyp, entry.Etyp, entry.Eref, entry.Low, entry.High,
			entry.Elsz, entry.Size)
	}

	fmt.Println()
}
