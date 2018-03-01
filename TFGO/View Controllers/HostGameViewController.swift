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

class CusAnnotation: NSObject, MKAnnotation {
    let title: String?
    let subtitle: String?
    let discipline: String
    let coordinate: CLLocationCoordinate2D
    var count = 0
    
    init(title: String, subtitle: String?, discipline: String, coordinate: CLLocationCoordinate2D) {
        self.title = title
        self.subtitle = subtitle
        self.discipline = discipline
        self.coordinate = coordinate
        
        super.init()
    }
}

class HostGameViewController: UITableViewController, UITextFieldDelegate, CLLocationManagerDelegate, MKMapViewDelegate {
    
    /* New game object */
    private let game = Game()
    
    @IBOutlet weak var host_map: MKMapView!
    
    let manager = CLLocationManager()
    
    var initialized = false  // boolean set to true after the first tracking of user's position
    
    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        
        // we want the most recent position of our user
        let location = locations [0]
        
        var region:MKCoordinateRegion
        
        let myLocation:CLLocationCoordinate2D = CLLocationCoordinate2DMake(location.coordinate.latitude, location.coordinate.longitude)
        
        // update the user's location information
        gameState.getUser().setLocation(to: myLocation.latitude, to: myLocation.longitude)
        
        // we only want to set the span the first time we locate the user;
        // we want to keep the current span otherwise
        if (initialized == false) {
            let span:MKCoordinateSpan = MKCoordinateSpanMake(0.0015, 0.0015)
            region = MKCoordinateRegionMake(myLocation, span)
            host_map.isRotateEnabled = false
            initialized = true
        
            host_map.setRegion(region, animated: false)
            self.host_map.showsUserLocation = true
            
            initialized = true
            
            host_map.delegate = self
        }
        else {
            region = host_map.region
            region.center = myLocation
        }
        
        // update the region of the map with the appropriate information
        host_map.setRegion(region, animated: false)
        self.host_map.showsUserLocation = true
        
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Initial game setup
        gameState.getUser().makeHost()
        game.setMode(to: .cp)
        
        let center = CLLocationCoordinate2DMake(41.794409, -87.595241)
        let span = MKCoordinateSpanMake(0.1, 0.1)
        let region = MKCoordinateRegionMake(center, span)
        self.host_map.setRegion(region, animated: true)
        
        
        let gestureRecognizer = UITapGestureRecognizer(target: self, action:#selector(handleLongPress))
        gestureRecognizer.delegate = self as? UIGestureRecognizerDelegate
        host_map.addGestureRecognizer(gestureRecognizer)
        
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
    /*
     func handleTap(gestureReconizer: UILongPressGestureRecognizer) {
     
     let location = gestureReconizer.locationc(in: host_map)
     let coordinate = host_map.convert(location,toCoordinateFrom: host_map)
     
     // Add annotation:
     let annotation = MKPointAnnotation()
     annotation.coordinate = coordinate
     annotation.title = "corner"
     annotation.subtitle = "tbd"
     host_map.addAnnotation(annotation)
     }
     */
    
    //@IBOutlet weak var host_map: MKMapView!
    var testpoint = CLLocationCoordinate2D(latitude: 1, longitude: 0)
    
    var boundaries = [MKMapPoint]()
    
    @objc func handleLongPress (gestureRecognizer: UITapGestureRecognizer) {
        print("test if running")
        //if gestureRecognizer.state == UIGestureRecognizerState.began {
        
        let touchPoint: CGPoint = gestureRecognizer.location(in: host_map)
        let newCoordinate: CLLocationCoordinate2D = host_map.convert(touchPoint, toCoordinateFrom: host_map)
        testpoint = newCoordinate
        print(testpoint)
        game.addBounary(to: MKMapPoint(x: testpoint.latitude, y: testpoint.longitude))
        //print(testpoint)
        addAnnotationOnLocation(pointedCoordinate: newCoordinate)
        //}
    }
    
    var annotations = [CusAnnotation]()
    
    //    mapview
    
    func mapView(_ mapView: MKMapView, didSelect view: MKAnnotationView) {
        //        print(view.annotation)
        let a = view.annotation as! CusAnnotation
        //        print(a.discipline)
        a.count += 1
        //        print(annotations.index(of: a))
        if a.count > 1 {
            let p = MKMapPoint(x: a.coordinate.longitude, y: a.coordinate.latitude)
            for (index, b) in boundaries.enumerated() {
                if b.x == p.x && b.y == p.y {
                    boundaries.remove(at: index)
                }
            }
            annotations.remove(at: annotations.index(of: a)!)
            mapView.removeAnnotation(a)
            print("called")
        }
        
    }
    var counter = 0
    func addAnnotationOnLocation(pointedCoordinate: CLLocationCoordinate2D) {
        //var title1 = "Boundary Point"
        //title1 +=
        let annotation = CusAnnotation(title: "Boundary Point" + String((counter % 4 + 1)), subtitle: "", discipline: "\(pointedCoordinate.latitude),\(pointedCoordinate.longitude)", coordinate: pointedCoordinate)
        //        annotation.discipline =
        annotations.append(annotation)
        if annotations.count <= 4 {
            boundaries.append(MKMapPoint(x: pointedCoordinate.latitude, y: pointedCoordinate.longitude))
            print(pointedCoordinate)
            counter += 1
            game.setBoundaries(boundaries)
        }
        if annotations.count > 4 {
            host_map.removeAnnotation(annotations.first!)
            annotations.removeFirst()
            print(counter % 4)
            boundaries[counter % 4] = MKMapPoint(x: pointedCoordinate.latitude, y: pointedCoordinate.longitude)
            print(pointedCoordinate)
            print(boundaries[0])
            print(boundaries[1])
            print(boundaries[2])
            print(boundaries[3])
            counter += 1
            game.setBoundaries(boundaries)
        }
        host_map.addAnnotation(annotation)
    }
    
    /* Text field variables */
    @IBOutlet weak var nameField: UITextField!
    @IBOutlet weak var descriptionField: UITextField!
    @IBOutlet weak var passwordField: UITextField!
    
    /* Hide keyboard when the return key is pressed */
    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        return textField.resignFirstResponder()
    }
    
