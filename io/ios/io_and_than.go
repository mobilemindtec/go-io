package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/types"
	"log"
	"reflect"
)

type IOAndThan[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect types.IOEffect
	f          func() types.IORunnable
	debug      bool
	debugInfo  *types.IODebugInfo
}

func NewAndThan[A any](f func() types.IORunnable) *IOAndThan[A] {
	return &IOAndThan[A]{f: f}
}

func (this *IOAndThan[A]) Lift() *types.IO[A] {
	return types.NewIO[A]().Effects(this)
}

func (this *IOAndThan[A]) TypeIn() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAndThan[A]) TypeOut() reflect.Type {
	return reflect.TypeFor[A]()
}

func (this *IOAndThan[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAndThan[A]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAndThan[A]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAndThan[A]) String() string {
	return fmt.Sprintf("AndThan(%v)", this.value.String())
}

func (this *IOAndThan[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAndThan[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAndThan[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAndThan[A]) UnsafeRun() types.IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else {
			runnableIO := this.f()
			this.value = runtime.New[A](runnableIO).UnsafeRun()
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}