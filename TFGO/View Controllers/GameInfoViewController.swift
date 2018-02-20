//
//  GameInfoViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class GameInfoViewController: UITableViewController, UICollectionViewDelegate, UICollectionViewDataSource {
    
    @IBOutlet weak var gameNameLbl: UILabel!
    @IBOutlet weak var gamemodeLbl: UILabel!
    @IBOutlet weak var descriptionLbl: UILabel!
    @IBOutlet weak var statusLbl: UILabel!
    
    private var game = gameState.getCurrentGame()
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        gameNameLbl.text = game.getName()
        switch game.getMode() {
        case .cp:
            gamemodeLbl.text = "◆ - Standard"
        case .multi:
            gamemodeLbl.text = "⇥ - Payload"
        case .payload:
            gamemodeLbl.text = "❖ - Multipoint"
        }
        descriptionLbl.text = game.getDescription()
        statusLbl.text = "Waiting for players - [\(game.getPlayers().count)/\(game.getMaxPlayers())]"
    }
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return game.getPlayers().count
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "Icon", for: indexPath) as! IconViewCell
        cell.label.text = game.getPlayers()[indexPath.row].getIcon()
        return cell
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "infoToWaiting" {
            if gameState.getConnection().sendData(data: JoinGameMsg(IDtoJoin: game.getID()!)).isSuccess {
                if MsgFromServer().parse() {
                    if game.isValid() {
                        return true
                    } else {
                        let alertController = UIAlertController(title: "This game is no longer valid", message:
                            "Please select a different game", preferredStyle: UIAlertControllerStyle.alert)
                        alertController.addAction(UIAlertAction(title: "It was hubris to try and join", style: UIAlertActionStyle.default,handler: nil))
                        self.present(alertController, animated: true, completion: nil)
                    }
                }
            }
            return false
        }
        return true
    }
}