    /* Gamemode information variables */
    @IBOutlet private var gamemodeButtons: [UIButton]!
    @IBOutlet weak var gamemodeLabel: UILabel!
    private var gamemodes = [Gamemode.cp, .payload, .multi]
    @IBOutlet weak var objectiveCell: UITableViewCell!
    
    /* chooseGamemode() - Action on any gamemode button */
    /* Sets gamemode to the button clicked and displays any additional information if required */
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
    
    /* Game information labels */
    @IBOutlet weak var playerLabel: UILabel!
    @IBOutlet weak var timeLabel: UILabel!
    @IBOutlet weak var pointLabel: UILabel!
    @IBOutlet weak var objectiveLabel: UILabel!
    
    /* step() - Action on any stepper */
    /* Changes the value of a game information label according to a stepper press */
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
    
    /* Password variables */
    private var usePassword = false
    @IBOutlet weak var passwordCell: UITableViewCell!
    
    /* passwordSwitch() - Action on password switch */
    /* Shows or hides the password table cell when needed */
    @IBAction func passwordSwitch(_ sender: UISwitch) {
        usePassword = sender.isOn
        passwordCell.isHidden = !sender.isOn
    }
    
    /* shouldPerformSegue() - override */
    /* Segue only if the game was able to send the game to the server */
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "Create Game" {
            if game.validBoundaries() && gameState.setCurrentGame(to: game) {
                if gameState.getConnection().sendData(data: CreateGameMsg(game: game)).isSuccess {
                    return true
                }
            }
        }
        return false
    }
    
    /* checkGame() - Action on Create Game button */
     /* Checks if the game is valid and causes an alert if it is not */
    @IBAction func checkGame(_ sender: UIButton) {
        if game.setName(to: nameField.text ?? "") {}
        if game.setDescription(to: descriptionField.text ?? "") {}
        if game.setPassword(to: passwordField.text ?? "") {}
        if !game.isValid() || !game.validBoundaries() {
            let alertController = UIAlertController(title: "Invalid Game", message:
                "Please ensure that you have filled all required fields", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "I hope you can forgive me", style: UIAlertActionStyle.default,handler: nil))
            self.present(alertController, animated: true, completion: nil)
        }
    }
}

