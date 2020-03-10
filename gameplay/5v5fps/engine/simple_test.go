package engine

import (
	"sync"
	"testing"
	"time"

	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
	"github.com/stretchr/testify/assert"
)

func fillTestData(e *SimpleEngine) []Player {
	wp1 := weapon.SecondaryOne
	p1 := Player{UserID: 1, HP: 100, Team: util.Team1, SecondaryWeapon: &wp1}
	p1.EquippedWeapon = p1.SecondaryWeapon
	wp2 := weapon.SecondaryOne
	p2 := Player{UserID: 2, HP: 100, Team: util.Team1, SecondaryWeapon: &wp2}
	p2.EquippedWeapon = p2.SecondaryWeapon
	wp3 := weapon.SecondaryOne
	p3 := Player{UserID: 3, HP: 100, Team: util.Team1, SecondaryWeapon: &wp3}
	p3.EquippedWeapon = p3.SecondaryWeapon
	wp4 := weapon.SecondaryOne
	p4 := Player{UserID: 4, HP: 100, Team: util.Team1, SecondaryWeapon: &wp4}
	p4.EquippedWeapon = p4.SecondaryWeapon
	wp5 := weapon.SecondaryOne
	p5 := Player{UserID: 5, HP: 100, Team: util.Team1, SecondaryWeapon: &wp5}
	p5.EquippedWeapon = p5.SecondaryWeapon
	wp6 := weapon.SecondaryOne
	p6 := Player{UserID: 6, HP: 100, Team: util.Team2, SecondaryWeapon: &wp6}
	p6.EquippedWeapon = p6.SecondaryWeapon
	wp7 := weapon.SecondaryOne
	p7 := Player{UserID: 7, HP: 100, Team: util.Team2, SecondaryWeapon: &wp7}
	p7.EquippedWeapon = p7.SecondaryWeapon
	wp8 := weapon.SecondaryOne
	p8 := Player{UserID: 8, HP: 100, Team: util.Team2, SecondaryWeapon: &wp8}
	p8.EquippedWeapon = p8.SecondaryWeapon
	wp9 := weapon.SecondaryOne
	p9 := Player{UserID: 9, HP: 100, Team: util.Team2, SecondaryWeapon: &wp9}
	p9.EquippedWeapon = p9.SecondaryWeapon
	wp10 := weapon.SecondaryOne
	p10 := Player{UserID: 10, HP: 100, Team: util.Team2, SecondaryWeapon: &wp10}
	p10.EquippedWeapon = p10.SecondaryWeapon
	e.Init(p1)
	e.Init(p2)
	e.Init(p3)
	e.Init(p4)
	e.Init(p5)
	e.Init(p6)
	e.Init(p7)
	e.Init(p8)
	e.Init(p9)
	e.Init(p10)

	return []Player{p1, p2, p3, p4, p5, p6, p7, p8, p9, p10}
}

func getPlayers(e *SimpleEngine) []Player {
	players := make([]Player, len(e.players))
	for k, v := range e.players {
		players[k-1] = *v
	}
	return players
}

