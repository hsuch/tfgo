//
//  Game.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import MapKit

enum Gamemode: String {
    case cp = "SingleCapture"
    case payload = "Payload"
    case multi = "MultiCapture"
}

class Player {
    
    private var name: String
    private var icon: String
    private var team: String
    private var loc = CLLocation(latitude: 0.0, longitude: 0.0)
    private var orientation: Float
    private var weapon: String
    private var status: String
    private var health: Int
    private var armor: Int
    private var host: Bool = false
    
    func getName() -> String {
        return name
    }
    
    func getWeapon() -> String {
        return weapon
    }
    
    func isHost() -> Bool {
        return host
    }
    
    func makeHost() {
        host = true
    }
    
    func setName(to name: String) {
        self.name = name
    }
    
    func getLocation() -> CLLocation {
        return loc
    }
    
    func setLocation(to x: Double, to y: Double) {
        self.loc = CLLocation(latitude: x, longitude: y)
    }
    
    func getIcon() -> String {
        return icon
    }
    
    func setIcon(to icon: String) {
        self.icon = icon
    }
    
    func getStatus() -> String {
        return status
    }
    
    func setStatus(to status: String) {
        self.status = status
    }
    
    func getHealth() -> Int {
        return health
    }
    
    func setHealth(to health: Int) {
        self.health = health
    }
    
    func getArmor() -> Int {
        return armor
    }
    
    func setArmor(to armor: Int) {
        self.armor = armor
    }
    
    func getOrientation() -> Float {
        return orientation
    }
    
    func setOrientation(to orientation: Float) {
        self.orientation = orientation
    }
    
    func getTeam() -> String {
        return team
    }
    
    func setTeam(to team: String) {
        self.team = team
    }
    
    func isValid() -> Bool {
        return name != "" && icon.count == 1
    }
    
    
    init(name: String, icon: String) {
        self.name = name
        self.icon = icon
        self.orientation = 0
        self.weapon = "Sword" // later
        self.status = "" // later
        self.health = 100 // later
        self.armor = 0 // later
        self.team = ""
    }
}

public class Objective {
    
    private var xLoc: Double
    private var yLoc: Double
    private var radius: Double
    private var occupants: [String]
    private var owner: String
    private var progress: Int
    
    func getXLoc() -> Double? {
        return xLoc
    }
    
    func getYLoc() -> Double? {
        return yLoc
    }
    
    func getRadius() -> Double? {
        return radius
    }
    
    func setOccupants(to occupants: [String]) {
        self.occupants = occupants
    }
    
    func getOwner() -> String {
        return owner
    }
    
    func setOwner(to owner: String) {
        self.owner = owner
    }
    
    func getProgress() -> Int {
        return progress
    }
    
    func setProgress(to progress: Int) {
        self.progress = progress
    }
    
    
    init(x: Double, y: Double, radius: Double) {
        self.radius = radius
        self.occupants = []
        self.owner = "Neutral"
        self.progress = 0
        self.xLoc = x
        self.yLoc = y
    }
    
}

public class Game {
    
    private var ID: String?
    private var name: String?
    private var mode: Gamemode = .cp
    private var maxTime: Int = 2
    private var maxPoints: Int = 10
    private var maxPlayers: Int = 2
    private var redPoints: Int = 0
    private var bluePoints: Int = 0
    private var description: String
    private var password: String?
    private var startTime: String = ""
    
    private var objectives: [Objective]
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
    
    func getMode() -> Gamemode {
        return mode
    }
    
    func setMode(to mode: Gamemode) {
        self.mode = mode
    }
    
    func getTimeLimit() -> Int {
        return maxTime
    }
    
    func setTimeLimit(to time: Int) {
        self.maxTime = time
    }
    
    func getMaxPoints() -> Int {
        return maxPoints
    }
    
    func setMaxPoints(to points: Int) {
        self.maxPoints = points
    }
    
    func getRedPoints() -> Int {
        return redPoints
    }
    
    func setRedPoints(to redPoints: Int) {
        self.redPoints = redPoints
    }
    
    func getBluePoints() -> Int {
        return bluePoints
    }
    
    func setBluePoints(to bluePoints: Int) {
        self.bluePoints = bluePoints
    }
    
    func getMaxPlayers() -> Int {
        return maxPlayers
    }
    
    func setMaxPlayers(to players: Int) {
        self.maxPlayers = players
    }
    
    func getDescription() -> String {
        return description
    }
    
    func getStartTime() -> String {
        return startTime
    }
    
    func setStartTime(to startTime: String) {
        self.startTime = startTime
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
    
    func getObjectives() -> [Objective] {
        return objectives
    }
    
    func addPlayer(toGame player: Player) {
        players.append(player)
    }
    
    func removePlayer(index: Int) {
        players.remove(at: index)
    }
    
    func addObjective(toObjective objective: Objective) {
        objectives.append(objective)
    }
    
    func getBoundaries() -> [MKMapPoint] {
        return boundaries
    }
    
    func setBoundaries(_ points: [MKMapPoint]) {
        boundaries = points
    }
    
    func hasPlayer(name: String) -> Bool {
        for player in players {
            if name == player.getName() {
                return true
            }
        }
        return false
    }
    
    func findPlayerIndex(name: String) -> Int {
        var i = 0;
        for player in players {
            if name == player.getName() {
                return i
            }
            i = i + 1
        }
        return -1
    }
    
    func findObjectiveIndex(x: Double, y: Double) -> Int {
        var i = 0
        for objective in objectives {
            if x == objective.getXLoc(), y == objective.getYLoc() {
                return i
            }
            i = i + 1
        }
        return -1
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
        objectives = []
    }
}
