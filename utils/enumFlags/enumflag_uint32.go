package enumFlags

type Bits32 uint32

//Set : set the Flag to 1.
func (b *Bits32) Set(flag Bits32) { *b = *b | flag }

//Clear : set the Flag to 0.
func (b *Bits32) Clear(flag Bits32) { *b = *b &^ flag }

//Toggle : inverts the Flag.
func (b *Bits32) Toggle(flag Bits32) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits32) Has(flag Bits32) bool { return (*b)&flag != 0 }
