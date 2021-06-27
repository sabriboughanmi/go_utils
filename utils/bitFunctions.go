package utils


const Uint16One uint16 = 1
const ByteOne byte = 1


//Uint16SetBit Sets the bit at pos in the Uint16.
func Uint16SetBit(n uint16, index uint) uint16 {
	n |= (Uint16One << index)
	return n
}

//ByteSetBit Sets the bit at pos in the Byte.
func ByteSetBit(n byte, index uint) byte {
	n |= (1 << index)
	return n
}

//Uint16ClearBit Clears the bit at pos for a Uint16.
func Uint16ClearBit(n uint16, index uint) uint16 {
	mask := ^(Uint16One << index)
	n &= mask
	return n
}

//ByteClearBit Clears the bit at pos for a byte.
func ByteClearBit(n byte, index uint) byte {
	mask := ^(ByteOne << index)
	n &= mask
	return n
}

//Uint16IsBitSet return true if the bit at index is 1
func Uint16IsBitSet(n uint16, index uint) bool {
	val := n & (Uint16One << index)
	return val > 0
}

//ByteIsBitSet return true if the bit at index is 1
func ByteIsBitSet(n byte, index uint) bool {
	val := n & (ByteOne << index)
	return val > 0
}
