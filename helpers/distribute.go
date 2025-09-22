package helpers


func Distribute(paths [][]string, ants int) [][]string {
	var l int
	var lenpath int
	var tour []int
	var tours []int
	var newslice [][]string
	exist := false
	for i := 0; i < len(paths); i++ {
		if newslice == nil {
			lenpath = len(paths[i]) - 1
			l = lenpath + ants - 1
			tours = append(tours, l)
			newslice = append(newslice, paths[i])
		} else  {
			for j := 0; j < len(newslice); j++ {
				n := ants/(i + 1)
				lenpath = len(paths[i]) - 1
				l = lenpath + n - 1
				tours = append(tours, l)
				tour = append(tour, l)
			}
			newslice = append(newslice, paths[i])
			max := max(tour)
			for _, c := range tours {
				if c == max {
					exist = true
					break
				}
			}
			tour = nil
			if exist {
				break
			}

		}
	}
	return newslice
}

func max(s []int) int {
    max := s[0]
    for _, v := range s {
        if v > max {
            max = v
        }
    }
    return max
}



/*
tours = longueur(Path1) + nb_fourmis -1
tours = 4 + 6 -1 = 9 tours
*/

/*
max(5,5)
*/


/*
tours = longueur(Path1) + nb_fourmis -1
tours = 4 + 6 -1 = 9 tours
*/

/*
max(5,5)
*/
