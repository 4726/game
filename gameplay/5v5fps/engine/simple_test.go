package engine

import (
	"testing"

	"github.com/4726/game/gameplay/5v5fps/util"
	"github.com/4726/game/gameplay/5v5fps/weapon"
	"github.com/stretchr/testify/assert"
)

func fillTestData(e *SimpleEngine) []Player {
	p1 := Player{UserID: 1, HP: 100, Team: util.Team1, SecondaryWeapon: &weapon.SecondaryOne}
	p1.EquippedWeapon = p1.SecondaryWeapon
	p2 := Player{UserID: 2, HP: 100, Team: util.Team1, SecondaryWeapon: &weapon.SecondaryOne}
	p2.EquippedWeapon = p2.SecondaryWeapon
	p3 := Player{UserID: 3, HP: 100, Team: util.Team1, SecondaryWeapon: &weapon.SecondaryOne}
	p3.EquippedWeapon = p3.SecondaryWeapon
	p4 := Player{UserID: 4, HP: 100, Team: util.Team1, SecondaryWeapon: &weapon.SecondaryOne}
	p4.EquippedWeapon = p4.SecondaryWeapon
	p5 := Player{UserID: 5, HP: 100, Team: util.Team1, SecondaryWeapon: &weapon.SecondaryOne}
	p5.EquippedWeapon = p5.SecondaryWeapon
	p6 := Player{UserID: 6, HP: 100, Team: util.Team2, SecondaryWeapon: &weapon.SecondaryOne}
	p6.EquippedWeapon = p6.SecondaryWeapon
	p7 := Player{UserID: 7, HP: 100, Team: util.Team2, SecondaryWeapon: &weapon.SecondaryOne}
	p7.EquippedWeapon = p7.SecondaryWeapon
	p8 := Player{UserID: 8, HP: 100, Team: util.Team2, SecondaryWeapon: &weapon.SecondaryOne}
	p8.EquippedWeapon = p8.SecondaryWeapon
	p9 := Player{UserID: 9, HP: 100, Team: util.Team2, SecondaryWeapon: &weapon.SecondaryOne}
	p9.EquippedWeapon = p9.SecondaryWeapon
	p10 := Player{UserID: 10, HP: 100, Team: util.Team2, SecondaryWeapon: &weapon.SecondaryOne}
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
	var players []Player
	for _, v := range e.players {
		players = append(players, *v)
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUp(1)

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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUpLeft(1)

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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()

	e.MoveUpRight(1)

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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()
	e.MoveLeft(2)

	for i := 0; i < players[0].EquippedWeapon.Ammo {
		e.Shoot(1, util.Vector3{-2, 0, 0})
	}

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
		for range ch{}
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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()
	e.MoveLeft(2)

	e.Shoot(1, util.Vector3{-1, 0, 0})

	updatedP1 := players[0]
	updatedP1.EquippedWeapon.Ammo--
	updatedP2 := players[1]
	updatedP2.Position = util.Vector3{-1, 0, 0}
	updatedP2.HP = updatedP2.HP - players[0].EquippedWeapon.Damage
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
		for range ch{}
	}()

	players := fillTestData(e)
	e.Start()

	actualPlayers := e.All(1)
	assert.ElementsMatch(t, players, actualPlayers)
}

func TestPickupWeaponNone(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestPickupWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestChannel(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestState(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestDropWeaponNone(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestDropWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchPrimaryWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchPrimaryWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchSecondaryWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchSecondaryWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchKnifeWeaponAlready(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSwitchKnifeWeapon(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestStart(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestBuyNotEnoughMoney(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestBuyWeaponDoesNotExist(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestBuy(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestSetOrientation(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestReloadFullAmmo(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}

func TestReload(t *testing.T) {
	e := NewSimpleEngine()
	ch := e.Channel()
	go func() {
		for range ch{}
	}()

	players := fillTestData(e)
}