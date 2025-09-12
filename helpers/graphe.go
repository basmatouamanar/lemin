package helpers

func Graphe(Rooms []*Room) map[string][]string {
	graphe := make(map[string][]string)
	for _, Room := range Rooms {
		for _, l := range Room.linker {
			graphe[Room.Name] = append(graphe[Room.Name], l.Name)
		}
	}
	return graphe
}