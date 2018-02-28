//
//  HostGameViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit
import MapKit
import CoreLocation

class HostGameViewController: UITableViewController, UITextFieldDelegate, CLLocationManagerDelegate {
    
    private let game = Game()
    
    @IBOutlet weak var host_map: MKMapView!
    
    let manager = CLLocationManager()
    
    var initialized = false  // boolean set to true after the first tracking of user's position
    
    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        
        
        if (initialized == false) {
            // we want the most recent position of our user
            let location = locations [0]
        
            var region:MKCoordinateRegion
        
            let myLocation = CLLocationCoordinate2DMake(location.coordinate.latitude, location.coordinate.longitude)
        
            let span:MKCoordinateSpan = MKCoordinateSpanMake(0.0015, 0.0015)
            region = MKCoordinateRegionMake(myLocation, span)
            host_map.isRotateEnabled = false
            initialized = true
        
            host_map.setRegion(region, animated: false)
            self.host_map.showsUserLocation = true
            
            let myLat = myLocation.latitude
            let myLon = myLocation.longitude
            
            // These are hardcoded boundaries for the purpose of testing iteration 1 code
            game.setBoundaries([MKMapPointMake(myLat + 0.1, myLon + 0.1), MKMapPointMake(myLat + 0.1, myLon - 0.1), MKMapPointMake(myLat - 0.1, myLon + 0.1), MKMapPointMake(myLat - 0.1, myLon - 0.1)])
            
            initialized = true
        }
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        game.setMode(to: .cp)
        
        let center = CLLocationCoordinate2DMake(41.794409, -87.595241)
        let span = MKCoordinateSpanMake(0.1, 0.1)
        let region = MKCoordinateRegionMake(center, span)
        self.host_map.setRegion(region, animated: true)
        
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyBest
        manager.requestWhenInUseAuthorization()
        manager.startUpdatingLocation()
        manager.startUpdatingHeading()
        // Uncomment the following line to preserve selection between presentations
        // self.clearsSelectionOnViewWillAppear = false
        
        // Uncomment the following line to display an Edit button in the navigation bar for this view controller.
        // self.navigationItem.rightBarButtonItem = self.editButtonItem
    }
    
    @IBOutlet weak var nameField: UITextField!
    @IBOutlet weak var descriptionField: UITextField!
    @IBOutlet weak var passwordField: UITextField!

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        return textField.resignFirstResponder()
    }
    
    @IBOutlet private var gamemodeButtons: [UIButton]!
    
    @IBOutlet weak var gamemodeLabel: UILabel!
    
    private var gamemodes = [Gamemode.cp, .payload, .multi]
    
    @IBOutlet weak var objectiveCell: UITableViewCell!
    
    @IBAction private func chooseGamemode(_ sender: UIButton) {
        if let modeIndex = gamemodeButtons.index(of: sender) {
            gamemodeLabel.text = gamemodes[modeIndex].rawValue
            game.setMode(to: gamemodes[modeIndex])
            if gamemodes[modeIndex] == .multi {
                objectiveCell.isHidden = false
            } else {
                objectiveCell.isHidden = true
            }
        }
    }
    
    @IBOutlet weak var playerLabel: UILabel!
    @IBOutlet weak var timeLabel: UILabel!
    @IBOutlet weak var pointLabel: UILabel!
    @IBOutlet weak var objectiveLabel: UILabel!
    
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
        case 50:
            game.setMaxObjectives(to: value)
            objectiveLabel.text = "\(value) Control Points"
        default:
            break;
        }
    }
    
    
    private var usePassword = false
    @IBOutlet weak var passwordCell: UITableViewCell!
    
    @IBAction func passwordSwitch(_ sender: UISwitch) {
        usePassword = sender.isOn
        passwordCell.isHidden = !sender.isOn
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "Create Game" {
            if gameState.setCurrentGame(to: game) {
                if gameState.getConnection().sendData(data: CreateGameMsg(game: game)).isSuccess {
                    gameState.getUser().makeHost()
                    return true
                }
            }
        }
        //Give incomplete gamestate alert
        return false
    }
    
    @IBAction func checkGame(_ sender: UIButton) {
        if game.setName(to: nameField.text ?? ""), game.setDescription(to: descriptionField.text ?? "") {}
        if usePassword {
            if game.setPassword(to: passwordField.text ?? "") {}
        }
        if !game.isValid() {
            let alertController = UIAlertController(title: "Invalid Game", message:
                "Please ensure that you have filled all required fields", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "I hope you can forgive me", style: UIAlertActionStyle.default,handler: nil))
            self.present(alertController, animated: true, completion: nil)
        }
    }
}

