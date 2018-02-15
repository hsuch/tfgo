//
//  Network.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import SwiftSocket
import SwiftyJSON
import MapKit

var gameState = GameState()

class Connection {
    private var servadd: String = "128.135.175.185" // to be replaced with real server ip
    private var servport: Int32 = 9265
    private var client: TCPClient

    func sendData(data: Data) -> Result{
        return client.send(data: data)
    }

    func recvData() -> Data? {
        let response = client.read(1024*10)
        return Data.init(response!)
    }

    init() {
        client = TCPClient(address: servadd, port: servport)
        DispatchQueue.global(qos: .background).async {
            switch self.client.connect(timeout: 10) {
            case .failure:
                print("Connection failed")
            default:
                print("Successful connection")
            }
        }
    }
}

class MsgFromServer {
    private var type: String
    /* possible message types:
        PlayerListUpdate, AvailableGames, GameInfo, JoinGameError, GameStartInfo, GameUpdate, StatusUpdate
    */

    private var data: [String: Any]

    func getType() -> String {
        return type
    }

    /* parse(): convert data array into appropriate data struct depending on message type */
    func parse() -> Bool {
        switch type {
        case "PlayerListUpdate":
            parsePlayerListUpdate(data: data)
            return true
        case "AvailableGames":
            return parseAvailableGames(data: data)
        case "GameInfo":
            return parseGameInfo(data: data)
        case "JoinGameError":
            return parseJoinGameError(data: data)
        case "GameStartInfo":
            return parseGameStartInfo(data: data)
        case "GameUpdate":
            return parseGameUpdate(data: data)
        case "StatusUpdate":
            return parseStatusUpdate(data: data)
        case "TakeHit":
            return parseTakeHit(data: data)
        default:
            return false
        }
    }

    init() {
        let conn = gameState.getConnection()
        let received = conn.recvData()
            self.data = try! JSONSerialization.jsonObject(with: received!, options: []) as! [String: Any]
            let type = data.removeValue(forKey: "Type")
            self.type = type as! String
        

    }
}

/* Parsing functions: helper functions called by parse() to parse different messages */

func parsePlayerListUpdate(data: [String: Any]) {
    
    var players = [Player]()
    
    if let info = data["Data"] as? [[String: Any]] {
        for player in info {
            if let name = player["Name"] as? String, let icon = player["Icon"] as? String {
                players.append(Player(name: name, icon: icon))
            }
        }
    }
    
    var current_players = gameState.getCurrentGame().getPlayers()
    var index = 0
    
    for current_player in current_players {
        let current_name = current_player.getName()
        for player in players {
            if current_name == player.getName() {
                gameState.getCurrentGame().removePlayer(index: index)
                index = index - 1
                break
            }
        }
        index = index + 1
    }
    
    current_players = gameState.getCurrentGame().getPlayers()
    
    for player in players {
        if gameState.getCurrentGame().hasPlayer(name: player.getName()) {
            gameState.getCurrentGame().addPlayer(toGame: player)
        }
    }
}

func parseAvailableGames(data: [String: Any]) -> Bool {

    if let info = data["Data"] as? [[String: Any]] {
        for game in info {
            if let id = game["ID"] as? String {
                if !gameState.hasGame(to: id) {
                    //TODO also get the game location and player list
                    if let name = game["Name"] as? String, let mode = game["Mode"] as? Gamemode {
                        let newGame = Game()
                        newGame.setID(to: id)
                        newGame.setName(to: name)   // these games will always give a valid name
                        newGame.setMode(to: mode)
                        gameState.addFoundGame(found: newGame)
                    }
                }
            }
        }
        return true
    }
    return false
}

