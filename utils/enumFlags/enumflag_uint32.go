package enumFlags

type Bits32 uint32

//Set : set the Flag to 1.
func (b *Bits32) Set(flag Bits32) { *b = *b | flag }

//SetMulti : set the Flags to 1.
func (b *Bits32) SetMulti(flags ...Bits32) {
	for _, flag := range flags {
		b.Set(flag)
	}
}

//Clear : set the Flag to 0.
func (b *Bits32) Clear(flag Bits32) { *b = *b &^ flag }

//ClearMulti : set the Flags to 0.
func (b *Bits32) ClearMulti(flags ...Bits32) {
	for _, flag := range flags {
		b.Clear(flag)
	}
}

//Toggle : inverts the Flag.
func (b *Bits32) Toggle(flag Bits32) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits32) Has(flag Bits32) bool { return (*b)&flag != 0 }
