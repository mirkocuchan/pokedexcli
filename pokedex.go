package main

type Pokedex struct{
	Pokemons map[string]Pokemon
} 

type Pokemon struct{
	Name string `json:"name"`
	URL  string `json:"url"`    
	Height int `json:"height"`
	Weight int	`json:"weight"`
	Stats []struct {
    	BaseStat int `json:"base_stat"`
    		Stat     struct {
        	Name string `json:"name"`
    		} `json:"stat"`
	}`json:"stats"`
	Types []struct {
        Type struct {
            Name string `json:"name"`
        } `json:"type"`
    } `json:"types"`
	BaseExperience int `json:"base_experience"`
}

func (p *Pokedex) Add(mon Pokemon) error{
	p.Pokemons[mon.Name] = mon
	return nil
}

func (p *Pokedex) Get(name string) (Pokemon, bool){
	pokemon, ok := p.Pokemons[name]
	return pokemon, ok
}

func NewPokedex() *Pokedex{
	p := Pokedex{
    	Pokemons: make(map[string]Pokemon),
	}

	return &p
}