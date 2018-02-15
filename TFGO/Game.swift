//
//  Game.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import MapKit

enum Gamemode: String{
    case cp = "Control Point"
    case payload = "Payload"
    case multi = "Multipoint"
}

class Player {
    
    private var name: String
    private var icon: String
    private var loc = CLLocation(latitude: 0.0, longitude: 0.0)
    private var orientation: Float
    private var weapon: String
    
    func getName() -> String {
        return name
    }
    
    func getWeapon() -> String {
        return weapon
    }
    
    func setName(to name: String) {
        self.name = name
    }
    
    func getLocation() -> CLLocation {
        return loc
    }
    
    func getIcon() -> String {
        return icon
    }
    
    func setIcon(to icon: String) {
        self.icon = icon
    }
    
    func getOrientation() -> Float {
        return orientation
    }
    
    func isValid() -> Bool {
        return name != "" && icon.count == 1
    }
    
    
    init(name: String, icon: String) {
        self.name = name
        self.icon = icon
        self.orientation = 0
        self.weapon = "" // later
    }
}

public class Game {
    
    private var ID: String?
    private var name: String?
    private var mode: Gamemode?
    private var maxTime: Int?
    private var maxPoints: Int?
    private var maxPlayers: Int?
    private var description: String
    private var password: String?
    
    private var players: [Player]
    private var boundaries: [MKMapPoint]
    
    func getID() -> String? {
        return ID
    }
    
    func setID(to id: String) {
        ID = id
    }
    
    func getName() -> String? {
        return name
    }
    
    func setName(to name: String) -> Bool {
        if validName(name) {
            self.name = name
            return true
        }
        return false
    }
    
    func getMode() -> Gamemode? {
        return mode
    }
    
    func setMode(to mode: Gamemode) {
        self.mode = mode
    }
    
    func getTimeLimit() -> Int? {
        return maxTime
    }
    
    func setTimeLimit(to time: Int) -> Bool {
        if validNumber(of: time) {
            self.maxTime = time
            return true
        }
        return false
    }
    
    func getMaxPoints() -> Int? {
        return maxPoints
    }
    
    func setMaxPoints(to points: Int) -> Bool {
        if validNumber(of: points) {
            self.maxPoints = points
            return true
        }
        return false
    }
    
    func getMaxPlayers() -> Int? {
        return maxPlayers
    }
    
    func setMaxPlayers(to players: Int) -> Bool {
        if validNumber(of: players) {
            self.maxPlayers = players
            return true
        }
        return false
    }
    
    func getDescription() -> String {
        return description
    }
    
    func setDescription(to description: String) -> Bool {
        if validDescription(description) {
            self.description = description
            return true
        }
        return false
    }
    
    func getPassword() -> String? {
        return password
    }
    
    func setPassword(to password: String) -> Bool {
        if validDescription(password) {
            self.password = password
            return true
        }
        return false
    }
    
    func getPlayers() -> [Player] {
        return players
    }
    
    func addPlayer(toGame player: Player) {
        players.append(player)
    }
    
    func getBoundaries() -> [MKMapPoint] {
        return boundaries
    }
    
    func setBoundaries(_ points: [MKMapPoint]) {
        boundaries = points
    }
    
    private func validName(_ name: String?) -> Bool {
        if let name = name, name != "", name.count < 30 {
            return true
        }
        return false
    }
    
    private func validNumber(of number: Int?) -> Bool {
        if let number = number, number >= 0 {
            return true
        }
        return false
    }
    
    private func validDescription(_ description: String) -> Bool {
        return description.count < 100
    }
    
    func isValid() -> Bool {
        if validName(name), validNumber(of: maxTime), validNumber(of: maxPoints), validNumber(of: maxPlayers), validDescription(description) {
            return true
        }
        return false
    }
    
    init() {
        description = ""
        players = []
        boundaries = []
    }
}
