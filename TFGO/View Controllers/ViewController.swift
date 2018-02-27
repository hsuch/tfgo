//
//  ViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 1/30/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

@IBDesignable
class ViewController: UIViewController {
    
    @IBOutlet weak var background: UIImageView!

    override func viewDidLoad() {
        super.viewDidLoad()
        
        background.image = UIImage(named: "bg2")
        self.navigationController?.setNavigationBarHidden(true, animated: true)
        // Do any additional setup after loading the view, typically from a nib.
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "temp" {
            let connection = gameState.getConnection()
            if connection.sendData(data: ShowGamesMsg()).isSuccess {
                if handleMsgFromServer(), gameState.findPublicGames().count > 0 {
                    let game = gameState.findPublicGames()[0]
                    if connection.sendData(data: JoinGameMsg(IDtoJoin: game.getID()!)).isSuccess {
                        return gameState.setCurrentGame(to: game)
                    }
                }
            }
            return false
        } else {
            return true
        }
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        self.navigationController?.setNavigationBarHidden(false, animated: true)
    }
    
}

