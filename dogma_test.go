package dogma

import "testing"
import "github.com/stretchr/testify/assert"

func BenchmarkRifter(b *testing.B) {
	Init()
	for n := 0; n < b.N; n++ {
		var ctx DogmaContext

		ctx.Init()
		ctx.SetShip(587)         // Rifter
		i := ctx.AddModule(3831) // Medium Shield Extender II
		ctx.GetShipAttribute(263)
		ctx.GetModuleAttribute(50, i)
		ctx.Destroy()
	}
}

func BenchmarkRifterReuseCtx(b *testing.B) {
	Init()
	var ctx DogmaContext
	ctx.Init()
	ctx.SetShip(587) // Rifter
	for n := 0; n < b.N; n++ {
		i := ctx.AddModule(3831) // Medium Shield Extender II
		ctx.GetShipAttribute(263)
		ctx.GetModuleAttribute(50, i)
		ctx.Strip()
	}
}

func TestDogmaContext_PowerLeft(t *testing.T) {
	Init()
	var ctx DogmaContext
	ctx.Init()
	ctx.SetShip(587) // Rifter
	assert.EqualValues(t, 51.25, ctx.PowerLeft())
	ctx.AddModule(3831) // Medium Shield Extender II
	assert.EqualValues(t, 28.75, ctx.PowerLeft())
}

func TestDogmaContext_Validate(t *testing.T) {
	Init()
	var ctx DogmaContext
	ctx.Init()
	ctx.SetShip(587) // Rifter
	powerLeft, cpuLeft := ctx.Validate()
	assert.EqualValues(t, 51.25, powerLeft)
	assert.EqualValues(t, 162.5, cpuLeft)
}

func TestDogmaContext_GetChargeAttributes(t *testing.T) {
	Init()
	var ctx DogmaContext
	ctx.Init()
	ctx.SetShip(587) // Rifter
	i := ctx.AddModule(10631) // Rocket Launcher II
	ctx.AddCharge(2514 , i) // Inferno Rocket
	//DamageTypeEm:        dmgMultiplier * f.ReadChargeAttribute(114, idx),
	//DamageTypeKinetic:   dmgMultiplier * f.ReadChargeAttribute(117, idx),
	//DamageTypeExplosive: dmgMultiplier * f.ReadChargeAttribute(116, idx),
	//DamageTypeThermal:   dmgMultiplier * f.ReadChargeAttribute(118, idx),
	attrs := ctx.GetChargeAttributes([]AttributeId{114, 116, 117, 118}, i)
	assert.EqualValues(t, 0.0, attrs[0])
	assert.EqualValues(t, 0.0, attrs[1])
	assert.EqualValues(t, 0.0, attrs[2])
	assert.EqualValues(t, 45.375, attrs[3])
}