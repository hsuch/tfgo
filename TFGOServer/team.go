package main

type Team struct {
	Name string
	Players map[string]*Player
	Base Location
	Points int
}