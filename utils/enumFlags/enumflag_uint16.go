package enumFlags

type Bits16 uint16

//Set : set the Flag to 1.
func (b *Bits16) Set(flag Bits16) { *b = *b | flag }

//Clear : set the Flag to 0.
func (b *Bits16) Clear(flag Bits16) { *b = *b &^ flag }

//Toggle : inverts the Flag.
func (b *Bits16) Toggle(flag Bits16) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits16) Has(flag Bits16) bool { return (*b)&flag != 0 }
