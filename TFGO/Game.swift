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

func randomColor() -> UIColor {
    let color = UIColor(red: CGFloat(arc4random_uniform(100))/200.0 + 0.5, green: CGFloat(arc4random_uniform(100))/200.0 + 0.5, blue: CGFloat(arc4random_uniform(100))/200.0 + 0.5, alpha: 1)
    return color
}

class Player {
    
    private var name: String
    private var icon: String
    private var team: String
    private var loc = CLLocation(latitude: 0.0, longitude: 0.0)
    private var orientation: Float
    private var weapon: Weapon
    private var weapons: [Weapon]
    private var pickups: [Pickup] = []
    private var status: String
    private var health: Int
    private var armor: Int
    private var host: Bool = false
    private var annotation: MKPointAnnotation = MKPointAnnotation()
    
    func getName() -> String {
        return name
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
    
    func addWeapon(to weapon: String) {
        self.weapons.append(weaponByName(name: name))
    }
    
    func getWeapon() -> Weapon {
        return weapon
    }
    
    func setWeapon(to weapon: Weapon) {
        self.weapon = weapon
    }
    
    func getWeaponsList() -> [Weapon] {
        return weapons
    }
    
    func updateAnnotation() {
        self.annotation.coordinate = loc.coordinate
        self.annotation.title = name
        self.annotation.subtitle = team
    }
    
    func getAnnotation() -> MKPointAnnotation {
        return annotation
    }
    
    func isValid() -> Bool {
        return name != "" && icon.count == 1
    }
    
    init(name: String, icon: String) {
        self.name = name
        self.icon = icon
        self.orientation = 0
        self.weapon = BeeSwarm() // later
        self.status = "" // later
        self.health = 100 // later
        self.armor = 0 // later
        self.team = ""
        self.weapons = [BeeSwarm(), Sword(), Shotgun()]
    }
}

public class Objective {
    
    private var xLoc: Double
    private var yLoc: Double
    private var radius: Double
    private var occupants: [String]
    private var owner: String
    private var progress: Int
    private var id: String
    
    func getXLoc() -> Double {
        return xLoc
    }
    
    func getYLoc() -> Double {
        return yLoc
    }
    
    func getRadius() -> Double {
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
    
    func getID() -> String {
        return id
    }
    
    func setID(to id: String) {
        self.id = id
    }
    
    init(x: Double, y: Double, radius: Double, id: String) {
        self.radius = radius
        self.occupants = []
        self.owner = "Neutral"
        self.progress = 0
        self.xLoc = x
        self.yLoc = y
        self.id = id
    }
    
}

public class Pickup {
    
    private var loc : MKMapPoint
    private var type : String
    private var amount : Int
    private var available : Bool = true
    private var annotation: MKPointAnnotation = MKPointAnnotation()
    
    init(loc: MKMapPoint, type: String, amount: Int) {
        self.loc = loc
        self.type = type
        self.amount = amount
        
        self.annotation.coordinate = CLLocationCoordinate2D(latitude: loc.x, longitude: loc.y)
        self.annotation.title = type
        self.annotation.subtitle = "Available"
    }
    
    func getX() -> Double {
        return loc.x
    }
    
    func getY() -> Double {
        return loc.y
    }
    
    func getAvailability() -> Bool {
        return available
    }
    
    func updateAnnotation() {
        if (available) {
            self.annotation.subtitle = "Available"
        }
        else {
            self.annotation.subtitle = "Unavailable"
        }
    }
    
    func getAnnotation() -> MKPointAnnotation {
        return annotation
    }
    
    func setAvailability(to availability: Bool) {
        self.available = availability
    }
    
}

public class Game {
    
    private var ID: String?
    private var name: String?
    private var mode: Gamemode = .cp
    private var maxTime: Int = 2
    private var maxPoints: Int = 10
    private var maxPlayers: Int = 2
    private var maxObjectives: Int = 2
    private var redPoints: Int = 0
    private var bluePoints: Int = 0
    private var description: String
    private var password: String?
    private var hasPassword: Bool = false
    private var startTime: [String] = []
    
