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
    private var player = gameState.getUser()
    
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
        self.game_map.showsUserLocation = false
        
    }
    
    func locationManager(_ manager: CLLocationManager, didUpdateHeading newHeading: CLHeading) {
        //update the user's orientation with respect to North every time it changes
        gameState.getUser().setOrientation(to: Float(newHeading.magneticHeading))
    }
    
    func addBoundary() {
        let bounds = game.getBoundaries()
        var polyBounds: [CLLocationCoordinate2D] = []
        for bound in bounds {
            let polyBound = CLLocationCoordinate2DMake(bound.x, bound.y)
            polyBounds.append(polyBound)
        }
        let polygon = MKPolygon(coordinates: polyBounds, count: polyBounds.count)
        game_map.add(polygon)
    }
    
    func addRadiusCircle(location: CLLocation, radius: CLLocationDistance) {
        let circle = MKCircle(center: location.coordinate, radius: radius)
        game_map.add(circle)
    }
    
    func mapView(_ mapView: MKMapView, rendererFor overlay: MKOverlay) -> MKOverlayRenderer! {
        if overlay is MKPolygon {
            let polygonView = MKPolygonRenderer(overlay: overlay)
            polygonView.strokeColor = #colorLiteral(red: 0.5568627715, green: 0.3529411852, blue: 0.9686274529, alpha: 1)
            polygonView.lineWidth = 2.0
            
            return polygonView
        }
        else if overlay is MKCircle {
            let circleRenderer = MKCircleRenderer(overlay: overlay)
            circleRenderer.fillColor = #colorLiteral(red: 0.8039215803, green: 0.8039215803, blue: 0.8039215803, alpha: 1)
            circleRenderer.strokeColor = #colorLiteral(red: 0.2549019754, green: 0.2745098174, blue: 0.3019607961, alpha: 1)
            circleRenderer.lineWidth = 1.0
            circleRenderer.alpha = 0.5
            
            return circleRenderer
        }
        
        return nil
    }
    
    func mapView(_ mapView: MKMapView, viewFor annotation: MKAnnotation) -> MKAnnotationView? {
        let annotationView = MKAnnotationView(annotation: annotation, reuseIdentifier: "rahhhh")
       // annotationView.transform = CGAffineTransform(scaleX: 0.1, y: 0.1)
        annotationView.canShowCallout = true
        if let title = annotation.title, let subtitle = annotation.subtitle {
            if title == "OBJECTIVE" {
                if subtitle == "Neutral" {
                    annotationView.image = UIImage(named: "cap_gray")
                }
                else if subtitle == "Blue" {
                    annotationView.image = UIImage(named: "cap_blue")
                }
                else {
                    annotationView.image = UIImage(named: "cap_red")
                }
            }
            else if title == "RED BASE" {
                annotationView.image = UIImage(named: "base_red")
            }
            else if title == "BLUE BASE" {
                annotationView.image = UIImage(named: "base_blue")
            }
            else if subtitle == "Red" {
                annotationView.image = UIImage(named: "player_red")
                annotationView.transform = CGAffineTransform(scaleX: 0.8, y: 0.8)
            }
            else if subtitle == "Blue" {
                annotationView.image = UIImage(named: "player_blue")
                annotationView.transform = CGAffineTransform(scaleX: 0.8, y: 0.8)
            }
            else {
                annotationView.image = UIImage(named: "pickup")
            }
        }
        return annotationView
    }
    
    private var status = (100, 0)
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
        
        game_map.delegate = self
        
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
        
        // now we make a pin on the map for each objective
        for objective in game.getObjectives() {

            game_map.addAnnotation(objective.getAnnotation())
            
            let center = CLLocation(latitude: objective.getXLoc(), longitude: objective.getYLoc())
            addRadiusCircle(location: center, radius: objective.getRadius())
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
        
        // we also want to draw circles representing each base's range
        addRadiusCircle(location: CLLocation(latitude: redBaseLoc.x, longitude: redBaseLoc.y), radius: gameState.getCurrentGame().getRedBaseRad())
        addRadiusCircle(location: CLLocation(latitude: blueBaseLoc.x, longitude: blueBaseLoc.y), radius: gameState.getCurrentGame().getBlueBaseRad())
        
        // next, we need pins for the pickups in the game
        let pickupAnnotations = gameState.getCurrentGame().getPickupAnnotations()
        game_map.addAnnotations(pickupAnnotations)
        
        // now we set pins for each player in the game
        gameState.getCurrentGame().updatePlayerAnnotations()
        let playerAnnotations = gameState.getCurrentGame().getPlayerAnnotations()
        game_map.addAnnotations(playerAnnotations)
        
        // finally, we draw the games boundaries
        addBoundary()
        
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
    
    /* Status bar variables */
    @IBOutlet weak var armorBar: UIProgressView!
    @IBOutlet weak var healthBar: UIProgressView!
    @IBOutlet weak var clipBar: UIProgressView!
    
    /* viewWillAppear() - override */
    /* Set status bar values */
    override func viewWillAppear(_ animated: Bool) {
        armorBar.layer.cornerRadius = 3.0
        armorBar.layer.borderWidth = 3.0
        armorBar.layer.borderColor = #colorLiteral(red: 0.8078431487, green: 0.02745098062, blue: 0.3333333433, alpha: 1)
        healthBar.layer.cornerRadius = 3.0
        healthBar.layer.borderWidth = 3.0
        healthBar.layer.borderColor = #colorLiteral(red: 0.3333333433, green: 0.3333333433, blue: 0.3333333433, alpha: 1)
        clipBar.layer.cornerRadius = 3.0
        clipBar.layer.borderWidth = 3.0
        clipBar.layer.borderColor = #colorLiteral(red: 0.7254902124, green: 0.4784313738, blue: 0.09803921729, alpha: 1)
    }
    
    /* talkShitGetHit() */
    /* Updates the status bar fills and alerts the user to when they are attacked */
    private func talkShitGetHit() {
        if status.0 < gameState.getUserHealth() || status.1 < gameState.getUserArmor() {
            status = (gameState.getUserHealth(), gameState.getUserArmor())
            AudioServicesPlaySystemSound(kSystemSoundID_Vibrate)
//            let alertController = UIAlertController(title: "Temp", message: "You were hit", preferredStyle: UIAlertControllerStyle.alert)
//            alertController.addAction(UIAlertAction(title: "Ouch", style: UIAlertActionStyle.default,handler: nil))
//            present(alertController, animated: true, completion: nil)
        }
        armorBar.setProgress(Float(self.status.1)/100, animated: true)
        healthBar.setProgress(Float(self.status.0)/100, animated: true)
    }
    
    /* Upper display labels */
    @IBOutlet weak var clock: UILabel!
    @IBOutlet weak var redScore: UILabel!
    @IBOutlet weak var blueScore: UILabel!
    
    /* tick() - Timer function */
    /* Updates the clock according to the start time value sent by the server */
    private func tick() {
        let curtime = Date()
        
        if(curtime < starttime) { // game has not started yet
            let diff = calendar.dateComponents([.minute, .second], from: curtime, to: starttime)
            clock.text = "-" + String(format: "%02d", diff.minute!) + ":" + String(format: "%02d", diff.second!)
            statusLabel.text = "Head to the \(player.getTeam()) team's base!"
            statusLabel.textColor = #colorLiteral(red: 0, green: 0, blue: 0, alpha: 1)
        }
        else if(curtime > starttime.addingTimeInterval(Double(game.getTimeLimit()) * 60.0)) { // time up
            clock.text = "00:00"
        }
        else { // game started
            let diff = calendar.dateComponents([.minute, .second], from: curtime, to: starttime.addingTimeInterval(Double(game.getTimeLimit()) * 60.0))
            clock.text = String(format: "%02d", diff.minute!) + ":" + String(format: "%02d", diff.second!)
            if statusLabel.textColor == #colorLiteral(red: 0, green: 0, blue: 0, alpha: 1) {
                statusLabel.text = ""
            }
        }
    }
    
    /* Timer variables */
    private var clipTimer = Timer()
    private var reloading = false
    private var timeLeft: TimeInterval = 0
    
    /* fireButton() - Action for when fire button is pressed */
    /* Sends a fire message to the server if the weapon is not currently being reloaded */
    @IBAction func fireButton(_ sender: UIButton) {
        if !reloading {
            if gameState.getConnection().sendData(data: FireMsg()).isSuccess {
                let weapon = gameState.getUser().getWeapon()
                if weapon.clipFill > 1 {
                    weapon.clipFill -= 1
                    clipBar.setProgress(Float(weapon.clipFill)/Float(weapon.clipSize), animated: true)
                } else if weapon.clipFill == 1 {
                    reloading = true
                    weapon.clipFill -= 1
                    timeLeft = weapon.clipReload
                    clipBar.setProgress(Float(weapon.clipFill)/Float(weapon.clipSize), animated: true)
                    clipTimer = Timer.scheduledTimer(timeInterval: 1/5, target: self,   selector: (#selector(GameViewController.clipUpdate)), userInfo: nil, repeats: true)
                } else {
                    reloading = true
                }
            }
        }
    }
    
    /* clipUpdate() - Timer function */
    /* Update the clip's status bar based on how full it is */
    @objc func clipUpdate() {
        timeLeft -= 1/5
        let weapon = gameState.getUser().getWeapon()
        weapon.clipFill += Int(0.2 * Float(weapon.clipSize))
        clipBar.setProgress(Float(weapon.clipFill)/Float(weapon.clipSize), animated: true)
        if weapon.clipFill == weapon.clipSize {
            reloading = false
            clipTimer.invalidate()
        }
    }
    
    var updateTimer = Timer()
    
    /* runTimer() */
    /* Start the update timer */
    func runTimer() {
        updateTimer = Timer.scheduledTimer(timeInterval: 0.25, target: self,   selector: (#selector(GameViewController.update)), userInfo: nil, repeats: true)
    }
    
    /* update() - Timer Function */
    /* Sends a location update message to the server and checks various quicly changing values */
    @objc func update() {
        if gameState.getConnection().sendData(data: LocUpMsg()).isSuccess {
            print(gameState.getUser().getLocation())
        }
        
        redScore.text = "\(game.getRedPoints())"
        blueScore.text = "\(game.getBluePoints())"
        tick()
        
        if game.getGameOver() {
            endGame()
        }
        
        handleStatus()
        
        // update the locations of other players on the map and the status of the pickups
        gameState.getCurrentGame().updatePlayerAnnotations()
        gameState.getCurrentGame().updatePickupAnnotations()
        
        // update objective status in the case that the owner changes
        let objectives = gameState.getCurrentGame().getObjectives()
        for objective in objectives {
            if objective.getRedraw() {
                game_map.removeAnnotation(objective.getAnnotation())
                objective.updateAnnotation()
                game_map.addAnnotation(objective.getAnnotation())
            }
        }
    }
    
    @IBOutlet weak var statusLabel: UILabel!
    
    /* handleStatus() - Timer function */
    /* Controls the status values being displayed */
    private func handleStatus() {
        switch player.getStatus() {
        case "OutOfBounds":
            statusLabel.text = "Out of Bounds!"
            statusLabel.textColor = randomColor()
            break
        case "BackInBounds":
            statusLabel.text = ""
            break
        case "Respawn":
            statusLabel.text = "Return to Base"
            statusLabel.textColor = randomColor()
            player.setHealth(to: 0)
            break
        case "Respawned":
            statusLabel.text = ""
            healthBar.setProgress(1, animated: true)
            break
        default:
            break
        }
    }
    
    /* endGame() */
    /* Processes the end of game when it receives the end of game message */
    private func endGame() {
        let victory = (game.getBluePoints() > game.getRedPoints()) ? "Blue" : "Red"
        updateTimer.invalidate()
        let actionController = UIAlertController(title: "Game Over", message:
            "\(victory) team victory!", preferredStyle: UIAlertControllerStyle.actionSheet)
        actionController.addAction(UIAlertAction(title: "Let me leave", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            self.performSegue(withIdentifier: "leaveGame", sender: nil)
        }))
        actionController.addAction(UIAlertAction(title: "I wanna stay", style: UIAlertActionStyle.cancel,handler: nil))
        self.present(actionController, animated: true, completion: nil)
    }
    
    /* leaveGame() - Action on leave game button */
    /* Checks to make sure a player wants to leave the game early and then sends the leave game message to the server */
    @IBAction func leaveGame(_ sender: UIButton) {
        let actionController = UIAlertController(title: nil, message:
            "Are you sure you want to leave?", preferredStyle: UIAlertControllerStyle.actionSheet)
        actionController.addAction(UIAlertAction(title: "Yep", style: UIAlertActionStyle.default,handler: {(alert: UIAlertAction!) -> Void in
            if gameState.getConnection().sendData(data: LeaveGameMsg()).isSuccess {
                self.performSegue(withIdentifier: "leaveGame", sender: nil)
            }
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
