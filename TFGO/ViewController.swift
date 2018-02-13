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
        
//        background.image = UIImage(named: "dark-triangles")
        // Do any additional setup after loading the view, typically from a nib.
    }

    var state = GameState()
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        if let identifier = segue.identifier {
            switch identifier {
            case "Host Game":
                if let hostVC = segue.destination as? HostGameViewController {
                    hostVC.state = state
                }
            case "Game Lobby":
                if let lobbyVC = segue.destination as? GameLobbyViewController {
                    lobbyVC.state = state
                }
            default:
                break;
            }
        }
    }
}