func parseGameInfo(data: [String: Any]) -> Bool {
    
    if let info = data["Data"] as? [String: Any] {
        if let desc = info["Description"] as? String, let playerNum = info["PlayerLimit"] as? Int, let pointLim = info["PointLimit"] as? Int, let timeLim = info["TimeLimit"] as? String {
            gameState.getCurrentGame().setMaxPlayers(to: playerNum)
            gameState.getCurrentGame().setDescription(to: desc)
            gameState.getCurrentGame().setMaxPoints(to: pointLim)
            
        }
        if let players = info["PlayerList"] as? [[String: Any]] {
            for player in players {
                if let name = player["Name"] as? String, let icon = player["Icon"] as? String {
                    gameState.getCurrentGame().addPlayer(toGame: Player(name: name, icon: icon))
                }
            }
            return true
        }
    }
    return false
}

func parseJoinGameError(data: [String: Any]) -> Bool {
    if let error = data["Data"] as? String {
        let alert = UIAlertController(title: error, message: "Please join a different game", preferredStyle: .alert)
        alert.addAction(UIAlertAction(title: "Return", style: .cancel, handler: nil))
        
        //self.present(alert, animated: true)  TODO figure out how to send message without extension UIViewController
        
        return true
    }
    return false
}

func parseGameStartInfo(data: [String: Any]) -> Bool {
    
    //TODO add bases for iteration 2
    if let info = data["Data"] as? [String: Any] {
        if let players = info["PlayerList"] as? [[String: Any]] {
            for player in players {
                if let name = player["Name"] as? String, let team = player["Team"] as? String {
                    let index = gameState.getCurrentGame().findPlayerIndex(name: name)
                    if index > -1 {
                        gameState.getCurrentGame().getPlayers()[index].setTeam(to: team)
                    }
                }
            }
        }
        if let objectives = info["Objectives"] as? [[String: Any]] {
            for objective in objectives {
                if let loc = objective["Location"] as? [String: Any], let radius = objective["Radius"] as? Double {
                    if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                        gameState.getCurrentGame().addObjective(toObjective: Objective(x: x, y: y, radius: radius))
                    }
                }
            }
        }
        return true
    }
    return false
}

func parseGameUpdate(data: [String: Any]) -> Bool {
    
    if let info = data["Data"] as? [String: Any] {
        if let players = info["PlayerList"] as? [[String: Any]], let points = info["Points"] as? [String: Any], let objectives = info["Objectives"] as? [[String: Any]] {
            
            for player in players {
                if let name = player["Name"] as? String, let orientation = player["Orientation"] as? Float, let loc = player["Location"] as? [String: Any] {
                    
                    let index = gameState.getCurrentGame().findPlayerIndex(name: name)
                    if index > -1 {
                        gameState.getCurrentGame().getPlayers()[index].setOrientation(to: orientation)
                        if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                            gameState.getCurrentGame().getPlayers()[index].setLocation(to: x, to: y)
                        }
                    }
                }
            }
            
            if let redPoints = points["Red"] as? Int, let bluePoints = points["Blue"] as? Int {
                gameState.getCurrentGame().setRedPoints(to: redPoints)
                gameState.getCurrentGame().setBluePoints(to: bluePoints)
            }
            
            for objective in objectives {
                if let loc = objective["Location"] as? [String: Any], let occupants = objective["Occupying"] as? [String], let owner = objective["BelongsTo"] as? String, let progress = objective["Progress"] as? Int {
                    
                    if let x = loc["X"] as? Double, let y = loc["Y"] as? Double {
                        let objIndex = gameState.getCurrentGame().findObjectiveIndex(x: x, y: y)
                        if objIndex > -1 {
                            gameState.getCurrentGame().getObjectives()[objIndex].setOwner(to: owner)
                            gameState.getCurrentGame().getObjectives()[objIndex].setProgress(to: progress)
                            gameState.getCurrentGame().getObjectives()[objIndex].setOccupants(to: occupants)
                        }
                    }
                }
            }
        }
        return true
    }
    return false
}

