package enumFlags

type Bits16 uint16

//Set : set the Flag to 1.
func (b *Bits16) Set(flag Bits16) { *b = *b | flag }

//SetMulti : set the Flags to 1.
func (b *Bits16) SetMulti(flags ...Bits16) {
	for _, flag := range flags {
		b.Set(flag)
	}
}

//Clear : set the Flag to 0.
func (b *Bits16) Clear(flag Bits16) { *b = *b &^ flag }

//ClearMulti : set the Flags to 0.
func (b *Bits16) ClearMulti(flags ...Bits16) {
	for _, flag := range flags {
		b.Clear(flag)
	}
}

//Toggle : inverts the Flag.
func (b *Bits16) Toggle(flag Bits16) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits16) Has(flag Bits16) bool { return (*b)&flag != 0 }
