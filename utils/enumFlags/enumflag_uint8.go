package enumFlags

type Bits8 uint8

//Set : set the Flag to 1.
func (b *Bits8) Set(flag Bits8) { *b = *b | flag }

//Clear : set the Flag to 0.
func (b *Bits8) Clear(flag Bits8) { *b = *b &^ flag }

//Toggle : inverts the Flag.
func (b *Bits8) Toggle(flag Bits8) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits8) Has(flag Bits8) bool { return (*b)&flag != 0 }
