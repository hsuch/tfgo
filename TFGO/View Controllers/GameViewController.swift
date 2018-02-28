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

class GameViewController: UIViewController, CLLocationManagerDelegate, MKMapViewDelegate {
    
    @IBOutlet weak var game_map: MKMapView!
    
    private var game = gameState.getCurrentGame()
    
    let manager = CLLocationManager() // used to track the user's location
    
    var initialized = false  // boolean set to true after the first tracking of user's position
    
    var playerLocs: [MKPointAnnotation] = []
    
    /* time variables */
    let calendar = Calendar.current
    var startcomponents = DateComponents()
    var startArr: [String] = []
    var starttime = Date()
    
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
    
    func game_map(game_map: MKMapView, rendererForOverlay overlay: MKOverlay) -> MKOverlayRenderer {
        if overlay is MKCircle {
            var circleRenderer = MKCircleRenderer(overlay: overlay)
            circleRenderer.fillColor = UIColor.blue
            circleRenderer.strokeColor = UIColor.red
            circleRenderer.lineWidth = 1
            return circleRenderer
        }
        return MKOverlayRenderer(overlay: overlay)
    }
    
    func addRadiusCircle(location: CLLocation, radius: CLLocationDistance) {
        self.game_map.delegate = self
        var circle = MKCircle(center: location.coordinate, radius: radius)
        self.game_map.add(circle)
    }
    
    private var status = (100, 0)
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
        
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyBest
        manager.requestWhenInUseAuthorization()
        manager.startUpdatingLocation()
        manager.startUpdatingHeading()
        
        // set start time variables
        startcomponents = DateComponents()
        startArr = game.getStartTime()
        startcomponents.year = Int(startArr[0])
        startcomponents.month = Int(startArr[1])
        startcomponents.day = Int(startArr[2])
        startcomponents.hour = Int(startArr[3])
        startcomponents.minute = Int(startArr[4])
        startcomponents.second = Int(startArr[5])
        starttime = calendar.date(from: startcomponents)!
        
        // used to differentiate between objectives if we have more than one
        var objectiveNumber = 1
        
        // now we make a pin on the map for each objective
        for objective in game.getObjectives() {
            let annotation = MKPointAnnotation()
            annotation.coordinate = CLLocationCoordinate2D(latitude: objective.getXLoc(), longitude: objective.getYLoc())
            annotation.title = "OBJECTIVE"
            annotation.subtitle = String(objectiveNumber)
            game_map.addAnnotation(annotation)
            
            let center = CLLocation(latitude: objective.getXLoc(), longitude: objective.getYLoc())
            addRadiusCircle(location: center, radius: objective.getRadius())
            
            objectiveNumber = objectiveNumber + 1
        }
        
        // next, we set pins for the bases
        let redBaseLoc = gameState.getCurrentGame().getRedBaseLoc()
        let blueBaseLoc = gameState.getCurrentGame().getBlueBaseLoc()
        let rbAnnotation = MKPointAnnotation()
        let bbAnnotation = MKPointAnnotation()
        rbAnnotation.coordinate = CLLocationCoordinate2D(latitude: redBaseLoc.x, longitude: redBaseLoc.y)
        rbAnnotation.title = "RED BASE"
        game_map.addAnnotation(rbAnnotation)
        bbAnnotation.coordinate = CLLocationCoordinate2D(latitude: blueBaseLoc.x, longitude: blueBaseLoc.y)
        bbAnnotation.title = "BLUE BASE"
        game_map.addAnnotation(bbAnnotation)
        
        // now we set pins for each player in the game
        gameState.getCurrentGame().updatePlayerAnnotations()
        let annotations = gameState.getCurrentGame().getPlayerAnnotations()
        game_map.addAnnotations(annotations)
        
        // finally, we draw the games boundaries
        let gameBounds = game.getBoundaries()
        let drawBounds = UIBezierPath.init()
        var first = true
        
        for bound in gameBounds {
            if (first) {
                drawBounds.move(to: CGPoint(x: bound.x, y: bound.y))
                first = false
            }
            else {
                drawBounds.addLine(to: CGPoint(x: bound.x, y: bound.y))
            }
        }
        // connect the last and first boundaries, then draw the boundaries
        drawBounds.move(to: CGPoint(x: gameBounds[0].x, y: gameBounds[0].y))
        drawBounds.stroke()
        
