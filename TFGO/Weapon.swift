//
//  Weapon.swift
//  TFGO
//
//  Created by Sam Schlang on 2/28/18.
//  Copyright © 2018 University of Chicago. All rights reserved.
//

import Foundation


/* Weapon "Abstract" Class */
 /* Contains information for weapon firing and mapping weapons to their names */
 /* Initializes to the null weapon, which is overridden by its children */
class Weapon {
    var name: String
    var clipSize: Int
    var clipFill: Int
    var shotReload: TimeInterval
    var clipReload: TimeInterval
    
    init () {
        name = ""
        clipSize = 0
        clipFill = 0
        shotReload = 0
        clipReload = 0
    }
}

/* weaponByName(String) -> Weapon */
/* Returns a new weapon based on the filename of the weapon image */
func weaponByName(name: String) -> Weapon {
    switch name {
    case "Sword":
        return Sword()
    case "Shotgun":
        return Shotgun()
    case "Pistol":
        return Pistol()
    case "Blaster":
        return Blaster()
    case "Crossbow":
        return Crossbow()
    case"Rifle":
        return Rifle()
    case "Boomerange":
        return Boomerang()
    case "Lightsaber":
        return Lightsaber()
    case "Spear":
        return Spear()
    case "Banhammer":
        return Banhammer()
    default:
        return BeeSwarm()
    }
}

class Sword: Weapon {
    
    override init() {
        super.init()
        name = "Sword"
        clipSize = 50
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Shotgun: Weapon {
    
    override init() {
        super.init()
        name = "Shotgun"
        clipSize = 2
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Pistol: Weapon {
    
    override init() {
        super.init()
        name = "Pistol"
        clipSize = 6
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Blaster: Weapon {
    
    override init() {
        super.init()
        name = "Blaster"
        clipSize = 20
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Crossbow: Weapon {
    
    override init() {
        super.init()
        name = "Crossbow"
        clipSize = 4
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Rifle: Weapon {
    
    override init() {
        super.init()
        name = "Rifle"
        clipSize = 3
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Boomerang: Weapon {
    
    override init() {
        super.init()
        name = "Boomerang"
        clipSize = 1
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Lightsaber: Weapon {
    
    override init() {
        super.init()
        name = "Lightsaber"
        clipSize = 100
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Spear: Weapon {
    
    override init() {
        super.init()
        name = "Spear"
        clipSize = 75
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class Banhammer: Weapon {
    
    override init() {
        super.init()
        name = "BanHammer"
        clipSize = 200
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}

class BeeSwarm: Weapon {
    
    override init() {
        super.init()
        name = "BeeSwarm"
        clipSize = 5
        clipFill = clipSize
        shotReload = 1
        clipReload = 1
    }
}
