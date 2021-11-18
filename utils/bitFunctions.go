package utils

const Uint16One uint16 = 1
const Uint32One uint32 = 1
const Uint64One uint64 = 1
const IntOne int = 1
const ByteOne byte = 1

///////////////////////////// SetBit \\\\\\\\\\\\\\\\\\\\\\\\\\\\

//ByteSetBit Sets the bit at pos in the Byte.
func ByteSetBit(n byte, index uint) byte {
	n |= (1 << index)
	return n
}

//Uint16SetBit Sets the bit at pos in the Uint16.
func Uint16SetBit(n uint16, index uint) uint16 {
	n |= (Uint16One << index)
	return n
}

//Uint32SetBit Sets the bit at pos in the Uint32.
func Uint32SetBit(n uint32, index uint) uint32 {
	n |= (Uint32One << index)
	return n
}


//IntSetBit Sets the bit at pos in the int.
func IntSetBit(n int, index uint) int {
	n |= (IntOne << index)
	return n
}

//Uint64SetBit Sets the bit at pos in the Uint32.
func Uint64SetBit(n uint64, index uint) uint64 {
	n |= (Uint64One << index)
	return n
}

//////////////////////////////// ClearBit \\\\\\\\\\\\\\\\\\\\\\\\\\\\

//ByteClearBit Clears the bit at pos for a byte.
func ByteClearBit(n byte, index uint) byte {
	mask := ^(ByteOne << index)
	n &= mask
	return n
}

//Uint16ClearBit Clears the bit at pos for a Uint16.
func Uint16ClearBit(n uint16, index uint) uint16 {
	mask := ^(Uint16One << index)
	n &= mask
	return n
}

//Uint32ClearBit Clears the bit at pos for a Uint32.
func Uint32ClearBit(n uint32, index uint) uint32 {
	mask := ^(Uint32One << index)
	n &= mask
	return n
}

//IntClearBit Clears the bit at pos for a Uint32.
func IntClearBit(n int, index uint) int {
	mask := ^(IntOne << index)
	n &= mask
	return n
}

//Uint64ClearBit Clears the bit at pos for a Uint32.
func Uint64ClearBit(n uint64, index uint) uint64 {
	mask := ^(Uint64One << index)
	n &= mask
	return n
}

//////////////////////////////// IsBitSet \\\\\\\\\\\\\\\\\\\\\\\\\\\\

//ByteIsBitSet return true if the bit at index is 1
func ByteIsBitSet(n byte, index uint) bool {
	val := n & (ByteOne << index)
	return val > 0
}

//Uint16IsBitSet return true if the bit at index is 1
func Uint16IsBitSet(n uint16, index uint) bool {
	val := n & (Uint16One << index)
	return val > 0
}

//Uint32IsBitSet return true if the bit at index is 1
func Uint32IsBitSet(n uint32, index uint) bool {
	val := n & (Uint32One << index)
	return val > 0
}

//IntIsBitSet return true if the bit at index is 1
func IntIsBitSet(n int, index uint) bool {
	val := n & (IntOne << index)
	return val > 0
}

//Uint64IsBitSet return true if the bit at index is 1
func Uint64IsBitSet(n uint64, index uint) bool {
	val := n & (Uint64One << index)
	return val > 0
}
