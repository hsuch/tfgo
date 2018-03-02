//
//  State.swift
//  TFGO
//
//  Created by Sam Schlang on 2/20/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import MapKit

var gameState = GameState()

class GameState {
    
    private var currentGame: Game?
    private var foundGames: [Game] = []
    private var user: Player = Player(name: "", icon:"", id: "")
    private var connection = Connection()
    
    
    func getUserName() -> String {
        return user.getName()
    }
    
    func setUserName(to name: String) {
        user.setName(to: name)
    }
    
    func getUserIcon() -> String {
        return user.getIcon()
    }
    
    func setUserIcon(to icon: String) {
        user.setIcon(to: icon)
    }
    
    func getUserStatus() -> String {
        return user.getStatus()
    }
    
    func setUserStatus(to status: String) {
        user.setStatus(to: status)
    }
    
    func getUserHealth() -> Int {
        return user.getHealth()
    }
    
    func setUserHealth(to health: Int) {
        user.setHealth(to: health)
    }
    
    func getUserArmor() -> Int {
        return user.getArmor()
    }
    
    func setUserArmor(to armor: Int) {
        user.setArmor(to: armor)
    }
    
    func getUserWeapon() -> Weapon {
        return user.getWeapon()
    }
    
    func getUser() -> Player {
        return user
    }
    
    func getUserLocation() -> CLLocation {
        return user.getLocation()
    }
    
    func getUserOrientation() -> Float {
        return user.getOrientation()
    }
    
    /* Do not call unless a game exists!!! */
    func getCurrentGame() -> Game {
        return currentGame!
    }
    
    func addFoundGame(found: Game) {
        foundGames.append(found)
    }
    
    func hasGame(to id: String) -> Bool {
        for game in foundGames {
            if id == game.getID() {
                return true
            }
        }
        return false
    }
    
    func findGame(to id: String) -> Game! {
        for game in foundGames {
            if id == game.getID() {
                return game
            }
        }
        return nil
    }
    
    func setFoundGames(to games: [Game]) {
        self.foundGames = games
    }
    
    func setCurrentGame(to game: Game) -> Bool {
        let valid = game.isValid()
        if valid {
            currentGame = game
        }
        return valid
    }
    
    func getConnection() -> Connection {
        return connection
    }
    
    func findPublicGames() -> [Game] {
        var publicGames: [Game] = []
        for game in foundGames {
            if game.getPassword() == nil {
                publicGames.append(game)
            }
        }
        return publicGames
    }
    
    func findPrivateGames() -> [Game] {
        var privateGames: [Game] = []
        for game in foundGames {
            if game.getPassword() != nil {
                privateGames.append(game)
            }
        }
        return privateGames
    }
    
    func getFoundGames() -> [Game] {
        return foundGames
    }
    
    func getDistanceFromGame(game: Game) -> Float {
        return 0.0
    }
    
    init() {
        
    }
}