        runTimer()
        DispatchQueue.global(qos: .userInitiated).async {
            while true {
                if handleMsgFromServer() {
                    DispatchQueue.main.async {
                        self.talkShitGetHit()
                    }
                }
            }
        }
    }
    
    @IBOutlet weak var armorBar: UIProgressView!
    @IBOutlet weak var healthBar: UIProgressView!
    
    override func viewWillAppear(_ animated: Bool) {
        armorBar.layer.cornerRadius = 2.0
        armorBar.layer.borderWidth = 3.0
        armorBar.layer.borderColor = #colorLiteral(red: 0.8078431487, green: 0.02745098062, blue: 0.3333333433, alpha: 1)
        healthBar.layer.cornerRadius = 2.0
        healthBar.layer.borderWidth = 3.0
        healthBar.layer.borderColor = #colorLiteral(red: 0.3333333433, green: 0.3333333433, blue: 0.3333333433, alpha: 1)
    }
    
    private func talkShitGetHit() {
        if status != (gameState.getUserHealth(), gameState.getUserArmor()) {
            status = (gameState.getUserHealth(), gameState.getUserArmor())
            AudioServicesPlaySystemSound(kSystemSoundID_Vibrate)
            armorBar.setProgress(Float(self.status.0)/100, animated: true)
            healthBar.setProgress(Float(self.status.1)/100, animated: true)
            let alertController = UIAlertController(title: "Temp", message: "You were hit", preferredStyle: UIAlertControllerStyle.alert)
            alertController.addAction(UIAlertAction(title: "Ouch", style: UIAlertActionStyle.default,handler: nil))
            present(alertController, animated: true, completion: nil)
        }
    }
    
    @IBOutlet weak var clock: UILabel!
    @IBOutlet weak var redScore: UILabel!
    @IBOutlet weak var blueScore: UILabel!
    
    private func tick() {
        let curtime = Date()
        
        if(curtime < starttime) { // game has not started yet
            let diff = calendar.dateComponents([.minute, .second], from: curtime, to: starttime)
            clock.text = "-" + String(diff.minute!) + ":" + String(diff.second!)
        }
        else { // game started
            let diff = calendar.dateComponents([.minute, .second], from: curtime, to: starttime.addingTimeInterval(Double(game.getTimeLimit()) * 60.0))
            clock.text = String(diff.minute!) + ":" + String(diff.second!)
        }
    }
    
    @IBAction func fireButton(_ sender: UIButton) {
        if gameState.getConnection().sendData(data: FireMsg()).isSuccess {
            //Put on Cooldown. Not necessary for Iteration 1
        }
    }
    
    var updateTimer = Timer()
    
    func runTimer() {
        updateTimer = Timer.scheduledTimer(timeInterval: 1, target: self,   selector: (#selector(GameViewController.update)), userInfo: nil, repeats: true)
    }
    
    @objc func update() {
        if gameState.getConnection().sendData(data: LocUpMsg()).isSuccess {
            print(gameState.getUser().getLocation())
        }
        
        redScore.text = "\(game.getRedPoints())"
        blueScore.text = "\(game.getBluePoints())"
        tick()
        
        // update the locations of other players on the map
        gameState.getCurrentGame().updatePlayerAnnotations()
        
        
//        // first we remove all the previous player annotations
//        for playerLoc in playerLocs {
//            let annotation = playerLoc as MKAnnotation
//            self.game_map.removeAnnotation(annotation)
//        }
//
//        // then we build a new list of player annotations
//        let playerList = gameState.getCurrentGame().getPlayers()
//        for player in playerList {
//            if player.getName() != gameState.getUser().getName() {
//                let annotation = MKPointAnnotation()
//                let loc = player.getLocation().coordinate
//                annotation.coordinate = CLLocationCoordinate2D(latitude: loc.latitude, longitude: loc.longitude)
//                annotation.title = player.getName()
//                annotation.subtitle = player.getTeam()
//                game_map.addAnnotation(annotation)
//                playerLocs.append(annotation)
//            }
//        }
        
    }
    
    @IBAction func leaveGame(_ sender: UIButton) {
        let actionController = UIAlertController(title: nil, message:
            "Are you sure you want to leave?", preferredStyle: UIAlertControllerStyle.actionSheet)
        actionController.addAction(UIAlertAction(title: "Yep", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            self.performSegue(withIdentifier: "leaveGame", sender: nil)
        }))
        actionController.addAction(UIAlertAction(title: "Nope", style: UIAlertActionStyle.cancel,handler: nil))
        self.present(actionController, animated: true, completion: nil)
    }
    
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        if segue.identifier != "inventory" {
            updateTimer.invalidate()
        }
    }
    
}
