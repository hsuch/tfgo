//
//  GameViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/9/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit
import MapKit
import CoreLocation
import AudioToolbox.AudioServices

class GameViewController: UIViewController, CLLocationManagerDelegate {
    
    @IBOutlet weak var game_map: MKMapView!
    
    private var game = gameState.getCurrentGame()
    
    let manager = CLLocationManager() // used to track the user's location
    
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
            let span:MKCoordinateSpan = MKCoordinateSpanMake(0.01, 0.01)
            region = MKCoordinateRegionMake(myLocation, span)
            game_map.isRotateEnabled = false
            initialized = true
        }
        else {
            region = game_map.region
            region.center = myLocation
        }

        // update the region of the map with the appropriate information
        game_map.setRegion(region, animated: false)
        self.game_map.showsUserLocation = true
        
    }
    
    func locationManager(_ manager: CLLocationManager, didUpdateHeading newHeading: CLHeading) {
        //update the user's orientation with respect to North every time it changes
        gameState.getUser().setOrientation(to: Float(newHeading.magneticHeading))
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
        
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyBest
        manager.requestWhenInUseAuthorization()
        manager.startUpdatingLocation()
        manager.startUpdatingHeading()
        
        // we want to mark objective locations with pins
        let game = gameState.getCurrentGame()
        let annotation = MKPointAnnotation()
        
        // used to differentiate between objectives if we have more than one
        var objectiveNumber = 1
        
        // now we make a pin on the map for each objective
        for objective in game.getObjectives() {
            annotation.coordinate = CLLocationCoordinate2D(latitude: objective.getXLoc(), longitude: objective.getYLoc())
            annotation.title = "OBJECTIVE"
            annotation.subtitle = String(objectiveNumber)
            game_map.addAnnotation(annotation)
            
            objectiveNumber = objectiveNumber + 1
        }
        
        runTimer()
        DispatchQueue.global(qos: .userInitiated).async {
            while true {
                if handleMsgFromServer() {
                    DispatchQueue.main.async {
                        if self.currentHealth != gameState.getUserHealth() {
                            self.currentHealth = gameState.getUserHealth()
                            AudioServicesPlaySystemSound(kSystemSoundID_Vibrate)
                            let alertController = UIAlertController(title: "Temp", message: "You were hit", preferredStyle: UIAlertControllerStyle.alert)
                            alertController.addAction(UIAlertAction(title: "Ouch", style: UIAlertActionStyle.default,handler: nil))
                            self.present(alertController, animated: true, completion: nil)
                        }
                    }
                }
            }
        }
    }
    
    @IBOutlet weak var clock: UILabel!
    @IBOutlet weak var redScore: UILabel!
    @IBOutlet weak var blueScore: UILabel!
    
    private func tick() {
//        let calender = Calendar.current
//        let date = calender.dateComponents([.year,.month,.day,.hour,.minute,.second], from: Date())
//
//
//        let year = date.year
//        let month = date.month
//        let day = date.day
//        let hour = date.hour
//        let minute = date.minute
//        let second = date.second
    }
    
    @IBAction func fireButton(_ sender: UIButton) {
        if gameState.getConnection().sendData(data: FireMsg()).isSuccess {
            //Put on Cooldown. Not necessary for Iteration 1
        }
        redScore.text = "\(game.getRedPoints())"
        blueScore.text = "\(game.getBluePoints())"
        tick()
    }
    
    var updateTimer = Timer()
    
    func runTimer() {
        updateTimer = Timer.scheduledTimer(timeInterval: 1, target: self,   selector: (#selector(GameViewController.update)), userInfo: nil, repeats: true)
    }
    
    private var currentHealth = 100
    
    @objc func update() {
        if gameState.getConnection().sendData(data: LocUpMsg()).isSuccess {
            print(gameState.getUser().getLocation())
        }
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        updateTimer.invalidate()
    }
    
}