    private var objectives: [Objective]
    private var players: [Player]
    private var pickups: [Pickup] = []
    private var boundaries: [MKMapPoint]
    private var location: MKMapPoint
    private var blueBaseLoc : MKMapPoint
    private var blueBaseRad : Double = 0.0
    private var redBaseLoc : MKMapPoint
    private var redBaseRad : Double = 0.0
    
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
    
    func getMaxObjectives() -> Int {
        return maxObjectives
    }
    
    func setMaxObjectives(to objectives: Int) {
        self.maxObjectives = objectives
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
    
    func getStartTime() -> [String] {
        return startTime
    }
    
    func stringSplit(to time: String) -> [String] {
        let result = time.components(separatedBy: ["-", " ", ":"])
        return result
    }
    
    func setStartTime(to startTime:String) {
        let time: [String] = stringSplit(to: startTime)
        self.startTime = time
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
            self.hasPassword = true
            return true
        }
        return false
    }
    
    func isPrivate() -> Bool {
        return hasPassword
    }
    
    func setHasPassword(to hasPassword: Bool) {
        self.hasPassword = hasPassword
    }
    
    func getPlayers() -> [Player] {
        return players
    }
    
    func getObjectives() -> [Objective] {
        return objectives
    }
    
    func setPlayers(toGame players: [Player]) {
        self.players = players
    }
    
    func getPickups() -> [Pickup] {
        return pickups
    }
    
    func setPickups(toPickup pickups: [Pickup]) {
        self.pickups = pickups
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
    
    func addBounary(to point: MKMapPoint) {
        boundaries.append(point)
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
    
    func findObjectiveIndex(id: String) -> Int {
        var i = 0
        for objective in objectives {
            if id == objective.getID() {
                return i
            }
            i = i + 1
        }
        return -1
    }
    
    func findPickupIndex(x: Double, y: Double) -> Int {
        var i = 0
        for pickup in pickups {
            if x == pickup.getX(), y == pickup.getY() {
                return i
            }
            i = i + 1
        }
        return -1
    }
    
    func setLocation(to loc: MKMapPoint) {
        location = loc
    }
    
    func getLocation() -> MKMapPoint {
        return location
    }
    
    func getBlueBaseLoc() -> MKMapPoint {
        return blueBaseLoc
    }
    
    func setBlueBaseLoc(to loc: MKMapPoint) {
        self.blueBaseLoc = loc
    }
    
    func getBlueBaseRad() -> Double {
        return blueBaseRad
    }
    
    func setBlueBaseRad(to radius: Double) {
        self.blueBaseRad = radius
    }
    
    func getRedBaseLoc() -> MKMapPoint {
        return redBaseLoc
    }
    
    func setRedBaseLoc(to loc: MKMapPoint) {
        self.redBaseLoc = loc
    }
    
    func getRedBaseRad() -> Double {
        return redBaseRad
    }
    
    func setRedBaseRad(to radius: Double) {
        self.redBaseRad = radius
    }
    
    func updatePlayerAnnotations() {
        for player in players {
            player.updateAnnotation()
        }
    }
    
    func updatePickupAnnotations() {
        for pickup in pickups {
            pickup.updateAnnotation()
        }
    }
    
    func getPlayerAnnotations() -> [MKPointAnnotation] {
        var annotations: [MKPointAnnotation] = []
        for player in players {
            annotations.append(player.getAnnotation())
        }
        return annotations
    }
    
    func getPickupAnnotations() -> [MKPointAnnotation] {
        var annotations: [MKPointAnnotation] = []
        for pickup in pickups {
            annotations.append(pickup.getAnnotation())
        }
        return annotations
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
    
    func validBoundaries() -> Bool {
        return boundaries.count == 4
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
        location = MKMapPoint()
        blueBaseLoc = MKMapPoint()
        redBaseLoc = MKMapPoint()
    }
}
