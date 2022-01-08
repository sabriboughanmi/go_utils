package enumFlags

type Bits8 uint8

//Set : set the Flag to 1.
func (b *Bits8) Set(flag Bits8) { *b = *b | flag }

//SetMulti : set the Flags to 1.
func (b *Bits8) SetMulti(flags ...Bits8) {
	for _, flag := range flags {
		b.Set(flag)
	}
}

//Clear : set the Flag to 0.
func (b *Bits8) Clear(flag Bits8) { *b = *b &^ flag }

//ClearMulti : set the Flags to 0.
func (b *Bits8) ClearMulti(flags ...Bits8) {
	for _, flag := range flags {
		b.Clear(flag)
	}
}

//Toggle : inverts the Flag.
func (b *Bits8) Toggle(flag Bits8) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits8) Has(flag Bits8) bool { return (*b)&flag != 0 }
