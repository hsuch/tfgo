//
//  Network.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import SwiftSocket

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
    func parse() -> AnyObject {
        switch type {
        case "PlayerListUpdate":
            return parsePlayerListUpdate(data: data)
        default:
            break
        }
    }

    init(conn: Connection) {
        let received = conn.recvData()
        self.data = try! JSONSerialization.jsonObject(with: received!, options: []) as! [String: Any]
        let type = data.removeValueForKey("Type")
        self.type = type

    }
}

/* Parsing functions: helper functions called by parse() to parse different messages */

func parsePlayerListUpdate(data: [String: Any]) -> Player {

}

class MsgToServer {
    private var action: String
    /* possible message actions:
        case CreateGame, ShowGames, ShowGameInfo, JoinGame, StartGame, LocationUpdate, Fire
    */

    private var data: [String: Any]

    /* toJson(): convert action type and data array into server-readable json */
    func toJson() -> Data {
        let retval = Data.init() //todo
        return retval
    }

    init(action: String, data: [String: Any]) {
        self.action = action
        self.data = data
    }
}

/* Message generators: the following functions generate messages that can be directly sent to the server via Connection.sendData()*/

func CreateGameMsg(game: Game) -> Data {
    // todo
    return Data.init()
}
func ShowGamesMsg() -> Data {
    let payload = ["Name": gameState.getUserName(), "Icon": gameState.getUserIcon()]
    return MsgToServer(action: "ShowGames", data: payload).toJson()
}
func ShowGameInfo(IDtoShow: String) -> Data {
    return MsgToServer(action: "ShowGameInfo", data: ["GameID": IDtoShow]).toJson()
}
func JoinGameMsg(IDtoJoin: String) -> Data {
    return MsgToServer(action: "JoinGame", data: ["GameID": IDtoJoin]).toJson()
}
func StartGameMsg() -> Data {
    return MsgToServer(action: "StartGame", data: [:]).toJson()
}
func LocUpMsg() -> Data {
    // todo, take location from this client's player
    return MsgToServer(action: "LocationUpdate", data: [:]).toJson()
}
func FireMsg() -> Data {
    // todo, take orientation and weapon from this client's player
    return MsgToServer(action: "Fire", data: [:]).toJson()

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
    
    func getUser() -> Player {
        return user
    }
    
    /* Do not call unless a game exists!!! */
    func getCurrentGame() -> Game {
        return currentGame!
    }
    
    
    func setCurrentGame(to game: Game) -> Bool {
        let valid = game.isValid()
        if valid {
            currentGame = game
        }
        return valid
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