func TestInit(t *testing.T) {
	e := NewSimpleEngine()

	players := fillTestData(e)
	actualPlayers := getPlayers(e)

	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveLeft(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveLeft(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{-1, 0, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveLeftUp(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveLeftUp(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{-1, 1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveLeftDown(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveLeftDown(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{-1, -1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveRight(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveRight(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{1, 0, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveRightUp(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveRightUp(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{1, 1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveRightDown(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveRightDown(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{1, -1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveUp(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUp(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{0, 1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveUpLeft(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUpLeft(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{-1, 1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveUpRight(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUpRight(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{1, 1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveDown(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveDown(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{0, -1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveDownLeft(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveDownLeft(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{-1, -1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestMoveDownRight(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveDownRight(1)

	expectedPlayer := players[0]
	expectedPlayer.Position = util.Vector3{1, -1, 0}
	players[0] = expectedPlayer

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestShootNoAmmo(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	ammoCount := players[0].EquippedWeapon.Ammo
	for i := 0; i < ammoCount; i++ {
		e.Shoot(1, util.Vector3{-1, 0, 0})
	}

	e.MoveLeft(2)
	e.Shoot(1, util.Vector3{-1, 0, 0})

	updatedP1 := players[0]
	updatedP1.EquippedWeapon.Ammo = 0
	updatedP2 := players[1]
	updatedP2.Position = util.Vector3{-1, 0, 0}
	players[0] = updatedP1
	players[1] = updatedP2

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestShootMiss(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()
	e.MoveLeft(2)

	e.Shoot(1, util.Vector3{-2, 0, 0})

	updatedP1 := players[0]
	updatedP1.EquippedWeapon.Ammo--
	updatedP2 := players[1]
	updatedP2.Position = util.Vector3{-1, 0, 0}
	players[0] = updatedP1
	players[1] = updatedP2

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestShootHit(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()
	e.MoveLeft(2)

	e.Shoot(1, util.Vector3{-1, 0, 0})

	updatedP1 := players[0]
	updatedP1.EquippedWeapon.Ammo--
	updatedP2 := players[1]
	updatedP2.Position = util.Vector3{-1, 0, 0}
	updatedP2.HP = updatedP2.HP - int(updatedP1.EquippedWeapon.Damage)
	players[0] = updatedP1
	players[1] = updatedP2

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestShootTeamKill(t *testing.T) {

}

func TestShootKill(t *testing.T) {

}

func TestAll(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	actualPlayers := e.All(1)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestPickupWeaponDoesNotExist(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.PickupWeapon(1, 10)
	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestPickupWeaponNone(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.PickupWeapon(1, 2)
	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestPickupWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	wp := weapon.PrimaryOne
	e.players[2].PrimaryWeapon = &wp
	e.players[2].EquippedWeapon = e.players[2].PrimaryWeapon
	e.Start()
	e.MoveLeft(2)
	droppedWeapon := e.players[2].PrimaryWeapon
	e.DropWeapon(2)
	e.MoveRight(2)

	e.MoveLeft(1)
	e.PickupWeapon(1, 2)

	updatedP1 := players[0]
	updatedP1.Position = util.Vector3{-1, 0, 0}
	updatedP1.PrimaryWeapon = droppedWeapon
	updatedP2 := players[1]
	updatedP2.EquippedWeapon = nil
	players[0] = updatedP1
	players[1] = updatedP2

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

// func TestChannel(t *testing.T) {
// 	e := NewSimpleEngine()
// 	ch := e.Channel()
// 	go func() {
// 		for range ch{}
// 	}()

// 	players := fillTestData(e)
// }

// func TestState(t *testing.T) {
// 	e := NewSimpleEngine()
// 	ch := e.Channel()
// 	go func() {
// 		for range ch{}
// 	}()

// 	players := fillTestData(e)
// }

func TestDropWeaponNone(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].SecondaryWeapon = nil
	e.players[1].EquippedWeapon = nil
	players = getPlayers(e)
	e.Start()

	e.DropWeapon(1)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestDropWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.DropWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = nil
	updatedP1.SecondaryWeapon = nil
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchPrimaryWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].PrimaryWeapon = &weapon.PrimaryOne
	players = getPlayers(e)
	e.Start()

	e.SwitchPrimaryWeapon(1)
	e.SwitchPrimaryWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.PrimaryWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchPrimaryWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].PrimaryWeapon = &weapon.PrimaryOne
	players = getPlayers(e)

	e.SwitchPrimaryWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.PrimaryWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchSecondaryWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.SwitchSecondaryWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.SecondaryWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchSecondaryWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].KnifeWeapon = &weapon.KnifeOne
	e.players[1].EquippedWeapon = e.players[1].KnifeWeapon
	players = getPlayers(e)
	e.Start()

	e.SwitchSecondaryWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.SecondaryWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchKnifeWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].KnifeWeapon = &weapon.KnifeOne
	players = getPlayers(e)
	e.Start()

	e.SwitchKnifeWeapon(1)
	e.SwitchKnifeWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.KnifeWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSwitchKnifeWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].KnifeWeapon = &weapon.KnifeOne
	players = getPlayers(e)
	e.Start()

	e.SwitchKnifeWeapon(1)

	updatedP1 := players[0]
	updatedP1.EquippedWeapon = updatedP1.KnifeWeapon
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestStart(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	var chanMessages []interface{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					wg.Done()
				}
				chanMessages = append(chanMessages, msg)
			case <-time.After(time.Second * 3):
				wg.Done()
			}
		}
	}()

	players := fillTestData(e)
	e.Start()

	wg.Wait()
	expectedChanMessages := []interface{}{
		MatchStart{},
		RoundStart{},
	}
	assert.Equal(t, expectedChanMessages, chanMessages)
	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestBuyNotEnoughMoney(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].Money = int(weapon.PrimaryOne.Price - 1)
	players = getPlayers(e)
	e.Start()

	e.Buy(1, weapon.PrimaryOne.ID)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestBuyWeaponDoesNotExist(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].Money = int(weapon.PrimaryOne.Price)
	players = getPlayers(e)
	e.Start()

	e.Buy(1, 10)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestBuy(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].Money = int(weapon.PrimaryOne.Price)
	e.Start()

	e.Buy(1, weapon.PrimaryOne.ID)

	updatedP1 := players[0]
	updatedP1.PrimaryWeapon = &weapon.PrimaryOne
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestSetOrientation(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.SetOrientation(1, util.Vector3{0, 90, 0})

	updatedP1 := players[0]
	updatedP1.Orientation = util.Vector3{0, 90, 0}
	players[0] = updatedP1

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestReloadFullAmmo(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.Reload(1)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestReloadKnifeWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.players[1].KnifeWeapon = &weapon.KnifeOne
	e.players[1].EquippedWeapon = e.players[1].KnifeWeapon
	players = getPlayers(e)
	e.Start()

	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Reload(1)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestReload(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch {
		}
	}()

	players := fillTestData(e)
	e.Start()

	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Shoot(1, util.Vector3{-1, 0, 0})
	e.Reload(1)

	actualPlayers := getPlayers(e)
	assert.ElementsMatch(t, players, actualPlayers)
}
