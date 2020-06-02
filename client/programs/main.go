/*
* Project: FLOTEA - Decentralized passenger transport system
* Copyright (c) 2020 Flotea, All Rights Reserved
* For conditions of distribution and use, see copyright notice in LICENSE
*/

package programs

var CurrentWorking = make(map[string]bool)

func IsCurrentlyWorking() bool {

	for _, isWorking := range CurrentWorking {
		if isWorking {
			return true
		}
	}
	return false
}
