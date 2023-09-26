package fuota

import (
	"errors"
)

func GenerateXORFragments(data []byte, fragmentSize int) ([][]byte, error) {
	if len(data)%fragmentSize != 0 {
		return nil, errors.New("length of data must be a multiple of the given fragment-size")
	}

	// fragment the data into rows
	var dataRows [][]byte
	for i := 0; i < len(data)/fragmentSize; i++ {
		offset := i * fragmentSize
		dataRows = append(dataRows, data[offset:offset+fragmentSize])
	}

	redundancyPackets := len(dataRows) / 2

	if(len(dataRows) %2 != 0){
        redundancyPackets++
    }
	

	for y := 0; y < redundancyPackets; y++ {
		s := make([]byte, fragmentSize)

		for m := 0; m < fragmentSize; m++ {
		    if (y*2)+1 == len(data)/fragmentSize {
		        s[m] = dataRows[(y*2)][m]
		    } else {
		        s[m] = dataRows[y*2][m] ^ dataRows[(y*2)+1][m]
		    }
			
		}

		dataRows = append(dataRows, s)
	}

	return dataRows, nil
}
