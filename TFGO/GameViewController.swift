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

class GameViewController: UIViewController, CLLocationManagerDelegate {
    
    
    @IBOutlet weak var game_map: MKMapView!
    
    let manager = CLLocationManager()
    
    var initialized = false  // boolean set to true after the first tracking of user's position
    
    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        
        // we want the most recent position of our user
        let location = locations [0]
        
        var region:MKCoordinateRegion
        
        let myLocation:CLLocationCoordinate2D = CLLocationCoordinate2DMake(location.coordinate.latitude, location.coordinate.longitude)
        
        gameState.getUser().setLocation(to: myLocation.latitude, to: myLocation.longitude)
        
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

        game_map.setRegion(region, animated: false)
        self.game_map.showsUserLocation = true
    }
    
    func locationManager(_ manager: CLLocationManager, didUpdateHeading newHeading: CLHeading) {
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
        
        runTimer()
        DispatchQueue.global(qos: .userInitiated).async {
            while true {
                if MsgFromServer().parse() {
                    DispatchQueue.main.async {
                        print("fuck")
                        if self.currentHealth != gameState.getUserHealth() {
                            self.currentHealth = gameState.getUserHealth()
                            let alertController = UIAlertController(title: "Temp", message: "You were hit", preferredStyle: UIAlertControllerStyle.alert)
                            alertController.addAction(UIAlertAction(title: "Ouch", style: UIAlertActionStyle.default,handler: nil))
                            self.present(alertController, animated: true, completion: nil)
                        }
                    }
                }
            }
        }
    }
    
    
    @IBAction func fireButton(_ sender: UIButton) {
        if gameState.getConnection().sendData(data: FireMsg()).isSuccess {
            //Put on Cooldown
        }
    }
    
    private var updateTimer = Timer()
    
    func runTimer() {
        updateTimer = Timer.scheduledTimer(timeInterval: 1, target: self,   selector: (#selector(GameViewController.update)), userInfo: nil, repeats: true)
    }
    
    private var currentHealth = 100
    
    @objc func update() {
        print("why")
        if gameState.getConnection().sendData(data: LocUpMsg()).isSuccess {
            print(gameState.getUser().getLocation())

        }
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        updateTimer.invalidate()
    }
    
    /*
     // MARK: - Navigation
     
     // In a storyboard-based application, you will often want to do a little preparation before navigation
     override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
     // Get the new view controller using segue.destinationViewController.
     // Pass the selected object to the new view controller.
     }
     */
    
}
