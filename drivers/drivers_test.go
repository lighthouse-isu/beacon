package drivers

import (
    "testing"

    "github.com/lighthouse/beacon/structs"

    "github.com/lighthouse/beacon/drivers/unknown"
)


func IsGood() bool {
    return true
}

func IsBad() bool {
    return false
}

var DriverA = &structs.Driver {
    Name: "a",
    IsApplicable: IsGood,
}

var DriverB = &structs.Driver {
    Name: "b",
    IsApplicable: IsBad,
}

var DriverC = &structs.Driver {
    Name: "c",
    IsApplicable: IsGood,
}


func TestFind(t *testing.T) {
    result := Find("b", []*structs.Driver{
        DriverA,
        DriverB,
        DriverC,
    })

    if result != DriverB {
        t.Fail()
    }
}

func TestFindUnkown(t *testing.T) {
    result := Find("something that's not real", []*structs.Driver{
        DriverA,
        DriverB,
        DriverC,
    })

    if result != unknown.Driver {
        t.Fail()
    }
}

func TestGuess(t *testing.T) {
    result := Guess([]*structs.Driver{
        DriverA,
        DriverB,
        DriverC,
    })
    if result != DriverA {
        t.Fail()
    }

    result = Guess([]*structs.Driver{
        DriverB,
        DriverB,
        DriverC,
    })
    if result != DriverC {
        t.Fail()
    }
}

func TestGuessUknown(t *testing.T) {
    result := Guess([]*structs.Driver{
        DriverB,
        DriverB,
        DriverB,
    })
    if result != unknown.Driver {
        t.Fail()
    }
}

func TestDecide(t *testing.T) {
    Defaults = []*structs.Driver{
        DriverA,
        DriverB,
        DriverC,
    }

    result := Decide()
    if result != DriverA {
        t.Fail()
    }

    *Preferred = "c"

    result = Decide()
    if result != DriverC {
        t.Fail()
    }
}
