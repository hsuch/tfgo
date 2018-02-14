//
//  Network.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation
import SwiftSocket

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
        return String(bytes: response, encoding: .utf8)?.data(using: .utf8)
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
    
    /* parse(): convert data array into appropriate data struct depending on message type */
    func parse() {
        //todo
    }
    
    init(conn: Connection) {
        let received = conn.recvData()
        self.data = try! JSONSerialization.jsonObject(with: received!, options: []) as! [String: Any]
        let type = data["Type"] as! String
        self.type = type
        
    }
}

class MsgToServer {
    private var action: String
    /* possible message actions:
        case CreateGame, ShowGames, ShowGameInfo, JoinGame, LocationUpdate, Fire
    */
    
    private var data: [String: Any]
    
    /* toJson(): convert action type and data array into server-readable json */
    func toJson() {
        // todo
    }
    
    init(action: String, data: [String: Any]) {
        self.action = action
        self.data = data
    }
}

class GameState {
    
    private var currentGame: Game?
    private var foundGames: [Game] = []
    
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
