package dogma

// #cgo LDFLAGS: -ldogma
// #include <dogma.h>
// #include <dogma-extra.h>
//
// int dogma_get_charge_attributes_batch(dogma_context_t* ctx, dogma_key_t index, dogma_attributeid_t attributeid[],
//                                       double out[], int size) {
//   int i;
//   for (i = 0; i < size; ++i) {
//     double v;
//     dogma_get_charge_attribute(ctx, index, attributeid[i], &v);
//     out[i] = v;
//   }
//   return DOGMA_OK;
// }
//
// int dogma_validate(dogma_context_t* ctx, double* power, double *cpu) {
//   dogma_get_power_left(ctx, power);
//   dogma_get_cpu_left(ctx, cpu);
//   return DOGMA_OK;
// }
//
// int dogma_get_power_left(dogma_context_t* ctx, double* out) {
//   double power_output, power_left;
//   dogma_get_ship_attribute(ctx, 11, &power_output);
//   dogma_get_ship_attribute(ctx, 15, &power_left);
//   *out = power_output - power_left;
//   return DOGMA_OK;
// }
//
// int dogma_get_cpu_left(dogma_context_t* ctx, double* out) {
//   double cpu_output, cpu_left;
//   dogma_get_ship_attribute(ctx, 48, &cpu_output);
//   dogma_get_ship_attribute(ctx, 49, &cpu_left);
//   *out = cpu_output - cpu_left;
//   return DOGMA_OK;
// }
import "C"

// https://blog.golang.org/c-go-cgo

const (
	StateOnline     = 17
	StateActive     = 31
	StateOverloaded = 63
)

type DogmaContext struct {
	ctx  *C.dogma_context_t
	mods []ModIdx
}

type TypeId uint32
type AttributeId uint16
type ModIdx uint8
type AttributeValue float32

func (c *DogmaContext) Init() {
	if err := C.dogma_init_context(&c.ctx); err != 0 {
		panic("oops")
	}
}

func (c *DogmaContext) SetShip(t TypeId) {
	if err := C.dogma_set_ship(c.ctx, C.dogma_typeid_t(t)); err != 0 {
		panic("oops")
	}
}

func (c *DogmaContext) AddModule(t TypeId) ModIdx {
	var i C.dogma_key_t
	if err := C.dogma_add_module_s(c.ctx, C.dogma_typeid_t(t), &i, StateOnline); err != 0 {
		panic("oops")
	}
	idx := ModIdx(i)
	c.mods = append(c.mods, idx)
	return idx
}

func (c *DogmaContext) AddCharge(t TypeId, idx ModIdx) {
	if err := C.dogma_add_charge(c.ctx, C.dogma_key_t(idx), C.dogma_typeid_t(t)); err != 0 {
		panic("oops")
	}
}

func (c *DogmaContext) GetShipAttribute(t AttributeId) AttributeValue {
	var value C.double
	if err := C.dogma_get_ship_attribute(c.ctx, C.dogma_attributeid_t(t), &value); err != 0 {
		panic("oops")
	}
	return AttributeValue(value)
}

func (c *DogmaContext) GetModuleAttribute(t AttributeId, i ModIdx) AttributeValue {
	var value C.double
	if err := C.dogma_get_module_attribute(c.ctx, C.dogma_key_t(i), C.dogma_attributeid_t(t), &value); err != 0 {
		panic("oops")
	}
	return AttributeValue(value)
}

func (c *DogmaContext) GetChargeAttribute(t AttributeId, i ModIdx) AttributeValue {
	var value C.double
	if err := C.dogma_get_charge_attribute(c.ctx, C.dogma_key_t(i), C.dogma_attributeid_t(t), &value); err != 0 {
		panic("oops")
	}
	return AttributeValue(value)
}

func (c *DogmaContext) GetChargeAttributes(attrIds []AttributeId, idx ModIdx) []AttributeValue {
	size := len(attrIds)
	attrs := make([]C.dogma_attributeid_t, size)
	for i := 0; i < size; i++ {
		attrs[i] = C.dogma_attributeid_t(attrIds[i])
	}
	values := make([]C.double, size)
	if err := C.dogma_get_charge_attributes_batch(c.ctx, C.dogma_key_t(idx), &attrs[0], &values[0], C.int(size)); err != 0 {
		panic("oops")
	}
	results := make([]AttributeValue, size)
	for i := 0; i < size; i++ {
		results[i] = AttributeValue(values[i])
	}
	return results
}

func (c *DogmaContext) Destroy() {
	if err := C.dogma_free_context(c.ctx); err != 0 {
		panic("oops")
	}
}

func (c *DogmaContext) Strip() {
	for _, i := range c.mods {
		c.RemoveModule(i)
	}
	c.mods = c.mods[:0]
}

func (c *DogmaContext) RemoveModule(idx ModIdx) {
	if err := C.dogma_remove_module(c.ctx, C.dogma_key_t(idx)); err != 0 {
		panic("oops")
	}
}

func Init() {
	if err := C.dogma_init(); err != 0 {
		panic("oops")
	}
}

// performance optimization functions

func (c *DogmaContext) PowerLeft() AttributeValue {
	var value C.double
	if err := C.dogma_get_power_left(c.ctx, &value); err != 0 {
		panic("oops")
	}
	return AttributeValue(value)
}

func (c *DogmaContext) CpuLeft() AttributeValue {
	var value C.double
	if err := C.dogma_get_cpu_left(c.ctx, &value); err != 0 {
		panic("oops")
	}
	return AttributeValue(value)
}

func (c *DogmaContext) Validate() (AttributeValue, AttributeValue) {
	var power, cpu C.double
	if err := C.dogma_validate(c.ctx, &power, &cpu); err != 0 {
		panic("oops")
	}
	return AttributeValue(power), AttributeValue(cpu)
}
