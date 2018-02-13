//
//  HostGameViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit

class HostGameViewController: UITableViewController, UITextFieldDelegate {
    
    var state: GameState?
    private let game = Game()
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        if game.setMaxPlayers(to: 2) {}
        if game.setMaxPoints(to: 10) {}
        if game.setTimeLimit(to: 2) {}
        game.setMode(to: .cp)
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
    
//    func textFieldDidEndEditing(_ textField: UITextField, reason: UITextFieldDidEndEditingReason) {
//        if reason == .committed, let text = textField.text {
//            if textField == nameField {
//                if !game.setName(to: text) {
//                    textField.text = ""
//                    //give invalid name message
//                }
//            } else if textField == descriptionField {
//                if !game.setDescription(to: text) {
//                    textField.text = ""
//                    //give invalid description message
//                }
//            } else if textField == passwordField {
//                if usePassword {
//                    if !game.setPassword(to: text) {
//                        textField.text = ""
//                        //give invalid password message
//                    }
//                }
//            }
//        }
//    }

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
            if game.setMaxPlayers(to: value) {
                playerLabel.text = "\(value) Players"
            }
        case 525600:
            if game.setTimeLimit(to: value) {
                timeLabel.text = "\(value) Minutes"
            }
        case 100000:
            if game.setMaxPoints(to: value) {
                pointLabel.text = "\(value) Points"
            }
        default:
            break;
        }
    }
    
    
    private var usePassword = false
    
    @IBAction func passwordSwitch(_ sender: UISwitch) {
        usePassword = sender.isOn
        if usePassword {
            passwordField.backgroundColor = #colorLiteral(red: 1.0, green: 1.0, blue: 1.0, alpha: 1.0)
        } else {
            passwordField.backgroundColor = #colorLiteral(red: 0.7540688515, green: 0.7540867925, blue: 0.7540771365, alpha: 1)
        }
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "Create Game" {
            if state?.setCurrentGame(to: game) ?? false {
                return true
            }
        }
        //Give incomplete gamestate alert
        return false
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        if let identifier = segue.identifier {
            if identifier == "Create Game" {
                if let waitingVC = segue.destination as? WaitingViewController {
                    waitingVC.state = state
                }
            }
        }
    }
    
    

    // MARK: - Table view data source

    /*
    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "reuseIdentifier", for: indexPath)

        // Configure the cell...

        return cell
    }
    */

    /*
    // Override to support conditional editing of the table view.
    override func tableView(_ tableView: UITableView, canEditRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the specified item to be editable.
        return true
    }
    */

    /*
    // Override to support editing the table view.
    override func tableView(_ tableView: UITableView, commit editingStyle: UITableViewCellEditingStyle, forRowAt indexPath: IndexPath) {
        if editingStyle == .delete {
            // Delete the row from the data source
            tableView.deleteRows(at: [indexPath], with: .fade)
        } else if editingStyle == .insert {
            // Create a new instance of the appropriate class, insert it into the array, and add a new row to the table view
        }    
    }
    */

    /*
    // Override to support rearranging the table view.
    override func tableView(_ tableView: UITableView, moveRowAt fromIndexPath: IndexPath, to: IndexPath) {

    }
    */

    /*
    // Override to support conditional rearranging of the table view.
    override func tableView(_ tableView: UITableView, canMoveRowAt indexPath: IndexPath) -> Bool {
        // Return false if you do not want the item to be re-orderable.
        return true
    }
    */

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destinationViewController.
        // Pass the selected object to the new view controller.
    }
    */

}
