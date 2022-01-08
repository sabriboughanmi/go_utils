package enumFlags

type Bits64 uint64

//Set : set the Flag to 1.
func (b *Bits64) Set(flag Bits64) { *b = *b | flag }

//SetMulti : set the Flags to 1.
func (b *Bits64) SetMulti(flags ...Bits64) {
	for _, flag := range flags {
		b.Set(flag)
	}
}

//Clear : set the Flag to 0.
func (b *Bits64) Clear(flag Bits64) { *b = *b &^ flag }

//ClearMulti : set the Flags to 0.
func (b *Bits64) ClearMulti(flags ...Bits64) {
	for _, flag := range flags {
		b.Clear(flag)
	}
}

//Toggle : inverts the Flag.
func (b *Bits64) Toggle(flag Bits64) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits64) Has(flag Bits64) bool { return (*b)&flag != 0 }
