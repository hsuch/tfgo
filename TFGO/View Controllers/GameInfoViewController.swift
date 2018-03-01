//
//  GameInfoViewController.swift
//  TFGO
//
//  Created by Sam Schlang on 2/12/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import UIKit
import MapKit
import CoreLocation

class GameInfoViewController: UITableViewController, UICollectionViewDelegate, UICollectionViewDataSource, CLLocationManagerDelegate, MKMapViewDelegate {
    
    @IBOutlet weak var game_map: MKMapView!
    @IBOutlet weak var gameNameLbl: UILabel!
    @IBOutlet weak var gamemodeLbl: UILabel!
    @IBOutlet weak var descriptionLbl: UILabel!
    @IBOutlet weak var statusLbl: UILabel!
    @IBOutlet weak var minutesLbl: UILabel!
    @IBOutlet weak var pointsLbl: UILabel!
    
    private var game = gameState.getCurrentGame()
    
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
    
    func mapView(_ mapView: MKMapView, rendererFor overlay: MKOverlay) -> MKOverlayRenderer! {
        if overlay is MKPolygon {
            let polygonView = MKPolygonRenderer(overlay: overlay)
            polygonView.strokeColor = UIColor.purple
            polygonView.lineWidth = 2.0
            
            return polygonView
        }
        
        return nil
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        gameNameLbl.text = game.getName()
        switch game.getMode() {
        case .cp:
            gamemodeLbl.text = "◆ - Standard"
        case .payload:
            gamemodeLbl.text = "⇥ - Payload"
        case .multi:
            gamemodeLbl.text = "❖ - Multipoint"
        }
        if game.getDescription() == "" {
            descriptionLbl.text = "None"
        } else {
            descriptionLbl.text = game.getDescription()
        }
        statusLbl.text = "Waiting for players - [\(game.getPlayers().count)/\(game.getMaxPlayers())]"
        minutesLbl.text = "\(game.getTimeLimit()) Minutes"
        pointsLbl.text = "\(game.getMaxPoints()) Points"
        
        game_map.delegate = self
        
        // now we have to set up the GameInfo Map
        manager.delegate = self
        manager.desiredAccuracy = kCLLocationAccuracyBest
        manager.requestWhenInUseAuthorization()
        manager.startUpdatingLocation()
        manager.startUpdatingHeading()
        
        // used to differentiate between objectives if we have more than one
        var objectiveNumber = 1
        
        // now we make a pin on the map for each objective
        for objective in game.getObjectives() {
            let annotation = MKPointAnnotation()
            annotation.coordinate = CLLocationCoordinate2D(latitude: objective.getXLoc(), longitude: objective.getYLoc())
            annotation.title = "OBJECTIVE"
            annotation.subtitle = String(objectiveNumber)
            game_map.addAnnotation(annotation)
            
//            let center = CLLocation(latitude: objective.getXLoc(), longitude: objective.getYLoc())
//            addRadiusCircle(location: center, radius: objective.getRadius())
            
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
        
        // finally, we draw the boundaries
        addBoundary()
    }
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return game.getPlayers().count
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "Icon", for: indexPath) as! IconViewCell
        cell.label.text = game.getPlayers()[indexPath.row].getIcon()
        cell.backgroundColor = randomColor()
        cell.layer.cornerRadius = 8.0
        cell.clipsToBounds = true
        return cell
    }
    
    override func shouldPerformSegue(withIdentifier identifier: String, sender: Any?) -> Bool {
        if identifier == "infoToWaiting" {
            if gameState.getConnection().sendData(data: JoinGameMsg(IDtoJoin: game.getID()!, password: "")).isSuccess {
                if handleMsgFromServer() {
                    if game.isValid() {
                        return true
                    } else {
                        let alertController = UIAlertController(title: "This game is no longer valid", message:
                            "Please select a different game", preferredStyle: UIAlertControllerStyle.alert)
                        alertController.addAction(UIAlertAction(title: "It was hubris to try and join", style: UIAlertActionStyle.default,handler: nil))
                        self.present(alertController, animated: true, completion: nil)
                    }
                }
            }
            return false
        }
        return true
    }
}
