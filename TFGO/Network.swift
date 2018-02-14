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
    private var address: String?
    private var port: UInt32?
    let client = TCPClient(address: address, port: port)
    
    func sendData(data: Data) -> Bool {
        let result = client.send(data)
        return result
    }
    
    func recvData() -> [Data] {
        guard let response = client.read(1024*10)
        else { return nil }
        return String(bytes: response, encoding: .utf8)
    }
    
    init(){
        client.connect(timeout: 10) // this should probably have success and failure case but whatever
    }
}

class ServerMessage {
    private enum type: String {
        case PlayerListUpdate, AvailableGames, GameInfo, JoinGame, GameStartInfo, GameUpdate, StatusUpdate
    }
    private var data: JSONSerialization
    
    func parse() {
        
    }
    
    init(conn: Connection) {
        received = conn.recvData()
        data = try? JSONSerialization.jsonObject(with: received, options: [])
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
