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
    private var servadd: String = "www.google.com" // to be replaced with real server ip
    private var servport: Int32 = 80
    private var client: TCPClient

    func sendData(data: Data) -> Result{
        return client.send(data: data)
    }

    func recvData() -> Data? {
        guard let response = client.read(1024*10)
        else { return nil }
        return Data.init(response)
    }

    init(){
        client = TCPClient(address: servadd, port: servport)
        client.connect(timeout: 10) // this should probably have success and failure case but whatever
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
//        case "AvailableGames":
//            return parseAvailableGames(data: data)
//        case "GameInfo":
//            return parseGameInfo(data: data)
        case "JoinGameError":
            return parseJoinGameError(data: data)
//        case "GameStartInfo":
//            return parseGameStartInfo(data: data)
//        case "GameUpdate":
//            return parseGameUpdate(data: data)
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
    
    for c_player in current_players {
        let c_name = c_player.getName()
        for player in players {
            if c_name == player.getName() {
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

//func parseAvailableGames(data: [String: Any]) -> Bool {
//
//}
//
//func parseGameInfo(data: [String: Any]) -> Bool {
//    let json = JSON(data)
//}
//
func parseJoinGameError(data: [String: Any]) -> Bool {
    if let error = data["Data"] as? String {
        let alert = UIAlertController(title: error, message: "Please join a different game", preferredStyle: .alert)
        alert.addAction(UIAlertAction(title: "Return", style: .cancel, handler: nil))
        
        //self.present(alert, animated: true)  TODO figure out how to send message without extension UIViewController
        
        return true
    }
    return false
}
//
//func parseGameStartInfo(data: [String: Any]) -> Bool {
//    let json = JSON(data)
//}

//func parseGameUpdate(data: [String: Any]) -> Bool {
//
//}

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
    let payload = ["Name": game.getName(), "Password": game.getPassword(), "Description": game.getDescription(), "PlayerLimit": game.getMaxPlayers(), "PointLimit": game.getMaxPoints(), "TimeLimit": game.getTimeLimit(), "Mode": game.getMode().rawValue, "Boundaries": [], "Host": host] as [String: Any]
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
    let payload = ["Location": location, "Orientation": gameState.getUserOrientation()]
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

