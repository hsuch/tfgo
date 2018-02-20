//
//  GameLobbyViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/10/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class LobbyCustomViewCell: UITableViewCell {
    @IBOutlet weak var gamemodeLabel: UILabel!
    @IBOutlet weak var gameNameLabel: UILabel!
    @IBOutlet weak var gameDistanceLabel: UILabel!
    @IBOutlet weak var userCollection: UICollectionView!
}

class GameLobbyViewController: UIViewController, UITableViewDelegate, UICollectionViewDelegate, UICollectionViewDataSource {
    
    private var hasChosenGame = false
    private var showPublic = false

    override func viewDidLoad() {
        super.viewDidLoad()
        
        runTimer()

        DispatchQueue.global(qos: .background).async {
            if  gameState.getConnection().sendData(data: ShowGamesMsg()).isSuccess {
                while !self.hasChosenGame {
                    if MsgFromServer().parse() { }
                }
            }
        }
    }
    
    private var gamesPrivate: [Game] = []
    private var gamesPublic: [Game] = []
    private var gamesList: [Game] {
        get {
            return showPublic ? gamesPublic : gamesPrivate
        }
    }
    
    @IBAction func segmentButtonPress(_ sender: UISegmentedControl) {
        showPublic = !showPublic
    }
    
    private var updateTimer = Timer()
    
    private func runTimer() {
        updateTimer = Timer.scheduledTimer(timeInterval: 3, target: self,   selector: (#selector(GameLobbyViewController.updateGames)), userInfo: nil, repeats: true)
    }
    
    @IBOutlet weak var table: UITableView!
    
    @objc private func updateGames() {
        gamesPrivate = []
        gamesPublic = []
        for game in gameState.getFoundGames() {
            if game.getPassword() == nil {
                gamesPublic.append(game)
            } else {
                gamesPrivate.append(game)
            }
        }
        
        table.numberOfRows(inSection: gamesList.count)
        table.reloadData()
    }
    
    private var currentGame: Game?
    
    func tableView(_ tableView: UITableView, willDisplay cell: UITableViewCell, forRowAt indexPath: IndexPath) {
        let cell = tableView.dequeueReusableCell(withIdentifier: "Game", for: indexPath) as! LobbyCustomViewCell
        let game = gamesList[indexPath.row]
        currentGame = game
        switch game.getMode() {
        case .cp:
            cell.gamemodeLabel.text = "◆"
        case .multi:
            cell.gamemodeLabel.text = "⇥"
        case .payload:
            cell.gamemodeLabel.text = "❖"
        }
        cell.gameNameLabel.text = game.getName()!
        cell.gameDistanceLabel.text = "\(gameState.getDistanceFromGame(game: game)) units away"
        cell.userCollection.reloadData()
    }
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return currentGame!.getPlayers().count
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "icon", for: indexPath) as! IconViewCell
        cell.label.text = currentGame!.getPlayers()[indexPath.row].getIcon()
        return cell
    }

}
