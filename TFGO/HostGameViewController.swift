//
//  HostGameViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit
import MapKit

class HostGameViewController: UITableViewController, UITextFieldDelegate {
    
    private let game = Game()
    
    @IBOutlet weak var host_map: MKMapView!
    
    override func viewDidLoad() {
        super.viewDidLoad()

        game.setMode(to: .cp)
        
        let center = CLLocationCoordinate2DMake(41.794409, -87.595241)
        let span = MKCoordinateSpanMake(0.1, 0.1)
        let region = MKCoordinateRegionMake(center, span)
        self.host_map.setRegion(region, animated: true)
        // Uncomment the following line to preserve selection between presentations
        // self.clearsSelectionOnViewWillAppear = false
        
        // Uncomment the following line to display an Edit button in the navigation bar for this view controller.
        // self.navigationItem.rightBarButtonItem = self.editButtonItem
    }
    
    
    
    @IBOutlet weak var nameField: UITextField!
    @IBOutlet weak var descriptionField: UITextField!
    @IBOutlet weak var passwordField: UITextField!
    
    func textField(_ textField: UITextField, shouldChangeCharactersIn range: NSRange, replacementString string: String) -> Bool {
        var name = game.getName() ?? ""
        name.append(string)
        if textField == nameField {
            if !game.setName(to: name) {
                //give invalid name message
                return false
            }
        } else if textField == descriptionField {
            if !game.setDescription(to: name) {
                //give invalid description message
                return false
            }
        } else if textField == passwordField {
            if usePassword {
                if !game.setPassword(to: name) {
                    //give invalid password message
                    return false
                }
            }
        }
        return true
    }
    
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        return textField.resignFirstResponder()
    }
    
    @IBOutlet private var gamemodeButtons: [UIButton]!
    
    @IBOutlet weak var gamemodeLabel: UILabel!
    
    private var gamemodes = [Gamemode.cp, .payload, .multi]
    
    @IBAction private func chooseGamemode(_ sender: UIButton) {
        if let modeIndex = gamemodeButtons.index(of: sender) {
            gamemodeLabel.text = gamemodes[modeIndex].rawValue
            game.setMode(to: gamemodes[modeIndex])
        }
    }
    
    @IBOutlet weak var playerLabel: UILabel!
    @IBOutlet weak var timeLabel: UILabel!
    @IBOutlet weak var pointLabel: UILabel!
    
    @IBAction func step(_ sender: UIStepper) {
        let value = Int(sender.value)
        switch sender.maximumValue {
        case 10000:
            game.setMaxPlayers(to: value)
            playerLabel.text = "\(value) Players"
        case 525600:
            game.setTimeLimit(to: value)
            timeLabel.text = "\(value) Minutes"
        case 100000:
            game.setMaxPoints(to: value)
            pointLabel.text = "\(value) Points"
        default:
            break;
        }
    }
    
    
    private var usePassword = false
    
    @IBAction func passwordSwitch(_ sender: UISwitch) {
        usePassword = sender.isOn
        if usePassword {
            passwordField.backgroundColor =  #colorLiteral(red: 1.0, green: 1.0, blue: 1.0, alpha: 1.0)
        } else {
            passwordField.backgroundColor =  #colorLiteral(red: 0.7540688515, green: 0.7540867925, blue: 0.7540771365, alpha: 1)
        }
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "Create Game" {
            if gameState.setCurrentGame(to: game) {
                if gameState.getConnection().sendData(data: CreateGameMsg(game: game)).isSuccess {
                    return true
                }
            }
        }
        //Give incomplete gamestate alert
        return false
    }
    
    @IBAction func checkGame(_ sender: UIButton) {
        if !game.isValid() {
            let alertController = UIAlertController(title: "Invalid Game", message:
                "Please ensure that you have filled all required fields", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "I hope you can forgive me", style: UIAlertActionStyle.default,handler: nil))
            self.present(alertController, animated: true, completion: nil)
        }
    }
}