func parseStatusUpdate(data: [String: Any]) -> Bool {
    if let status = data["Data"] as? String {
        gameState.setUserStatus(to: status)
        return true
    }
    return false
}

func parseTakeHit(data: [String: Any]) -> Bool {
    if let info = data["Data"] as? [String: Any] {
        if let health = info["Health"] as? Int, let armor = info["Armor"] as? Int {
            gameState.setUserHealth(to: health)
            gameState.setUserArmor(to: armor)
            return true
        }
    }
    return false
}

class MsgToServer {
    private var action: String
    /* possible message actions:
        case CreateGame, ShowGames, ShowGameInfo, JoinGame, StartGame, LocationUpdate, Fire
    */

    private var data: [String: Any]

    /* toJson(): convert action type and data array into server-readable json */
    func toJson() -> Data {
        let message = ["Action": action, "Data": data] as [String : Any]
        let retval = try! JSONSerialization.data(withJSONObject: message, options: JSONSerialization.WritingOptions.prettyPrinted)
        return retval
    }

    init(action: String, data: [String: Any]) {
        self.action = action
        self.data = data
    }
}

/* Message generators: the following functions generate messages that can be directly sent to the server via Connection.sendData()*/

func CreateGameMsg(game: Game) -> Data {
    // still needs to take boundaries from game page
    let host = ["Name": gameState.getUserName(), "Icon": gameState.getUserIcon()] as [String: Any]
    let minutes = game.getTimeLimit()
    let timelimit = "0h" + "\(minutes)" + "m0s"
    let payload = ["Name": game.getName()!, "Password": game.getPassword() ?? "", "Description": game.getDescription(), "PlayerLimit": game.getMaxPlayers(), "PointLimit": game.getMaxPoints(), "TimeLimit": timelimit, "Mode": game.getMode().rawValue, "Boundaries": [], "Host": host] as [String: Any]
    return MsgToServer(action: "CreateGame", data: payload).toJson()
}
func ShowGamesMsg() -> Data {
    let payload = ["Name": gameState.getUserName(), "Icon": gameState.getUserIcon()]
    return MsgToServer(action: "ShowGames", data: payload).toJson()
}
func ShowGameInfoMsg(IDtoShow: String) -> Data {
    return MsgToServer(action: "ShowGameInfo", data: ["GameID": IDtoShow]).toJson()
}
func JoinGameMsg(IDtoJoin: String) -> Data {
    return MsgToServer(action: "JoinGame", data: ["GameID": IDtoJoin]).toJson()
}
func StartGameMsg() -> Data {
    return MsgToServer(action: "StartGame", data: [:]).toJson()
}
func LocUpMsg() -> Data {
    let location = ["X": gameState.getUserLocation().coordinate.latitude, "Y": gameState.getUserLocation().coordinate.longitude]
    let payload = ["Location": location, "Orientation": gameState.getUserOrientation()] as [String : Any]
    return MsgToServer(action: "LocationUpdate", data: payload).toJson()
}
func FireMsg() -> Data {
    // todo, take orientation and weapon from this client's player
    let weapon = gameState.getUserWeapon()
    let direction = gameState.getUserOrientation()
    let payload = ["Weapon": weapon, "Direction": direction] as [String: Any]
    return MsgToServer(action: "Fire", data: payload).toJson()

}



class GameState {
    
    private var currentGame: Game?
    private var foundGames: [Game] = []
    private var user: Player = Player(name: "", icon:"")
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
    
    func getUserWeapon() -> String {
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
    
    func findPublicGames() -> [Game]{
        var publicGames: [Game] = []
        for game in foundGames {
            if game.getPassword() == nil {
                publicGames.append(game)
            }
        }
        return publicGames
    }
    
    func findPrivateGames() -> [Game]{
        var privateGames: [Game] = []
        for game in foundGames {
            if game.getPassword() != nil {
                privateGames.append(game)
            }
        }
        return privateGames
    }
    
    init() {
        
    }
}

