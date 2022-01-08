package enumFlags

type Bits64 uint64

//Set : set the Flag to 1.
func (b *Bits64) Set(flag Bits64) { *b = *b | flag }

//Clear : set the Flag to 0.
func (b *Bits64) Clear(flag Bits64) { *b = *b &^ flag }

//Toggle : inverts the Flag.
func (b *Bits64) Toggle(flag Bits64) { *b = *b ^ flag }

//Has : returns true if Flag is 1.
func (b *Bits64) Has(flag Bits64) bool { return (*b)&flag != 0 }
