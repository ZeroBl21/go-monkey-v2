package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
	FreeScope    SymbolScope = "FREE"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:       map[string]Symbol{},
		FreeSymbols: []Symbol{},
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer

	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol

	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{
		Name:  original.Name,
		Index: len(s.FreeSymbols) - 1,
		Scope: FreeScope,
	}

	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]

	if !ok && s.Outer != nil {
		symbol, ok = s.Outer.Resolve(name)
		if !ok {
			return symbol, ok
		}

		if symbol.Scope == GlobalScope || symbol.Scope == BuiltinScope {
			return symbol, ok
		}

		free := s.defineFree(symbol)
		return free, true
	}

	return symbol, ok
}
